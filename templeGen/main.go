package main

import (
	"flag"
	"log"

	"github.com/captncraig/temple/gen"
	"gopkg.in/fsnotify.v1"
)

var (
	flagPkg   = flag.String("pkg", "main", "package declaration for generated file")
	flagVar   = flag.String("var", "templates", "variable name for embedded template map")
	flagFile  = flag.String("o", "", "Output file name")
	flagDir   = flag.String("dir", "templates", "Directory containing templates to embed")
	flagWatch = flag.Bool("w", false, "watch directory for changes and regenerate")
)

func main() {
	flag.Parse()
	if *flagFile == "" {
		flag.PrintDefaults()
		log.Fatal("output file name required")
	}

	params := gen.GenParameters{
		*flagPkg, *flagVar, *flagFile, *flagDir,
	}

	err := gen.Generate(params)
	if err != nil {
		log.Fatal(err)
	}

	if !*flagWatch {
		return
	}

}
