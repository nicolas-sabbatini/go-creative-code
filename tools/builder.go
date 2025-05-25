package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
)

//go:embed template/template.html
var templateFile string

type Flags struct {
	target  string
	os      string
	arch    string
	help    bool
	options bool
}

func parseFlags() *Flags {
	help := flag.Bool("h", false, "Print help")
	options := flag.Bool("o", false, "Print target options")
	target := flag.String("target", "", "Name of the build target - required (all will build everything)")
	os := flag.String("os", "linux", "Target plataform: js - linux - darwin - windows (def: linux)")
	arch := flag.String("arch", "amd64", "Target architecture: wasm - amd64 - arm (def: amd64)")
	flag.Parse()
	return &Flags{
		help:    *help,
		options: *options,
		target:  *target,
		os:      *os,
		arch:    *arch,
	}
}

func printHelp() {
	fmt.Println(`	> Flags:
		-h: Print help
		-o: Print target options
		-target: Name of the build target - required (all will build everything)
		-os: Target plataform: js - linux - darwin - windows (def: linux)
		-arch: Target architecture: wasm - amd64 - arm (def: amd64)

	> Example:
		"go run tools/builder.go -target=all -plataform=js -arch=wasm"
		Build all targets to web
    `)
}

func help(flags *Flags) {
	fmt.Println("== Project builder ==")
	if flags.help {
		printHelp()
		os.Exit(0)
	}
}

func readProjectDir() []string {
	fmt.Println("> Reading project")
	entries, err := os.ReadDir("./cmd")
	if err != nil {
		log.Fatal(err)
	}
	options := make([]string, 0)
	for _, e := range entries {
		options = append(options, e.Name())
	}
	return options
}

func options(flags *Flags, op []string) {
	if flags.options {
		fmt.Println("> Options")
		fmt.Println("  - all")
		for _, o := range op {
			fmt.Printf("  - %s\n", o)
		}
		os.Exit(0)
	}
}

func selectBuild(flags *Flags, op []string) []string {
	if flags.target == "" {
		log.Error("The `-target` flag is required")
		printHelp()
		os.Exit(1)
	}
	if flags.target == "all" {
		return op
	}
	exist := false
	for _, o := range op {
		if o == flags.target {
			exist = true
		}
	}
	if !exist {
		log.Errorf("The `-target=%s` is not do not exist", flags.target)
		log.Error("Options:")
		log.Error("  - all")
		for _, o := range op {
			log.Errorf("  - %s", o)
		}
		printHelp()
		os.Exit(1)
	}
	return []string{flags.target}
}

func tidy() {
	fmt.Println("> Tidy the modules")
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Can't tidy the modules, error: %s", err)
	}
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	log.Fatal(err)
	return false
}

func buildDesktop(flags *Flags, targetName string, outputFolder string) {
	outputName := targetName
	if flags.os == "windows" {
		outputName = outputName + ".exe"
	}
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		fmt.Sprintf("%s/%s", outputFolder, outputName),
		fmt.Sprintf("./cmd/%s/%s.go", targetName, targetName))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", flags.os))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", flags.arch))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Can't build %s-%s, error: %s", flags.os, flags.arch, err)
	}
}

func camelCaseToTitleCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	s = re.ReplaceAllString(s, `${1} ${2}`)
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}

func buildJs(flags *Flags, targetName string, outputFolder string) {
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		fmt.Sprintf("%s/%s.wasm", outputFolder, targetName),
		fmt.Sprintf("./cmd/%s/%s.go", targetName, targetName),
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", flags.os))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", flags.arch))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Can't build %s-%s, error: %s", flags.os, flags.arch, err)
	}
	htmlTemplate, err := template.New("template").Parse(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	indexFile, err := os.Create(fmt.Sprintf("%s/index.html", outputFolder))
	defer indexFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = htmlTemplate.Execute(indexFile, struct {
		Title   string
		BinName string
	}{
		Title:   camelCaseToTitleCase(targetName),
		BinName: fmt.Sprintf("%s.wasm", targetName),
	})
	if err != nil {
		log.Fatal(err)
	}
	goRootCmd := exec.Command("go", "env", "GOROOT")
	goRootOutput, err := goRootCmd.Output()
	if err != nil {
		log.Fatalf("Failed to get GOROOT: %s", err)
	}
	cmd = exec.Command(
		"cp",
		fmt.Sprintf("%s/lib/wasm/wasm_exec.js", strings.TrimSpace(string(goRootOutput))),
		outputFolder,
	)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Can't copy the go `wasm_exec.js` file, error: %s", err)
	}
}

func build(flags *Flags, targets []string) {
	baseFolder := fmt.Sprintf("./target/%s-%s", flags.os, flags.arch)
	exist := fileExist("./target")
	fmt.Println("> Crating target directory")
	if !exist {
		os.Mkdir("./target", 0775)
	}
	exist = fileExist(baseFolder)
	if !exist {
		os.Mkdir(baseFolder, 0775)
	}
	for _, t := range targets {
		fmt.Printf("> Building (%s-%s) %s\n", flags.os, flags.arch, t)
		targetFolder := fmt.Sprintf("%s/%s", baseFolder, t)
		err := os.RemoveAll(targetFolder)
		if err != nil {
			log.Fatalf("Can't clean the previews build %s", err)
		}
		os.Mkdir(targetFolder, 0775)
		if flags.os == "linux" || flags.os == "windows" || flags.os == "darwin" {
			buildDesktop(flags, t, targetFolder)
		} else if flags.os == "js" {
			buildJs(flags, t, targetFolder)
		} else {
			log.Fatalf("Unknown platform %s-%s", flags.os, flags.arch)
		}
	}
	fmt.Printf("> Finish compiling, all binares are in `%s`\n", baseFolder)
}

func main() {
	flags := parseFlags()
	help(flags)
	op := readProjectDir()
	options(flags, op)
	targets := selectBuild(flags, op)
	tidy()
	build(flags, targets)
}
