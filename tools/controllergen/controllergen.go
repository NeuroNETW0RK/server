package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

var handlerName string

func init() {
	handler := flag.String("controller", "", "请输入需要生成的 controller 名称\n")
	flag.Parse()

	handlerName = strings.ToLower(*handler)
}

func main() {
	// 可以自己指定handlerName
	handlerName = "task_resource"
	fs := token.NewFileSet()
	filePath := fmt.Sprintf("./internal/neuronetserver/controller/v1/%s", handlerName)
	parsedFile, err := decorator.ParseFile(fs, filePath+fmt.Sprintf("/%s.go", handlerName), nil, 0)
	if err != nil {
		log.Fatalf("parsing package: %s: %s\n", filePath, err)
	}

	//files, _ := ioutil.ReadDir(filePath)
	//if len(files) > 1 {
	//	log.Fatalf("请先确保 %s 目录中，有且仅有 handler.go 一个文件。", filePath)
	//}

	dst.Inspect(parsedFile, func(n dst.Node) bool {
		decl, ok := n.(*dst.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			typeSpec, _ok := spec.(*dst.TypeSpec)
			if !_ok {
				continue
			}

			var interfaceType *dst.InterfaceType
			if interfaceType, ok = typeSpec.Type.(*dst.InterfaceType); !ok {
				continue
			}

			for _, v := range interfaceType.Methods.List {
				if len(v.Names) > 0 {

					filepath := "./internal/neuronetserver/controller/v1/" + handlerName
					filename := fmt.Sprintf("%s/func_%s.go", filepath, strings.ToLower(v.Names[0].String()))
					if pathExists(filename) {
						continue
					}

					funcFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
					if err != nil {
						fmt.Printf("create and open func file error %v\n", err.Error())
						continue
					}

					if funcFile == nil {
						fmt.Printf("func file is nil \n")
						continue
					}

					fmt.Println("  └── file : ", filename)

					funcContent := fmt.Sprintf("package %s\n\n", handlerName)
					funcContent += fmt.Sprintf("func (c *controller) %s(ctx *gin.Context) {}", v.Names[0].String())

					funcFile.WriteString(funcContent)
					funcFile.Close()
				}

			}

			//domainRequestFilePath := "./domain/" + handlerName + "/request.go"
			//requestFile, err := os.OpenFile(domainRequestFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
			//if err != nil {
			//	fmt.Printf("create and open func file error %v\n", err.Error())
			//}
			//if requestFile == nil {
			//	fmt.Printf("func file is nil \n")
			//}
			//fmt.Printf("gen %v request file", handlerName)
			//var funcContent string
			//funcContent = fmt.Sprintf("package %s\n\n", handlerName)
			//for _, v := range interfaceType.Methods.List {
			//	if len(v.Names) > 0 {
			//		if v.Names[0].String() == "i" {
			//			continue
			//		}
			//		funcContent += fmt.Sprintf("\n\ntype %sRequest struct {}\n\n", LcFirst(v.Names[0].String()))
			//	}
			//}
			//requestFile.WriteString(funcContent)
			//requestFile.Close()
		}
		return true
	})
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
