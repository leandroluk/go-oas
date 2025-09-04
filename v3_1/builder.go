package oas

import (
	"encoding/json"
	"fmt"
	"strings"
)

func joinAndTrim(parts ...string) string {
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		cleaned = append(cleaned, strings.TrimSpace(p))
	}
	return strings.Join(cleaned, " ")
}

// =======================
// Root Builder
// =======================

type Builder struct {
	doc *Document
}

func NewBuilder() *Builder {
	return &Builder{
		doc: &Document{
			OpenAPI:    "3.1.0",
			Info:       Info{},
			Paths:      make(Paths),
			Components: &Components{},
		},
	}
}

func (b *Builder) SetTitle(title string) *Builder {
	b.doc.Info.Title = title
	return b
}

func (b *Builder) SetVersion(version string) *Builder {
	b.doc.Info.Version = version
	return b
}

func (b *Builder) SetSummary(parts ...string) *Builder {
	summary := joinAndTrim(parts...)
	b.doc.Info.Summary = &summary
	return b
}

func (b *Builder) SetDescription(parts ...string) *Builder {
	description := joinAndTrim(parts...)
	b.doc.Info.Description = &description
	return b
}

func (b *Builder) SetTermsOfService(parts ...string) *Builder {
	termsOfService := joinAndTrim(parts...)
	b.doc.Info.TermsOfService = &termsOfService
	return b
}

func (b *Builder) SetContact(contact *Contact) *Builder {
	b.doc.Info.Contact = contact
	return b
}

func (b *Builder) SetLicense(license License) *Builder {
	b.doc.Info.License = &license
	return b
}

func (b *Builder) AddServer(url, description string) *Builder {
	b.doc.Servers = append(b.doc.Servers, Server{
		URL:         url,
		Description: &description,
	})
	return b
}

func (b *Builder) AddSchema(name string, schema Schema) *Builder {
	if b.doc.Components.Schemas == nil {
		b.doc.Components.Schemas = make(map[string]SchemaOrRef)
	}
	b.doc.Components.Schemas[name] = SchemaOrRef{Schema: &schema}
	return b
}

func (b *Builder) AddSecurityScheme(name string, scheme SecurityScheme) *Builder {
	if b.doc.Components.SecuritySchemes == nil {
		b.doc.Components.SecuritySchemes = make(map[string]SecuritySchemeOrRef)
	}
	b.doc.Components.SecuritySchemes[name] = SecuritySchemeOrRef{Scheme: &scheme}
	return b
}

func (b *Builder) Security(req SecurityRequirement) *Builder {
	b.doc.Security = append(b.doc.Security, req)
	return b
}

func (b *Builder) AddTag(tag Tag) *Builder {
	b.doc.Tags = append(b.doc.Tags, tag)
	return b
}

func (b *Builder) ExternalDocs(desc, url string) *Builder {
	b.doc.ExternalDocs = &ExternalDocumentation{Description: &desc, URL: url}
	return b
}

func (b *Builder) Path(path string) *PathBuilder {
	piRef, ok := b.doc.Paths[path]
	if !ok || piRef.PathItem == nil {
		pi := &PathItem{}
		b.doc.Paths[path] = PathItemOrRef{PathItem: pi}
		return &PathBuilder{builder: b, path: path, item: pi}
	}
	return &PathBuilder{builder: b, path: path, item: piRef.PathItem}
}

func (b *Builder) Build() *Document {
	return b.doc
}

func (b *Builder) JSON() ([]byte, error) {
	return json.MarshalIndent(b.doc, "", "  ")
}

// =======================
// Path Builder
// =======================

type PathBuilder struct {
	builder *Builder
	path    string
	item    *PathItem
}

func (pb *PathBuilder) DonePath() *Builder {
	return pb.builder
}

