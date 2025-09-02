package oas

import (
	"encoding/json"
	"errors"
	"sort"
)

// ========== Utilidades para unions / helpers ==========

// Reference representa um $ref simples.
type Reference struct {
	Ref string `json:"$ref"`
}

func isRefObject(raw map[string]any) (string, bool) {
	if v, ok := raw["$ref"]; ok {
		if s, ok2 := v.(string); ok2 {
			return s, true
		}
	}
	return "", false
}

// StringOrStringArray aceita "type" como string ou []string (OpenAPI 3.1 / JSON Schema 2020-12).
type StringOrStringArray struct {
	One  *string
	Many []string
}

func (t *StringOrStringArray) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		t.One = &s
		return nil
	}
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		t.Many = arr
		return nil
	}
	return errors.New("StringOrStringArray: valor deve ser string ou []string")
}

func (t StringOrStringArray) MarshalJSON() ([]byte, error) {
	if t.One != nil {
		return json.Marshal(*t.One)
	}
	return json.Marshal(t.Many)
}

// SchemaOrRef: ou um objeto Schema, ou um $ref.
type SchemaOrRef struct {
	Schema *Schema
	Ref    *Reference
}

func (s *SchemaOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		s.Ref = &Reference{Ref: ref}
		return nil
	}
	var sch Schema
	if err := json.Unmarshal(b, &sch); err != nil {
		return err
	}
	s.Schema = &sch
	return nil
}

func (s SchemaOrRef) MarshalJSON() ([]byte, error) {
	if s.Ref != nil {
		return json.Marshal(s.Ref)
	}
	return json.Marshal(s.Schema)
}

// Items: JSON Schema 2020-12 permite schema único ou lista de schemas.
type Items struct {
	Single *SchemaOrRef
	List   []SchemaOrRef
}

func (it *Items) UnmarshalJSON(b []byte) error {
	var one SchemaOrRef
	if err := json.Unmarshal(b, &one); err == nil && (one.Ref != nil || one.Schema != nil) {
		it.Single = &one
		return nil
	}
	var arr []SchemaOrRef
	if err := json.Unmarshal(b, &arr); err == nil {
		it.List = arr
		return nil
	}
	return errors.New("Items: esperado schema único ou lista de schemas")
}

func (it Items) MarshalJSON() ([]byte, error) {
	if it.Single != nil {
		return json.Marshal(it.Single)
	}
	return json.Marshal(it.List)
}

// AdditionalProperties: bool ou schema (OpenAPI 3.1 / JSON Schema).
type AdditionalProperties struct {
	Allows *bool
	Schema *SchemaOrRef
}

func (ap *AdditionalProperties) UnmarshalJSON(b []byte) error {
	var bo bool
	if err := json.Unmarshal(b, &bo); err == nil {
		ap.Allows = &bo
		return nil
	}
	var sr SchemaOrRef
	if err := json.Unmarshal(b, &sr); err == nil && (sr.Ref != nil || sr.Schema != nil) {
		ap.Schema = &sr
		return nil
	}
	return errors.New("AdditionalProperties: esperado boolean ou schema")
}

func (ap AdditionalProperties) MarshalJSON() ([]byte, error) {
	if ap.Allows != nil {
		return json.Marshal(*ap.Allows)
	}
	return json.Marshal(ap.Schema)
}

// MapStringAny é útil pra extensions "x-*" e payloads variados.
type MapStringAny map[string]any

// ========== Núcleo do Documento ==========

// Document (OpenAPI root)
type Document struct {
	OpenAPI           string                   `json:"openapi"` // "3.1.x"
	Info              Info                     `json:"info"`
	JSONSchemaDialect *string                  `json:"jsonSchemaDialect,omitempty"`
	Servers           []Server                 `json:"servers,omitempty"`
	Paths             Paths                    `json:"paths,omitempty"`
	Webhooks          map[string]PathItemOrRef `json:"webhooks,omitempty"`
	Components        *Components              `json:"components,omitempty"`
	Security          []SecurityRequirement    `json:"security,omitempty"`
	Tags              []Tag                    `json:"tags,omitempty"`
	ExternalDocs      *ExternalDocumentation   `json:"externalDocs,omitempty"`
}

