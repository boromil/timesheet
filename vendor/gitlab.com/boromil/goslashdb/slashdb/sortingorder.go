package slashdb

// Asc sets the ascending format for the given field
// i.e. 'a' -> 'a', '-b' -> 'b'
func Asc(field string) string {
	if len(field) > 0 && field[0] == '-' {
		return field[1:]
	}
	return field
}

// Desc sets the descending format for the given field
// i.e. '-a' -> '-a', 'b' -> '-b'
func Desc(field string) string {
	if len(field) > 0 && field[0] != '-' {
		return "-" + field
	}
	return field
}