func (pb *PathBuilder) Get(summary string) *OperationBuilder {
	return pb.addOp("get", summary)
}
func (pb *PathBuilder) Post(summary string) *OperationBuilder {
	return pb.addOp("post", summary)
}
func (pb *PathBuilder) Put(summary string) *OperationBuilder {
	return pb.addOp("put", summary)
}
func (pb *PathBuilder) Delete(summary string) *OperationBuilder {
	return pb.addOp("delete", summary)
}
func (pb *PathBuilder) Patch(summary string) *OperationBuilder {
	return pb.addOp("patch", summary)
}

func (pb *PathBuilder) addOp(method, summary string) *OperationBuilder {
	op := &Operation{Summary: &summary, Responses: make(Responses)}
	switch method {
	case "get":
		pb.item.Get = op
	case "post":
		pb.item.Post = op
	case "put":
		pb.item.Put = op
	case "delete":
		pb.item.Delete = op
	case "patch":
		pb.item.Patch = op
	}
	return &OperationBuilder{pathBuilder: pb, method: method, op: op}
}

// =======================
// Operation Builder
// =======================

type OperationBuilder struct {
	pathBuilder *PathBuilder
	method      string
	op          *Operation
}

func (ob *OperationBuilder) AddTag(tag string) *OperationBuilder {
	ob.op.Tags = append(ob.op.Tags, tag)
	return ob
}

func (ob *OperationBuilder) SetSummary(s string) *OperationBuilder {
	ob.op.Summary = &s
	return ob
}

func (ob *OperationBuilder) SetDescription(d string) *OperationBuilder {
	ob.op.Description = &d
	return ob
}

func (ob *OperationBuilder) SetExternalDocs(desc, url string) *OperationBuilder {
	ob.op.ExternalDocs = &ExternalDocumentation{Description: &desc, URL: url}
	return ob
}

func (ob *OperationBuilder) SetOperationID(id string) *OperationBuilder {
	ob.op.OperationID = &id
	return ob
}

func (ob *OperationBuilder) SetParameters(params ...ParameterOrRef) *OperationBuilder {
	ob.op.Parameters = append(ob.op.Parameters, params...)
	return ob
}

func (ob *OperationBuilder) SetRequestBody(rb RequestBodyOrRef) *OperationBuilder {
	ob.op.RequestBody = &rb
	return ob
}

func (ob *OperationBuilder) AddSecurity(req SecurityRequirement) *OperationBuilder {
	ob.op.Security = append(ob.op.Security, req)
	return ob
}

func (ob *OperationBuilder) SetResponses(resps Responses) *OperationBuilder {
	for k, v := range resps {
		ob.op.Responses[k] = v
	}
	return ob
}

func (ob *OperationBuilder) AddServer(url, description string) *OperationBuilder {
	if ob.op.Servers == nil {
		ob.op.Servers = make([]Server, 0)
	}
	ob.op.Servers = append(ob.op.Servers, Server{
		URL:         url,
		Description: &description,
	})
	return ob
}

func (ob *OperationBuilder) DoneOp() *PathBuilder {
	return ob.pathBuilder
}

func (ob *OperationBuilder) SetDeprecated() *OperationBuilder {
	val := true
	ob.op.Deprecated = &val
	return ob
}

// ---------------- Params -----------------

func (ob *OperationBuilder) ParamQuery(name, typ, desc string, required bool) *OperationBuilder {
	return ob.addParam(name, InQuery, typ, desc, required)
}
func (ob *OperationBuilder) ParamPath(name, typ, desc string) *OperationBuilder {
	return ob.addParam(name, InPath, typ, desc, true)
}
func (ob *OperationBuilder) ParamHeader(name, typ, desc string, required bool) *OperationBuilder {
	return ob.addParam(name, InHeader, typ, desc, required)
}
func (ob *OperationBuilder) ParamCookie(name, typ, desc string, required bool) *OperationBuilder {
	return ob.addParam(name, InCookie, typ, desc, required)
}

func (ob *OperationBuilder) addParam(name string, in ParameterIn, typ string, desc string, required bool) *OperationBuilder {
	schema := Schema{Type: &StringOrArray{One: &typ}}
	param := Parameter{
		Name:        name,
		In:          in,
		Description: &desc,
		Required:    &required,
		Schema:      &SchemaOrRef{Schema: &schema},
	}
	ob.op.Parameters = append(ob.op.Parameters, ParameterOrRef{Param: &param})
	return ob
}

// ---------------- Request Body -----------------

func (ob *OperationBuilder) RequestJSON(schema SchemaOrRef, required bool) *OperationBuilder {
	rb := RequestBody{
		Content: map[string]MediaType{
			"application/json": {Schema: &schema},
		},
		Required: &required,
	}
	ob.op.RequestBody = &RequestBodyOrRef{Body: &rb}
	return ob
}

// ---------------- Responses -----------------

func (ob *OperationBuilder) ResponseJSON(status int, desc string, schema SchemaOrRef) *OperationBuilder {
	resp := Response{
		Description: desc,
		Content: map[string]MediaType{
			"application/json": {Schema: &schema},
		},
	}
	code := fmt.Sprintf("%d", status)
	ob.op.Responses[code] = ResponseOrRef{Resp: &resp}
	return ob
}

func (ob *OperationBuilder) ResponseText(status int, desc string) *OperationBuilder {
	resp := Response{
		Description: desc,
		Content: map[string]MediaType{
			"text/plain": {Schema: &SchemaOrRef{Schema: &Schema{Type: &StringOrArray{One: Ptr("string")}}}},
		},
	}
	code := fmt.Sprintf("%d", status)
	ob.op.Responses[code] = ResponseOrRef{Resp: &resp}
	return ob
}

func (ob *OperationBuilder) ResponseWithHeaders(status int, desc string, schema SchemaOrRef, headers map[string]Header) *OperationBuilder {
	resp := Response{
		Description: desc,
		Content: map[string]MediaType{
			"application/json": {Schema: &schema},
		},
	}
	if len(headers) > 0 {
		resp.Headers = make(map[string]HeaderOrRef)
		for k, v := range headers {
			copy := v
			resp.Headers[k] = HeaderOrRef{Header: &copy}
		}
	}
	code := fmt.Sprintf("%d", status)
	ob.op.Responses[code] = ResponseOrRef{Resp: &resp}
	return ob
}

// ---------------- Examples & Links -----------------

func (ob *OperationBuilder) Example(status int, mediaType string, example Example) *OperationBuilder {
	code := fmt.Sprintf("%d", status)
	if resp, ok := ob.op.Responses[code]; ok && resp.Resp != nil {
		if resp.Resp.Content == nil {
			resp.Resp.Content = make(map[string]MediaType)
		}
		mt := resp.Resp.Content[mediaType]
		if mt.Examples == nil {
			mt.Examples = make(map[string]ExampleOrRef)
		}
		mt.Examples["example"] = ExampleOrRef{Example: &example}
		resp.Resp.Content[mediaType] = mt
		ob.op.Responses[code] = resp
	}
	return ob
}

func (ob *OperationBuilder) Link(status int, name string, link Link) *OperationBuilder {
	code := fmt.Sprintf("%d", status)
	if resp, ok := ob.op.Responses[code]; ok && resp.Resp != nil {
		if resp.Resp.Links == nil {
			resp.Resp.Links = make(map[string]LinkOrRef)
		}
		resp.Resp.Links[name] = LinkOrRef{Link: &link}
		ob.op.Responses[code] = resp
	}
	return ob
}

// ---------------- Callbacks -----------------

func (ob *OperationBuilder) Callback(name string, cb Callback) *OperationBuilder {
	if ob.op.Callbacks == nil {
		ob.op.Callbacks = make(map[string]CallbackOrRef)
	}
	ob.op.Callbacks[name] = CallbackOrRef{Callback: &cb}
	return ob
}

// =======================
// Helpers
// =======================

func Ptr[T any](t T) *T { return &t }
