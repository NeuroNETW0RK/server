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

var serviceName string

type argsInfo struct {
	argsTypes []string
	funcArgs  string
}

type resultsInfo struct {
	resultName []string
	funcResult string
}

func init() {
	service := flag.String("service", "", "请输入需要生成的 service 名称\n")
	flag.Parse()

	serviceName = strings.ToLower(*service)
}

func main1() {
	// 可以自己指定serviceName
	serviceName = "task"
	fs := token.NewFileSet()
	filePath := fmt.Sprintf("./internal/neuronetserver/service/v1/%s", serviceName)
	parsedFile, err := decorator.ParseFile(fs, filePath+fmt.Sprintf("/%s.go", serviceName), nil, 0)
	if err != nil {
		log.Fatalf("parsing package: %s: %s\n", filePath, err)
	}

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

					funcParams := v.Type.(*dst.FuncType)

					resultList := funcParams.Results.List

					argsList := funcParams.Params.List
					argsInfo := parseArgs(argsList)
					resultsInfo := parseResult(resultList)

					filepath := "./internal/neuronetserver/service/v1/" + serviceName
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

					funcContent := fmt.Sprintf("package %s\n\n", serviceName)
					funcContent += "import (\n"
					funcContent += `"github.com/gin-gonic/gin"`
					funcContent += "\n"
					funcContent += fmt.Sprintf(`"neuronet/domain/%s"`, serviceName)
					funcContent += "\n)\n\n"
					funcContent += fmt.Sprintf("func (s *service) %s(%s) %s {}", v.Names[0].String(), argsInfo.funcArgs, resultsInfo.funcResult)
					fmt.Println(funcContent)
					funcFile.WriteString(funcContent)
					funcFile.Close()
				}
			}

			//domainRequestFilePath := "./domain/" + serviceName + "/dto.go"
			//requestFile, err := os.OpenFile(domainRequestFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
			//if err != nil {
			//	fmt.Printf("create and open func file error %v\n", err.Error())
			//}
			//if requestFile == nil {
			//	fmt.Printf("func file is nil \n")
			//}
			//fmt.Printf("gen %v dto file", serviceName)
			//var funcContent string
			//funcContent = fmt.Sprintf("package %s\n\n", serviceName)
			//for _, v := range interfaceType.Methods.List {
			//	if len(v.Names) > 0 {
			//		if v.Names[0].String() == "i" {
			//			continue
			//		}
			//		funcParams := v.Type.(*dst.FuncType)
			//		resultList := funcParams.Results.List
			//		argsList := funcParams.Params.List
			//
			//		argsInfo := parseArgs(argsList)
			//		resultsInfo := parseResult(resultList)
			//
			//		for _, args := range argsInfo.argsTypes {
			//			if args != "" {
			//				funcContent += fmt.Sprintf("type %s struct {}", args)
			//			}
			//			funcContent += "\n\n"
			//		}
			//		for _, result := range resultsInfo.resultName {
			//			if result != "" {
			//				funcContent += fmt.Sprintf("\n\ntype %s struct {}\n\n", result)
			//			}
			//			funcContent += "\n\n"
			//		}
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

func parseArgs(argsList []*dst.Field) *argsInfo {
	var (
		argsInfo  = new(argsInfo)
		ctxStr    = "c *gin.Context"
		argsStr   string
		argsTypes []string
	)

	if len(argsList) == 1 {
		argsInfo.funcArgs = ctxStr
		argsInfo.argsTypes = []string{""}
		return argsInfo
	}
	argsList = argsList[1:]
	for _, i := range argsList {
		argsName := i.Names[0].Name
		argsType := i.Type.(*dst.SelectorExpr).Sel.Name
		argsStr += argsName + " " + serviceName + "." + argsType
		argsTypes = append(argsTypes, argsType)
	}
	argsInfo.funcArgs = ctxStr + ", " + argsStr
	argsInfo.argsTypes = argsTypes
	return argsInfo
}

func parseResult(resultList []*dst.Field) *resultsInfo {
	var (
		resultsInfo = new(resultsInfo)
		replyTypes  []string
		replyStr    string
	)

	if len(resultList) == 1 {
		reply := resultList[0].Type.(dst.Expr).(*dst.Ident).Name
		resultsInfo.funcResult = reply
		resultsInfo.resultName = []string{""}
		return resultsInfo
	}
	resultList = resultList[:len(resultList)-1]
	for _, i := range resultList {
		replyType := i.Type.(*dst.StarExpr).X.(*dst.SelectorExpr).Sel.Name
		replyTypes = append(replyTypes, replyType)
		replyStr += "(*" + serviceName + "." + replyType + ", " + "error)"
	}
	resultsInfo.resultName = replyTypes
	resultsInfo.funcResult = replyStr
	return resultsInfo
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
