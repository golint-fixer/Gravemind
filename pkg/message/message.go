package message

import (
	"sort"
	"strconv"
	"strings"

	"github.com/mvdan/xurls"
)

type Message struct {
	Room         string `json:"room"`
	RawContent   string `json:"body"`
	Content      []MessagePart
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

func (m *Message) ParseTags() {
	m.Username = strings.Title(m.Login)
	m.Color = "#33CC33"

	for _, tag := range strings.Split(m.Tags, ";") {
		var key, val string
		i := strings.Index(tag, "=")
		if i > 0 {
			key, val = tag[:i], tag[i+1:]
		} else {
			key = tag
		}
		switch key {
		case "display-name":
			m.Username = val
		case "user-id":
			m.UserId, _ = strconv.Atoi(val)
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
					id, _ := strconv.Atoi(parts[0])
					for _, point := range strings.Split(parts[1], ",") {
						p := strings.Split(point, "-")
						start, _ := strconv.Atoi(p[0])
						end, _ := strconv.Atoi(p[1])
						m.emotes = append(m.emotes, emote{id, start, end + 1})
					}
				}
			}
		}
	}
}

func (m *Message) ParseMessage() {
	sort.Sort(m.emotes)
	urls := xurls.Relaxed.FindAllStringIndex(m.RawContent, -1)

	emoticons := append(m.emotes, emote{0, len(m.RawContent), 9999})
	urls = append(urls, []int{len(m.RawContent), 9999})

	bodyIndex, emoteIndex, urlIndex := 0, 0, 0
	for bodyIndex < len(m.RawContent) {
		emoticon := emoticons[emoteIndex]
		url := emote{-1, urls[urlIndex][0], urls[urlIndex][1]}

		// Prefer emotes if there's overlap
		if emoticon.Start < url.End {
			if emoticon.Start > bodyIndex {
				m.Content = append(m.Content, MessagePartString{m.RawContent[bodyIndex:emoticon.Start]})
			}
			if emoticon.End <= len(m.RawContent) {
				m.Content = append(m.Content, MessagePartEmote{m.RawContent[emoticon.Start:emoticon.End], emoticon.Id})
			}
			bodyIndex = emoticon.End
		} else {
			if url.Start > bodyIndex {
				m.Content = append(m.Content, MessagePartString{m.RawContent[bodyIndex:url.Start]})
			}
			if url.End <= len(m.RawContent) {
				m.Content = append(m.Content, MessagePartUrl{m.RawContent[url.Start:url.End]})
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
