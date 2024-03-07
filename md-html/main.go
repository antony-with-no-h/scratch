package main

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	sample = `# Heading
Paragraph text

` + "```" + `erlang
code block
` + "```"
)

/*
   def block_code(self, code, lang):
       return textwrap.dedent('''\
           <ac:structured-macro ac:name="code" ac:schema-version="1">
               <ac:parameter ac:name="language">{l}</ac:parameter>
               <ac:plain-text-body><![CDATA[{c}]]></ac:plain-text-body>
           </ac:structured-macro>
       ''').format(c=code, l=lang or '')
*/

func confluenceCodeBlock(w io.Writer, c *ast.CodeBlock, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf(`<ac:structured-macro ac:name="code" ac:schema-version="1">
<ac:parameter ac:name="language">%s</ac:parameter>
<ac:plain-text-body<![CDATA[%s]]></ac:plain-text-body></ac:structured-macro>`, theme(c), c.Literal))
	}
}

func Hook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if codeBlock, ok := node.(*ast.CodeBlock); ok {
		confluenceCodeBlock(w, codeBlock, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func main() {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(sample))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: Hook}
	renderer := html.NewRenderer(opts)

	//out := markdown.Render(doc, renderer)
	fmt.Printf("%s\n", markdown.Render(doc, renderer))

	//fmt.Printf("\n", out)

	//fmt.Print("--- AST tree:\n")
	//ast.Print(os.Stdout, doc)
	//fmt.Print("\n")
}

func theme(c *ast.CodeBlock) string {
	// aliases from github.com/github-linguist/linguist/blob/master/lib/linguist/languages.yml
	syntaxMap := map[string]string{
		"actionscript 3":  "actionscript3",
		"actionscript3":   "actionscript3",
		"as3":             "actionscript3",
		"applescript":     "applescript",
		"scpt":            "applescript",
		"sh":              "bash",
		"shell-script":    "bash",
		"bash":            "bash",
		"zsh":             "bash",
		"csharp":          "c#",
		"cake":            "c#",
		"cakescript":      "c#",
		"cpp":             "cpp",
		"css":             "css",
		"cfm":             "coldfusion",
		"cfml":            "coldfusion",
		"coldfusion html": "coldfusion",
		"delphi":          "delphi",
		"objectpascal":    "delphi",
		"diff":            "diff",
		"erlang":          "erl",
		"erl":             "erl",
		"groovy":          "groovy",
		"html":            "xml",
		"xhtml":           "xml",
		"xml":             "xml",
		"rss":             "xml",
		"xsd":             "xml",
		"wsdl":            "xml",
		"java":            "java",
		"js":              "js",
		"node":            "js",
		"php":             "php",
		"perl":            "perl",
		"none":            "text",
		"fundamental":     "text",
		"plain text":      "text",
		"powershell":      "powershell",
		"posh":            "powershell",
		"pwsh":            "powershell",
		"python":          "py",
		"jruby":           "ruby",
		"macruby":         "ruby",
		"rake":            "ruby",
		"rb":              "ruby",
		"rbx":             "ruby",
		"sql":             "sql",
		"sass":            "sass",
		"scala":           "scala",
		"visual basic":    "vb",
		"vbnet":           "vb",
		"vb .net":         "vb",
		"vb.net":          "vb",
		"yaml":            "yml",
		"yml":             "yml",
		"json":            "json",
		"geojson":         "json",
		"jsonl":           "json",
		"topojson":        "json",
	}

	syntax := fmt.Sprintf("%s", c.Info)
	if lang, ok := syntaxMap[syntax]; ok {
		return lang
	}

	return "text"
}
