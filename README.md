# config

[![travis build](https://travis-ci.org/flowpl/config.svg)](https://travis-ci.org/flowpl/config)
[![Coverage Status](https://coveralls.io/repos/github/flowpl/config/badge.svg?branch=master)](https://coveralls.io/github/flowpl/config?branch=master)

A simple hierarchical config system that can read configuration from a multitude of sources.
Currently JSON files and environment variables are supported to load configuration from.

## Install

```shell
go get github.com/flowpl/config
```

## Dependencies

none

## Usage

#### Normal usage

```go

import github.com/flowpl/config

// create a config struct that reads variables from the environment
cnf := config.NewConfig("app", config.NewEnvironmentLoader())

// loads a value and convert it to bool from $APP_BOOLVARIABLENAME
var boolVariable bool, err error = cnf.GetBool("boolVariableName")

// loads a value and convert to string from $APP_STRING
var stringVariable string, err error = cnf.GetString("string")

// loads a value and converts to int from $APP_INT_VAR
var integer int, err error = cnf.GetInt("int_var")

// loads a value and converts to float from $APP_FLOAT_VAR
var f float, err error = cnf.GetFloat("Float_Var")

```

If the environment variable does not contain a value that can be converted into the data type requested, an error is returned.
There is a little bit of leeway, though.
- integers can be loaded as float
- valid 'true' values for bool are: "yes", "1", "true", "y". This is case insensitive.
- valid 'false' values for bool are: "no", "0", "false", "n". This is also case insensitive.

#### Reading variables in single value context

```go

import github.com/flowpl/config

// create a config struct that reads variables from the environment
cnf := config.NewConfig("app", config.NewEnvironmentLoader())

// these functions return the same values as their non-May counterparts but they don't return errors.
var boolVariable bool = cnf.MayGetBool("boolVariableName")
var stringVariable string = cnf.MayGetString("string")
var integer int = cnf.MayGetInt("int_var")
var f float = cnf.MayGetFloat("Float_Var")

```

#### Creating Config Hierarchies

```go

import github.com/flowpl/config

// create a new config with an environment loader
cnf := config.NewConfig("app", config.NewEnvironmentLoader())

level2 := cnf.ChildConfig("level2")

// load the value from $APP_LEVEL2_BOOL
var level2BoolVal bool, level2Err error := level2.GetBool("bool")

// load the value from $APP_BOOL
var boolVal bool, err error := cnf.GetBool("bool")

```

## Loaders

the config library supports configurable loaders to support a multitude of configuration sources.
Currently only JSON files and environment variables are supported.

#### EnvironmentLoader
```go
loader := config.NewEnvironmentLoader()
```

refer to the usage examples above to see how values are loaded from the environment

#### JSONLoader
```go
loader, err := config.NewJsonLoader("path/to/file.json")
```

Example JSON:

```json
{
    "key1": "value1",
    "level2": {
        "key2": 123.5
    }
}
```

Accessing these values with the config packages looks like this:

```go

import github.com/flowpl/config

jsonLoader, loaderErr := config.NewJsonLoader("path/to/file.json")
if loaderErr != nil {
    panic(loaderErr.Error()
}

cnf := config.NewConfig("app", jsonLoader)


// access key1
key1 := cnf.MayGetString("key1")

// access key1
level2 := cnf.ChildConfig("level2")
key2 := level2.GetFloat("key2")

```


## Error handling

When using the normal interface to read variables error handling should by pretty straight forward.
However when using the single value interface errors are obviously not available directly. 
Fortunately the config library stores all errors that happened while loading variables. 
The error list can be retrieved by calling `cnf.GetErrors()`.
 
## panicking

In most cases it is not sensible to start or continue an application if errors happened while loading configuration.
The config package provides the convenient `cnf.TriggerErrorPanic()` that can be called after interacting with config and 
panics if config errors occurred and does nothing otherwise. 
The panic message contains all error messages concatenated into one string.

## License

MIT. See LICENSE
