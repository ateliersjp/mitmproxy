package main

import (
	"io"
	"net/url"
	"strings"

	"github.com/ateliersjp/http"
	"github.com/benhoyt/goawk/interp"
	"github.com/benhoyt/goawk/parser"
)

const (
	CMD                  = "awk"
)

type Sed struct {
	Program *parser.Program
	Config  interp.Config
}

func NewSed(cmd, src string) *Sed {
	prog, err := parser.ParseProgram([]byte(src), nil)
	if err != nil {
		return nil
	}
	sed := &Sed{}
	sed.Program = prog
	if cmd == "csv" {
		sed.Config.InputMode = interp.CSVMode
	} else if cmd == "tsv" {
		sed.Config.InputMode = interp.TSVMode
	}
	return sed
}

func (sed *Sed) TransformHeaders(src []string) (dst []string) {
	return src
}

func (sed *Sed) TransformBody(src io.Reader) (dst io.Reader) {
	sed.Config.Stdin = src
	dst, sed.Config.Output = io.Pipe()
	go func() {
		interp.ExecProgram(sed.Program, &sed.Config)
		sed.Config.Output.(io.WriteCloser).Close()
	}()
	return
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
	if cmd, src, ok := strings.Cut(cutRequestURI(m), "="); ok {
		_, cmd, _ := strings.Cut(cmd, ":")
		src, _ := url.QueryUnescape(src)
		return NewSed(cmd, src)
	}
	return nil
}
