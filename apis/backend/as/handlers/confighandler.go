package handlers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type ASStoreHandler struct {
}

func (r *ASStoreHandler) DryRunCreateFn(ctx context.Context, key types.NamespacedName, obj runtime.Object, dryrun bool) (runtime.Object, error) {
	/*
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return obj, err
		}
		targetKey, err := config.GetTargetKey(accessor.GetLabels())
		if err != nil {
			return obj, err
		}
		cfg := obj.(*config.Config)
		obj, _, err = r.Handler.SetIntent(ctx, targetKey, cfg, true, dryrun)
		if err != nil {
			return obj, err
		}
	*/
	return obj, nil
}
func (r *ASStoreHandler) DryRunUpdateFn(ctx context.Context, key types.NamespacedName, obj, old runtime.Object, dryrun bool) (runtime.Object, error) {
	/*
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return obj, err
		}
		targetKey, err := config.GetTargetKey(accessor.GetLabels())
		if err != nil {
			return obj, err
		}
		cfg := obj.(*config.Config)
		obj, _, err = r.Handler.SetIntent(ctx, targetKey, cfg, true, dryrun)
		if err != nil {
			return obj, err
		}
	*/
	return obj, nil
}
func (r *ASStoreHandler) DryRunDeleteFn(ctx context.Context, key types.NamespacedName, obj runtime.Object, dryrun bool) (runtime.Object, error) {
	/*
		accessor, err := meta.Accessor(obj)
		if err != nil {
			return obj, err
		}
		targetKey, err := config.GetTargetKey(accessor.GetLabels())
		if err != nil {
			return obj, err
		}
		cfg := obj.(*config.Config)
		return r.Handler.DeleteIntent(ctx, targetKey, cfg, dryrun)
	*/
	return obj, nil
}
