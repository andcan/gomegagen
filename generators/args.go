package generators

import (
	"k8s.io/gengo/types"
	"strings"
)

type CustomArgs struct {
	WhitelistStructs []string
	BlacklistStructs []string

	BlacklistFields []string
}

func (CustomArgs) hasMatch(name string, t *types.Type) bool {
	if strings.Contains(name, ".") {
		ss := strings.Split(name, ".")
		pkg := strings.Join(ss[:len(ss)-1], ".")
		name := ss[len(ss)-1]
		return pkg == t.Name.Package && name == t.Name.Name
	}
	return name == t.Name.Name
}

func (c *CustomArgs) IsWhitelistedStruct(t *types.Type) bool {
	for _, name := range c.WhitelistStructs {
		if c.hasMatch(name, t) {
			return true
		}
	}
	return false
}

func (c *CustomArgs) IsBlacklistedStruct(t *types.Type) bool {
	for _, name := range c.WhitelistStructs {
		if c.hasMatch(name, t) {
			return false
		}
	}
	return true
}

func (c *CustomArgs) ShouldGenerate(t *types.Type) bool {
	if 0 != len(c.WhitelistStructs) {
		return c.IsWhitelistedStruct(t)
	} else {
		return c.IsBlacklistedStruct(t)
	}
}
