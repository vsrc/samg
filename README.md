# SAM Generator

AWS SAM boilerplate generator command line utility

Motivation behind creating this tool is to help automate repetitive steps a developer has to do when creating api with AWS lambda. In addition to already existing amazing aws sam cli tools, for each new function and stack, developer needs to:

1. Create new folder/s
2. Copy paste template code
3. Copy paste any commonly used code snippets
4. Add configuration for the new function / stack in `template.yaml`

This tool does all of that for you with one command. Two if you need to create a stack.

## How to install

To install it, run:

```sh
  go get github.com/vsrc/samg
```

Command will download and compile the source code, then put it in the `bin` folder inside your `GOPATH`.

### Requirement

In order to properly work this tool requires templates. You can download `templates` folder from from this repo and edit them if needed. These files are generated with aws sam cli sample command `sam init` and contain simple demo functionality. 

As you are developing your api, you might start using your own set libraries / packages / code snippets in most of your functions and it makes sense to add them to the template code so it gets copied into every new function you generate afterwards.


## How to use

To see all commands and options available run 

```sh
  samg -h
```

or for specific command

```sh
  samg [COMMAND] -h
```

### To create new stack

```sh
  samg stack [command options]
```

Options that you can add to this command:
```txt
  --path value, -p value      path where new stack will be created (required)
  --parent value, -a value    path to folder where parent stack configuration file is in which we will put the reference to this new stack (default: parent folder of the provided path)
  --template value, -t value  path to template file for this new stack (default: templates/stack/template.yaml)
  --help, -h                  show help (default: false)
```

For example command:

```sh
  samg stack -p v1/
```

will create new folder `v1` and put `template.yaml` file inside, ready to be populated.

NOTE: when deploying via aws sam cli, deployment will create new separate API for each stack. It makes sense to keep same logical units inside same stack. 


### To create new function
```sh
  samg function [command options]
```

Options that you can add to this command:
```txt
  --path value, -p value      path where new function will be created (required)
   --parent value, -a value    path to folder where parent stack configuration file is in which we will put the reference to this new function (default: parent folder of the provided path)
   --template value, -t value  path to template folder for this new function (default: templates/fn/)
   --name value, -n value      optional name for this function (default: last part of provided path)
   --url value, -u value       url for the api for this function (default: same as provided path)
   --method value, -m value    http method for this function (default: GET)
   --help, -h                  show help (default: false)
```

For example command:

```sh
  samg fn -p v1/user/register
```

will create new folder on the path `v1/user/register/` and put two template files (`go.mod` and `main.go`) inside, ready to be used.


### To test locally

Use aws sam cli as usual, navigate to the function directory and run:


```sh
  sam build # Build package before starting local api
```

```sh
  sam local start-api # Start local api
```


### To deploy

Use aws sam cli as usual, navigate to the folder where `template.yaml` is either for function or for stack and run:

```sh
  sam deploy --guided --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND
```

This command will guide you through first time configuration setup and if you choose to save configuration, for every subsequent run of this command you can omit `--guided` flag.
