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
The available functions inside the template are extended by the functions defined in the FuncMap variable.

Additionally to the functionality provide via the text/template package there are contexts.

Context

A context inside a template is a folder or a file context.
A context is started by a line with the prefix ">>>" and ends by a line with the prefix "<<<".
Each context has a name that follows after ">>>" or "<<<" until the end of the line.
For example the context ">>>a" is ended by "<<<a" and the name of the context would be "a".

If the name of the context ends with a slash (/), the context is a folder context otherwise it
is a file context.

The name of a file context defines the name of the file into which the content
of the context will be saved. The folder of the file is defined by the surrounding folder contexts.
The outermost folder context is the baseDir parameter of the Run function (defaults to the current working directory in the CLI tool).

The following would create the file "fileZ.txt" inside the folder "[baseDir]/folder1/folderA". Any missing directories a created
on the fly.

    >>>folder1/
    >>>folderA/
    >>>fileZ.txt
    Hello World
    <<<fileZ.txt
    <<<folderA/
    <<<folder1/

The placeholders inside the template are organized as a json object. When the Run function is called, the
json objects is first applied to the template and then the folders and files are created as defined in the
applied template. That makes it possible to use placeholders as parts of folder or file names.

Most of the time this package will be used via the scaffold command sub package.

It can be installed via

  go get gopkg.in/metakeule/scaffold.v1/cmd/scaffold

Run

	scaffold help

to see the available options.

For a complete example have a look at the example directory.
The make.sh file contains the needed CLI command.
*/
package scaffold
