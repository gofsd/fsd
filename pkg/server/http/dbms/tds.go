package dbms

// DB root structure
type DB struct {
	TestCollection Example
	files          File
}

// Example is a tested collection structure
type Example struct {
	ID      uint
	UserID  uint
	GroupID uint
	Name    string
	Text    []byte
	User    File
}
