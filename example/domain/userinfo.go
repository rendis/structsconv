package domain

type UserInfo struct {
	Name          string
	LastName      string
	FullName      string
	Age           int
	ExternalValue string
	Addresses     []Address
	Description   DescriptionDomain
}
