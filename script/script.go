package main

//
// ./script --pkg-path [pkg] --version [ver] --component [comp] --controller [opdir]]
import (
	"crd/snapgen"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"strings"
)

const crdfile = "/crd.go"

func ExtractCommentTags(marker string, lines []string) map[string][]string {
	out := map[string][]string{}
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		if !strings.HasPrefix(line, marker) {
			continue
		}
		// TODO: we could support multiple values per key if we split on spaces
		kv := strings.SplitN(line[len(marker):], "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = append(out[kv[0]], kv[1])
		} else if len(kv) == 1 {
			out[kv[0]] = append(out[kv[0]], "")
		}
	}
	return out
}

type parsedFile struct {
	name string
	file *ast.File
}

type specStatus struct {
	spec   types.Type
	status types.Type
}

func main() {
	pkgParsed := []parsedFile{}

	endLineCommentGroup := make(map[int]*ast.CommentGroup)

	ctxt := snapgen.SnapGenContext{
		ConfigCrdMap: make(map[types.Type]snapgen.CrdDetail),
		StateCrdMap:  make(map[types.Type]snapgen.CrdDetail),
	}
	op := flag.String("output-package", "", "Package Path for generated files")
	pp := flag.String("pkg-path", "", "Package path to generate the file")
	comp := flag.String("group", "", "CRD component")
	ver := flag.String("version", "", "CRD version")
	contdir := flag.String("controller", "", "Output directory where controller files are generated")

	flag.Parse()
	ctxt.OutputPkg = *op
	ctxt.PkgPath = *pp
	ctxt.Component = *comp
	ctxt.Version = *ver
	ctxt.ControllerDir = *contdir

	if ctxt.PkgPath[len(ctxt.PkgPath)-1] != '/' {
		ctxt.PkgPath += "/"
	}

	if ctxt.OutputPkg[len(ctxt.PkgPath)-1] != '/' {
		ctxt.OutputPkg += "/"
	}

	dir, e := os.Getwd()
	if e != nil {
		fmt.Println("Error in getting working directory")
	}
	index := strings.LastIndex(dir, "/")

	fmt.Println(dir)
	ctxt.Dirprefix = dir[:index]
	pkgindex := strings.LastIndex(ctxt.PkgPath, "pkg")
	ctxt.PkgRoot = ctxt.PkgPath[pkgindex:]

	crdpath := ctxt.Dirprefix + "/" + ctxt.PkgRoot + ctxt.Component + "/" + ctxt.Version + crdfile

	src, err := ioutil.ReadFile(crdpath)

	if err != nil {
		fmt.Println("Couldn't open the type file", crdpath)
		return
	}
	fset := token.NewFileSet()

	// parse each file will set it's fset.
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		fmt.Println("Error in opening file")
		return
	}
	// printer.Fprint(os.Stdout, fset, f)

	pkgParsed = append(pkgParsed, parsedFile{crdpath, f})

	files := make([]*ast.File, len(pkgParsed))

	for i := range pkgParsed {
		files[i] = pkgParsed[i].file
	}
	for _, commentG := range f.Comments {
		pos := fset.Position(commentG.End())
		endLineCommentGroup[pos.Line] = commentG
	}

	c := types.Config{
		IgnoreFuncBodies: true,
		// Note that importAdapter can call b.importPackage which calls this
		// method. So there can't be cycles in the import graph.
		Importer: importer.Default(),
		Error: func(err error) {
			fmt.Println("error in parsing scope", err)
		},
	}
	//	newfs := token.NewFileSet()
	pkgDir := ctxt.PkgPath + ctxt.Component + "/" + ctxt.Version
	pkg, err := c.Check(pkgDir, fset, files, nil)
	if err != nil {
		fmt.Println("error in processing package", pkgDir)

	}
	s := pkg.Scope()
	for _, n := range s.Names() {
		obj := s.Lookup(n)
		tn, ok := obj.(*types.TypeName)
		if ok {
			position := fset.Position(obj.Pos())
			if c, ok := endLineCommentGroup[position.Line-3]; ok {
				l := strings.Split(strings.TrimRight(c.Text(), "\n"), "\n")
				values := ExtractCommentTags("+", l)
				if val, ok := values["gencrd"]; ok {
					if val[0] == "config" {
						ctxt.ConfigCrdMap[tn.Type()] = snapgen.CrdDetail{
							Name:       tn.Name(),
							NameLower:  strings.ToLower(tn.Name()),
							NamePlural: strings.ToLower(tn.Name() + "s"),
						}
					}

				}
			} else {
				// Look for state/action
				if c, ok := endLineCommentGroup[position.Line-2]; ok {
					l := strings.Split(strings.TrimRight(c.Text(), "\n"), "\n")
					values := ExtractCommentTags("+", l)
					if val, ok := values["gencrd"]; ok {
						if val[0] == "state" {
							ctxt.StateCrdMap[tn.Type()] = snapgen.CrdDetail{
								Name:       tn.Name(),
								NameLower:  strings.ToLower(tn.Name()),
								NamePlural: strings.ToLower(tn.Name() + "s"),
							}
						}
					}

				}

			}
		}
	}
	// Remove debug
	fmt.Println("crds we care about")
	for k, v := range ctxt.ConfigCrdMap {
		fmt.Println(k, v)
	}
	for k, v := range ctxt.StateCrdMap {
		fmt.Println(k, v)
	}
	snapgen.GenerateController(&ctxt)
}
