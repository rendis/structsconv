package domain

type Address struct {
	Street     string
	Number     string
	City       string
	Zip        string
	Primary    bool
	References []string
	Country    string
}
