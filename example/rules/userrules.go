package rules

import (
	"github.com/rendis/structsconv"
	"github.com/rendis/structsconv/example/domain"
	"github.com/rendis/structsconv/example/dto"
)

// GetUserDtoToUserDomainRules dto.UserDto{} -> domain.UserDomain{}
func (d *RulesDefinitions) GetUserDtoToUserDomainRules() *structsconv.RulesDefinition {
	var rules = structsconv.RulesSet{}

	// Name association
	rules["ID"] = "UserID"
	rules["NickName"] = "UserName"
	rules["Nicks"] = "Aliases"
	rules["Info"] = "UserInfo"

	// Set field as ignored
	rules["IgnorableField"] = nil

	return &structsconv.RulesDefinition{
		Rules:  rules,
		Source: dto.UserDto{},
		Target: domain.UserDomain{},
	}
}
