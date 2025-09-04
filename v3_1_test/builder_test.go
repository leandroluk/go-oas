package oas_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	oas "github.com/leandroluk/go-oas/v3_1"
)

func TestBuilder_FullSpec(t *testing.T) {
	builder := oas.NewBuilder().
		SetTitle("Minha API").
		SetVersion("1.0.0").
		SetSummary("Este é", "um summary com", "multiplas linhas").
		SetDescription("Exemplo de API completa", " gerada em runtime").
		SetTermsOfService("Este é", "um terms of service com", "multiplas linhas").
		SetContact(oas.NewContact("nome", "http://example.com", "a@a.com")).
		SetLicense(*oas.NewLicense("licença", "http://example.com", "a@a.com")).
		AddServer("http://localhost:8080", "Servidor local").
		AddSchema("User", oas.Schema{
			Type: &oas.StringOrArray{One: oas.Ptr("object")},
			Properties: map[string]oas.SchemaOrRef{
				"id":   {Schema: &oas.Schema{Type: &oas.StringOrArray{One: oas.Ptr("integer")}}},
				"name": {Schema: &oas.Schema{Type: &oas.StringOrArray{One: oas.Ptr("string")}}},
			},
			Required: []string{"id", "name"},
		}).
		AddSecurityScheme("bearerAuth", oas.SecurityScheme{
			Type:         "http",
			Scheme:       oas.Ptr("bearer"),
			BearerFormat: oas.Ptr("JWT"),
		}).
		Security(oas.SecurityRequirement{"bearerAuth": {}}).
		AddTag(oas.Tag{Name: "users", Description: oas.Ptr("Operações de usuários")}).
		ExternalDocs("Docs externos", "https://example.com/docs").
		Path("/users").
		Get("Lista usuários").
		AddTag("users").
		ParamQuery("limit", "integer", "Número de resultados", false).
		ResponseJSON(200, "OK", oas.SchemaOrRef{
			Schema: &oas.Schema{
				Type:  &oas.StringOrArray{One: oas.Ptr("array")},
				Items: &oas.Items{Single: &oas.SchemaOrRef{Ref: &oas.Reference{Ref: "#/components/schemas/User"}}},
			},
		}).
		Example(200, "application/json", oas.Example{
			Summary: oas.Ptr("Exemplo de lista de usuários"),
			Value: []map[string]any{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		}).
		DoneOp().
		Post("Cria usuário").
		AddTag("users").
		RequestJSON(oas.SchemaOrRef{Ref: &oas.Reference{Ref: "#/components/schemas/User"}}, true).
		ResponseJSON(201, "Criado", oas.SchemaOrRef{Ref: &oas.Reference{Ref: "#/components/schemas/User"}}).
		Link(201, "GetUserById", oas.Link{
			OperationID: oas.Ptr("getUser"),
			Parameters: map[string]any{
				"id": "$response.body#/id",
			},
		}).
		DoneOp().
		DonePath()

	// constrói o documento
	doc := builder.Build()
	data, err := json.MarshalIndent(doc, "", "  ")
	require.NoError(t, err)

	// printa pra debug
	fmt.Println(string(data))

	// valida que campos básicos existem
	require.Equal(t, "3.1.0", doc.OpenAPI)
	require.Equal(t, "Minha API", doc.Info.Title)
	require.Equal(t, "Este é um summary com multiplas linhas", *doc.Info.Summary)
	require.Equal(t, "Exemplo de API completa gerada em runtime", *doc.Info.Description)
	require.Equal(t, "Este é um terms of service com multiplas linhas", *doc.Info.TermsOfService)
	require.NotNil(t, doc.Info.Contact)
	require.NotNil(t, doc.Info.License)
	require.Contains(t, doc.Paths, "/users")
	require.NotNil(t, doc.Components.Schemas["User"])
	require.NotNil(t, doc.Components.SecuritySchemes["bearerAuth"])
}

func TestBuilder_ExtraCoverage(t *testing.T) {
	b := oas.NewBuilder().
		SetTitle("API Extra").
		SetVersion("1.0.0").
		AddTag(oas.Tag{Name: "extra"}).
		ExternalDocs("Docs", "https://example.com").
		AddSecurityScheme("apiKeyAuth", oas.SecurityScheme{
			Type: "apiKey",
			Name: oas.Ptr("X-API-Key"),
			In:   oas.InHeader,
		}).
		Security(oas.SecurityRequirement{"apiKeyAuth": {}})

	// PUT /items
	b.Path("/items").
		Put("Atualiza item").
		SetSummary("Sumário").
		SetDescription("Descrição").
		SetDeprecated().
		SetExternalDocs("Op docs", "https://op.example.com").
		AddSecurity(oas.SecurityRequirement{"apiKeyAuth": {}}).
		ParamPath("id", "string", "ID do item").
		ParamHeader("X-Custom", "string", "Header custom", false).
		ParamCookie("session", "string", "Sessão", true).
		RequestJSON(oas.SchemaOrRef{
			Schema: &oas.Schema{Type: &oas.StringOrArray{One: oas.Ptr("object")}},
		}, true).
		ResponseText(204, "Sem conteúdo").
		ResponseStatus(218, "This is fine").
		ResponseWithHeaders(200, "OK com header",
			oas.SchemaOrRef{Schema: &oas.Schema{Type: &oas.StringOrArray{One: oas.Ptr("string")}}},
			map[string]oas.Header{
				"X-Rate-Limit": {Description: oas.Ptr("limite")},
			},
		).
		Example(200, "application/json", oas.Example{
			Summary: oas.Ptr("Exemplo vazio"),
			Value:   map[string]any{"foo": "bar"},
		}).
		Link(200, "next", oas.Link{
			OperationID: oas.Ptr("getNext"),
		}).
		Callback("onEvent", oas.Callback{
			"{$request.body#/url}": oas.PathItemOrRef{
				PathItem: &oas.PathItem{
					Post: &oas.Operation{
						Responses: oas.Responses{
							"200": oas.ResponseOrRef{Resp: &oas.Response{Description: "callback ok"}},
						},
					},
				},
			},
		}).
		DoneOp().
		DonePath()

	// Força JSON()
	data, err := b.JSON()
	require.NoError(t, err)

	// valida parse reverso
	var parsed oas.Document
	require.NoError(t, json.Unmarshal(data, &parsed))
	require.Equal(t, "3.1.0", parsed.OpenAPI)
	require.Contains(t, parsed.Paths, "/items")

	// cobre helper oas.Ptr
	require.Equal(t, "hello", *oas.Ptr("hello"))
}

func TestBuilder_DeletePatchAndExampleNilContent(t *testing.T) {
	b := oas.NewBuilder().SetTitle("API").SetVersion("1.0.0")

	// cria um DELETE com resposta JSON
	pathBuilder := b.Path("/items/{id}")
	pathBuilder.
		Delete("Remove item").
		ResponseJSON(200, "deleted", oas.SchemaOrRef{}).
		DoneOp().
		DonePath()

	// força o Content da resposta a ser nil
	op := b.Build().Paths["/items/{id}"].PathItem.Delete
	resp := op.Responses["200"].Resp
	resp.Content = nil

	// chama Example -> deve cair no branch Content==nil
	b.Path("/items/{id}").
		Delete("Remove item again").
		Example(200, "application/json", oas.Example{
			Summary: oas.Ptr("forced"),
			Value:   map[string]any{"foo": "bar"},
		}).
		DoneOp().
		DonePath()

	// adiciona PATCH em outro path
	b.Path("/items").
		Patch("Patch item").
		ResponseText(200, "ok").
		DoneOp().
		DonePath()

	// valida JSON final tem delete + patch
	data, err := b.JSON()
	require.NoError(t, err)
	s := string(data)
	require.Contains(t, s, `"delete"`)
	require.Contains(t, s, `"patch"`)
}

func TestBuilder_Path_MergeAndExampleNilContent(t *testing.T) {
	b := oas.NewBuilder().SetTitle("API").SetVersion("1.0.0")

	// Começa no mesmo path
	pb := b.Path("/items/{id}")

	// DELETE com response 200
	ob := pb.
		Delete("Remove item").
		ResponseJSON(200, "deleted", oas.SchemaOrRef{})

	// Força o branch: zera o Content dessa mesma response ANTES de chamar Example
	op := b.Build().Paths["/items/{id}"].PathItem.Delete
	resp := op.Responses["200"].Resp
	resp.Content = nil

	// Agora Example deve cair no branch Content == nil
	pb = ob.
		Example(200, "application/json", oas.Example{
			Summary: oas.Ptr("forced"),
			Value:   map[string]any{"ok": true},
		}).
		DoneOp()

	// Adiciona PATCH no MESMO path (merge, não sobrescreve)
	pb.
		Patch("Patch item").
		ResponseText(200, "ok").
		DoneOp().
		DonePath()

	// JSON deve conter delete e patch
	data, err := b.JSON()
	require.NoError(t, err)

	var doc oas.Document
	require.NoError(t, json.Unmarshal(data, &doc))

	pi := doc.Paths["/items/{id}"].PathItem
	require.NotNil(t, pi.Delete)
	require.NotNil(t, pi.Patch)

	// Garante que o Example foi registrado no mediaType "application/json"
	r := pi.Delete.Responses["200"].Resp
	require.NotNil(t, r)
	require.NotNil(t, r.Content)

	mt, ok := r.Content["application/json"]
	require.True(t, ok, "application/json mediaType not found on DELETE 200 response")
	require.NotNil(t, mt.Examples)
	_, has := mt.Examples["example"]
	require.True(t, has, "example key not found in mediaType examples")
}

func TestBuilder_NewSetters(t *testing.T) {
	b := oas.NewBuilder().SetTitle("API Setters").SetVersion("2.0.0")

	// cria um GET e usa os novos setters
	b.Path("/setters").
		Get("Operation with setters").
		SetOperationID("op123").
		SetParameters(
			oas.ParameterOrRef{Param: &oas.Parameter{Name: "foo", In: oas.InQuery}},
			oas.ParameterOrRef{Param: &oas.Parameter{Name: "bar", In: oas.InHeader}},
		).
		SetRequestBody(oas.RequestBodyOrRef{
			Body: &oas.RequestBody{
				Description: oas.Ptr("body desc"),
				Required:    oas.Ptr(true),
			},
		}).
		SetResponses(oas.Responses{
			"400": {Resp: &oas.Response{Description: "bad request"}},
			"500": {Resp: &oas.Response{Description: "internal error"}},
		}).
		AddServer("http://op.example.com", "op server").
		DoneOp().
		DonePath()

	data, err := b.JSON()
	require.NoError(t, err)

	var doc oas.Document
	require.NoError(t, json.Unmarshal(data, &doc))

	op := doc.Paths["/setters"].PathItem.Get
	require.NotNil(t, op)
	require.Equal(t, "op123", *op.OperationID)
	require.Len(t, op.Parameters, 2)
	require.NotNil(t, op.RequestBody)
	require.Contains(t, op.Responses, "400")
	require.Contains(t, op.Responses, "500")
	require.Len(t, op.Servers, 1)
	require.Equal(t, "http://op.example.com", op.Servers[0].URL)
}
