package {{replace .Name "_" "."}}

type {{camelCase1 .Name}} struct {
{{range .Fields}}
	{{camelCase1 .Name}} {{.Type}}
{{end}}
}