package main

import (
	"encoding/json"
	"fmt"

	"github.com/rendis/structsconv"
	"github.com/rendis/structsconv/example/domain"
	"github.com/rendis/structsconv/example/dto"
	"github.com/rendis/structsconv/example/rules"
)

func main() {
	structsconv.RegisterRulesDefinitions(rules.GetUserInfoDtoToUserInfoDomainRules())
	structsconv.RegisterSetOfRulesDefinitions(&rules.RulesDefinitions{})

	var uDto = &dto.UserDto{
		UserID:   15369764,
		UserName: "John",
		Email:    "john.test@doe.org",
		Password: "de_*wwe?-QW.",
		Aliases:  []string{"john", "johnny", "big john"},
		UserInfo: dto.UserInfoDto{
			FirstName: "John",
			LastName:  "Doe",
			Age:       25,
			Addresses: []dto.AddressDto{
				{
					Street:     "Main Street",
					City:       "New York",
					Number:     "123",
					ZipCode:    "10001",
					Primary:    true,
					References: []string{"Near the shoe store", "Near the bank"},
				},
				{
					Street:     "Second Street",
					City:       "New York",
					Number:     "456",
					ZipCode:    "10002",
					References: []string{"Near the supermarket", "Near the food store"},
				},
			},
		},

		FavoriteBooks: map[string]dto.BookDto{
			"Sci-Fi": {
				Author: "Cixin Liu",
				Genre:  "Sci-Fi",
				Title:  "El problema de los 3 cuerpos",
			},
			"Fantasy": {
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings",
				Genre:  "Fantasy",
			},
		},
		Top5Movies: [5]dto.MovieDto{
			{
				Title: "The Lord of the Rings",
				Genre: "Fantasy",
				Year:  1954,
			},
			{
				Title: "The Hobbit",
				Genre: "Fantasy",
				Year:  1937,
			},
			{
				Title: "The Matrix",
				Genre: "Sci-Fi",
				Year:  1999,
			},
			{
				Title: "The Dark Knight",
				Genre: "Sci-Fi",
				Year:  2008,
			},
			{
				Title: "The Lord of the Rings",
				Genre: "Fantasy",
				Year:  1954,
			},
		},
	}
	var uDomain = &domain.UserDomain{}

	structsconv.Map(uDto, uDomain, "hello", 3.14, "word")

	// Format the output as JSON
	vJSON, _ := json.MarshalIndent(uDomain, "", "  ")
	fmt.Printf("%s\n", string(vJSON))
}
