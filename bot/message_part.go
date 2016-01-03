package main

import (
	"bytes"
	"fmt"
	"html/template"
)

type MessagePartType int

const (
	STRING MessagePartType = iota
	EMOTE
	URL
)

type MessageParts []MessagePart

func (m MessageParts) HTML() template.HTML {
	var buf bytes.Buffer
	for _, p := range m {
		buf.WriteString(string(p.HTML()))
	}
	return template.HTML(buf.String())
}

type MessagePart interface {
	Type() MessagePartType
	String() string
	HTML() template.HTML
}

type MessagePartString struct {
	content string
}

func (mp MessagePartString) Type() MessagePartType {
	return STRING
}
func (mp MessagePartString) String() string {
	return mp.content
}
func (mp MessagePartString) HTML() template.HTML {
	return template.HTML(template.HTMLEscapeString(mp.content))
}

type MessagePartEmote struct {
	content string
	emote   int
}

func (mp MessagePartEmote) Type() MessagePartType {
	return EMOTE
}
func (mp MessagePartEmote) String() string {
	return mp.content
}
func (mp MessagePartEmote) URL() string {
	return fmt.Sprintf("http://static-cdn.jtvnw.net/emoticons/v1/%d/3.0", mp.emote)
}
func (mp MessagePartEmote) HTML() template.HTML {
	return template.HTML(fmt.Sprintf(`<img src="%s">`, mp.URL()))
}

type MessagePartUrl struct {
	content string
}

func (mp MessagePartUrl) Type() MessagePartType {
	return URL
}
func (mp MessagePartUrl) String() string {
	return mp.content
}
func (mp MessagePartUrl) URL() string {
	return "?????????" // TODO
}
func (mp MessagePartUrl) HTML() template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, mp.URL(), mp.content))
}

type emote struct {
	Id    int
	Start int
	End   int
}

type emotes []emote

func (e emotes) Len() int           { return len(e) }
func (e emotes) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e emotes) Less(i, j int) bool { return e[i].Start < e[j].Start }
