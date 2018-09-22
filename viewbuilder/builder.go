package viewbuilder

import(
"encoding/json"
"errors"
"fmt"
"io/ioutil"
"os"
"path/filepath"
)

var viewdir  string = "api"
var indexdir string = "index"
var contentdir string = "content"


type Explorer struct{
  journey_start int
}

type view struct{
  info   os.FileInfo
  source string
}

var Views = []view{}



// Explore walks through all directories searching untracked content
func (e Explorer) explore(d string){
  // checks everything in directory d
  //dircontents := os.Readdir
  //dirnames := os.Readdirnames
}

// FindViews returns jsons in api folder
func FindViews(dir string) (Views []view, err error){
  flist, err := ioutil.ReadDir(dir)
  for _, d := range flist{
    if !d.IsDir(){
      if filepath.Ext(d.Name()) == ".json" {
        abs, _:= filepath.Abs( filepath.Join(dir, d.Name()))
        Views = append( Views, view{ d, abs })
        println(d.Name())
      }
    }
  }
  return
}

// 
func IndexView(v view) (err error) {
  var content []byte
  vbytes, err := ioutil.ReadFile(v.source)
  var vjson interface{}
  if err == nil{
    err = json.Unmarshal(vbytes, vjson)
  }
  if err != nil{
    return err
  }
  if vjson, ok := vjson.(map[string]interface{}); ok{
    if val, ok := vjson["dir"]; ok && val != nil {
      if val, ok := val.(string); ok{
        // Restrict dirs to be searched to dir and its subdirectories
        content, err = findContentKeyValued(filepath.Join(contentdir, val), vjson)
      }else{
        err = errors.New("dir in view json is not a string")
      }
    }else{
      content, err = findContentKeyValued(contentdir, vjson)
    }
  }else{
    errors.New("View rule is not map[string]interface{}")
  }
  err = ioutil.WriteFile(filepath.Join(indexdir, filepath.Base(v.source)), content, 0666)
  return
}

// findContentKeyValued find content in dir that matches the rules
func findContentKeyValued(dir string, rules interface{}) (content []byte, err error){
  var viewbytes []byte
  flist, _ := ioutil.ReadDir(dir)
  for _, d := range flist{
    if !d.IsDir(){
      if filepath.Ext(d.Name()) == ".json" {
        abs, _ := filepath.Abs( filepath.Join(dir, d.Name()) )
        viewbytes, err = ioutil.ReadFile(abs)
        var vjson interface{}
        if err == nil{
          err = json.Unmarshal(viewbytes, vjson)
        }
        if err != nil{
          return content, err
        }
        match, err := matchKeyValued(vjson, rules)
        if err == nil{
          c, err := json.Marshal(match)
          if err == nil{
            content = append(content, c...)
          }
        }
      }
    }
  }
  return content, nil
}

// matchKeyValued object against rule. Returns rule filled with object values and an error if found.
func matchKeyValued(o interface{}, rule interface{}) (match interface{}, err error){
  object, ook := o.(map[string]interface{})
  rules, ruleok := rule.(map[string]interface{})
  errMessage := ""
  if !ook || !ruleok {
    return rule, errors.New("Object or Rule is not a map string interface{}")
  }
  for k, v := range rules{
    if k!= "dir"{
      if v != nil{
        if object[k] != v{
          errMessage += fmt.Sprintf("Object key %s does not match the rule.\n", k)
        }
      }else{
          rules[k] = object[k]
      }
    }
  }
  if errMessage != "" {
    return rule, errors.New(errMessage)
  }else{
    return rule, nil
  }
}
