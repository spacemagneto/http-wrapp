package client

// Selector is the strategy interface for choosing a proxy from a candidate list.
//
// Implementations receive only the healthy entries pre-filtered by the pool,
// so they do not need to re-check health themselves.
// A nil or empty slice must return nil.
type Selector interface {
	Select(args []*Entry) *Entry
}
