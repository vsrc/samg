package main

import (
	"fmt"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
)

// StackNode struct for template yaml config file
type StackNode struct {
  Type string `yaml:"Type,omitempty"`
  Properties Properties `yaml:"Properties,omitempty"`
}

func stackCommand() *cli.Command {

  var path string
  var parentPath string
  var templatePath string

  stack := cli.Command{
      Name:    "stack",
      Aliases: []string{"s"},
      Usage:   "create new stack",
      Action:  func(c *cli.Context) error { 

        if templatePath == "" {
          templatePath = "templates/stack/template.yaml"
        }

        _ = mkdirWith(path)

        err := copyFileContents(templatePath, path + "/template.yaml")
        if err != nil {
          return err
        }

        if parentPath == "" {
          parentPath, err = getDefaultParentPath(path)
          if err != nil {
            panic(err)
          }
        }

        // edit parent template file to include this new template
        err = updateStackParentWith(parentPath, path)
        if err != nil {
          return err
        }

        fmt.Println("created new stack on path ", path)

        return nil

      },
  }

  stack.Flags = []cli.Flag {
    &cli.StringFlag{
      Name: "path",
      Aliases: []string{"p"},
      Usage: "path where new stack will be created (required)",
      Destination: &path,
      Required: true,
    },
    &cli.StringFlag{
      Name: "parent",
      Aliases: []string{"a"},
      Usage: "path to folder where parent stack configuration file is in which we will put the reference to this new stack",
      Destination: &parentPath,
      DefaultText: "parent folder of the provided path",
    },
    &cli.StringFlag{
      Name: "template",
      Aliases: []string{"t"},
      Usage: "path to template file for this new stack",
      Destination: &templatePath,
      DefaultText: "templates/stack/template.yaml",
    },
  }

  return &stack
}

func updateStackParentWith(parentPath string, childPath string) (err error) {

  newStackName := generateChildNodekNameWith(childPath, "Stack")

  parentChildPath, err := filepath.Rel(parentPath, childPath)
  
  if err != nil {
    panic(err)
  }

  parentChildFileLocation := strings.ReplaceAll(parentChildPath + "/template.yaml", "//", "/") 


  newNode := StackNode {
    Type: "AWS::Serverless::Application",
    Properties: Properties {
      Location: parentChildFileLocation,
    },
  }

  return updateParentWith(parentPath, newStackName, newNode)

}

func generateChildNodekNameWith(path string, suffix string) (name string) {
  // v1/user/ -> V1UserStack

  // var name string

  pathSlice := strings.Split(path, "/")

  for _, el := range pathSlice {
    name = name + strings.Title(strings.ToLower(el))
  }

  name = name + suffix
  
  return
}