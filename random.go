package client

import "math/rand/v2"

type Random struct{}

func (r *Random) Next(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	return entries[rand.IntN(len(entries))]
}
