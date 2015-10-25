package message

import "fmt"

type MessagePartType int

const (
	STRING MessagePartType = iota
	EMOTE
	URL
)

type MessagePart interface {
	Type() MessagePartType
	String() string
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
func (mp MessagePartEmote) Url() string {
	return fmt.Sprintf("http://static-cdn.jtvnw.net/emoticons/v1/%d/3.0", mp.emote)
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

type emote struct {
	Id    int
	Start int
	End   int
}

type emotes []emote

func (e emotes) Len() int           { return len(e) }
func (e emotes) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e emotes) Less(i, j int) bool { return e[i].Start < e[j].Start }
