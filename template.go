package main

const template = `
{{define "INVOCATION"}}type {{.Name}}Invocation struct {{"{"}}{{if .Params}}
	Parameters struct {
{{range .Params}}		{{.StructDef}}
{{end}}	}{{end}}{{if .Results}}
	Results struct {
{{range .Results}}		{{.StructDef}}
{{end}}	}{{end}}
}

{{end}}

{{define "FAKE"}}type Fake{{.Name}} struct {
{{range .Methods}}	{{.Name}}Hook func({{.FormatParamsDeclaration}}) ({{.FormatResultsDeclaration}})
{{end}}
{{range .Methods}}	{{.Name}}Calls []*{{.Name}}Invocation
{{end}}}

{{end}}

{{define "METHOD"}}func (a *Fake{{.InterfaceName}}) {{.Name}}({{.FormatParamsDeclaration}}) ({{.FormatResultsDeclaration}}) {
	invocation := new({{.Name}}Invocation)

{{if .Params}}{{range .Params}}	invocation.Parameters.{{.CapitalName}} = {{.Name}}
{{end}}{{end}}
{{if .Results}}	{{.FormatResultsCall}} = a.{{.Name}}Hook({{.FormatParamsCall}})
{{else}}	a.{{.Name}}Hook({{.FormatParamsCall}})
{{end}}
{{if .Results}}	{{range .Results}}invocation.Results.{{.CapitalName}} = {{.Name}}
{{end}}{{end}}
	a.{{.Name}}Calls = append(a.{{.Name}}Calls, invocation)

	return {{.FormatResultsCall}}
}

{{end}}
`
