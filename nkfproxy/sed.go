package main

import (
	"io"
	"bufio"
	"strings"

	"github.com/ateliersjp/http"
	"go4.org/strutil"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

const (
	BUF_SIZE             = 4096
	CMD                  = "charset="
)

type Sed struct {
	src     encoding.Encoding
	dst     encoding.Encoding
}

func (s *Sed) Transform(r io.Reader) io.Reader {
	if t := s.Transformer(); t != nil {
		return transform.NewReader(r, t)
	}
	return r
}

func (s *Sed) Transformer() (t transform.Transformer) {
	if s.src != nil {
		if s.dst == nil {
			t = s.src.NewDecoder()
		}
	} else {
		if s.dst != nil {
			t = s.dst.NewEncoder()
		}
	}
	return
}

func (s *Sed) DetectFrom(m *http.Msg) {
	if enc := detectFromHeader(m); enc != nil {
		s.src = enc
	} else if enc := detectFromBody(m); enc != nil {
		s.src = enc
	}
}

func (s *Sed) Invert() *Sed {
	return &Sed{
		src: s.dst,
		dst: s.src,
	}
}

func cutRequestURI(m *http.Msg) {
	if len(m.Headers) > 0 {
		method, path, _ := strings.Cut(m.Headers[0], " /")
		_, path, _ = strings.Cut(path, "/")
		m.Headers[0] = method + " /" + path
	}
}

func detectFromHeader(m *http.Msg) (enc encoding.Encoding) {
	for _, line := range m.Headers {
		if enc = detect(line); enc != nil {
			break
		}
	}
	return
}

func detectFromBody(m *http.Msg) (enc encoding.Encoding) {
	if r, ok := m.Body.(*bufio.Reader); ok {
		if data, err := r.Peek(BUF_SIZE); err == nil || err == io.EOF {
			enc = detect(string(data))
		}
	}
	return
}

func detect(data string) (enc encoding.Encoding) {
	if _, charset, ok := strings.Cut(data, CMD); ok {
		charset = strings.TrimPrefix(charset, `"`)
		if strutil.HasPrefixFold(charset, "Shift_JIS") {
			enc = japanese.ShiftJIS
		} else if strutil.HasPrefixFold(charset, "EUC-JP") {
			enc = japanese.EUCJP
		} else if strutil.HasPrefixFold(charset, "ISO-2022-JP") {
			enc = japanese.ISO2022JP
		}
	}
	return
}

func getSed(m *http.Msg) (s *Sed) {
	s = &Sed{}
	if enc := detectFromHeader(m); enc != nil {
		s.dst = enc
	}
	cutRequestURI(m)
	return
}
