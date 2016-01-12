package main

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/mvdan/xurls"
)

type Message struct {
	Room         string `json:"room"`
	RawContent   string `json:"body"`
	Content      MessageParts
	Login        string `json:"nick"`
	Username     string
	UserId       int
	UserType     string
	IsTurbo      bool
	IsSubscriber bool
	Color        string
	IsAction     bool

	Tags   string `json:"tags"`
	emotes emotes
}

func (m *Message) String() string {
	var buf bytes.Buffer
	buf.WriteString("<")
	buf.WriteString(m.Username)
	buf.WriteString(m.Room)
	buf.WriteString("> ")
	if m.IsAction {
		buf.WriteString("/me ")
	}
	buf.WriteString(m.Content.String())
	return buf.String()
}

var tagEscaping = map[byte]byte{
	':':  ';',
	's':  ' ',
	'\\': '\\',
	'r':  '\r',
	'n':  '\n',
}

func (m *Message) ParseTags() {
	m.Username = strings.Title(m.Login)
	m.Color = "#33CC33"

	for _, tag := range strings.Split(m.Tags, ";") {
		// Get the data
		var key, val string
		i := strings.Index(tag, "=")
		if i > 0 {
			key, val = tag[:i], tag[i+1:]
		} else {
			key = tag
		}

		// Unescape the data
		j, v := 0, []byte(val)
		for i := 0; i < len(v); i++ {
			if v[i] == '\\' {
				i++
				c, ok := tagEscaping[v[i]]
				if !ok {
					break
				}
				v[j] = c
			} else {
				v[j] = v[i]
			}
			j++
		}
		val = string(v[0:j])

		// Utilize the data
		var err error
		switch key {
		case "display-name":
			if val != "" {
				m.Username = val
			}
		case "user-id":
			m.UserId, err = strconv.Atoi(val)
			logErr(err)
		case "user-type":
			m.UserType = val
		case "turbo":
			m.IsTurbo = val == "1"
		case "subscriber":
			m.IsSubscriber = val == "1"
		case "color":
			m.Color = val
		case "emotes":
			if val != "" {
				for _, emote_data := range strings.Split(val, "/") {
					parts := strings.Split(emote_data, ":")
					id, err := strconv.Atoi(parts[0])
					logErr(err)
					for _, point := range strings.Split(parts[1], ",") {
						p := strings.Split(point, "-")
						start, err := strconv.Atoi(p[0])
						logErr(err)
						end, err := strconv.Atoi(p[1])
						logErr(err)
						m.emotes = append(m.emotes, emote{id, start, end + 1})
					}
				}
			}
		}
	}
}

func (m *Message) ParseMessage() {
	content := m.RawContent
	if len(content) > 8 && content[0:8] == "\x01ACTION " && content[len(content)-1] == '\x01' {
		content = content[8 : len(content)-1]
		m.IsAction = true
	}

	sort.Sort(m.emotes)
	urls := xurls.Relaxed.FindAllStringIndex(content, -1)

	emoticons := append(m.emotes, emote{0, len(content), 9999})
	urls = append(urls, []int{len(content), 9999})

	bodyIndex, emoteIndex, urlIndex := 0, 0, 0
	for bodyIndex < len(content) {
		emoticon := emoticons[emoteIndex]
		url := emote{-1, urls[urlIndex][0], urls[urlIndex][1]}

		// Prefer emotes if there's overlap
		if emoticon.Start < url.End {
			if emoticon.Start > bodyIndex {
				m.Content = append(m.Content, MessagePartString{content[bodyIndex:emoticon.Start]})
			}
			if emoticon.End <= len(content) {
				m.Content = append(m.Content, MessagePartEmote{content[emoticon.Start:emoticon.End], emoticon.Id})
			}
			bodyIndex = emoticon.End
		} else {
			if url.Start > bodyIndex {
				m.Content = append(m.Content, MessagePartString{content[bodyIndex:url.Start]})
			}
			if url.End <= len(content) {
				m.Content = append(m.Content, MessagePartUrl{content[url.Start:url.End]})
			}
			bodyIndex = url.End
		}

		// Cycle through the lists until they're valid again
		for emoteIndex < len(emoticons) && emoticons[emoteIndex].Start < bodyIndex {
			emoteIndex++
		}
		for urlIndex < len(urls) && urls[urlIndex][0] < bodyIndex {
			urlIndex++
		}
	}
}
