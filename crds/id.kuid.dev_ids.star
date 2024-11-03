def getEndpointID(nodeID, port, endpoint, name):
  return {
    "partition": nodeID.get("partition", ""),
    "region": nodeID.get("region", ""),
    "site": nodeID.get("site", ""),
    "node": nodeID.get("node", ""),
    "port": str(int(port)),
    "endpoint": str(int(endpoint)),
    "name": name,
  }

def genEndpointIDString(epID):
  nodeIDstr = genNodeIDString(epID)

  epIDstr = nodeIDstr + "." + str(int(epID.get("port", 0))) + "." + str(int(epID.get("endpoint", 0)))
  epShortName = epID.get("name", "")
  if epShortName == "" or epShortName == "interface":
    return epIDstr 
  return epIDstr + "." + epShortName

def getEndpointKeys():
  return {
    "partition": True,
    "region": True,
    "site": True,
    "node": True,
    "port": True,
    "endpoint": True,
    "name": True,
  }