package oas_test

import (
	"encoding/json"
	"testing"

	oas "github.com/leandroluk/go-oas/v3_1"
	"github.com/stretchr/testify/require"
)

func TestStringOrStringArray_JSON(t *testing.T) {
	// string única
	{
		var s oas.StringOrArray
		require.NoError(t, json.Unmarshal([]byte(`"foo"`), &s))
		require.Equal(t, "foo", *s.One)
		out, _ := json.Marshal(s)
		require.Equal(t, `"foo"`, string(out))
	}
	// array
	{
		var s oas.StringOrArray
		require.NoError(t, json.Unmarshal([]byte(`["a","b"]`), &s))
		require.Equal(t, []string{"a", "b"}, s.Many)
		out, _ := json.Marshal(s)
		require.Contains(t, string(out), `"a"`)
	}
	// inválido
	{
		var s oas.StringOrArray
		err := json.Unmarshal([]byte(`123`), &s)
		require.Error(t, err)
	}
}

func TestSchemaOrRef_JSON(t *testing.T) {
	// ref
	{
		var sr oas.SchemaOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/schemas/User"}`), &sr))
		require.NotNil(t, sr.Ref)
		out, _ := json.Marshal(sr)
		require.Contains(t, string(out), `"$ref"`)
	}
	// schema válido
	{
		var sr oas.SchemaOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"type":"string"}`), &sr))
		require.NotNil(t, sr.Schema)
		out, _ := json.Marshal(sr)
		require.Contains(t, string(out), `"type"`)
	}
	// schema inválido
	{
		var sr oas.SchemaOrRef
		err := json.Unmarshal([]byte(`{"type":123}`), &sr)
		require.Error(t, err)
	}
	// json não objeto
	{
		var sr oas.SchemaOrRef
		err := json.Unmarshal([]byte(`123`), &sr)
		require.Error(t, err)
	}
}

func TestItems_JSON(t *testing.T) {
	// single schema
	{
		var it oas.Items
		require.NoError(t, json.Unmarshal([]byte(`{"type":"string"}`), &it))
		require.NotNil(t, it.Single)
		out, _ := json.Marshal(it)
		require.Contains(t, string(out), `"type"`)
	}
	// lista
	{
		var it oas.Items
		require.NoError(t, json.Unmarshal([]byte(`[{"type":"string"},{"type":"integer"}]`), &it))
		require.Len(t, it.List, 2)
		out, _ := json.Marshal(it)
		require.Contains(t, string(out), `"integer"`)
	}
	// erro
	{
		var it oas.Items
		err := json.Unmarshal([]byte(`123`), &it)
		require.Error(t, err)
	}
}

func TestAdditionalProperties_JSON(t *testing.T) {
	// bool
	{
		var ap oas.AdditionalProperties
		require.NoError(t, json.Unmarshal([]byte(`true`), &ap))
		require.NotNil(t, ap.Allows)
		out, _ := json.Marshal(ap)
		require.Equal(t, "true", string(out))
	}
	// schema
	{
		var ap oas.AdditionalProperties
		require.NoError(t, json.Unmarshal([]byte(`{"type":"string"}`), &ap))
		require.NotNil(t, ap.Schema)
		out, _ := json.Marshal(ap)
		require.Contains(t, string(out), `"type"`)
	}
	// erro
	{
		var ap oas.AdditionalProperties
		err := json.Unmarshal([]byte(`123`), &ap)
		require.Error(t, err)
	}
}

func TestPathItemOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var pir oas.PathItemOrRef
		err := json.Unmarshal([]byte(`123`), &pir)
		require.Error(t, err)
	}
	// ref
	{
		var pir oas.PathItemOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/pathItems/Foo"}`), &pir))
		require.NotNil(t, pir.Ref)
		out, _ := json.Marshal(pir)
		require.Contains(t, string(out), `"$ref"`)
	}
	// válido
	{
		var pir oas.PathItemOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"get":{"responses":{"200":{"description":"ok"}}}}`), &pir))
		require.NotNil(t, pir.PathItem)
	}
	// inválido
	{
		var pir oas.PathItemOrRef
		err := json.Unmarshal([]byte(`{"get":123}`), &pir)
		require.Error(t, err)
	}
}

func TestParameterOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var pr oas.ParameterOrRef
		err := json.Unmarshal([]byte(`123`), &pr)
		require.Error(t, err)
	}
	// ref
	{
		var pr oas.ParameterOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/parameters/Bar"}`), &pr))
		require.NotNil(t, pr.Ref)
	}
	// válido
	{
		var pr oas.ParameterOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"name":"q","in":"query","required":true,"schema":{"type":"string"}}`), &pr))
		require.NotNil(t, pr.Param)
		out, _ := json.Marshal(pr)
		require.Contains(t, string(out), `"name":"q"`)
	}
	// inválido
	{
		var pr oas.ParameterOrRef
		err := json.Unmarshal([]byte(`{"name":123}`), &pr)
		require.Error(t, err)
	}
	// branch Marshal com Ref
	{
		pr := oas.ParameterOrRef{Ref: &oas.Reference{Ref: "#/components/parameters/Foo"}}
		out, err := json.Marshal(pr)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/parameters/Foo")
	}
	// branch Marshal com Param
	{
		pr := oas.ParameterOrRef{Param: &oas.Parameter{Name: "q", In: oas.InQuery}}
		out, err := json.Marshal(pr)
		require.NoError(t, err)
		require.Contains(t, string(out), `"q"`)
	}
}

func TestRequestBodyOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var rb oas.RequestBodyOrRef
		err := json.Unmarshal([]byte(`123`), &rb)
		require.Error(t, err)
	}
	// ref
	{
		var rb oas.RequestBodyOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/requestBodies/Baz"}`), &rb))
		require.NotNil(t, rb.Ref)
	}
	// válido
	{
		var rb oas.RequestBodyOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"content":{"application/json":{"schema":{"type":"string"}}}}`), &rb))
		require.NotNil(t, rb.Body)
	}
	// inválido
	{
		var rb oas.RequestBodyOrRef
		err := json.Unmarshal([]byte(`{"content":123}`), &rb)
		require.Error(t, err)
	}
	// Marshal branch Ref
	{
		r := oas.RequestBodyOrRef{Ref: &oas.Reference{Ref: "#/components/requestBodies/Foo"}}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/requestBodies/Foo")
	}
	// Marshal branch Body
	{
		r := oas.RequestBodyOrRef{
			Body: &oas.RequestBody{
				Content: map[string]oas.MediaType{
					"application/json": {
						Schema: &oas.SchemaOrRef{Schema: &oas.Schema{
							Type: &oas.StringOrArray{One: oas.Ptr("string")},
						}},
					},
				},
			},
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		require.Contains(t, string(out), "application/json")
	}
}

func TestResponseOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var rr oas.ResponseOrRef
		err := json.Unmarshal([]byte(`123`), &rr)
		require.Error(t, err)
	}
	// ref
	{
		var rr oas.ResponseOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/responses/Foo"}`), &rr))
		require.NotNil(t, rr.Ref)
	}
	// válido
	{
		var rr oas.ResponseOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"description":"ok"}`), &rr))
		require.NotNil(t, rr.Resp)
	}
	// inválido
	{
		var rr oas.ResponseOrRef
		err := json.Unmarshal([]byte(`{"description":123}`), &rr)
		require.Error(t, err)
	}
	// Marshal com Ref
	{
		r := oas.ResponseOrRef{Ref: &oas.Reference{Ref: "#/components/responses/Foo"}}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/responses/Foo")
	}
	// Marshal com Resp
	{
		r := oas.ResponseOrRef{
			Resp: &oas.Response{Description: "ok"},
		}
		out, err := json.Marshal(r)
		require.NoError(t, err)
		require.Contains(t, string(out), `"ok"`)
	}
}

func TestHeaderOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var h oas.HeaderOrRef
		err := json.Unmarshal([]byte(`123`), &h)
		require.Error(t, err)
	}
	// ref
	{
		var h oas.HeaderOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/headers/X-Rate-Limit"}`), &h))
		require.NotNil(t, h.Ref)
	}
	// válido
	{
		var h oas.HeaderOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"description":"limite"}`), &h))
		require.NotNil(t, h.Header)
	}
	// inválido
	{
		var h oas.HeaderOrRef
		err := json.Unmarshal([]byte(`{"description":123}`), &h)
		require.Error(t, err)
	}
	// Marshal com Ref
	{
		h := oas.HeaderOrRef{Ref: &oas.Reference{Ref: "#/components/headers/X-Rate-Limit"}}
		out, err := json.Marshal(h)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/headers/X-Rate-Limit")
	}
	// Marshal com Header
	{
		h := oas.HeaderOrRef{Header: &oas.Header{Description: oas.Ptr("limite")}}
		out, err := json.Marshal(h)
		require.NoError(t, err)
		require.Contains(t, string(out), `"limite"`)
	}
}

func TestExampleOrRef_JSON(t *testing.T) {
	// inválido
	{
		var ex oas.ExampleOrRef
		err := json.Unmarshal([]byte(`123`), &ex)
		require.Error(t, err)
	}
	// ref
	{
		var ex oas.ExampleOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"$ref":"#/components/examples/MyExample"}`), &ex))
		require.NotNil(t, ex.Ref)
		out, _ := json.Marshal(ex)
		require.Contains(t, string(out), "#/components/examples/MyExample")
	}
	// válido
	{
		var ex oas.ExampleOrRef
		require.NoError(t, json.Unmarshal([]byte(`{"summary":"um exemplo","value":{"foo":"bar"}}`), &ex))
		require.NotNil(t, ex.Example)
	}
	// inválido dentro de Example
	{
		var ex oas.ExampleOrRef
		err := json.Unmarshal([]byte(`{"summary":123}`), &ex)
		require.Error(t, err)
	}
	// Marshal com Ref
	{
		ex := oas.ExampleOrRef{Ref: &oas.Reference{Ref: "#/components/examples/Other"}}
		out, err := json.Marshal(ex)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/examples/Other")
	}
	// Marshal com Example
	{
		ex := oas.ExampleOrRef{
			Example: &oas.Example{Summary: oas.Ptr("manual"), Value: map[string]any{"x": 1}},
		}
		out, err := json.Marshal(ex)
		require.NoError(t, err)
		require.Contains(t, string(out), `"manual"`)
	}
}

