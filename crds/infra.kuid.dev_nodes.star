load("id.kuid.dev.ids.star", "getNodeKeys")

def getNodeKeys():
  return {
    "partition": True,
    "region": True,
    "site": True,
    "node": True,
  }

def getNodeID(partition, region, site, node):
  return {
    "partition": partition,
    "region": region,
    "site": site,
    "node": node,
  }

def getNodeID(self):
  nodeKeys = getNodeKeys()
  spec = getSpec(self)
  nodeID = {}
  for key, val in spec.items():
    if key in nodeKeys:
      nodeID[key] = val
  return nodeID

def genNodeIDString(nodeID):
  return nodeID["partition"] + "." + nodeID["region"] + "." + nodeID["site"] + "." + nodeID["node"]


def getSpec(self):
  return self.get("spec", {})

def getPartition(self):
  spec = getSpec(self)
  return spec.get("partition", "")

def getProvider(self):
  spec = getSpec(self)
  return spec.get("provider", "")

def getPlatformType(self):
  spec = getSpec(self)
  return spec.get("platformType", "")

def getStatus(self):
  return self.get("status", {})

def getNode(name, namespace, spec):
  return {
    "apiVersion": "infra.kuid.dev/v1alpha1",
    "kind": "Node",
    "metadata": {
      "name": name,
      "namespace": namespace,
    },
    "spec": spec,
  }