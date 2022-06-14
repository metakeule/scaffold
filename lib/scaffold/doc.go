// Copyright (c) 2015 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package scaffold provides file and directory generation based on templates.

A template must be UTF8 without byte order marker and have \n (linefeed) as line terminator.
It has a head and a body, separated by an empty line:

    1. head (must not contain an empty line)
    2. empty line
    3. body

The head might contain anything but empty lines, but it is recommended to put some annotated
json string as example for the usage into it. Also authorship of the template and contact infos can be
put there.

The syntax of the body is a superset of the Go text/template package (http://golang.org/pkg/text/template).
The available functions inside the body are extended by the functions defined in the FuncMap variable.

Additionally to the functionality provide via the text/template package there are contexts.

Context

The template body can have folder and file contexts.
A context is started by a line with the prefix ">>>" and ends by a line with the prefix "<<<".
Each context has a name that follows after ">>>" or "<<<" until the end of the line.
For example the context ">>>a" is ended by "<<<a" and the name of the context would be "a".

If the name of the context ends with a slash (/), the context is a folder context otherwise it
is a file context.

The name of a file context defines the name of the file into which the content
of the context will be saved. The folder of the file is defined by the surrounding folder contexts.
The name of a folder context defines the name of the folder inside which the inner folders and files (as defined
by the inner contexts) are saved.

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

The placeholders inside the body are organized as a json object / map. When the Run function is called, the
json objects is mixed to the template and after that the folders and files are created as defined in the
result. That makes it possible to use placeholders as parts of folder or file names.

Escaping of double curly braces and dollar chars

Curly braces and dollar chars are part of syntax of the go template engine and there
is no syntax to replace them from inside the template. However, the scaffold package has
included some helper functions that return them.

   {{doubleCurlyOpen}}{{dollar}}{{doubleCurlyClose}}

will result in the string "{{$}}".

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
