package twitch

import "time"

type timestamp string

func (ts timestamp) Time() time.Time {
	t, err := time.Parse("", ts)
	if err != nil {
		return time.Time{}
	}
	return t
}
