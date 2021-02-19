# SAM Generator

AWS SAM boilerplate generator command line utility



## How to install

To install it, run:

```sh
  go get github.com/vsrc/samg
```

Command will download and compile the source code, then put it in the `bin` folder inside your `GOPATH`.



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
