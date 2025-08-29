# go-oas

<img align="right" width="180px" src="https://raw.githubusercontent.com/leandroluk/go-oas/refs/heads/master/assets/go-oas.png">

[![Build Status](https://github.com/leandroluk/go-oas/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/leandroluk/go-oas/actions)  
[![Coverage Status](https://img.shields.io/codecov/c/github/leandroluk/go-oas/master.svg)](https://codecov.io/gh/leandroluk/go-oas)  
[![Go Report Card](https://goreportcard.com/badge/github.com/leandroluk/go-oas)](https://goreportcard.com/report/github.com/leandroluk/go-oas)  
[![Go Doc](https://godoc.org/github.com/leandroluk/go-oas?status.svg)](https://pkg.go.dev/github.com/leandroluk/go-oas)  
[![Release](https://img.shields.io/github/release/leandroluk/go-oas.svg?style=flat-square)](https://github.com/leandroluk/go-oas/releases)  

Uma biblioteca em Go para **modelagem, serialização e construção de documentos OpenAPI 3.1**, com suporte a `$ref`, schemas JSON, builders fluentes e integração com **Gin**.

---

## Instalação

```sh
go get github.com/leandroluk/go-oas
```

---

## Exemplo Rápido

```go
package main

import (
    "encoding/json"
    "fmt"
    oas "github.com/leandroluk/go-oas/v3_1"
)

func main() {
    b := oas.NewBuilder().
        Info("API Example", "1.0.0").
        Path("/items").
        Get().
            Summary("Lista itens").
            Response(200, "ok", oas.SchemaOrRef{
                Schema: &oas.Schema{Type: &oas.StringOrStringArray{One: oas.Ptr("string")}},
            }).
            Done().
        Done()

    doc := b.Build()
    data, _ := json.MarshalIndent(doc, "", "  ")
    fmt.Println(string(data))
}
```

Saída (resumida):

```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "API Example",
    "version": "1.0.0"
  },
  "paths": {
    "/items": {
      "get": {
        "summary": "Lista itens",
        "responses": {
          "200": {
            "description": "ok",
            "content": {
              "application/json": {
                "schema": { "type": "string" }
              }
            }
          }
        }
      }
    }
  }
}
```

---

## Integração com Gin

```go
import "github.com/gin-gonic/gin"

r := gin.Default()

r.GET("/items", func(c *gin.Context) {
    c.JSON(200, []string{"item1", "item2"})
})

// Em paralelo você pode expor o documento OpenAPI:
r.GET("/openapi.json", func(c *gin.Context) {
    doc := b.Build()
    c.JSON(200, doc)
})
```

---

## Estrutura do Projeto

```
v3_1/
  builder.go    # Builder fluente para criar documentos OAS
  struct.go     # Definições das structs OpenAPI 3.1
v3_1_test/
  builder_test.go
  struct_test.go
```

---

## Testes

```sh
make test       # roda os testes sem cobertura
make test.ci    # roda os testes com cobertura
make test.html  # abre o relatório de cobertura
```

---

## Status

- ✅ Suporte completo ao OpenAPI 3.1 (Document, Schema, Components, Paths, etc.)  
- ✅ Builders fluentes para criar specs programaticamente  
- ✅ Integração simples com Gin  
- ✅ 100% de cobertura de testes em `struct.go`  
- 🚧 Futuro: suporte às versões 2.0 e 3.0 para migração de specs

---

## Licença

MIT © [Leandro Luk](https://github.com/leandroluk)
