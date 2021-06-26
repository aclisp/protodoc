package main

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

const ProtoTpl = `
{{- /* ------------------------------------------------------------- */ -}}

{{define "service"}}
## Service {{.ServiceName}}

{{.Comment}}
{{range .Infs}}
### Method {{.ServiceName}}.{{.MethodName}}
{{if .IsWebSocket}}
WebSocket {{.Typ}}
{{end}}
> {{.HTTPMethod}} {{.URLPath}} <br/>
{{- if not .IsWebSocket}}
> Content-Type: application/json <br/>
> Authorization: Bearer (token) <br/>
{{- end}}

{{.Comment}}

{{if .Req.Empty}}Request is empty
{{else}}Request parameters
{{template "fields" .Req.Params}}{{end}}
{{if .Res.Empty}}Response is empty
{{else}}Response parameters
{{template "fields" .Res.Params}}{{end}}
{{end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

{{define "enum"}}
### enum {{.Name}}

{{.Comment}}

Constants

|   Value   |   Name    |  Description |
| --------- | --------- | ------------ |
{{- range .Constants}}
| {{.Val}}  | {{.Name}} | {{.Comment}} |
{{- end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

{{- define "object"}}
### object {{.Name}}

{{.Comment}}

{{if .Empty}}It has no attributes
{{else}}Attributes
{{template "fields" .Attrs}}{{end}}
{{- end}}

{{- /* ------------------------------------------------------------- */ -}}

{{define "fields"}}
|   Name    |   Type    |  Description |
| --------- | --------- | ------------ |
{{- range .}}
| {{.Name}} | {{.TypeHRef}} | {{.Comment}} |
{{- end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

{{define "toc"}}
{{- range .Services}}
* [Service {{.ServiceName}}]({{.HRef}})
{{- range .Infs}}
    * [Method {{.ServiceName}}.{{.MethodName}}]({{.HRef}})
{{- end}}
{{- end}}
* [Enums](#enums)
{{- range .Enums}}
    * [Enum {{.Name}}]({{.HRef}})
{{- end}}
* [Objects](#objects)
{{- range .Objects}}
    * [Object {{.Name}}]({{.HRef}})
{{- end}}
{{end}}

{{- /* ------------------------------------------------------------- */ -}}

# API Protocol

Table of Contents
{{template "toc" .}}

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

// HRef generates a cross reference ID used in markdown
func (s Service) HRef() string {
	return "#service-" + strings.ToLower(s.ServiceName)
}

func (e Endpoint) HRef() string {
	return "#method-" + strings.ToLower(e.ServiceName) + strings.ToLower(e.MethodName)
}

func (o Object) HRef() string {
	return "#object-" + href(o.Name)
}

func (e Enum) HRef() string {
	return "#enum-" + href(e.Name)
}

func (r Request) Empty() bool {
	return len(r.Params) == 0
}

func (r Response) Empty() bool {
	return len(r.Params) == 0
}

func (o Object) Empty() bool {
	return len(o.Attrs) == 0
}

func (e Endpoint) IsWebSocket() bool {
	return e.Typ != Unary
}
