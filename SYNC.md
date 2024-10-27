# sync startegy

- create/delete index lock/unlock
    - when creating/deleting an index -> claims are called which cause a lock issue
        - we have a single lock on index create/delete and claim create/delete
    - we introduce a special dryrin flag since the delete has very few options if []string{"recursion"} is set we treat this as a special behavior

## differences

1. deleteIndex
- sync: claims and entries get deleted 
- async: claims from index get deleted, but other not -> reconciler should trigger the status

2. restart app -> restore
- sync/async: create index need to be triggered (reconciler needed)

3. status of index
- sync/async: reconciler update status with status subresource (reconciler needed)

4. order is important
- sync: when something fails the client need to handle this
- async: reconciler tries to recover

5. 