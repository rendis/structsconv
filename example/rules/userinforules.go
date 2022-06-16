package rules

import (
	"fmt"
	"github.com/rendis/structsconv"
	"github.com/rendis/structsconv/example/domain"
	"github.com/rendis/structsconv/example/dto"
)

// GetUserInfoDtoToUserInfoDomainRules dto.UserInfoDto -> domain.UserInfo
func GetUserInfoDtoToUserInfoDomainRules() *structsconv.RulesDefinition {
	var rules = structsconv.RulesSet{}

	// Name association
	rules["Name"] = "FirstName"

	// Request current source
	rules["FullName"] = func(i dto.UserInfoDto) string {
		return i.FirstName + " " + i.LastName
	}

	// Request current source and root parent source
	rules["Description"] = func(i dto.UserInfoDto, u dto.UserDto) domain.DescriptionDomain {
		return domain.DescriptionDomain{
			ID:       u.UserID,
			FullName: i.FirstName + " " + i.LastName,
			Age:      i.Age,
			Email:    u.Email,
		}
	}

	// Request current source and first argument of type "string" (hello)
	rules["ExternalValue"] = func(i dto.UserInfoDto, s string) string {
		return fmt.Sprintf("ExternalValue composed from '%s' and '%s'", i.FirstName, s)
	}

	return &structsconv.RulesDefinition{
		Rules:  rules,
		Source: dto.UserInfoDto{},
		Target: domain.UserInfo{},
	}
}