// Info
type Info struct {
	Title          string   `json:"title"`
	Version        string   `json:"version"`
	Summary        *string  `json:"summary,omitempty"`
	Description    *string  `json:"description,omitempty"`
	TermsOfService *string  `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Contact struct {
	Name  *string `json:"name,omitempty"`
	URL   *string `json:"url,omitempty"`
	Email *string `json:"email,omitempty"`
}

func NewContact(name string, url string, email string) *Contact {
	return &Contact{Name: Ptr(name), URL: Ptr(url), Email: Ptr(email)}
}

type License struct {
	Name string  `json:"name"`
	ID   *string `json:"identifier,omitempty"` // OAS 3.1
	URL  *string `json:"url,omitempty"`
}

func NewLicense(name string, urlAndEmail ...string) *License {
	var id *string = nil
	var url *string = nil
	if len(urlAndEmail) >= 1 {
		id = Ptr(urlAndEmail[0])
	}
	if len(urlAndEmail) >= 2 {
		url = Ptr(urlAndEmail[1])
	}
	return &License{Name: name, ID: id, URL: url}
}

// Servers
type Server struct {
	URL         string                    `json:"url"`
	Description *string                   `json:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description *string  `json:"description,omitempty"`
}

// ========== Paths / PathItem / Operation ==========

type Paths map[string]PathItemOrRef

// PathItemOrRef: $ref ou PathItem.
type PathItemOrRef struct {
	Ref      *Reference
	PathItem *PathItem
}

func (p *PathItemOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		p.Ref = &Reference{Ref: ref}
		return nil
	}
	var pi PathItem
	if err := json.Unmarshal(b, &pi); err != nil {
		return err
	}
	p.PathItem = &pi
	return nil
}

func (p PathItemOrRef) MarshalJSON() ([]byte, error) {
	if p.Ref != nil {
		return json.Marshal(p.Ref)
	}
	return json.Marshal(p.PathItem)
}

// PathItem
type PathItem struct {
	Summary     *string          `json:"summary,omitempty"`
	Description *string          `json:"description,omitempty"`
	Get         *Operation       `json:"get,omitempty"`
	Put         *Operation       `json:"put,omitempty"`
	Post        *Operation       `json:"post,omitempty"`
	Delete      *Operation       `json:"delete,omitempty"`
	Options     *Operation       `json:"options,omitempty"`
	Head        *Operation       `json:"head,omitempty"`
	Patch       *Operation       `json:"patch,omitempty"`
	Trace       *Operation       `json:"trace,omitempty"`
	Servers     []Server         `json:"servers,omitempty"`
	Parameters  []ParameterOrRef `json:"parameters,omitempty"`
}

// Operation
type Operation struct {
	Tags         []string                 `json:"tags,omitempty"`
	Summary      *string                  `json:"summary,omitempty"`
	Description  *string                  `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation   `json:"externalDocs,omitempty"`
	OperationID  *string                  `json:"operationId,omitempty"`
	Parameters   []ParameterOrRef         `json:"parameters,omitempty"`
	RequestBody  *RequestBodyOrRef        `json:"requestBody,omitempty"`
	Responses    Responses                `json:"responses"`
	Callbacks    map[string]CallbackOrRef `json:"callbacks,omitempty"`
	Deprecated   *bool                    `json:"deprecated,omitempty"`
	Security     []SecurityRequirement    `json:"security,omitempty"`
	Servers      []Server                 `json:"servers,omitempty"`
}

type ExternalDocumentation struct {
	Description *string `json:"description,omitempty"`
	URL         string  `json:"url"`
}

// ========== Parameters / RequestBody / MediaType / Encoding ==========

type ParameterIn string

const (
	InQuery  ParameterIn = "query"
	InHeader ParameterIn = "header"
	InPath   ParameterIn = "path"
	InCookie ParameterIn = "cookie"
)

type ParameterStyle string

// (mantemos valores comuns; a 3.1 valida por texto)
const (
	StyleForm           ParameterStyle = "form"
	StyleSimple         ParameterStyle = "simple"
	StyleMatrix         ParameterStyle = "matrix"
	StyleLabel          ParameterStyle = "label"
	StyleSpaceDelimited ParameterStyle = "spaceDelimited"
	StylePipeDelimited  ParameterStyle = "pipeDelimited"
	StyleDeepObject     ParameterStyle = "deepObject"
)

type Parameter struct {
	Name            string                  `json:"name"`
	In              ParameterIn             `json:"in"`
	Description     *string                 `json:"description,omitempty"`
	Required        *bool                   `json:"required,omitempty"`
	Deprecated      *bool                   `json:"deprecated,omitempty"`
	AllowEmptyValue *bool                   `json:"allowEmptyValue,omitempty"` // only query
	Style           *ParameterStyle         `json:"style,omitempty"`
	Explode         *bool                   `json:"explode,omitempty"`
	AllowReserved   *bool                   `json:"allowReserved,omitempty"`
	Schema          *SchemaOrRef            `json:"schema,omitempty"`
	Example         any                     `json:"example,omitempty"`
	Examples        map[string]ExampleOrRef `json:"examples,omitempty"`
	Content         map[string]MediaType    `json:"content,omitempty"`
}

type ParameterOrRef struct {
	Param *Parameter
	Ref   *Reference
}

func (p *ParameterOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		p.Ref = &Reference{Ref: ref}
		return nil
	}
	var pr Parameter
	if err := json.Unmarshal(b, &pr); err != nil {
		return err
	}
	p.Param = &pr
	return nil
}

func (p ParameterOrRef) MarshalJSON() ([]byte, error) {
	if p.Ref != nil {
		return json.Marshal(p.Ref)
	}
	return json.Marshal(p.Param)
}

// RequestBody
type RequestBody struct {
	Description *string              `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    *bool                `json:"required,omitempty"`
}

