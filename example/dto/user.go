package dto

type UserDto struct {
	UserID        int64
	UserName      string
	Email         string
	Password      string
	UserInfo      UserInfoDto
	Aliases       []string
	FavoriteBooks map[string]BookDto
	Top5Movies    [5]MovieDto
}
