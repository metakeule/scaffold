// Copyright (c) 2015 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package scaffold provides file and directory generation based on templates.

A template consists of 3 parts:

1. help section (must not contain an empty line)
2. an empty line
3. core template

The help section might contain anything but empty lines, but it is recommended to put some annotated
json string as example for the usage into it. Also authorship of the template and contact infos can be
put there.

The syntax of the core template is a superset of the Go text/template package (http://golang.org/pkg/text/template).
The following additional functions are made available inside the template

replace":    Replace,
camelCase1": CamelCase1,
camelCase2": CamelCase2,
title      corresponds to http://golang.org/pkg/strings/#Title
toLower    corresponds to http://golang.org/pkg/strings/#ToLower
toUpper    corresponds to http://golang.org/pkg/strings/#ToUpper
trim       corresponds to http://golang.org/pkg/strings/#Trim

with the added feature of files and directory markers.

Example for a template

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
*/
package scaffold
