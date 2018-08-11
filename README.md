# scaffold

scaffolding via go templates

[![Build Status Travis](https://secure.travis-ci.org/metakeule/scaffold.png)](http://travis-ci.org/metakeule/scaffold) [![Documentation](http://godoc.org/gopkg.in/metakeule/scaffold.v1?status.png)](http://godoc.org/gopkg.in/metakeule/scaffold.v1) 

Installation
============

`go get github.com/metakeule/scaffold/cmd/scaffold`

Usage
=====

given the following template (file `models.templ`)

```go
{
    "Models": [
        {
            "Name": "",
            "Fields": [
                {"Name": "", "Type": ""}
            ]
        }
    ]
}

>>>models/
{{range .Models}}
>>>{{toLower .Name}}/
>>>model.go
package {{replace .Name "_" "."}}

type {{camelCase1 .Name}} struct {
{{range .Fields}}
    {{camelCase1 .Name}} {{.Type}}
{{end}}
}

<<<model.go
<<<{{toLower .Name}}/
{{end}}
<<<models/

```

and the following json (file `models.json`)

```json
{
    "Models": [
        {
            "Name": "person",
            "Fields": [
                {
                    "Name": "first_name",
                    "Type": "string"
                },
                {
                    "Name": "last_name",
                    "Type": "string"
                }
            ]
        },
        {
            "Name": "address",
            "Fields": [
                {
                    "Name": "street_no",
                    "Type": "string"
                },
                {
                    "Name": "city",
                    "Type": "string"
                }
            ]
        }
    ]
}
```

when running the command

```sh
scaffold -t=models.templ < models.json
```

the following directory structure would be build:

```sh
models
├── address
│   └── model.go
└── person
    └── model.go
```

where `models/address/model.go` contains

```go
package address

type Address struct {

    StreetNo string

    City string

}

```

and `models/person/model.go` contains

```go
package person

type Person struct {

    FirstName string

    LastName string

}

```