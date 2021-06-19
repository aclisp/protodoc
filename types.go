package main

import (
	"log"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// a proto file is a parsing unit
type protoFile struct {
	// a proto file could have multiple service
	services []service
	// a proto file should have multiple object
	objects []object
	// a proto file should have multiple enum
	enums []enum
}

// a proto file could have multiple service
type service struct {
	// service only use comment placed at the beginning
	comment string
	// the package name of the proto file
	packageName string
	// my name
	serviceName string
	// a service has multiple endpoint
	infs []endpoint
}

// enndpoint is also called method or interface
type endpoint struct {
	// the package name of the proto file
	packageName string
	// the enclosing service name
	serviceName string
	// my name
	methodName string
	// the url path where api gateway resolves
	urlPath string
	// the http method, which is always POST
	httpMethod string
	// comment placed at the beginning, as well as inline-comment placed at the ending
	comment string

	typ rpcType
	req request
	res response
}

type rpcType int

const (
	unary = iota
	serverStreaming
	clientStreaming
	bidirectionalStreaming
)

type request struct {
	params []field
	typ    string
}

type response struct {
	params []field
	typ    string
}

type field struct {
	// comment placed at the beginning, as well as inline-comment placed at the ending
	comment string
	name    string
	typ     string
	repeat  bool
	// the enclosing type name of this field
	enclosing string
}

// object is user-defined field type
type object struct {
	// only use comment placed at the beginning
	comment string
	// refers by field.typ
	name string
	// attributes
	attrs []field
}

// enum is user-defined type which has one of a pre-defined list of values
type enum struct {
	// only use comment placed at the beginning
	comment string
	// refers by field.typ
	name string
	// constants
	constants []enumField
}

type enumField struct {
	// comment placed at the beginning, as well as inline-comment placed at the ending
	comment string
	name    string
	val     string
	// the enclosing type name of this field
	enclosing string
}

func (t rpcType) String() string {
	switch t {
	case unary:
		return "unary"
	case clientStreaming:
		return "client-streaming"
	case serverStreaming:
		return "server-streaming"
	case bidirectionalStreaming:
		return "bidirectional-streaming"
	}
	return "unknown"
}

func extract(pp *parser.Proto) (pf protoFile) {
	// find all services in proto body
	for _, x := range pp.ProtoBody {
		if service, ok := x.(*parser.Service); ok {
			s := extractService(service, pp)
			pf.services = append(pf.services, s)
		}
	}
	pf.objects, pf.enums = composeObjectsAndEnums(pf, pp)
	return pf
}

// extract our service from parser
func extractService(ps *parser.Service, pp *parser.Proto) (s service) {
	s.comment = composeHeadComment(ps.Comments)
	s.packageName = extractPackageName(pp)
	s.serviceName = ps.ServiceName
	s.infs = composeInterfaces(s, ps, pp)
	return s
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

func composeInterfaces(s service, ps *parser.Service, pp *parser.Proto) []endpoint {
	eps := make([]endpoint, 0, len(ps.ServiceBody))
	for _, x := range ps.ServiceBody {
		var ep endpoint
		ep.packageName = s.packageName
		ep.serviceName = s.serviceName
		if rpc, ok := x.(*parser.RPC); ok {
			ep.methodName = rpc.RPCName
			ep.comment = composeHeadAndInlineComment(rpc.Comments, rpc.InlineComment, "\n")
			ep.typ = extractRPCType(rpc)
			ep.req = extractRPCRequest(rpc.RPCRequest, pp)
			ep.res = extractRPCResponse(rpc.RPCResponse, pp)
		}
		ep.validate()
		eps = append(eps, ep)
	}
	return eps
}

func (e *endpoint) validate() {
	e.urlPath = "/" + e.packageName + "/" + e.serviceName + "/" + e.methodName
	e.httpMethod = "POST"
	if e.typ != unary {
		e.httpMethod = "GET" // websocket uses GET
	}
}

func extractRPCType(rpc *parser.RPC) rpcType {
	if rpc.RPCRequest.IsStream && rpc.RPCResponse.IsStream {
		return bidirectionalStreaming
	}
	if rpc.RPCRequest.IsStream {
		return clientStreaming
	}
	if rpc.RPCResponse.IsStream {
		return serverStreaming
	}
	return unary
}

func extractRPCRequest(rr *parser.RPCRequest, pp *parser.Proto) (r request) {
	msg := findMessage(pp, rr.MessageType)
	r.params = composeFields(msg, msg.MessageName)
	r.typ = rr.MessageType
	return r
}

func extractRPCResponse(rr *parser.RPCResponse, pp *parser.Proto) (r response) {
	msg := findMessage(pp, rr.MessageType)
	r.params = composeFields(msg, msg.MessageName)
	r.typ = rr.MessageType
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

func composeFields(pm *parser.Message, enclosing string) []field {
	fs := make([]field, 0, len(pm.MessageBody))
	for _, x := range pm.MessageBody {
		if pf, ok := x.(*parser.Field); ok {
			var f field
			f.comment = composeHeadAndInlineComment(pf.Comments, pf.InlineComment, " ")
			f.name = pf.FieldName
			f.typ = pf.Type
			f.repeat = pf.IsRepeated
			f.enclosing = enclosing
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

// Type resolves f's type to a readable format in the scope of protoFile
func (f field) Type(pf protoFile) (r string) {
	if f.typ == "" {
		return "(nil)"
	}
	isScalar := func() (string, bool) {
		_, ok := scalarTypes[f.typ]
		return f.typ, ok
	}
	isEnum := func() (string, bool) {
		scopes := strings.Split(f.enclosing, ".")
		for i := len(scopes); i >= 0; i-- {
			scope := strings.Join(scopes[:i], ".")
			qualified := f.typ
			if scope != "" {
				qualified = scope + "." + f.typ
			}
			for _, e := range pf.enums {
				if qualified == e.name {
					return e.name, true
				}
			}
		}
		return "", false
	}
	isObject := func() (string, bool) {
		scopes := strings.Split(f.enclosing, ".")
		for i := len(scopes); i >= 0; i-- {
			scope := strings.Join(scopes[:i], ".")
			qualified := f.typ
			if scope != "" {
				qualified = scope + "." + f.typ
			}
			for _, o := range pf.objects {
				if qualified == o.name {
					return o.name, true
				}
			}
		}
		return "", false
	}
	if typeName, ok := isScalar(); ok {
		r = typeName
	} else if typeName, ok := isEnum(); ok {
		r = "enum " + typeName
	} else if typeName, ok := isObject(); ok {
		r = "object " + typeName
	} else {
		r = "(" + f.typ + ")"
	}
	if f.repeat {
		r = "array of " + r
	}
	return r
}

// get all messages and enums in the proto, but exclude the messages used directly by interfaces
func composeObjectsAndEnums(pf protoFile, pp *parser.Proto) (objects []object, enums []enum) {
	excludes := make(map[string]bool)
	for _, s := range pf.services {
		for _, inf := range s.infs {
			excludes[inf.req.typ] = true
			excludes[inf.res.typ] = true
		}
	}

	var extractMessagesAndEnums func(body []parser.Visitee, enclosingName string)
	extractMessagesAndEnums = func(body []parser.Visitee, enclosingName string) {
		for _, x := range body {
			if enum, ok := x.(*parser.Enum); ok {
				e := composeEnum(enum, enclosingName)
				enums = append(enums, e)
			} else if msg, ok := x.(*parser.Message); ok {
				o := composeObject(msg, enclosingName)
				if excludes[o.name] {
					continue
				}
				objects = append(objects, o)
				extractMessagesAndEnums(msg.MessageBody, enclosingName+msg.MessageName+".")
			}
		}
	}

	extractMessagesAndEnums(pp.ProtoBody, ".")
	return objects, enums
}

func composeObject(pm *parser.Message, enclosingName string) (o object) {
	enclosingName = enclosingName[1:] // trim the first "."
	o.comment = composeHeadComment(pm.Comments)
	o.name = enclosingName + pm.MessageName
	o.attrs = composeFields(pm, o.name)
	return o
}

func composeEnum(pe *parser.Enum, enclosingName string) (e enum) {
	enclosingName = enclosingName[1:] // trim the first "."
	e.comment = composeHeadComment(pe.Comments)
	e.name = enclosingName + pe.EnumName
	e.constants = composeEnumFields(pe, e.name)
	return e
}

func composeEnumFields(pe *parser.Enum, enclosing string) []enumField {
	fs := make([]enumField, 0, len(pe.EnumBody))
	for _, x := range pe.EnumBody {
		if pf, ok := x.(*parser.EnumField); ok {
			var f enumField
			f.comment = composeHeadAndInlineComment(pf.Comments, pf.InlineComment, " ")
			f.name = pf.Ident
			f.val = pf.Number
			f.enclosing = enclosing
			fs = append(fs, f)
		}
	}
	return fs
}
