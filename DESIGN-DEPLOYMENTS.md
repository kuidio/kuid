# deployment options

- async design with a generic backend -> reocncilers needed
- async design with etcd -> reconcilers needed
- sync design with a generic backend
- sync design with choreo


async versus sync:
- Async need a reconciler
- Sync can have a special way to allocate before -> create/update/delete are special; watch/list are the same
- Sync cannot be done with etcd as the storage interface changed
- Sync -> We do a special init for both options -> main reason is to use a direct storage interface for saveALL and restore functions
-> we init all 3 resources together due to the storage

## How to select which objects should be rendered?

- name: ENABLE_BE_AS (group)
  value: "sync,badgerdb" | "true" 
- name: ENABLE_BE_VLAN (group)
  value: "sync" | "true"
- name: ENABLE_BE_IPAM (group)
  value: "sync" | "true"

## Select between sync and async also per group

-> using ENV flags


## open

- FieldSelector would be nice when listing the items to be able to check if they belong to the index
- Update? How to handle
- Conversion function for CRD(s) - how does this work ???? -> this will determine if we can keep choreo aligned or not

Choreo -> we need a conversion function conversion function


TODO:
- align the rest interface from choreo with the rest interface of 
