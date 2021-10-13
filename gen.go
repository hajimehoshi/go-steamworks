// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const version = "151"

type SteamAPI struct {
	CallbackStructs []SteamAPICallbackStruct `json:"callback_structs"`
	Consts          []SteamAPIConst          `json:"consts"`
	Enums           []SteamAPIEnum           `json:"enums"`
	Interfaces      []SteamAPIInterface      `json:"interfaces"`
	Structs         []SteamAPIStruct         `json:"structs"`
	Typedefs        []SteamAPITypedef        `json:"typedefs"`
}

type SteamAPICallbackStruct struct {
	ID     int             `json:"callback_id"`
	Enums  []SteamAPIEnum  `json:"enums"`
	Fields []SteamAPIField `json:"fields"`
	Struct string          `json:"struct"`
}

func goType(str string) string {
	// Function pointer
	if strings.Contains(str, "(*)") {
		return "uintptr"
	}

	var tokens []string
	var unsigned bool
	var long bool
	for _, token := range strings.Split(str, " ") {
		token = strings.TrimSpace(token)
		if idx := strings.Index(token, "::"); idx >= 0 {
			token = token[idx+2:]
		}
		switch token {
		case "unsigned":
			unsigned = true
			continue
		case "char":
			if unsigned {
				token = "uint8"
			} else {
				token = "char"
			}
			unsigned = false
			long = false
		case "short":
			if unsigned {
				token = "uint16"
			} else {
				token = "int16"
			}
			unsigned = false
			long = false
		case "int":
			if unsigned {
				token = "uint32"
			} else {
				token = "int32"
			}
			unsigned = false
			long = false
		case "long":
			// Assume there is no single 'long'.
			if long {
				if unsigned {
					token = "uint64"
				} else {
					token = "int64"
				}
				unsigned = false
				long = false
			} else {
				long = true
				continue
			}
		case "int32_t":
			token = "int32"
		case "int64_t":
			token = "int64"
		case "float":
			token = "float32"
		case "double":
			token = "float64"
		case "const":
			continue
		case "&":
			token = "*"
		case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "bool":
			// Do nothing
		case "void":
			// Do nothing
		default:
			if isIdent(token) {
				token = strings.Title(token)
			}
		}
		tokens = append(tokens, token)
	}
	for i := 0; i < len(tokens)/2; i++ {
		j := len(tokens) - i - 1
		tokens[i], tokens[j] = tokens[j], tokens[i]
	}

	var pointer int
	for i, t := range tokens {
		switch t {
		case "*":
			pointer++
			continue
		case "char":
			if pointer > 0 {
				tokens[i-1] = ""
				tokens[i] = "string"
			} else {
				tokens[i] = "byte"
			}
			pointer = 0
		case "void":
			if pointer > 0 {
				tokens[i-1] = ""
				tokens[i] = "uintptr"
			} else {
				tokens[i] = ""
			}
		}
		pointer = 0
	}

	return strings.Join(tokens, "")
}

func goTypeDefaultValue(t string) string {
	if strings.HasPrefix(t, "*") {
		return "nil"
	}
	switch t {
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return "0"
	case "string":
		return `""`
	case "bool":
		return "false"
	}
	return fmt.Sprintf("*new(%s)", t)
}

func (s *SteamAPICallbackStruct) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "const %s_ID = %d\n", s.Struct, s.ID)
	for _, e := range s.Enums {
		b.WriteString(e.String())
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "type %s struct {\n", s.Struct)
	for _, f := range s.Fields {
		fmt.Fprintf(&b, "\t%s %s\n", strings.Title(f.Name), goType(f.Type))
	}
	b.WriteString("}")
	return b.String()
}

type SteamAPIConst struct {
	Name string `json:"constname"`
	Type string `json:"consttype"`
	Val  string `json:"constval"`
}

func isIdent(str string) bool {
	if len(str) == 0 {
		return false
	}

	for i, r := range str {
		if i == 0 {
			if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' {
				continue
			}
		} else {
			if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '_' {
				continue
			}
		}
		return false
	}
	return true
}

func (s *SteamAPIConst) String() string {
	name := s.Name
	name = strings.Title(name)

	value := s.Val
	if strings.HasSuffix(value, "ull") {
		value = value[:len(value)-3]
	}
	if value == "( SteamItemInstanceID_t ) ~ 0" {
		// Assume that SteamItemInstanceID_t is uint64.
		value = "^uint64(0)"
	}
	if value == "( ( uint32 ) 'd' << 16U ) | ( ( uint32 ) 'e' << 8U ) | ( uint32 ) 'v'" {
		value = "uint32('d')<<16 | uint32('e')<<8 | uint32('v')"
	}

	tokens := strings.Split(value, "|")
	for i, token := range tokens {
		t := strings.TrimSpace(token)
		if isIdent(t) {
			t = strings.Title(t)
		}
		tokens[i] = t
	}
	value = strings.Join(tokens, " | ")

	return fmt.Sprintf("const %s = %s", name, value)
}

type SteamAPIEnum struct {
	Name   string          `json:"enumname"`
	FqName string          `json:"fqname"`
	Values []SteamAPIValue `json:"values"`
}

func (s *SteamAPIEnum) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "type %s int\n", s.Name)
	b.WriteString("const (\n")
	for _, v := range s.Values {
		fmt.Fprintf(&b, "\t%s %s = %s\n", strings.Title(v.Name), s.Name, v.Value)
	}
	b.WriteString(")")
	return b.String()
}