func TestLinkOrRef_JSON(t *testing.T) {
	// erro inicial (json inválido → cobre primeiro return err)
	{
		var l oas.LinkOrRef
		err := json.Unmarshal([]byte(`123`), &l)
		require.Error(t, err)
	}

	// com $ref (cobre branch Ref no Unmarshal e Marshal)
	{
		var l oas.LinkOrRef
		data := []byte(`{"$ref":"#/components/links/MyLink"}`)
		require.NoError(t, json.Unmarshal(data, &l))
		require.NotNil(t, l.Ref)
		out, err := json.Marshal(l)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/links/MyLink")
	}

	// com objeto Link válido (cobre branch Link no Unmarshal e Marshal)
	{
		var l oas.LinkOrRef
		data := []byte(`{"operationId":"getUser","parameters":{"id":"123"}}`)
		require.NoError(t, json.Unmarshal(data, &l))
		require.NotNil(t, l.Link)
		out, err := json.Marshal(l)
		require.NoError(t, err)
		require.Contains(t, string(out), `"getUser"`)
		require.Contains(t, string(out), `"id"`)
	}

	// objeto inválido (cobre segundo return err)
	{
		var l oas.LinkOrRef
		data := []byte(`{"operationId":123}`)
		err := json.Unmarshal(data, &l)
		require.Error(t, err)
	}
}

func TestCallbackOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var c oas.CallbackOrRef
		err := json.Unmarshal([]byte(`123`), &c)
		require.Error(t, err)
	}
	// ref
	{
		var c oas.CallbackOrRef
		data := []byte(`{"$ref":"#/components/callbacks/MyCb"}`)
		require.NoError(t, json.Unmarshal(data, &c))
		require.NotNil(t, c.Ref)
		out, err := json.Marshal(c)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/callbacks/MyCb")
	}
	// callback válido
	{
		var c oas.CallbackOrRef
		data := []byte(`{"onEvent":{"get":{"responses":{"200":{"description":"ok"}}}}}`)
		require.NoError(t, json.Unmarshal(data, &c))
		require.NotNil(t, c.Callback)
		out, err := json.Marshal(c)
		require.NoError(t, err)
		require.Contains(t, string(out), `"onEvent"`)
	}
	// callback inválido
	{
		var c oas.CallbackOrRef
		data := []byte(`{"onEvent":123}`)
		err := json.Unmarshal(data, &c)
		require.Error(t, err)
	}
}

func TestSecuritySchemeOrRef_JSON(t *testing.T) {
	// erro inicial
	{
		var s oas.SecuritySchemeOrRef
		err := json.Unmarshal([]byte(`123`), &s)
		require.Error(t, err)
	}

	// ref
	{
		var s oas.SecuritySchemeOrRef
		data := []byte(`{"$ref":"#/components/securitySchemes/MySec"}`)
		require.NoError(t, json.Unmarshal(data, &s))
		require.NotNil(t, s.Ref)
		out, err := json.Marshal(s)
		require.NoError(t, err)
		require.Contains(t, string(out), "#/components/securitySchemes/MySec")
	}

	// válido
	{
		var s oas.SecuritySchemeOrRef
		data := []byte(`{"type":"http","scheme":"bearer"}`)
		require.NoError(t, json.Unmarshal(data, &s))
		require.NotNil(t, s.Scheme)
		out, err := json.Marshal(s)
		require.NoError(t, err)
		require.Contains(t, string(out), `"bearer"`)
	}

	// inválido
	{
		var s oas.SecuritySchemeOrRef
		data := []byte(`{"type":123}`)
		err := json.Unmarshal(data, &s)
		require.Error(t, err)
	}
}

