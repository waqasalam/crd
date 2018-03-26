package main

import (
	"fmt"
	"go/types"
	"os"
	"strings"
	"text/template"
)

type CrdDetail struct {
	Name       string
	NameLower  string
	NamePlural string
	Spec       string
	Status     string
}

type SnapGenContext struct {
	PkgPath       string
	PkgRoot       string
	Component     string
	Version       string
	OutputPkg     string
	ControllerDir string
	Dirprefix     string
	ConfigCrdMap  map[types.Type]CrdDetail
	StateCrdMap   map[types.Type]CrdDetail
}

const (
	crdfile        = "/crdtmp.go"
	controllerfile = "/controllertmp.go"
	registerfile   = "/registertmp.go"
	utilfile       = "/util.go"
	snapgen        = "/script"
)

var fmap = template.FuncMap{
	"toLower": strings.ToLower,
	"title":   strings.Title,
}

const utilText = `package {{ .Version }}

import (
	"reflect"
)

func compare(v1, v2 reflect.Value) bool {
	change := false

	switch v1.Kind() {
	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			change = true
		}
	case reflect.Int:
		if v1.Int() != v2.Int() {
			change = true
		}
	case reflect.Int8:
		if int8(v1.Int()) != int8(v2.Int()) {
			change = true
		}
	case reflect.Int16:
		if int16(v1.Int()) != int16(v2.Int()) {
			change = true
		}
	case reflect.Int32:
		if int32(v1.Int()) != int32(v2.Int()) {
			change = true
		}
	case reflect.Uint:
		if v1.Uint() != v2.Uint() {
			change = true
		}
	case reflect.Uint8:
		if uint8(v1.Uint()) != uint8(v2.Uint()) {
			change = true
		}
	case reflect.Uint16:
		if uint16(v1.Uint()) != uint16(v2.Uint()) {
			change = true
		}
	case reflect.Uint32:
		if uint32(v1.Uint()) != uint32(v2.Uint()) {
			change = true
		}
	case reflect.Uint64:
		if uint64(v1.Uint()) != uint64(v2.Uint()) {
			change = true
		}
	case reflect.Float64:
		if v1.Float() != v2.Float() {
			change = true
		}
	case reflect.Slice:
		for i, n := 0, v1.Len(); i < n; i++ {
			if compare(v1.Index(i), v2.Index(i)) {
				change = true
			}
		}
	case reflect.String:
		if v1.String() != v2.String() {
			change = true
		}
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if compare(v1.Field(i), v2.Field(i)) {
				change = true
			}
		}

	}
	return change
}

func CompareObjAndDiff(x, y interface{}) []bool {

	v1 := reflect.ValueOf(x).Elem()
	v2 := reflect.ValueOf(y).Elem()

	if v1.Kind() != reflect.Struct {
		return nil
	}
	attrset := make([]bool, v1.NumField())

	for i, n := 0, v1.NumField(); i < n; i++ {
		if compare(v1.Field(i), v2.Field(i)) {
			attrset[i] = true
		}
	}

	return attrset
}
`

func GenerateController(input *SnapGenContext) {
	crdpath := input.Dirprefix + "/" + input.PkgRoot + input.Component + "/" + input.Version + crdfile
	crdHndl, err := os.OpenFile(crdpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("Open file failed", crdpath)
	}

	regpath := input.Dirprefix + "/" + input.PkgRoot + input.Component + "/" + input.Version + registerfile
	regHndl, err := os.OpenFile(regpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("Open file failed", crdpath)
	}

	utilPath := input.Dirprefix + "/" + input.PkgRoot + input.Component + "/" + input.Version + utilfile
	utilHndl, err := os.OpenFile(utilPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("Open file failed", utilHndl)
	}

	controllerDir := input.Dirprefix + "/" + input.ControllerDir
	if _, err := os.Stat(controllerDir); os.IsNotExist(err) {
		os.Mkdir(controllerDir, 0755)

	}
	controllerpath := controllerDir + controllerfile

	controller, err := os.OpenFile(controllerpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("Open file failed", controllerpath)
	}

	// Move to snapgen directory to execute the templates
	err = os.Chdir(input.Dirprefix + snapgen)
	if err != nil {
		fmt.Println("Couldn't open the snapgen directory", snapgen)
		panic(err)
	}

	tmpl, err := template.New("").Funcs(fmap).ParseGlob("*.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.ExecuteTemplate(crdHndl, "crd.tmpl", *input)
	if err != nil {
		fmt.Println("error in processing template", err)
	}

	err = tmpl.ExecuteTemplate(regHndl, "register.tmpl", *input)
	if err != nil {
		fmt.Println("error in processing template", err)
	}

	err = tmpl.ExecuteTemplate(controller, "controller.tmpl", *input)
	if err != nil {
		fmt.Println("error in processing template", err)
	}

	utilTmpl, err := template.New("util").Parse(utilText)
	if err != nil {
		fmt.Println("Error in processing template", err)
	}

	err = utilTmpl.Execute(utilHndl, *input)
	if err != nil {
		fmt.Println("Error in Executing template", err)
	}
	/*	err = tmpl.ExecuteTemplate(statusfile, "status.tmpl", *input)
		if err != nil {
			fmt.Println("error in processing template", err)
		}
	*/
}
