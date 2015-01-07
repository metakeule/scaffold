# scaffold

scaffolding via go templates

# Usage

Create a template e.g.

```
{ // here is the help / example area
    Models: [
        {
            Name: "",
            Fields: [
                {Name: "", Type: ""}
            ]
        }
    ]
} // end of the help area at the first empty line

{{range .Models}}
>>>models/
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
<<<models/
{{end}}
```

The real template starts after the first empty line. All above that line is just help text.

Inside the real template placeholders in go template syntax may be used and the following functions are available:

```go
var FuncMap = template.FuncMap{
    "replace":    Replace,
    "camelCase1": CamelCase1,
    "camelCase2": CamelCase2,
    "title":      strings.Title,
    "toLower":    strings.ToLower,
    "toUpper":    strings.ToUpper,
    "trim":       strings.Trim,
}
```

The template is used to create files and folders.
Each line that starts with `>>>` indicates that a file or a folder context begins. If the line ends with a slash `/` it is a folder context, otherwise a file.

Everything following is part of the file or folder up to a line starting with `<<<` that closes the file/folder.

Files and folders are always relative to the surrounding folder contexts.

The scaffolding takes place by using a template like this and apply some json string on it. In this example the json string would be something like this

```json
{ 
    Models: [
        {
            Name: "first_model",
            Fields: [
                {Name: "first_field", Type: "string"},
                {Name: "second_field", Type: "int"}
            ]
        },
        {
            Name: "second_model",
            Fields: [
                ...
            ]
        }
        ...
    ]
}
```

The expected structure of the template should be written inside the help block of the template.

If we take the following json input:

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

and then save the template as `models.templ` and the json input as `models.json`. Then after we have installed scaffold via `go get github.com/metakeule/scaffold/cmd/scaffold` we can run

```sh
cat models.json | scaffold models.templ
```

and then the files `models/person/model.go` and `models/address/model.go` are generated below the current working directory and any missing directories are created.

See example directory to recap what you have learned.

Simply by creating and sharing template files it should be possible to generate a fairly amount of boilerplate.

To show the help message, run `scaffold models.templ help`.
So to write the input json you could simply run `scaffold models.templ help > my_models.json` and then editing the generated json file.

Values in the input json should always be **snake-case**. They could easily be transformed to camel-case where needed via the provided `camelCase1` and `camelCase2` functions inside the template.
