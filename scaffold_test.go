package scaffold

import (
	"bytes"
	"strings"
	"testing"
)

func TestCamelCase1(t *testing.T) {

	tests := []struct {
		input, expected string
	}{
		{"field_name", "FieldName"},
		{"FieldName", "FieldName"},
		{"Field_name", "FieldName"},
		{"field_Name", "FieldName"},
		{"fieldname", "Fieldname"},
	}

	for _, test := range tests {

		if got, want := CamelCase1(test.input), test.expected; got != want {
			t.Errorf("CamelCase1(%v) = %v; want %v", test.input, got, want)
		}
	}

}

func TestCamelCase2(t *testing.T) {

	tests := []struct {
		input, expected string
	}{
		{"field_name", "fieldName"},
		{"FieldName", "FieldName"},
		{"Field_name", "FieldName"},
		{"field_Name", "fieldName"},
		{"fieldname", "fieldname"},
	}

	for _, test := range tests {

		if got, want := CamelCase2(test.input), test.expected; got != want {
			t.Errorf("CamelCase2(%#v) = %#v; want %#v", test.input, got, want)
		}
	}

}

func TestReplace(t *testing.T) {

	tests := []struct {
		input, src, dest, expected string
	}{
		{"field_name", "_", "*", "field*name"},
		{"field_name_x", "_", "*", "field*name*x"},
		{"_field_name_", "_", "*", "*field*name*"},
	}

	for _, test := range tests {

		if got, want := Replace(test.input, test.src, test.dest), test.expected; got != want {
			t.Errorf("Replace(%#v, %#v, %#v) = %#v; want %#v", test.input, test.src, test.dest, got, want)
		}
	}

}

var validHead = `{
	"Models": [
		{
			"Name": "",
			"Fields": [
				{"Name": "", "Type": ""}
			]
		}
	]
}`

var validBody = `{{range .Models}}
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
{{end}}`

var validJSON = `
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
`

var validTemplate = validHead + "\n\n" + validBody

func TestSplitTemplate(t *testing.T) {
	h, b := SplitTemplate(validTemplate)

	if h != validHead {
		t.Errorf("wrong head: %#v, expecting %#v", h, validHead)
	}

	if b != validBody {
		t.Errorf("wrong body: %#v, expecting %#v", b, validBody)
	}
}

func TestRun(t *testing.T) {

	tests := []struct {
		dir, body, json, expected string
	}{
		{
			"start",
			">>>file.txt\n<<<file.txt",
			`{}`,
			"start/file.txt\n",
		},
		{
			"start",
			">>>file1.txt\n<<<file1.txt\n>>>file2.txt\n<<<file2.txt",
			`{}`,
			"start/file1.txt\nstart/file2.txt\n",
		},
		{
			"start",
			">>>a/\n>>>b/\n>>>file1.txt\n<<<file1.txt\n>>>file2.txt\n<<<file2.txt\n<<<b/\n<<<a/\n",
			`{}`,
			"start/a/b/file1.txt\nstart/a/b/file2.txt\n",
		},
		{
			"a/dir",
			"{{range .Files}}>>>{{.Name}}.txt\n<<<{{.Name}}.txt\n{{end}}",
			`{"Files": [{"Name": "file1"},{"Name": "file2"}]}`,
			"a/dir/file1.txt\na/dir/file2.txt\n",
		},
		{
			"start/dir",
			validBody,
			validJSON,
			"start/dir/models/person/model.go\nstart/dir/models/address/model.go\n",
		},
	}

	for _, test := range tests {
		var log bytes.Buffer
		err := Run(test.dir, test.body, strings.NewReader(test.json), &log, true)
		if err != nil {
			t.Errorf("Run(%#v, %#v, %#v,...) returned error: %v", test.dir, test.body, test.json, err)
			continue
		}
		if got, want := log.String(), test.expected; got != want {
			t.Errorf("Run(%#v, %#v, %#v,...) = %#v; want %#v", test.dir, test.body, test.json, got, want)
		}
	}
}

func TestRunErrors(t *testing.T) {

	tests := []struct {
		dir, body, json string
	}{
		{
			"start",
			">>>file.txt\n<<<file.txt",
			`{`,
		},
		{
			"start",
			">>file.txt\n<<<file.txt",
			`{}`,
		},
		{
			"start",
			">>>file1.txt\n<<<file2.txt",
			`{}`,
		},
		{
			"start",
			">>>a/\n>>>b\n>>>file1.txt\n<<<file1.txt\n<<<a/\n<<<b\n",
			`{}`,
		},
		{
			"start",
			">>>file1.txt\nho\n>>>file2.txt\nhu<<<file2.txt\n<<<file1.txt\n",
			`{}`,
		},
		{
			"start",
			"{{range .x}}",
			`{}`,
		},
	}

	for _, test := range tests {
		var log bytes.Buffer
		err := Run(test.dir, test.body, strings.NewReader(test.json), &log, true)
		if err == nil {
			t.Errorf("Run(%#v, %#v, %#v,...) returned no error", test.dir, test.body, test.json)
		}
	}
}
