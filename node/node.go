package node

import (
  "os"
	"fmt"
  "strings"
	"time"
	"encoding/gob"
)

const ContentFolder string = "content"

// node is any content in the website
// they can be saved and loaded
type Node struct {
	Title    string
	PathId   string
	Modifyed time.Time
	Creation time.Time
	Filetype string
	Content  []byte
}

// Write changes the content and saves node to disk
func (s *Node) Write(d []byte) ( err error) {
	s.Content = d
	s.Modifyed = time.Now()
  err = s.Save()
  return
}

//getPathId returns an Id that can be used as path
func (s *Node) GetPathId() (pathId string, err error) {
	if s.PathId == "" {
		f, err := os.Open("/dev/urandom")
		if err != nil {
			return pathId, err
		}else{
      b := make([]byte, 32)
      f.Read(b)
      s.PathId = fmt.Sprintf("%x/%x/%x/%x/%x/%x/%x/%x",
        b[0:4], b[4:8], b[8:12], b[12:16], b[16:20], b[20:24], b[24:28],b[28:32])
    }
	}
  pathId = s.PathId
	return pathId, err
}

// Save saves a node to disk
func (s *Node) Save() error {
  d := strings.Join( strings.Split(s.PathId, "/")[0:7], "/")
  os.MkdirAll(ContentFolder+"/"+d, 0775)
	f, err := os.Create(ContentFolder+"/"+s.PathId)
	if err != nil {
		return err
	}
	defer f.Close()
  enc:=gob.NewEncoder(f)
  err=enc.Encode(s)
	if err != nil {
		return err
	}
	return nil
}

// Create a new node named title with content content
// pathId is created automatically and returned
func Create(title string, content []byte) (pathId string, err error) {
  node := Node{
    Title    : title,
    Modifyed : time.Now(),
    Creation : time.Now(),
    Filetype : "gob",
  }
  pathId, err = node.GetPathId()
  if err != nil {
    return "", err
  }else{
    err := node.Write(content)
    if err != nil {
      return pathId, err
    }
  }
  return pathId, nil
}

// Gzip content before saving
//fz := gzip.NewWriter(f)
//defer fz.Close()
//enc := gob.NewEncoder(fz)
//err = enc.Encode(s.addr)
