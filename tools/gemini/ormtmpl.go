package main

import "fmt"

var ormTmpl = fmt.Sprintf(`package model

import (
	{{if has . "time"}}"time"
	"github.com/Jarnpher553/gemini/model/orm"{{else}}"github.com/Jarnpher553/gemini/model/orm"{{end}}
)
{{range .}} 
type {{ title .Name}} struct { {{range .Fields}} 
	{{if ne .Name "id"}}{{title .Name}}	{{if .Nullable}}*{{end}}{{if eq .Type "uuid"}}uuid.GUID{{else if eq .Type "time"}}time.Time{{else}}{{.Type}}{{end}} {{if not .Nullable}}%s{{end}}{{else}}{{if eq .Type "int"}}orm.ModelInt{{else}}orm.ModelUUID{{end}}{{end}} {{end}}
} {{end}}
`, "`gorm:\"not null\"`")
