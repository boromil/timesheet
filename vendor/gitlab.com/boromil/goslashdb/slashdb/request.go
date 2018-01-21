package slashdb

import (
	"bytes"
)

// Request main SlashDB request object
type Request struct {
	Kind      string
	Parts     []Part
	Params    map[string]string
	Separator string
}

// NewDataRequest the data resource request constructor
func NewDataRequest(separator string) *Request {
	if separator == "" {
		separator = ","
	}

	return &Request{
		Kind:      "/db",
		Parts:     []Part{},
		Params:    map[string]string{},
		Separator: separator,
	}
}

// NewQueryRequest the query resource request constructor
func NewQueryRequest(separator string) *Request {
	if separator == "" {
		separator = ","
	}

	return &Request{
		Kind:      "/query",
		Parts:     []Part{},
		Params:    map[string]string{},
		Separator: separator,
	}
}

func (req *Request) String() string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(req.Kind)

	for _, part := range req.Parts {
		buf.WriteString(part.String())
	}

	buf.WriteString(".json")

	if len(req.Params) > 0 {
		buf.WriteString("?")
		for k, v := range req.Params {
			buf.WriteString(k + "=" + v + "&")
		}
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String()
}

// AddParts adds the user defined parts to the request
func (req *Request) AddParts(part ...Part) {
	for i := range part {
		part[i].Separator = req.Separator
	}
	req.Parts = append(req.Parts, part...)
}
