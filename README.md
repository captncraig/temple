##temple - intelligent template utilities for go

The `html/template` package is extremely powerful, but lacks many quality-of-life features that can make it tricky. This package adds some of these features:

- A way to statically embed templates into your binary for production use, while still allowing editing without restarting your app while developing.
- Rendering groups of templates together (Like rendering Header, Content, Footer in a single response)
- Master / Child template execution.
- Efficient In-memory template execution. Rendering to the http response seems reasonable, but errors midway can corrupt your response.

**Potential Additions**
- Type restrictions on context objects on a per-template basis. Make sure you never render a template with an unexpected type.

### Usage:
#### Embed your templates:

`go get github.com/captncraig/temple/templeGen`

`templeGen [pkg=myPackage] [var=myWebTemplates] -o=templates.go -dir=templates`

This will generate a file (templates.go) with a map of template names to base-64 encoded template data for every file in the `templates` directory. It will include all top level files regardless of extension.

#### Retrieve a template by name

```import (
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
}```

The flag approach is only one way to activate dev-mode. Other ways are certainly possible. In dev-mode, the templateManager will load all template files any time any template is requested. This will ensure you always serve the latest version of all templates. If devMode is false, it will rely on its embedded, pre-parsed versions of the templates for maximum portability and performance.
