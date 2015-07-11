package main

import (
	"flag"
	"log"
	
	"github.com/captncraig/temple"
)

var devMode = flag.Bool("dev", false, "activate dev mode for templates")

func main() {
	flag.Parse()
	templateManager, err := temple.New(*devMode, myTemplates, "templates")
	if err != nil {
		log.Fatal(err)
	}

	template, err := templateManager.GetTemplate("main.tpl")
	...
}
