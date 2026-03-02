package client

type Weighted struct {
}

func (Weighted) Next(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	return nil
}
