package dto

type AddressDto struct {
	Street     string
	City       string
	Number     string
	ZipCode    string
	Primary    bool
	References []string
}
