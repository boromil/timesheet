package slashdb

import (
	"strconv"
	"strings"
)

// SetSort sets the sort query param on the request
func (req *Request) SetSort(fields ...string) {
	req.Params["sort"] = strings.Join(fields, ",")
}

// RemoveSort removes the sort query param from the request
func (req *Request) RemoveSort() {
	delete(req.Params, "sort")
}

// DistinctOn turns on distinct for the request
func (req *Request) DistinctOn() {
	req.Params["distinct"] = ""
}

// DistinctOff turns off distinct for the request
func (req *Request) DistinctOff() {
	delete(req.Params, "distinct")
}

// SetLimit sets the limit query param on the request
func (req *Request) SetLimit(limit int) {
	req.Params["limit"] = strconv.Itoa(limit)
}

// RemoveLimit removes the limit query param from the request
func (req *Request) RemoveLimit() {
	delete(req.Params, "limit")
}

// SetOffset sets the offset query param on the request
func (req *Request) SetOffset(offset int) {
	req.Params["offset"] = strconv.Itoa(offset)
}

// RemoveOffset removes the offset query param from the request
func (req *Request) RemoveOffset() {
	delete(req.Params, "offset")
}

// StreamingOn turns on streaming for the request
func (req *Request) StreamingOn() {
	req.Params["stream"] = "true"
}

// StreamingOff turns off streaming for the request
func (req *Request) StreamingOff() {
	delete(req.Params, "stream")
}

// SetDepth sets the depth query param on the request
func (req *Request) SetDepth(depth int) {
	req.Params["depth"] = strconv.Itoa(depth)
}

// RemoveDepth removes the depth query param from the request
func (req *Request) RemoveDepth() {
	delete(req.Params, "depth")
}

// WantArrayOn turns on wantarray for the request
func (req *Request) WantArrayOn() {
	req.Params["wantarray"] = "true"
}

// WantArrayOff turns off wantarray for the request
func (req *Request) WantArrayOff() {
	delete(req.Params, "wantarray")
}

// SetSeparator sets the separator query param on the request
func (req *Request) SetSeparator(separator string) {
	req.Separator = separator
	req.Params["separator"] = separator
	for i := range req.Parts {
		req.Parts[i].Separator = separator
	}
}

// ResetSeparator resets the separator query param to the default ','
func (req *Request) ResetSeparator() {
	req.Separator = ","
	delete(req.Params, "separator")
	for i := range req.Parts {
		req.Parts[i].Separator = req.Separator
	}
}

// SetURLStringSub sets the URL string substitution query param on the request
func (req *Request) SetURLStringSub(urlStringSub string) {
	req.Params[urlStringSub] = "/"
}

// RemoveSetURLStringSub removes the URL string substitution query param from the request
func (req *Request) RemoveSetURLStringSub() {
	keys := []string{}
	for k, v := range req.Params {
		if v == "/" {
			keys = append(keys, k)
		}
	}
	for _, k := range keys {
		delete(req.Params, k)
	}
}

// SetWildcard sets the wildcard query param on the request
func (req *Request) SetWildcard(wildcard string) {
	req.Params["wildcard"] = wildcard
}

// RemoveWildcard removes the wildcard query param from the request
func (req *Request) RemoveWildcard() {
	delete(req.Params, "wildcard")
}

// NilVisibleOn turns on nil_visible for the request
func (req *Request) NilVisibleOn() {
	req.Params["nil_visible"] = "true"
}

// NilVisibleOff turns off nil_visible for the request
func (req *Request) NilVisibleOff() {
	delete(req.Params, "nil_visible")
}

// SetCardinality sets the cardinality query param on the request
func (req *Request) SetCardinality(cardinality int) {
	req.Params["cardinality"] = strconv.Itoa(cardinality)
}

// RemoveCardinality removes the cardinality query param from the request
func (req *Request) RemoveCardinality() {
	delete(req.Params, "cardinality")
}

// HeadersOn turns on headers for the request
func (req *Request) HeadersOn() {
	req.Params["headers"] = "true"
}

// HeadersOff turns off headers for the request
func (req *Request) HeadersOff() {
	delete(req.Params, "headers")
}

// SetCSVNullStr turns on csvNullStr for the request
func (req *Request) SetCSVNullStr(csvNullStr string) {
	req.Params["csvNullStr"] = csvNullStr
}

// RemoveCSVNullStr turns off csvNullStr for the request
func (req *Request) RemoveCSVNullStr() {
	delete(req.Params, "csvNullStr")
}
