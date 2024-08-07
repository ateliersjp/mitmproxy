package main

import (
	"net/url"
	"strings"

	"github.com/ateliersjp/http"
	"golang.org/x/text/transform"
	texttransform "github.com/tenntenn/text/transform"
)

const (
	CMD                  = "s"
)

type Sed struct {
	rawtt   map[string]string
	tt      texttransform.ReplaceByteTable
}

func NewSed() *Sed {
	return &Sed{
		rawtt   : make(map[string]string),
		tt      : make(texttransform.ReplaceByteTable, 0, 1),
	}
}

func (s *Sed) Transformer() transform.Transformer {
	return texttransform.ReplaceAll(s.tt)
}

func (s *Sed) AppendRule(before, after string) {
	s.rawtt[before] = after
	s.tt = append(s.tt, []byte(before), []byte(after))
	if uri, ok := strings.CutPrefix(before, "http://"); ok {
		s.tt = append(s.tt, []byte("https://" + uri), []byte(after))
	} else if uri, ok := strings.CutPrefix(before, "https://"); ok {
		s.tt = append(s.tt, []byte("http://" + uri), []byte(after))
	}
}

func (s *Sed) Invert() (dst *Sed) {
	dst = NewSed()
	for after, before := range s.rawtt {
		dst.AppendRule(before, after)
	}
	return
}

func cutRequestURI(m *http.Msg) (segment string) {
	if len(m.Headers) > 0 {
		method, path, _ := strings.Cut(m.Headers[0], " /")
		segment, path, _ = strings.Cut(path, "/")
		m.Headers[0] = method + " /" + path
	}
	return
}

func cutRequestURIFunc(m *http.Msg, f func(string) bool) (segment string, ok bool) {
	if len(m.Headers) > 0 {
		method, path, _ := strings.Cut(m.Headers[0], " /")
		segment, path, _ = strings.Cut(path, "/")
		if ok = f(segment); ok {
			m.Headers[0] = method + " /" + path
		}
	}
	return
}

func isCmd(segment string) bool {
	return segment == CMD
}

func getSed(m *http.Msg) (s *Sed) {
	s = NewSed()
	for {
		if _, ok := cutRequestURIFunc(m, isCmd); !ok {
			break
		}
		before, _ := url.QueryUnescape(cutRequestURI(m))
		after, _ := url.QueryUnescape(cutRequestURI(m))
		s.AppendRule(before, after)
	}
	return
}
