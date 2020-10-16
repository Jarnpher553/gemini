package main

const dtoTmpl = "package model\n" +
	"\n" +
	"import (\n" +
	"	{{if has . \"uuid\"}}\"github.com/Jarnpher553/gemini/uuid\"\n{{end}}" +
	"	{{if has . \"time\"}}\"time\"\n{{end}}" +
	")\n" +
	"\n" +
	"{{range .}}" +
	"type {{ title .Request.Name}} struct {\n" +
	"{{range .Request.Fields}}" +
	"	{{title .Name}}\t{{if eq .Type \"uuid\"}}uuid.GUID{{else if eq .Type \"time\"}}time.Time{{else}}{{.Type}}{{end}} `json:\"{{.Name}}\" form:\"{{.Name}}\" {{if .Required}}binding:\"required\"{{end}}`\n" +
	"{{end}}" +
	"}" +
	"\n" +
	"type {{ title .Response.Name}} struct {\n" +
	"{{range .Response.Fields}}" +
	"	{{title .Name}}\t{{if eq .Type \"uuid\"}}uuid.GUID{{else if eq .Type \"time\"}}time.Time{{else}}{{.Type}}{{end}} `json:\"{{.Name}}\"`\n" +
	"{{end}}" +
	"}" +
	"{{end}}"
