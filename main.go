package main

import (
	"flag"
	"go/parser"
	"log"
	"go/token"
	"go/doc"
	"regexp"
	"go/ast"
	"fmt"
	"strings"
	"text/template"
	"os"
	"path/filepath"
	"bytes"
	"reflect"
)

type ServiceGenData struct {
	Package string
	Service string
	Methods []string
	Packages map[string]*bool
}

type MethodGenData struct {
	Method string
	ArgsBlock string
	ResultBlock string
}

func main() {
	var hDir string
	flag.StringVar(&hDir, "s", "../../../service", "Handler dir")
	flag.Parse()

	regExpFile, err := regexp.Compile(`.*/jrpc2hh_.+\.go`)
	if err != nil {
		log.Fatal(err)
	}
	dFun := func(path string, f os.FileInfo, err error) error {
		if regExpFile.Match([]byte(path)) {
			os.Remove(path)
		}
		return nil
	}
	filepath.Walk(hDir, dFun)

	fs := token.NewFileSet()
	d, err := parser.ParseDir(fs, hDir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	regExpService, err := regexp.Compile("jrpc2hh:service\\n")
	if err != nil {
		log.Fatal(err)
	}
	regExpMethod, err := regexp.Compile("jrpc2hh:method\\n")
	if err != nil {
		log.Fatal(err)
	}
	sTmpl, err := template.ParseFiles("./templates/service.tmpl")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Can't parse service template, %s", err.Error()))
	}
	mTmpl, err := template.ParseFiles("./templates/method.tmpl")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Can't parse method template, %s", err.Error()))
	}

	for _, f := range d {
		docs := doc.New(f, hDir, 0)
		for _, t := range docs.Types {
			if regExpService.Match([]byte(t.Doc)) {
				// Place for finding function New for service
				log.Printf("S %s", t.Name)
				sgd := ServiceGenData{Package: f.Name, Service: t.Name, Methods: make([]string, 0), Packages:make(map[string]*bool)}
				for _, m := range t.Methods {
					if regExpMethod.Match([]byte(m.Doc)) {
						log.Printf("  M %s", m.Name)
						if len(m.Decl.Type.Params.List) != 2 {
							log.Fatalf("%s.%s bad interface, mast has 2 params", t.Name, m.Name)
						}
						if len(m.Decl.Type.Results.List) != 1 {
							log.Fatalf("%s.%s bad interface, mast has 1 result", t.Name, m.Name)
						}
						res, ok := m.Decl.Type.Results.List[0].Type.(*ast.Ident)
						if !ok {
							log.Fatalf("%s.%s can't cast res to *ast.Indentm, maybe you result doesn't implement error", t.Name, m.Name)
						}
						if res.Name != "error" {
							log.Fatalf("%s.%s result isn't error", t.Name, m.Name)
						}

						mgd := MethodGenData{Method:m.Name}

						if reflect.TypeOf(m.Decl.Type.Params.List[0].Type).String() == "*ast.SelectorExpr" {
							pack := m.Decl.Type.Params.List[0].Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
							model := m.Decl.Type.Params.List[0].Type.(*ast.SelectorExpr).Sel.Name
							if pack + model == "jModelsNilArgs" {
								mgd.ArgsBlock = `if reqBody.HasParams() {
			return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "That method of service can't has param", nil)
		}
		var args jModels.NilArgs`
							} else {
								mgd.ArgsBlock = fmt.Sprintf(`var args %s.%s
		if reqBody.HasParams() {
			err := json.Unmarshal(*reqBody.Params, &args)
			if err != nil {
				return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "Can't unmarshal params to args structure'", err.Error())
			}
		}`, pack, model)
								sgd.Packages["encoding/json"] = nil
								sgd.Packages["models"] = nil
							}
						} else {
							model := m.Decl.Type.Params.List[0].Type.(*ast.Ident).Name
							mgd.ArgsBlock = fmt.Sprintf(`var args %s
		if reqBody.HasParams() {
			err := json.Unmarshal(*reqBody.Params, &args)
			if err != nil {
				return nil, jModels.NewError(jModels.ErrorCodeInvalidParams, "Can't unmarshal params to args structure'", err.Error())
			}
		}`, model)
							sgd.Packages["encoding/json"] = nil
							sgd.Packages["models"] = nil
						}





						if reflect.TypeOf(m.Decl.Type.Params.List[1].Type).String() == "*ast.StarExpr" {
							var pack, model string
							if reflect.TypeOf(m.Decl.Type.Params.List[1].Type.(*ast.StarExpr).X).String() == "*ast.StarExpr" {
								pack = m.Decl.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name
								model = m.Decl.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name
								sgd.Packages["models"] = nil
							} else {
								pack = m.Decl.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name
								model = m.Decl.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name
								sgd.Packages["models"] = nil
							}
							mgd.ResultBlock = fmt.Sprintf("var res %s.%s", pack, model)
						} else {
							log.Fatal("Bad result, result must be star")
						}

						buf := bytes.NewBuffer(make([]byte, 0))
						mTmpl.Execute(buf, mgd)
						sgd.Methods = append(sgd.Methods, buf.String())
					}
				}

				file, err := os.OpenFile(fmt.Sprintf("%s/jrpc2hh_%s.go", hDir, strings.ToLower(t.Name)),
					os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
					0755)
				if err != nil {
					log.Fatalf(fmt.Sprintf("Can't open file for writing generated data, %s", err.Error()))
				}
				sTmpl.Execute(file, sgd)
			}
		}
	}
}
