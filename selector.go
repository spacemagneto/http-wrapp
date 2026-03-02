package client

type Selector interface {
	Select(args []*Entry) *Entry
}
