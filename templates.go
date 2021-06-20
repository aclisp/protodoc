package main

import (
	"bytes"
	"log"
	"text/template"
)

const ProtoTpl = `
{{- /* ------------------------------------------------------------- */ -}}

{{define "service"}}
## Service {{.ServiceName}}

{{.Comment}}
{{range .Infs}}
### Method {{.ServiceName}}.{{.MethodName}}

{{.Comment}}

Request
{{template "fields" .Req.Params}}
Response
{{template "fields" .Res.Params}}
{{end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

{{define "enum"}}
### enum {{.Name}}

{{.Comment}}

Constants

|   Value   |   Name    |   Comment    |
| --------- | --------- | ------------ |
{{- range .Constants}}
| {{.Val}}  | {{.Name}} | {{.Comment}} |
{{- end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

{{- define "object"}}
### object {{.Name}}

{{.Comment}}

Attributes
{{template "fields" .Attrs}}
{{- end}}

{{- /* ------------------------------------------------------------- */ -}}

{{define "fields"}}
|   Name    |   Type    |   Comment    |
| --------- | --------- | ------------ |
{{- range .}}
| {{.Name}} | {{.Type}} | {{.Comment}} |
{{- end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

# API Protocol

{{range .Services}}
{{template "service" .}}
{{end}}

## Enums

{{- range .Enums}}
{{template "enum" .}}
{{- end}}

## Objects

{{- range .Objects}}
{{template "object" .}}
{{- end}}
`

func (pf ProtoFile) generateMarkdown() string {
	proto := template.Must(template.New("proto").Parse(ProtoTpl))
	buf := bytes.Buffer{}
	if err := proto.Execute(&buf, pf); err != nil {
		log.Panicf("failed to execute template: %v", err)
	}
	return buf.String()
}
