package client

type Policy struct {
	httpStatusCodes map[int]struct{}

	errorCodes map[string]struct{}
}

func NewDefaultPolicy() *Policy {
	return &Policy{
		httpStatusCodes: CopySet(DefaultRetryHTTPStatusCodes),
		errorCodes:      CopySet(DefaultRetryErrorCodes),
	}
}
