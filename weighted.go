package client

import "math/rand/v2"

type WeightedRandom struct {
	randFloat64 func() float64
}

func NewWeightedRandom() *WeightedRandom {
	return &WeightedRandom{randFloat64: rand.Float64}
}

func (w *WeightedRandom) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	var sum float64
	for _, entry := range entries {
		sum += entry.stats.Weight()
	}

	r := w.randFloat64() * sum
	for _, entry := range entries {
		r -= entry.stats.Weight()
		if r <= 0 {
			return entry
		}
	}

	return entries[len(entries)-1]
}