type RequestBodyOrRef struct {
	Body *RequestBody
	Ref  *Reference
}

func (r *RequestBodyOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		r.Ref = &Reference{Ref: ref}
		return nil
	}
	var rb RequestBody
	if err := json.Unmarshal(b, &rb); err != nil {
		return err
	}
	r.Body = &rb
	return nil
}

func (r RequestBodyOrRef) MarshalJSON() ([]byte, error) {
	if r.Ref != nil {
		return json.Marshal(r.Ref)
	}
	return json.Marshal(r.Body)
}

// MediaType
type MediaType struct {
	Schema   *SchemaOrRef            `json:"schema,omitempty"`
	Example  any                     `json:"example,omitempty"`
	Examples map[string]ExampleOrRef `json:"examples,omitempty"`
	Encoding map[string]Encoding     `json:"encoding,omitempty"`
}

// Encoding
type Encoding struct {
	ContentType   *string                `json:"contentType,omitempty"`
	Headers       map[string]HeaderOrRef `json:"headers,omitempty"`
	Style         *ParameterStyle        `json:"style,omitempty"`
	Explode       *bool                  `json:"explode,omitempty"`
	AllowReserved *bool                  `json:"allowReserved,omitempty"`
}

// ========== Responses / Response / Header / Example / Link / Callback ==========

type Responses map[string]ResponseOrRef // inclui "default"

type Response struct {
	Description string                 `json:"description"`
	Headers     map[string]HeaderOrRef `json:"headers,omitempty"`
	Content     map[string]MediaType   `json:"content,omitempty"`
	Links       map[string]LinkOrRef   `json:"links,omitempty"`
}

type ResponseOrRef struct {
	Resp *Response
	Ref  *Reference
}

func (r *ResponseOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		r.Ref = &Reference{Ref: ref}
		return nil
	}
	var rr Response
	if err := json.Unmarshal(b, &rr); err != nil {
		return err
	}
	r.Resp = &rr
	return nil
}

func (r ResponseOrRef) MarshalJSON() ([]byte, error) {
	if r.Ref != nil {
		return json.Marshal(r.Ref)
	}
	return json.Marshal(r.Resp)
}

// Header
type Header struct {
	Description *string                 `json:"description,omitempty"`
	Required    *bool                   `json:"required,omitempty"`
	Deprecated  *bool                   `json:"deprecated,omitempty"`
	Style       *ParameterStyle         `json:"style,omitempty"`
	Explode     *bool                   `json:"explode,omitempty"`
	Schema      *SchemaOrRef            `json:"schema,omitempty"`
	Example     any                     `json:"example,omitempty"`
	Examples    map[string]ExampleOrRef `json:"examples,omitempty"`
	Content     map[string]MediaType    `json:"content,omitempty"`
}

type HeaderOrRef struct {
	Header *Header
	Ref    *Reference
}

