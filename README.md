##temple - intelligent template utilities for go

The `html/template` package is extremely powerful, but lacks many quality-of-life features that can make it tricky. This package adds some of these features:

- A way to statically embed templates into your binary for production use, while still allowing editing without restarting your app while developing.
- Rendering groups of templates together (Like rendering Header, Content, Footer in a single response)
- Master / Child template execution.
- Efficient In-memory template execution. Rendering to the http response seems reasonable, but errors midway can corrupt your response.

**Potential Additions**
- Type restrictions on context objects on a per-template basis. Make sure you never render a template with an unexpected type.

### Usage:
