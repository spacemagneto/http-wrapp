package client

type Rotation interface {
	Next(args []*Entry) *Entry
}
