package rules

import (
	"github.com/rendis/structsconv"
	"github.com/rendis/structsconv/example/domain"
	"github.com/rendis/structsconv/example/dto"
)

// GetAddressDtoToAddressDomainRules dto.AddressDto{} -> domain.Address{}
func (d *RulesDefinitions) GetAddressDtoToAddressDomainRules() *structsconv.RulesDefinition {
	var rules = structsconv.RulesSet{}

	// Name association
	rules["Zip"] = "ZipCode"

	// Constant
	rules["Country"] = func() string { return "CL" }

	return &structsconv.RulesDefinition{
		Rules:  rules,
		Source: dto.AddressDto{},
		Target: domain.Address{},
	}
}
