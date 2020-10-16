package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Object struct {
	Name        string
	ServiceList []string
}

type Service struct {
	Name string `mapstructure:"name"`
	Orm  Orm    `mapstructure:"orm"`
	Dto  Dto    `mapstructure:"dto"`
}

type Field struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	Nullable bool   `mapstructure:"nullable"`
}

type DtoField struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"`
	Required bool   `mapstructure:"required"`
}

type Dto struct {
	Request  DtoObj `mapstructure:"request"`
	Response DtoObj `mapstructure:"response"`
}

type Orm struct {
	Name   string  `mapstructure:"name"`
	Fields []Field `mapstructure:"fields"`
}

type DtoObj struct {
	Name   string     `mapstructure:"name"`
	Fields []DtoField `mapstructure:"fields"`
}

var name string

//var services string
var path string

func main() {
	rootCmd := &cobra.Command{
		Use:     "gemini",
		Short:   "Simple method of generate project",
		Example: "gemini -n demo -p .",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("start generate...")

			reader := strings.NewReader(gen)
			viper.SetConfigType("json")
			_ = viper.ReadConfig(reader)

			var s = make([]Service, 0)
			_ = viper.UnmarshalKey("services", &s)

			for i := range s {
				s[i].Name = strings.ToLower(s[i].Name) + "Service"
			}

			err := os.MkdirAll(filepath.Join(path, name), os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}

			_ = os.Chdir(filepath.Join(path, name))

			//新建文件夹
			//_ = os.MkdirAll("./logs", os.ModePerm)
			CreateDirectory()

			//新建文件
			f, _ := os.Create("config.toml")
			_, _ = f.WriteString(toml)

			f, _ = os.Create("go.mod")
			_, _ = f.WriteString("module " + name + "\n")

			f, _ = os.Create("main.go")
			mt := template.Must(template.New("main").Funcs(template.FuncMap{"title": strings.Title, "join": strings.Join, "name": Name, "services": Services}).Parse(mainTmpl))

			_ = mt.Execute(f, s)

			f, _ = os.Create("./model/dto.go")
			dt := template.Must(template.New("dto").Funcs(template.FuncMap{"title": strings.Title, "has": DtoHas}).Parse(dtoTmpl))
			_ = dt.Execute(f, []Dto{s[0].Dto})

			f, _ = os.Create("./model/orm.go")
			ot := template.Must(template.New("orm").Funcs(template.FuncMap{"title": strings.Title, "has": OrmHas}).Parse(ormTmpl))
			_ = ot.Execute(f, []Orm{s[0].Orm})

			Default(s)
		},
	}

	rootCmd.Flags().StringVarP(&name, "name", "n", "default", "your project name")
	//rootCmd.LocalFlags().StringVarP(&services, "services", "s", "default", "your services (delimit commas)")
	rootCmd.Flags().StringVarP(&path, "path", "p", "./", "your directory path")

	initCmd := &cobra.Command{
		Use:     "init",
		Short:   "Create a default template of gen file",
		Example: "gemini init",
		Run: func(cmd *cobra.Command, args []string) {
			f, _ := os.Create("./gen.json")
			_, _ = f.WriteString(gen)
		},
	}

	genCmd := &cobra.Command{
		Use:     "gen",
		Short:   "Generate project with gen file",
		Example: "gemini gen",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("start generate...")

			viper.SetConfigName("gen")
			viper.SetConfigType("json")
			viper.AddConfigPath("./")

			_, err := os.Stat("./gen.json")
			if err != nil {
				if !os.IsExist(err) {
					log.Fatalln(err)
				}
			}

			err = viper.ReadInConfig()
			if err != nil {
				log.Fatalln(err)
			}
			name = viper.GetString("name")
			path = viper.GetString("path")

			var s = make([]Service, 0)
			if err := viper.UnmarshalKey("services", &s); err != nil {
				log.Fatalln("gen.json format error")
			}

			for i := range s {
				s[i].Name = strings.ToLower(s[i].Name) + "Service"
			}

			err = os.MkdirAll(filepath.Join(path, name), os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}

			_ = os.Chdir(filepath.Join(path, name))

			//新建文件夹
			CreateDirectory()

			//新建文件
			f, _ := os.Create("config.yaml")
			_, _ = f.WriteString(toml)

			f, _ = os.Create("go.mod")
			_, _ = f.WriteString("module " + name + "\n")

			f, _ = os.Create("main.go")
			mt := template.Must(template.New("main").Funcs(template.FuncMap{"title": strings.Title, "join": strings.Join, "name": Name, "services": Services}).Parse(mainTmpl))

			_ = mt.Execute(f, s)

			f, _ = os.Create("./model/dto.go")
			dt := template.Must(template.New("dto").Funcs(template.FuncMap{"title": strings.Title, "has": DtoHas}).Parse(dtoTmpl))

			var dtoList []Dto
			for _, v := range s {
				dtoList = append(dtoList, v.Dto)
			}
			_ = dt.Execute(f, &dtoList)

			f, _ = os.Create("./model/orm.go")
			ot := template.Must(template.New("orm").Funcs(template.FuncMap{"title": strings.Title, "has": OrmHas}).Parse(ormTmpl))

			var ormList []Orm
			for _, v := range s {
				ormList = append(ormList, v.Orm)
			}

			_ = ot.Execute(f, &ormList)

			Default(s)
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func DtoHas(list []Dto, t string) bool {

	for _, v := range list {
		for _, vv := range v.Request.Fields {
			if vv.Type == t {
				return true
			}
		}
		for _, vv := range v.Response.Fields {
			if vv.Type == t {
				return true
			}
		}
	}
	return false
}

func OrmHas(list []Orm, t string) bool {

	for _, v := range list {
		for _, vv := range v.Fields {
			if vv.Type == t {
				return true
			}
		}
	}
	return false
}

func Name() string {
	return name
}

func Services(s []Service) string {
	var out []string
	for _, v := range s {
		out = append(out, v.Name)
	}
	return strings.Join(out, ", ")
}

func Default(s []Service) {
	f, _ := os.Create("./validators/phone.go")
	_, _ = f.WriteString(`package validators

import (
	"github.com/Jarnpher553/gemini/validator"
	"regexp"
)

func demo(v *validator.Validate, fl validator.FieldLevel) bool {
	if f, ok := fl.Field().Interface().(string); ok {
		reg := regexp.MustCompile("1\\d{10}")
		ret := reg.MatchString(f)
		return ret
	}
	return true
}
`)

	f, _ = os.Create("./validators/validator.go")
	_, _ = f.WriteString(`package validators

import (
	"github.com/Jarnpher553/gemini/validator"
)

func init() {
	validator.Register("demo", demo)
}
`)

	f, _ = os.Create("./middlewares/permission.go")
	_, _ = f.WriteString(`package middlewares

import (
	"github.com/Jarnpher553/gemini/service"
)

func Permission(code string) service.Middleware {
	return func(baseService service.IBaseService) service.HandlerFunc {
		return func(ctx *service.Ctx) {
			// do something here
		}
	}
}
`)

	st := template.Must(template.New("service").Funcs(template.FuncMap{"title": strings.Title, "name": Name, "trimSuffix": strings.TrimSuffix}).Parse(serviceTmpl))
	for _, v := range s {
		f, _ = os.Create("./services/" + v.Name + ".go")
		_ = st.Execute(f, v)
	}

	f, _ = os.Create("./schedules/schedule.go")
	_, _ = f.WriteString(`package schedules

import (
	"github.com/Jarnpher553/gemini/scheduler"
)

func init() {
	scheduler.Assign(scheduler.Every(5*scheduler.Second), demo)
}
`)

	f, _ = os.Create("./schedules/demo.go")
	_, _ = f.WriteString(`package schedules

import (
	"github.com/Jarnpher553/gemini/scheduler"
)

func demo(ops *scheduler.Options){
	//Doing something...
}
`)
	f, _ = os.Create("./errs/error.go")
	_, _ = f.WriteString(`package errs

import "github.com/Jarnpher553/gemini/erro"

const (
	ErrDemo = 5000
)

func init() {
	erro.Register(ErrDemo, "有误")
}`)
	f, _ = os.Create("./version/version.go")
	_, _ = f.WriteString(`package version

var (
	Release   = "unset"
	Commit    = "unset"
	BuildTime = "unset"
)
`)

	st = template.Must(template.New("version").Funcs(template.FuncMap{"name": Name}).Parse(versionTmpl))
	f, _ = os.Create("./services/versionService.go")
	_ = st.Execute(f, nil)

	st = template.Must(template.New("makefile").Funcs(template.FuncMap{"name": Name}).Parse(makefileTmpl))
	f, _ = os.Create("./makefile")
	_ = st.Execute(f, nil)

	log.Println("generate success")
}

func CreateDirectory() {
	_ = os.MkdirAll("./model", os.ModePerm)
	_ = os.MkdirAll("./services", os.ModePerm)
	_ = os.MkdirAll("./validators", os.ModePerm)
	_ = os.MkdirAll("./middlewares", os.ModePerm)
	_ = os.MkdirAll("./schedules", os.ModePerm)
	_ = os.MkdirAll("./errs", os.ModePerm)
	_ = os.MkdirAll("./version", os.ModePerm)
}
