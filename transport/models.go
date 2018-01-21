package transport

// User represents a User record
type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
}