type SteamAPIInterface struct {
	Accessors     []SteamAPIAccessor `json:"accessors"`
	ClassName     string             `json:"classname"`
	Fields        []SteamAPIField    `json:"fields"`
	Methods       []SteamAPIMethod   `json:"method"`
	VersionString string             `json:"version_string"`
}

type SteamAPIStruct struct {
	Fields  []SteamAPIField  `json:"fields"`
	Methods []SteamAPIMethod `json:"methods"`
	Struct  string           `json:"struct"`
}

func (s *SteamAPIStruct) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "type %s struct {\n", strings.Title(s.Struct))
	for _, f := range s.Fields {
		name := f.Name
		if f.Private {
			name = strings.ToLower(name[0:1]) + name[1:]
		} else {
			name = strings.Title(name)
		}
		fmt.Fprintf(&b, "\t%s %s\n", name, goType(f.Type))
	}
	b.WriteString("}")
	for _, m := range s.Methods {
		name := m.Name
		// A method might be an operator.
		if !isIdent(name) {
			tokens := strings.Split(m.NameFlat, "_")
			name = tokens[len(tokens)-1]
		}
		if name == "c_str" {
			name = "String"
		}
		b.WriteString("\n\n")
		var params []string
		for _, p := range m.Params {
			params = append(params, fmt.Sprintf("%s %s", p.Name, goType(p.Type)))
		}
		rt := goType(m.ReturnType)
		fmt.Fprintf(&b, "func (*%s) %s(%s) %s {\n", strings.Title(s.Struct), name, strings.Join(params, ", "), rt)
		fmt.Fprintf(&b, "\t// TODO: Invoke %s\n", m.NameFlat)
		if rt != "" {
			fmt.Fprintf(&b, "\treturn %s\n", goTypeDefaultValue(rt))
		}
		b.WriteString("}")
	}
	return b.String()
}

type SteamAPITypedef struct {
	Typedef string `json:"typedef"`
	Type    string `json:"type"`
}

func (s *SteamAPITypedef) String() string {
	switch s.Typedef {
	case "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64",
		"lint64", "ulint64", "intp", "uintp":
		return ""
	}
	return fmt.Sprintf("type %s %s", s.Typedef, goType(s.Type))
}

type SteamAPIField struct {
	Name    string `json:"fieldname"`
	Type    string `json:"fieldtype"`
	Private bool   `json:"private"`
}

type SteamAPIValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SteamAPIAccessor struct {
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	NameFlat string `json:"name_flat"`
}

type SteamAPIMethod struct {
	Name       string          `json:"methodname"`
	NameFlat   string          `json:"methodname_flat"`
	Params     []SteamAPIParam `json:"params"`
	ReturnType string          `json:"returntype"`
}

type SteamAPIParam struct {
	Name string `json:"paramname"`
	Type string `json:"paramtype"`
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "go-steamworks")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := processZip(dir); err != nil {
		return err
	}

	return nil
}

func processZip(dir string) error {
	zipfile, err := os.Open(fmt.Sprintf("steamworks_sdk_%s.zip", version))
	if err != nil {
		if os.IsNotExist(err) {
			const sdkURL = "https://partner.steamgames.com/downloads/steamworks_sdk_" + version + ".zip"
			return fmt.Errorf("steamworks_sdk.zip must exist; download it from %s with your Steamworks account", sdkURL)
		}
		return err
	}
	defer zipfile.Close()

	stat, err := zipfile.Stat()
	if err != nil {
		return err
	}
	r, err := zip.NewReader(zipfile, stat.Size())
	if err != nil {
		return err
	}

	for path, filename := range map[string]string{
		"sdk/redistributable_bin/linux32/libsteam_api.so": "libsteam_api.so",
		"sdk/redistributable_bin/linux64/libsteam_api.so": "libsteam_api64.so",
		"sdk/redistributable_bin/osx/libsteam_api.dylib":  "libsteam_api.dylib",
		"sdk/redistributable_bin/steam_api.dll":           "steam_api.dll",
		"sdk/redistributable_bin/win64/steam_api64.dll":   "steam_api64.dll",
	} {
		f, err := r.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var out io.Writer
		out, err = os.Create(filename)
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, f); err != nil {
			return err
		}
	}

	f, err := r.Open("sdk/public/steam/steam_api.json")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := writeAPI(f); err != nil {
		return err
	}

	return nil
}

const apiTmpl = `// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 The go-steamworks Authors

// Code generated by gen.go. DO NOT EDIT.

package steamworks

type CSteamID uint64
type CGameID uint64
{{range .Typedefs}}
{{.String}}{{end}}
{{range .Consts}}
{{.String}}{{end}}
{{range .Enums}}
{{.String}}
{{end}}
{{range .CallbackStructs}}
{{.String}}
{{end}}
{{range .Structs}}
{{.String}}
{{end}}
// TODO: Dump interfaces
`

func writeAPI(r io.Reader) error {
	var api SteamAPI
	dec := json.NewDecoder(r)
	if err := dec.Decode(&api); err != nil {
		return err
	}

	f, err := os.Create("api.go")
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("api.go").Parse(apiTmpl))
	if err := t.Execute(f, &api); err != nil {
		return err
	}

	var buf bytes.Buffer
	cmd := exec.Command("gofmt", "-s", "-w", "api.go")
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, buf.String())
	}

	return nil
}
