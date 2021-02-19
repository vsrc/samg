package main

import (
	"fmt"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
)

// FnNode struct for template yaml config file
type FnNode struct {
  Type string `yaml:"Type,omitempty"`
  Properties Properties `yaml:"Properties,omitempty"`
}

// Properties struct child node for StackNode and FnNode in template yaml config file
type Properties struct {
  Location  string `yaml:"Location,omitempty"`
  CodeURI   string `yaml:"CodeUri,omitempty"`
  Handler   string `yaml:"Handler,omitempty"`
  Runtime   string `yaml:"Runtime,omitempty"`
  Tracing   string `yaml:"Tracing,omitempty"`
  Events    Events `yaml:"Events,omitempty"`
}

// Events struct child node for Properties in template yaml config file
type Events struct {
  CatchAll CatchAll `yaml:"CatchAll,omitempty"`
}

// CatchAll struct child node for Events in template yaml config file
type CatchAll struct {
  Type string `yaml:"Type,omitempty"`
  Properties EventProperties `yaml:"Properties,omitempty"`
}

// EventProperties struct child node for CatchAll in template yaml config file
type EventProperties struct {
  Path    string `yaml:"Path,omitempty"`
  Method  string `yaml:"Method,omitempty"`
}

func fnCommand() *cli.Command {

  var path string
  var parentPath string
  var templatePath string
  var name string
  var url string
  var method string

  fn := cli.Command{
      Name:    "function",
      Aliases: []string{"fn"},
      Usage:   "create new aws lambda function",
      Action:  func(c *cli.Context) error { 

        path = strings.ReplaceAll(path + "/", "//", "/") 

        if templatePath == "" {
          templatePath = "templates/fn/"
        }

        // create new folder on the path given for this new function [needs path for new function]
        err := mkdirWith(path)
        if err != nil {
          return err
        }

        // copy go.mod.txt file into that new folder function
        // and rename go.mod.txt to go.mod
        err = copyFileContents(templatePath + "/go.mod.txt", path + "/go.mod")
        if err != nil {
          return err
        }


        // if name for module is not provided, generate it
        if name == "" {
          name = generateModuleName(path)
        }

        // edit go.mod.txt file - change MODULE_NAME [needs name for new function]
        err = replaceStringInFile(path + "/go.mod", "{MODULE_NAME}", name)

        
        // copy main.go.txt file into that new folder function
        // and rename main.go.txt to main.go
        err = copyFileContents(templatePath + "/main.go.txt", path + "/main.go")
        if err != nil {
          return err
        }
        
        // path to where parent configuration is      
        if parentPath == "" {
          parentPath, err = getDefaultParentPath(path)
          if err != nil {
            return err
          }
        }

        // url for api
        if url == "" {
          url = path
        }

        // url must start with leading slash
        url = strings.ReplaceAll("/" + url, "//", "/")

        // http method for this fn
        if method == "" {
          method = "GET"
        }

        // location of code in relation to the parent config file
        parentChildPath, err := filepath.Rel(parentPath, path)
        if err != nil {
          return err
        }
        codeLocation := strings.ReplaceAll(parentChildPath + "/", "//", "/") 

        // generate new configuration [needs parent stack path, url
        // ...  for api (optional), http method (optional, default get)]
        newNode := generateFnConfiguration(codeLocation, name, url, method)

        childNodekName := generateChildNodekNameWith(path, "Fn")
        // add information about new function into the parent stack [needs parent stack path, url
        err = updateParentWith(parentPath, childNodekName, newNode)
        if err != nil {
          return err
        }

        fmt.Println("new function template created at path " + path)

        return nil

      },
  }

  fn.Flags = []cli.Flag {
    &cli.StringFlag{
      Name: "path",
      Aliases: []string{"p"},
      Usage: "path where new function will be created (required)",
      Destination: &path,
      Required: true,
    },
    &cli.StringFlag{
      Name: "parent",
      Aliases: []string{"a"},
      Usage: "path to folder where parent stack configuration file is in which we will put the reference to this new function",
      Destination: &parentPath,
      DefaultText: "parent folder of the provided path",
    },
    &cli.StringFlag{
      Name: "template",
      Aliases: []string{"t"},
      Usage: "path to template folder for this new function",
      Destination: &templatePath,
      DefaultText: "templates/fn/",
    }, 
    &cli.StringFlag{
      Name: "name",
      Aliases: []string{"n"},
      Usage: "optional name for this function",
      Destination: &name,
      DefaultText: "last part of provided path",
    }, 
    &cli.StringFlag{
      Name: "url",
      Aliases: []string{"u"},
      Usage: "url for the api for this function",
      Destination: &url,
      DefaultText: "same as provided path",
    }, 
    &cli.StringFlag{
      Name: "method",
      Aliases: []string{"m"},
      Usage: "http method for this function",
      Destination: &method,
      DefaultText: "GET",
    },
  }

  return &fn
}

func generateFnConfiguration(path string, name string, url string, method string) FnNode {
  return FnNode {
    Type: "AWS::Serverless::Function",
    Properties: Properties {
      CodeURI: path,
      Handler: name,
      Runtime: "go1.x",
      Tracing: "Active",
      Events: Events {
        CatchAll: CatchAll{
          Type: "Api",
          Properties: EventProperties{
            Path: url,
            Method: method,
          },
        },
      },
    },
  }
}
