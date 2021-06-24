package main

import (
	"log"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// a proto file is a parsing unit
type ProtoFile struct {
	// a proto file could have multiple service
	Services []Service
	// a proto file should have multiple object
	Objects []Object
	// a proto file should have multiple enum
	Enums []Enum
}

// a proto file could have multiple Service
type Service struct {
	// service only use Comment placed at the beginning
	Comment string
	// the package name of the proto file
	PackageName string
	// my name
	ServiceName string
	// a service has multiple endpoint
	Infs []Endpoint
}

// enndpoint is also called method or interface
type Endpoint struct {
	// the package name of the proto file
	PackageName string
	// the enclosing service name
	ServiceName string
	// my name
	MethodName string
	// the url path where api gateway resolves
	URLPath string
	// the http method, which is always POST
	HTTPMethod string
	// Comment placed at the beginning, as well as inline-Comment placed at the ending
	Comment string

	Typ RPCType
	Req Request
	Res Response
}

type RPCType int

const (
	Unary = iota
	ServerStreaming
	ClientStreaming
	BidirectionalStreaming
)

type Request struct {
	Params []Field
	Typ    string
}

type Response struct {
	Params []Field
	Typ    string
}

type Field struct {
	// Comment placed at the beginning, as well as inline-Comment placed at the ending
	Comment string
	Name    string
	Typ     string
	Repeat  bool
	// the Enclosing type name of this field
	Enclosing string
	// reference to the enclosing proto file
	protoFile *ProtoFile
}

// Object is user-defined field type
type Object struct {
	// only use Comment placed at the beginning
	Comment string
	// refers by field.typ
	Name string
	// attributes
	Attrs []Field
}

// Enum is user-defined type which has one of a pre-defined list of values
type Enum struct {
	// only use Comment placed at the beginning
	Comment string
	// refers by field.typ
	Name string
	// Constants
	Constants []EnumField
}

type EnumField struct {
	// Comment placed at the beginning, as well as inline-Comment placed at the ending
	Comment string
	Name    string
	Val     string
	// the Enclosing type name of this field
	Enclosing string
}

func (t RPCType) String() string {
	switch t {
	case Unary:
		return "unary"
	case ClientStreaming:
		return "client-streaming"
	case ServerStreaming:
		return "server-streaming"
	case BidirectionalStreaming:
		return "bidirectional-streaming"
	}
	return "unknown"
}

func (pf *ProtoFile) composeFrom(pp *parser.Proto) {
	// find all services in proto body
	for _, x := range pp.ProtoBody {
		if service, ok := x.(*parser.Service); ok {
			pf.addService(service, pp)
		}
	}
	pf.addObjectsAndEnums(pp)
}

func (pf *ProtoFile) addService(ps *parser.Service, pp *parser.Proto) {
	var s Service
	s.Comment = composeHeadComment(ps.Comments)
	s.PackageName = extractPackageName(pp)
	s.ServiceName = ps.ServiceName
	s.Infs = pf.composeInterfaces(s, ps, pp)

	pf.Services = append(pf.Services, s)
}

func extractComment(pc *parser.Comment) string {
	if pc == nil {
		return ""
	}
	return strings.TrimSpace(strings.Join(pc.Lines(), "\n"))
}

func composeHeadComment(pcs []*parser.Comment) string {
	ss := make([]string, 0, len(pcs))
	for _, pc := range pcs {
		s := extractComment(pc)
		ss = append(ss, s)
	}
	return strings.Join(ss, " ")
}

func composeHeadAndInlineComment(pcs []*parser.Comment, pic *parser.Comment, sep string) string {
	head := composeHeadComment(pcs)
	inline := extractComment(pic)
	if head == "" {
		return inline
	}
	if inline == "" {
		return head
	}
	return head + sep + inline
}

func extractPackageName(pp *parser.Proto) string {
	for _, x := range pp.ProtoBody {
		if p, ok := x.(*parser.Package); ok {
			return p.Name
		}
	}
	return "(missed-package)"
}

func (pf *ProtoFile) composeInterfaces(s Service, ps *parser.Service, pp *parser.Proto) []Endpoint {
	eps := make([]Endpoint, 0, len(ps.ServiceBody))
	for _, x := range ps.ServiceBody {
		var ep Endpoint
		ep.PackageName = s.PackageName
		ep.ServiceName = s.ServiceName
		if rpc, ok := x.(*parser.RPC); ok {
			ep.MethodName = rpc.RPCName
			ep.Comment = composeHeadAndInlineComment(rpc.Comments, rpc.InlineComment, "\n")
			ep.Typ = extractRPCType(rpc)
			ep.Req = extractRPCRequest(rpc.RPCRequest, pp, pf)
			ep.Res = extractRPCResponse(rpc.RPCResponse, pp, pf)
		}
		ep.validate()
		eps = append(eps, ep)
	}
	return eps
}

func (e *Endpoint) validate() {
	e.URLPath = "/" + e.PackageName + "/" + e.ServiceName + "/" + e.MethodName
	e.HTTPMethod = "POST"
	if e.Typ != Unary {
		e.HTTPMethod = "GET" // websocket uses GET
	}
}

func extractRPCType(rpc *parser.RPC) RPCType {
	if rpc.RPCRequest.IsStream && rpc.RPCResponse.IsStream {
		return BidirectionalStreaming
	}
	if rpc.RPCRequest.IsStream {
		return ClientStreaming
	}
	if rpc.RPCResponse.IsStream {
		return ServerStreaming
	}
	return Unary
}

func extractRPCRequest(rr *parser.RPCRequest, pp *parser.Proto, pf *ProtoFile) (r Request) {
	msg := findMessage(pp, rr.MessageType)
	r.Params = composeFields(msg, msg.MessageName, pf)
	r.Typ = rr.MessageType
	return r
}

func extractRPCResponse(rr *parser.RPCResponse, pp *parser.Proto, pf *ProtoFile) (r Response) {
	msg := findMessage(pp, rr.MessageType)
	r.Params = composeFields(msg, msg.MessageName, pf)
	r.Typ = rr.MessageType
	return r
}

func findMessage(pp *parser.Proto, mt string) *parser.Message {
	for _, x := range pp.ProtoBody {
		if m, ok := x.(*parser.Message); ok && m.MessageName == mt {
			return m
		}
	}
	log.Panicf("proto doesn't has message %q", mt)
	return nil
}

func composeFields(pm *parser.Message, enclosing string, protoFile *ProtoFile) []Field {
	fs := make([]Field, 0, len(pm.MessageBody))
	for _, x := range pm.MessageBody {
		if pf, ok := x.(*parser.Field); ok {
			var f Field
			f.Comment = composeHeadAndInlineComment(pf.Comments, pf.InlineComment, " ")
			f.Name = pf.FieldName
			f.Typ = pf.Type
			f.Repeat = pf.IsRepeated
			f.Enclosing = enclosing
			f.protoFile = protoFile
			fs = append(fs, f)
		}
	}
	return fs
}

var scalarTypes = map[string]struct{}{
	"double":   {},
	"float":    {},
	"int32":    {},
	"int64":    {},
	"uint32":   {},
	"uint64":   {},
	"sint32":   {},
	"sint64":   {},
	"fixed32":  {},
	"fixed64":  {},
	"sfixed32": {},
	"sfixed64": {},
	"bool":     {},
	"string":   {},
	"bytes":    {},
}

func (f Field) isScalar() (typename string, ok bool) {
	_, ok = scalarTypes[f.Typ]
	return f.Typ, ok
}

func (f Field) isEnum() (typename string, ok bool) {
	scopes := strings.Split(f.Enclosing, ".")
	for i := len(scopes); i >= 0; i-- {
		scope := strings.Join(scopes[:i], ".")
		qualified := f.Typ
		if scope != "" {
			qualified = scope + "." + f.Typ
		}
		for _, e := range f.protoFile.Enums {
			if qualified == e.Name {
				return e.Name, true
			}
		}
	}
	return "", false
}

func (f Field) isObject() (typename string, ok bool) {
	scopes := strings.Split(f.Enclosing, ".")
	for i := len(scopes); i >= 0; i-- {
		scope := strings.Join(scopes[:i], ".")
		qualified := f.Typ
		if scope != "" {
			qualified = scope + "." + f.Typ
		}
		for _, o := range f.protoFile.Objects {
			if qualified == o.Name {
				return o.Name, true
			}
		}
	}
	return "", false
}

// Type resolves f's type to a readable format in the scope of protoFile
func (f Field) Type() (r string) {
	if f.Typ == "" {
		return "(nil)"
	}
	if typename, ok := f.isScalar(); ok {
		r = typename
	} else if typename, ok := f.isEnum(); ok {
		r = "enum " + typename
	} else if typename, ok := f.isObject(); ok {
		r = "object " + typename
	} else {
		r = "(" + f.Typ + ")"
	}
	if f.Repeat {
		r = "array of " + r
	}
	return r
}

// convert typename to href id
func href(typename string) string {
	return strings.Join(strings.Split(strings.ToLower(typename), "."), "")
}

// TypeHRef is Type in addition to a href
func (f Field) TypeHRef() (r string) {
	if f.Typ == "" {
		return "(nil)"
	}
	if typename, ok := f.isScalar(); ok {
		r = typename
	} else if typename, ok := f.isEnum(); ok {
		r = "[enum " + typename + "](#enum-" + href(typename) + ")"
	} else if typename, ok := f.isObject(); ok {
		r = "[object " + typename + "](#object-" + href(typename) + ")"
	} else {
		r = "(" + f.Typ + ")"
	}
	if f.Repeat {
		r = "array of " + r
	}
	return r
}

// add all messages and enums in the proto, but exclude the messages used directly by interfaces
func (pf *ProtoFile) addObjectsAndEnums(pp *parser.Proto) {
	excludes := make(map[string]bool)
	for _, s := range pf.Services {
		for _, inf := range s.Infs {
			excludes[inf.Req.Typ] = true
			excludes[inf.Res.Typ] = true
		}
	}

	var extractMessagesAndEnums func(body []parser.Visitee, enclosingName string)
	extractMessagesAndEnums = func(body []parser.Visitee, enclosingName string) {
		for _, x := range body {
			if enum, ok := x.(*parser.Enum); ok {
				e := composeEnum(enum, enclosingName)
				pf.Enums = append(pf.Enums, e)
			} else if msg, ok := x.(*parser.Message); ok {
				o := composeObject(msg, enclosingName, pf)
				if excludes[o.Name] {
					continue
				}
				pf.Objects = append(pf.Objects, o)
				extractMessagesAndEnums(msg.MessageBody, enclosingName+msg.MessageName+".")
			}
		}
	}

	extractMessagesAndEnums(pp.ProtoBody, ".")
}

func composeObject(pm *parser.Message, enclosingName string, pf *ProtoFile) (o Object) {
	enclosingName = enclosingName[1:] // trim the first "."
	o.Comment = composeHeadComment(pm.Comments)
	o.Name = enclosingName + pm.MessageName
	o.Attrs = composeFields(pm, o.Name, pf)
	return o
}

func composeEnum(pe *parser.Enum, enclosingName string) (e Enum) {
	enclosingName = enclosingName[1:] // trim the first "."
	e.Comment = composeHeadComment(pe.Comments)
	e.Name = enclosingName + pe.EnumName
	e.Constants = composeEnumFields(pe, e.Name)
	return e
}

func composeEnumFields(pe *parser.Enum, enclosing string) []EnumField {
	fs := make([]EnumField, 0, len(pe.EnumBody))
	for _, x := range pe.EnumBody {
		if pf, ok := x.(*parser.EnumField); ok {
			var f EnumField
			f.Comment = composeHeadAndInlineComment(pf.Comments, pf.InlineComment, " ")
			f.Name = pf.Ident
			f.Val = pf.Number
			f.Enclosing = enclosing
			fs = append(fs, f)
		}
	}
	return fs
}
