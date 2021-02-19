package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)


func mkdirWith(path string) error {

  if _, err := os.Stat(path); os.IsNotExist(err) {
      os.MkdirAll(path, os.ModePerm)
  }

  return nil
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src string, dst string) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return
    }
    err = out.Sync()
    return
}

func replaceStringInFile(path string, old string, new string) error {
  input, err := ioutil.ReadFile(path)
  if err != nil {
          return err
  }

  newContent := strings.ReplaceAll(string(input), old, new) 

  err = ioutil.WriteFile(path, []byte(newContent), 0644)
  if err != nil {
          return err
  }

  return nil
}

func generateModuleName(path string) string {
  parts := strings.Split(path, "/")
  return parts[len(parts) - 2]
}

func getDefaultParentPath(childPath string) (parentPath string, err error) {
  parentPath, err = filepath.Rel("", childPath + "/../")
  return
}



func updateParentWith(parentPath string, childName string, childConfig interface{}) (err error) {

  doesParentExist, err := doesFileExist(parentPath + "/template.yaml")
  if err != nil {
    return
  }

  // if parent does not exist do nothing
  if (!doesParentExist) {
    return nil
  }


  parentFile, err := ioutil.ReadFile(parentPath + "/template.yaml")

  if err != nil {
    return
  }

  parentData := yaml.MapSlice{}
  yaml.Unmarshal(parentFile, &parentData)

  // fetching the node to change and its position in the slice
  nodePos := -1

  var node yaml.MapSlice
  for i, n := range parentData {
    if n.Key == "Resources" {
        node = n.Value.(yaml.MapSlice)
        nodePos = i
        break
    }
  }


  // handle parent files with no resources param
  if nodePos < 0 {

    node = yaml.MapSlice{}

    parentData = append(parentData, yaml.MapItem{
      Key: "Resources",
      Value: yaml.MapSlice{},
    })

    nodePos = len(parentData) - 1
  }

  // append new data to the node
  node = append(node, yaml.MapItem{
    Key: childName,
    Value: childConfig,
  })

  parentData[nodePos].Value = node

  updatedParent, err := yaml.Marshal(&parentData)
  if err != nil {
    return
  }

  // overwrite parent configuration
  err = ioutil.WriteFile(parentPath + "/template.yaml", updatedParent, 0644)
  if err != nil {
    return
  }
  
  return
}

func doesFileExist(path string) (bool, error) {

  if _, err := os.Stat(path); err == nil {

    // path exists
    return true, nil
  } else if os.IsNotExist(err) {

    // path does *not* exist
    return false, nil
  } else {

    // Schrodinger: file may or may not exist. See err for details.
    // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existenc
    return false, err
  }
}