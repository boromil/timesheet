package slashdb

import "strings"

// Filter describes the request filtering condition (also ordering)
type Filter struct {
	Values map[string][]string
	Order  []string
}

// NewFilter the default constructor for the Filter object
func NewFilter(
	values map[string][]string,
	order []string,
) Filter {
	if values == nil {
		values = map[string][]string{}
	}
	if order != nil {
		order = []string{}
	}

	return Filter{
		Values: values,
		Order:  order,
	}
}

// Part describes a single request part/segment
type Part struct {
	Name      string
	Fields    []string
	Filter    Filter
	Separator string
}

// NewPart the default constructor for the Part object
func NewPart(
	name string,
	fields []string,
	filter Filter,
	separator string,
) Part {
	if separator == "" {
		separator = ","
	}

	return Part{
		Name:      name,
		Fields:    fields,
		Filter:    filter,
		Separator: separator,
	}
}

func prepareFilters(k string, values []string, sep string) string {
	if k == "" {
		return ""
	}
	if len(values) == 0 {
		return ""
	}
	value := strings.Join(values, sep)
	if value == "" {
		return ""
	}
	return "/" + k + "/" + value
}

func (part Part) String() string {
	name := part.Name
	if name != "" {
		name = "/" + name
	}

	sep := part.Separator
	var filters string
	if len(part.Filter.Order) != 0 {
		for _, k := range part.Filter.Order {
			filters += prepareFilters(k, part.Filter.Values[k], sep)
		}
	} else {
		for k, vs := range part.Filter.Values {
			filters += prepareFilters(k, vs, sep)
		}
	}

	var fields string
	if len(part.Fields) != 0 {
		fields = "/" + strings.Join(part.Fields, part.Separator)
	}

	return name + filters + fields
}