func (h *HeaderOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		h.Ref = &Reference{Ref: ref}
		return nil
	}
	var hd Header
	if err := json.Unmarshal(b, &hd); err != nil {
		return err
	}
	h.Header = &hd
	return nil
}

func (h HeaderOrRef) MarshalJSON() ([]byte, error) {
	if h.Ref != nil {
		return json.Marshal(h.Ref)
	}
	return json.Marshal(h.Header)
}

// Example
type Example struct {
	Summary       *string `json:"summary,omitempty"`
	Description   *string `json:"description,omitempty"`
	Value         any     `json:"value,omitempty"`
	ExternalValue *string `json:"externalValue,omitempty"`
}

type ExampleOrRef struct {
	Example *Example
	Ref     *Reference
}

func (e *ExampleOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		e.Ref = &Reference{Ref: ref}
		return nil
	}
	var ex Example
	if err := json.Unmarshal(b, &ex); err != nil {
		return err
	}
	e.Example = &ex
	return nil
}

func (e ExampleOrRef) MarshalJSON() ([]byte, error) {
	if e.Ref != nil {
		return json.Marshal(e.Ref)
	}
	return json.Marshal(e.Example)
}

// Link
type Link struct {
	OperationRef *string      `json:"operationRef,omitempty"`
	OperationID  *string      `json:"operationId,omitempty"`
	Parameters   MapStringAny `json:"parameters,omitempty"`
	RequestBody  any          `json:"requestBody,omitempty"`
	Description  *string      `json:"description,omitempty"`
	Server       *Server      `json:"server,omitempty"`
}

type LinkOrRef struct {
	Link *Link
	Ref  *Reference
}

func (l *LinkOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		l.Ref = &Reference{Ref: ref}
		return nil
	}
	var lk Link
	if err := json.Unmarshal(b, &lk); err != nil {
		return err
	}
	l.Link = &lk
	return nil
}

func (l LinkOrRef) MarshalJSON() ([]byte, error) {
	if l.Ref != nil {
		return json.Marshal(l.Ref)
	}
	return json.Marshal(l.Link)
}

// Callback
type Callback map[string]PathItemOrRef

type CallbackOrRef struct {
	Callback *Callback
	Ref      *Reference
}

func (c *CallbackOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		c.Ref = &Reference{Ref: ref}
		return nil
	}
	var cb Callback
	if err := json.Unmarshal(b, &cb); err != nil {
		return err
	}
	c.Callback = &cb
	return nil
}

func (c CallbackOrRef) MarshalJSON() ([]byte, error) {
	if c.Ref != nil {
		return json.Marshal(c.Ref)
	}
	return json.Marshal(c.Callback)
}

// ========== Components / Security / Tags ==========

