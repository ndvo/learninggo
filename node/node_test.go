package node_test

import (
  "testing"
  "learninggo/node"
  "regexp"
  "os"
)

// Test if GetPathId returns a uniq id fitting regex (\w{8}/){7}\w{8}
func TestGetPathId(t *testing.T){
  n:= node.Node{
  }
  id, err := n.GetPathId()
  if err != nil{
    t.Error(err)
  }
  r := regexp.MustCompile(`^(\w{8}/){7}\w{8}$`)
  if !r.MatchString(id) {
    t.Error("Id return does not conform to pattern")
  }
  if _, err := os.Stat(node.ContentFolder+"/"+id); !os.IsNotExist(err) {
    t.Error("Trying to create existing id")
  }
}

// Test if node.Create works appropriatelly
func TestCreate(t *testing.T){
  // Read some random data
  f, _ := os.Open("/dev/urandom")
  b := make([]byte, 32)
  f.Read(b)
  id, err:= node.Create("Nelson", b )
  if err!= nil{
    t.Error("Could not create node "+id+" - "+err.Error())
  }
  if _, err := os.Stat(node.ContentFolder+"/"+id); os.IsNotExist(err) {
    t.Error("Content was not created: "+err.Error())
  }else{
    os.Remove(node.ContentFolder+"/"+id)
  }
}

