package main

import (
	"io"
	"strings"
	"regexp"

	"github.com/ateliersjp/http"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	REGEXP_JS            = regexp.MustCompile("^(application|text)/(x-)?(java|ecma|j|live)script(1\\.[0-5])?$|^module$")
	REGEXP_JSON          = regexp.MustCompile("[/+]json$")
	REGEXP_XML           = regexp.MustCompile("[/+]xml$")
)

type Sed struct {
	mediaType   string
	cssMinify   minify.Minifier
	htmlMinify  minify.Minifier
	svgMinify   minify.Minifier
	jsMinify    minify.Minifier
	jsonMinify  minify.Minifier
	xmlMinify   minify.Minifier
}

func (sed *Sed) TransformHeaders(src []string) (dst []string) {
	for _, header := range src {
		if mediaType, ok := strings.CutPrefix(strings.ToLower(header), "content-type:"); ok {
			sed.mediaType = strings.TrimSpace(mediaType)
			break
		}
	}
	return src
}

func (sed *Sed) TransformBody(src io.Reader) (dst io.Reader) {
	m := minify.New()
	m.Add("text/css", sed.cssMinify)
	m.Add("text/html", sed.htmlMinify)
	m.Add("image/svg+xml", sed.svgMinify)
	m.AddRegexp(REGEXP_JS, sed.jsMinify)
	m.AddRegexp(REGEXP_JSON, sed.jsonMinify)
	m.AddRegexp(REGEXP_XML, sed.xmlMinify)
	r, w := io.Pipe()
	go func() {
		if err := m.Minify(sed.mediaType, w, src); err != nil {
			io.Copy(w, src)
		}
		w.Close()
	}()
	return r
}

func (sed *Sed) Transform(src io.Reader) (dst io.Reader) {
	return src
}

func cutRequestURI(m *http.Msg) (segment string) {
	if len(m.Headers) > 0 {
		method, path, _ := strings.Cut(m.Headers[0], " /")
		segment, path, _ = strings.Cut(path, "/")
		m.Headers[0] = method + " /" + path
	}
	return
}

func getSed(m *http.Msg) *Sed {
	if _, cmd, _ := strings.Cut(cutRequestURI(m), ":"); cmd == "keep" {
		return &Sed{
			cssMinify: &css.Minifier{
				KeepCSS2: true,
			},
			htmlMinify: &html.Minifier{
				KeepComments: true,
				KeepDefaultAttrVals: true,
				KeepDocumentTags: true,
				KeepEndTags: true,
				KeepQuotes: true,
				KeepWhitespace: true,
			},
			svgMinify: &svg.Minifier{
				KeepComments: true,
			},
			jsMinify: &js.Minifier{
				KeepVarNames: true,
			},
			jsonMinify: &json.Minifier{
				KeepNumbers: true,
			},
			xmlMinify: &xml.Minifier{
				KeepWhitespace: true,
			},
		}
	}
	return &Sed{
		cssMinify: &css.Minifier{},
		htmlMinify: &html.Minifier{},
		svgMinify: &svg.Minifier{},
		jsMinify: &js.Minifier{},
		jsonMinify: &json.Minifier{},
		xmlMinify: &xml.Minifier{},
	}
}
