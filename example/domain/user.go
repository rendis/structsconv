package domain

type UserDomain struct {
	ID             int64
	NickName       string
	Email          string
	Password       string
	Info           UserInfo
	Nicks          []string
	FavoriteBooks  map[string]BookDomain
	Top5Movies     [5]MovieDomain
	IgnorableField string
}