type Components struct {
	Schemas         map[string]SchemaOrRef         `json:"schemas,omitempty"`
	Responses       map[string]ResponseOrRef       `json:"responses,omitempty"`
	Parameters      map[string]ParameterOrRef      `json:"parameters,omitempty"`
	Examples        map[string]ExampleOrRef        `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBodyOrRef    `json:"requestBodies,omitempty"`
	Headers         map[string]HeaderOrRef         `json:"headers,omitempty"`
	SecuritySchemes map[string]SecuritySchemeOrRef `json:"securitySchemes,omitempty"`
	Links           map[string]LinkOrRef           `json:"links,omitempty"`
	Callbacks       map[string]CallbackOrRef       `json:"callbacks,omitempty"`
	PathItems       map[string]PathItemOrRef       `json:"pathItems,omitempty"` // OAS 3.1
}

type SecurityRequirement map[string][]string

// Security Schemes

type SecuritySchemeType string

const (
	SecAPIKey        SecuritySchemeType = "apiKey"
	SecHTTP          SecuritySchemeType = "http"
	SecMutualTLS     SecuritySchemeType = "mutualTLS"
	SecOAuth2        SecuritySchemeType = "oauth2"
	SecOpenIDConnect SecuritySchemeType = "openIdConnect"
)

type SecurityScheme struct {
	Type        SecuritySchemeType `json:"type"`
	Description *string            `json:"description,omitempty"`

	// apiKey
	Name *string     `json:"name,omitempty"`
	In   ParameterIn `json:"in,omitempty"` // "query" | "header" | "cookie"

	// http
	Scheme       *string `json:"scheme,omitempty"`
	BearerFormat *string `json:"bearerFormat,omitempty"`

	// oauth2
	Flows *OAuthFlows `json:"flows,omitempty"`

	// openIdConnect
	OpenIDConnectURL *string `json:"openIdConnectUrl,omitempty"`
}

type SecuritySchemeOrRef struct {
	Scheme *SecurityScheme
	Ref    *Reference
}

func (s *SecuritySchemeOrRef) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if ref, ok := isRefObject(raw); ok {
		s.Ref = &Reference{Ref: ref}
		return nil
	}
	var sc SecurityScheme
	if err := json.Unmarshal(b, &sc); err != nil {
		return err
	}
	s.Scheme = &sc
	return nil
}

func (s SecuritySchemeOrRef) MarshalJSON() ([]byte, error) {
	if s.Ref != nil {
		return json.Marshal(s.Ref)
	}
	return json.Marshal(s.Scheme)
}

// OAuth2 (simplificado conforme OAS 3.1)
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       *string           `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// Tags
type Tag struct {
	Name         string                 `json:"name"`
	Description  *string                `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}

// ========== JSON Schema (núcleo para 3.1) ==========

// Schema representa o dialeto JSON Schema 2020-12 na medida necessária para OAS 3.1.
// (Campos menos comuns podem ser adicionados no mesmo padrão.)
type Schema struct {
	// Meta
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Default     any     `json:"default,omitempty"`
	Deprecated  *bool   `json:"deprecated,omitempty"`
	ReadOnly    *bool   `json:"readOnly,omitempty"`
	WriteOnly   *bool   `json:"writeOnly,omitempty"`
	Examples    []any   `json:"examples,omitempty"`

	// Tipos
	Type  *StringOrStringArray `json:"type,omitempty"`
	Enum  []any                `json:"enum,omitempty"`
	Const any                  `json:"const,omitempty"`

	// Combinações
	AllOf []SchemaOrRef `json:"allOf,omitempty"`
	OneOf []SchemaOrRef `json:"oneOf,omitempty"`
	AnyOf []SchemaOrRef `json:"anyOf,omitempty"`
	Not   *SchemaOrRef  `json:"not,omitempty"`

	// Objetos
	Properties           map[string]SchemaOrRef `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties *AdditionalProperties  `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]SchemaOrRef `json:"patternProperties,omitempty"`
	MinProperties        *int                   `json:"minProperties,omitempty"`
	MaxProperties        *int                   `json:"maxProperties,omitempty"`

	// Arrays
	Items       *Items        `json:"items,omitempty"`
	PrefixItems []SchemaOrRef `json:"prefixItems,omitempty"` // JSON Schema 2020-12
	MinItems    *int          `json:"minItems,omitempty"`
	MaxItems    *int          `json:"maxItems,omitempty"`
	UniqueItems *bool         `json:"uniqueItems,omitempty"`
	Contains    *SchemaOrRef  `json:"contains,omitempty"`
	MinContains *int          `json:"minContains,omitempty"`
	MaxContains *int          `json:"maxContains,omitempty"`

	// Strings
	MinLength *int    `json:"minLength,omitempty"`
	MaxLength *int    `json:"maxLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`
	Format    *string `json:"format,omitempty"`

	// Numbers
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`

	// Misc (discriminator, xml, etc. — opcionais, comuns no mundo OAS)
	Discriminator *Discriminator `json:"discriminator,omitempty"`
	XML           *XML           `json:"xml,omitempty"`

	// Extensões "x-*" ficam livres no nível de uso (MapStringAny) quando necessário.
}

// Discriminator (OAS)
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

// XML (OAS)
type XML struct {
	Name      *string `json:"name,omitempty"`
	Namespace *string `json:"namespace,omitempty"`
	Prefix    *string `json:"prefix,omitempty"`
	Attribute *bool   `json:"attribute,omitempty"`
	Wrapped   *bool   `json:"wrapped,omitempty"`
}

// ========== Helpers opcionais ==========

// ValidateRequiredResponses garante que há pelo menos um response (padrão de sanity).
func (op *Operation) ValidateRequiredResponses() error {
	if op == nil {
		return nil
	}
	if len(op.Responses) == 0 {
		return errors.New("operation.responses não pode ser vazio")
	}
	// ordena chaves só pra debug
	keys := make([]string, 0, len(op.Responses))
	for k := range op.Responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return nil
}
