package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Script holds data model for script template
type Script struct {
	QueryString string
	ModuleName  string
	ModuleVar   string
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}
	// The doc string
	var queryStr = os.Args[1]
	// We need tge module name
	// if
	// input = http.Agent
	// moduleName = http
	var moduleName = strings.Split(
		queryStr, ".",
	)[0]
	// create new template
	t := template.New("script")
	// script template for generating doc
	var scriptTemplate = `
try{
  eval('try{ var {{.ModuleVar}} = require("{{.ModuleName}}"); console.log({{.QueryString}});} catch(e){console.log(e.message);}');
} catch(e){
  console.log(e.message);
}
`
	// parse our script template
	t, err := t.Parse(scriptTemplate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fill data
	var data = Script{
		QueryString: queryStr,
		ModuleName:  moduleName,
		ModuleVar:   moduleName,
	}
	// execute template
	var w = new(bytes.Buffer)
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// exec the script for printing doc
	printDoc(w)
}

// printUsage prints usage texts
func printUsage() {
	fmt.Println(`
USAGE of jsdoc:
 jsdoc <module>
 jsdoc <sym>[.<method>]                  
 jsdoc [<module>].<sym>[.<method>]

EXAMPLE USAGE of jsdoc:
 jsdoc http
 jsdoc http.Agent

CREATED BY:
 @AnikHasibul`)
}

func printDoc(script io.Reader) {
	cmd := exec.Command("node", "-")
	cmd.Stdin = script
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
