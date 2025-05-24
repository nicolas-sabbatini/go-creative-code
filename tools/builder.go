package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type Flags struct {
	target   string
	platform string
	arch     string
	help     bool
	options  bool
}

func printHelp() {
	fmt.Println(`== Project builder ==
	> Flags:
		-h: Print help
		-target: Name of the builder target - required (all whil build everything)
		-plataform: Target plataform: js - linux - darwin - windows (def: linux)
		-arch: Target architecture: wasm - amd64 - arm (def: amd64)

	> Example:
		"go run tools/builder.go -target=all -plataform=js -arch=wasm"
		Build all targets to web
    `)
}

func parseFlags() *Flags {
	help := flag.Bool("h", false, "Print help")
	options := flag.Bool("o", false, "Print options")
	target := flag.String("target", "", "Name of the target")
	platform := flag.String("platform", "linux", "Target plataform")
	arch := flag.String("arch", "amd64", "Target architecture")
	flag.Parse()
	return &Flags{
		help:     *help,
		options:  *options,
		target:   *target,
		platform: *platform,
		arch:     *arch,
	}
}

func readProjectDir() []string {
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

func main() {
	flags := parseFlags()
	if flags.help {
		printHelp()
		return
	}
	options := readProjectDir()

	if flags.target == "" {
		log.Error("The `-target` flag is required")
		printHelp()
		return
	}
	exist := false
	for _, o := range options {
		if o == flags.target {
			exist = true
		}
	}
	if !exist && flags.target != "all" {
		log.Errorf("The `-target=%s` is not do not exist", flags.target)
		log.Error("Options:")
		log.Error("  - all")
		for _, o := range options {
			log.Errorf("  - %s", o)
		}
		printHelp()
		return
	}
}