func TestOperation_ValidateRequiredResponses(t *testing.T) {
	// nil receiver → cobre return nil
	var opNil *oas.Operation
	require.NoError(t, opNil.ValidateRequiredResponses())

	// vazio → cobre erro responses
	op := &oas.Operation{Responses: oas.Responses{}}
	err := op.ValidateRequiredResponses()
	require.Error(t, err)

	// com responses → cobre caminho feliz e sort
	op = &oas.Operation{Responses: oas.Responses{
		"200": {Resp: &oas.Response{Description: "ok"}},
		"404": {Resp: &oas.Response{Description: "nope"}},
	}}
	err = op.ValidateRequiredResponses()
	require.NoError(t, err)
}

func TestNewObjectSchema(t *testing.T) {
	s := oas.NewObjectSchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "object", *s.Type.One)
	require.Nil(t, s.Type.Many)
}

func TestNewArraySchema(t *testing.T) {
	s := oas.NewArraySchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "array", *s.Type.One)
	require.Nil(t, s.Type.Many)
}

func TestNewStringSchema(t *testing.T) {
	s := oas.NewStringSchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "string", *s.Type.One)
	require.Nil(t, s.Type.Many)
}

func TestNewIntegerSchema(t *testing.T) {
	s := oas.NewIntegerSchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "integer", *s.Type.One)
	require.Nil(t, s.Type.Many)
}

func TestNewNumberSchema(t *testing.T) {
	s := oas.NewNumberSchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "number", *s.Type.One)
	require.Nil(t, s.Type.Many)
}

func TestNewBooleanSchema(t *testing.T) {
	s := oas.NewBooleanSchema()
	require.NotNil(t, s.Type)
	require.NotNil(t, s.Type.One)
	require.Equal(t, "boolean", *s.Type.One)
	require.Nil(t, s.Type.Many)
}
