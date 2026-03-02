package client

import "math/rand/v2"

type WeightedRandom struct {
}

func (WeightedRandom) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	var sum float64
	for _, entry := range entries {
		sum += entry.stats.Weight()
	}

	r := rand.Float64() * sum
	for _, entry := range entries {
		r -= entry.stats.Weight()
		if r <= 0 {
			return entry
		}
	}

	return entries[len(entries)-1]
}
