# deployment options

- asynchronous design with a generic backend
- synchronous design with a generic backend
- synchronous design with choreo


async versus sync:
- We could initialize the storage in such way that we use a special method for storing the objects -> useful for the backend object like ipam/as/vlan/etc
- Create/Update/Delete is Handled special
- List and Watch handled in the traditional way

## How to select which objects should be rendered?

        - name: ENABLE_BE_AS (group)
          value: "sync" | "true"
        - name: ENABLE_BE_VLAN (group)
          value: "sync" | "true"
        - name: ENABLE_BE_IPAM (group)
          value: "sync" | "true"

## Select between sync and async also per group

-> using ENV flags



## sync versus async

1. We do a special init for both options -> main reason is to use a direct storage interface for saveALL and restore functions
-> we init all 3 resources together due to the storage

2. The sync or async is decided on init



## store interface 

-> for tests it is better to use a regular storage interface
-> memory for tests
-> badgerdb or other for a real deployment

What about config maps ??


## open

- FieldSelector would be nice when listing the items to be able to check if they belong to the index
- Update? How to handle
- Conversion function for CRD(s) - how does this work ???? -> this will determine if we can keep the resources aligned or not

to be tested:
- getting storage
- how to walk over a list ? unstructured ?
- field manager is handled by
