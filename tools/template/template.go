package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"os"

	"github.com/charmbracelet/log"
)

type Project struct {
	Title   string
	BinName string
}

//go:embed template.html
var templateFile string

//go:embed wasm_exec.js
var wasmExecJs []byte

func printHelp() {
	fmt.Println(`
           -- GO WASM TEMPLATE --
  A small tool to bundle your wasm proyects

    Options:
      -title: Title of the HTML file
      -bin: Name of the bin file

    Flags:
      -h: Print this
    `)
}

func parseFlags() *Project {
	help := flag.Bool("h", false, "Print help")
	bin := flag.String("bin", "", "Name of the bin file")
	title := flag.String("title", "", "Title of the HTML file")
	flag.Parse()
	if *help || *title == "" || *bin == "" {
		return nil
	}
	return &Project{
		Title:   *title,
		BinName: *bin,
	}
}

func createTemplate() *template.Template {
	htmlTemplate, err := template.New("template").Parse(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	return htmlTemplate
}

func saveToDisk(htmlTemplate *template.Template, project *Project) {
	indexFile, err := os.Create("index.html")
	defer indexFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = htmlTemplate.Execute(indexFile, project)
	if err != nil {
		log.Fatal(err)
	}
	jsExecutor, err := os.Create("wasm_exec.js")
	defer jsExecutor.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = jsExecutor.Write(wasmExecJs)
}

func main() {
	project := parseFlags()
	if project == nil {
		printHelp()
		return
	}
	htmlTemplate := createTemplate()
	saveToDisk(htmlTemplate, project)
}
