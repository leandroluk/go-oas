# go-openapi

<img align="right" width="180px" src="https://raw.githubusercontent.com/leandroluk/go-openapi/refs/heads/master/assets/go-openapi.png">

[![Build Status](https://github.com/leandroluk/go-openapi/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/leandroluk/go-openapi/actions)  
[![Coverage Status](https://img.shields.io/codecov/c/github/leandroluk/go-openapi/main.svg)](https://codecov.io/gh/leandroluk/go-openapi)  
[![Go Report Card](https://goreportcard.com/badge/github.com/leandroluk/go-openapi)](https://goreportcard.com/report/github.com/leandroluk/go-openapi)  
[![Go Doc](https://godoc.org/github.com/leandroluk/go-openapi?status.svg)](https://pkg.go.dev/github.com/leandroluk/go-openapi)  
[![Release](https://img.shields.io/github/release/leandroluk/go-openapi.svg?style=flat-square)](https://github.com/leandroluk/go-openapi/releases)  

A schema description & validation library for Go, inspired by [hapi/joi](https://github.com/hapijs/joi).

---

## Contents
- [go-openapi](#go-openapi)
  - [Contents](#contents)
  - [Getting started](#getting-started)
  - [Usage](#usage)
  - [Available Schemas](#available-schemas)
  - [Examples](#examples)
    - [String Validation](#string-validation)
    - [Number Validation](#number-validation)
    - [Boolean Validation](#boolean-validation)
    - [Object Validation](#object-validation)
  - [Implementation Status](#implementation-status)
  - [About the Project](#about-the-project)
  - [Contributors](#contributors)
  - [License](#license)

---

## Getting started

1. Install the package:

```sh
go get github.com/leandroluk/go-openapi
```

2. Import it in your project:

```go
import "github.com/leandroluk/go-openapi"
```

3. Create and validate schemas:

```go
schema := joi.Object(map[string]joi.Schema{
    "username": joi.String().Min(3).Max(20).Required(),
    "age":      joi.Number().Min(18).Required(),
    "email":    joi.String().Regex(regexp.MustCompile(`.+@.+\..+`)),
})

value := map[string]any{
    "username": "john_doe",
    "age":      25,
    "email":    "john@example.com",
}

parsed, errs := schema.Validate("user", value)
if len(errs) > 0 {
    fmt.Println("Validation failed:", errs)
} else {
    fmt.Println("Validation passed:", parsed)
}
```

---

## Usage

- **Basic types**: `String`, `Number`, `Boolean`, `Object`
- **Rules**:  
  - String: `.Min()`, `.Max()`, `.Regex()`, `.Trim()`, `.Lowercase()`, `.Uppercase()`  
  - Number: `.Min()`, `.Max()`, `.Integer()`, `.Positive()`, `.Negative()`  
  - Boolean: `.True()`, `.False()`, `.Truthy()`, `.Falsy()`  
  - Object: `.Min()`, `.Max()`, `.Length()`, `.Unknown()`  

---

## Available Schemas
- ✅ String
- ✅ Number
- ✅ Boolean
- ✅ Object
- ⬜ Array *(coming soon)*

---

## Examples

### String Validation
```go
joi.String().Min(5).Max(10).Trim()
```

### Number Validation
```go
joi.Number().Integer().Positive()
```

### Boolean Validation
```go
joi.Boolean().Truthy("yes", "1").Falsy("no", "0")
```

### Object Validation
```go
joi.Object(map[string]joi.Schema{
    "id":   joi.Number().Required(),
    "name": joi.String().Min(3),
}).Unknown(false)
```

---

## Implementation Status
- [x] String rules
- [x] Number rules
- [x] Boolean rules
- [x] Object rules
- [ ] Array rules
- [ ] Custom extensions

---

## About the Project
This project was inspired by [joi](https://github.com/hapijs/joi) for Node.js and aims to bring a similar developer experience to Go.

---

## Contributors
Thanks to all the people who contribute! [[Contribute](CONTRIBUTING.md)]

---

## License
MIT License – see [LICENSE](LICENSE) file for details.
