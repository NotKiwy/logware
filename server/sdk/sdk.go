package sdk

import (
	"context"
	"strings"
	"unicode/utf8"
)

type Cfg map[string]string

type Choice struct {
	V     string `json:"v"`
	Label string `json:"label"`
}

type Opt struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Def     string   `json:"def"`
	Choices []Choice `json:"choices"`
	Text    bool     `json:"text,omitempty"`
	Max     int      `json:"max,omitempty"`
	Hide    bool     `json:"hide,omitempty"`
}

func (o Opt) Has(__v string) bool {
	if o.Text {
		return utf8.RuneCountInString(__v) <= o.Max && !strings.ContainsAny(__v, "\n\r")
	}
	for _, _c := range o.Choices {
		if _c.V == __v {
			return true
		}
	}
	return false
}

const (
	Reply = "reply"
	Name  = "name"
)

type Perm struct {
	Right string `json:"right"`
	Why   string `json:"why"`
}

type Msg struct {
	Text string
	Kind string
	Out  bool
}

type Bot interface {
	Edit(context.Context, string) error
	Reply(context.Context, string) error
	Note(context.Context, string) error
}

type Plugin struct {
	ID     string
	Name   string
	Desc   string
	Long   string
	Cat    string
	Ver    string
	Check  bool
	Icon   string
	Color  string
	Author string
	Perms  []Perm
	Opts   []Opt
	On     func(context.Context, Bot, Msg, Cfg) error
}

func (p *Plugin) Can(__right string) bool {
	for _, _r := range p.Perms {
		if _r.Right == __right {
			return true
		}
	}
	return false
}
