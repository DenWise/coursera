package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
	"text/template"
)

// код писать тут

var (
	commentStr string
	prefix     = "// apigen:api"
)

type Params struct {
	Url    string
	Auth   bool
	Method string
}

type TplSrv struct {
	StructName string
	Handlers   map[string]string
}

type Templates struct {
	Templates []*TplSrv
}

type TplWrap struct {
	StructName string
	MethodName string
}

var (
	serveTpl = template.Must(template.New("serveTpl").Parse(`
	func (s *{{.StructName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path { 
		{{range $k, $v := .Handlers}}
		case "{{$v}}":
			s.wrapper{{$k}}(w, r)
		{{end}}
		default:
			http.Error(w, "Method not found", 404)
		}
	}`))
	wrapperTpl = template.Must(template.New("wrapperTpl").Parse(`
	func (s *{{.StructName}}) wrapper{{.MethodName}}(w http.ResponseWriter, r *http.Request) {
		
	}
`))
)

type Method struct {
	Name   string
	Params *Params
}

func main() {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Printf("%#v\n", node.Name)
	fmt.Fprintln(out, "package "+node.Name.Name)
	fmt.Fprintln(out)
	fmt.Fprintln(out)

	methods := make([]*Method, 0)
	urls := make(map[string]string, 0)
	handlers := make(map[string][]string)

	for _, d := range node.Decls {
		f, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if f.Doc == nil {
			continue
		}

		needCodegen := false
		for _, comment := range f.Doc.List {
			needCodegen = needCodegen || strings.HasPrefix(comment.Text, prefix)
			commentStr = comment.Text
		}
		if !needCodegen {
			continue
		}

		params := &Params{}
		genParamStr := strings.TrimPrefix(commentStr, prefix)
		json.Unmarshal([]byte(genParamStr), params)

		structName := f.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Obj.Name
		handlers[structName] = append(handlers[structName],f.Name.Name)
		urls[f.Name.Name] = params.Url

		methods = append(methods, &Method{
			Name:   f.Name.Name,
			Params: params,
		})
	}

	tpls := new(Templates)

	for s, f := range handlers {
		tpl := &TplSrv{
			StructName: s,
			Handlers: make(map[string]string),
		}
		for _, fn := range f {
			tpl.Handlers[fn] = urls[fn]
		}
		tpls.Templates = append(tpls.Templates,tpl)
	}

	for _, t := range tpls.Templates {
		err = serveTpl.Execute(out, t)
		if err != nil {
			log.Fatal(err)
		}
	}

}
