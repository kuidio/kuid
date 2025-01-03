package generic

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/runtime"
)

func New(
	indexKind string,
	claimKind string,
	indexObjectFn func(runtime.Object) (backend.IndexObject, error),
	claimObjectFn func(runtime.Object) (backend.ClaimObject, error),
	entryObjectFn func(runtime.Object) (backend.EntryObject, error),
	entryFromCacheFn func(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject,
) bebackend.Backend {

	cache := bebackend.NewCache[*CacheInstanceContext]()

	return &be{
		cache:            cache,
		indexKind:        indexKind,
		claimKind:        claimKind,
		indexObjectFn:    indexObjectFn,
		claimObjectFn:    claimObjectFn,
		entryObjectFn:    entryObjectFn,
		entryFromCacheFn: entryFromCacheFn,
	}
}

type be struct {
	cache            bebackend.Cache[*CacheInstanceContext]
	m                sync.RWMutex
	indexKind        string
	claimKind        string
	indexObjectFn    func(runtime.Object) (backend.IndexObject, error)
	claimObjectFn    func(runtime.Object) (backend.ClaimObject, error)
	entryObjectFn    func(runtime.Object) (backend.EntryObject, error)
	entryFromCacheFn func(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject
	// added later
	//entryStorage *registry.Store
	//claimStorage *registry.Store
	bestorage BackendStorage
}

func (r *be) PrintEntries(ctx context.Context, index string) {
	entries, _ := r.listEntries(ctx, store.ToKey(index))
	for _, entry := range entries {
		uobj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(entry)
		fmt.Println("entry", uobj)
	}
}

func (r *be) AddStorageInterfaces(bes any) error {
	bestorage, ok := bes.(BackendStorage)
	if !ok {
		return fmt.Errorf("AddStorageInterfaces did not supply a generic BackendStorage interface, got: %s", reflect.TypeOf(bes).Name())
	}
	r.bestorage = bestorage
	return nil
}

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, obj runtime.Object) error {
	r.m.Lock()
	defer r.m.Unlock()
	index, err := r.indexObjectFn(obj)
	if err != nil {
		return err
	}
	ctx = bebackend.InitIndexContext(ctx, "create", index)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := index.GetKey()

	log.Debug("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// if the Cache is not initialized -> restore the cache
	// this happens upon initialization or backend restart
	if _, err := r.cache.Get(ctx, key); err != nil {
		// if it does not exist create the cache
		cacheInstanceCtx := NewCacheInstanceContext(index.GetTree(), index.GetType())
		r.cache.Create(ctx, key, cacheInstanceCtx)
	}

	if !r.cache.IsInitialized(ctx, key) {
		if err := r.restore(ctx, index); err != nil {
			log.Error("cannot restore index", "error", err.Error())
			index.SetConditions(condition.Failed(err.Error()))
			return err
		}
		log.Debug("restored")
		index.SetConditions(condition.Ready())
		obj = index

		if err := r.cache.SetInitialized(ctx, key); err != nil {
			return err
		}
	}
	log.Debug("update Index claims", "object", obj)
	return r.updateIndexClaims(ctx, index)
}

// DeleteIndex deletes a backend index
func (r *be) DeleteIndex(ctx context.Context, obj runtime.Object) error {
	r.m.Lock()
	defer r.m.Unlock()
	objidx, err := r.indexObjectFn(obj)
	if err != nil {
		return err
	}
	ctx = bebackend.InitIndexContext(ctx, "delete", objidx)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := objidx.GetKey()

	log.Debug("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// delete the data from the backend
	if err := r.destroy(ctx, key); err != nil {
		log.Error("cannot delete Index", "error", err.Error())
		return err
	}
	log.Debug("destroyed")
	r.cache.Delete(ctx, key)

	log.Debug("finished")
	return nil
}

func (r *be) Claim(ctx context.Context, obj runtime.Object, recursion bool) error {
	if !recursion {
		r.m.Lock()
		defer r.m.Unlock()
	}
	claim, err := r.claimObjectFn(obj)
	if err != nil {
		return err
	}

	ctx = bebackend.InitClaimContext(ctx, "create", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey())
	if err != nil {
		return err
	}
	if !r.cache.IsInitialized(ctx, claim.GetKey()) {
		return fmt.Errorf("cache not initialized")
	}

	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		return err
	}
	if err := a.Validate(ctx, claim); err != nil {
		return err
	}
	if err := a.Apply(ctx, claim); err != nil {
		return err
	}
	// store the resources in the backend
	if err := r.saveAll(ctx, claim.GetKey()); err != nil {
		return err
	}
	obj = claim
	return nil
}

func (r *be) Release(ctx context.Context, obj runtime.Object, recursion bool) error {
	if !recursion {
		r.m.Lock()
		defer r.m.Unlock()
	}
	claim, err := r.claimObjectFn(obj)
	if err != nil {
		return err
	}

	ctx = bebackend.InitClaimContext(ctx, "delete", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	cacheCtx, err := r.cache.Get(ctx, claim.GetKey())
	if err != nil {
		return err
	}
	if !r.cache.IsInitialized(ctx, claim.GetKey()) {
		return fmt.Errorf("cache not initialized")
	}

	a, err := getApplicator(ctx, cacheCtx, claim)
	if err != nil {
		return err
	}
	if err := a.Delete(ctx, claim); err != nil {
		return err
	}

	return r.saveAll(ctx, claim.GetKey())
}

func getApplicator(_ context.Context, cacheInstanceCtx *CacheInstanceContext, claim backend.ClaimObject) (Applicator, error) {
	claimType := claim.GetClaimType()
	var a Applicator
	switch claimType {
	case backend.ClaimType_DynamicID:
		a = &dynamicApplicator{name: string(claimType), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case backend.ClaimType_StaticID:
		a = &staticApplicator{name: string(claimType), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	case backend.ClaimType_Range:
		a = &rangeApplicator{name: string(claimType), applicator: applicator{cacheInstanceCtx: cacheInstanceCtx}}
	default:
		return nil, fmt.Errorf("invalid addressing, got: %s", string(claimType))
	}

	return a, nil
}
