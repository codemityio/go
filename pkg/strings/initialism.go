package strings

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

type Initialism struct {
	list map[string]struct{}
}

// ToPascal convert value to PascalCase.
func (i *Initialism) ToPascal(value string) string {
	value = strings.ReplaceAll(value, "_", "-")

	parts := strings.Split(value, "-")

	for idx, part := range parts {
		if _, ok := i.list[part]; ok {
			parts[idx] = strings.ToUpper(part)
		} else if len(part) > 0 {
			parts[idx] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(parts, "")
}

// ToCamel convert value to camelCase.
func (i *Initialism) ToCamel(value string) string {
	value = strings.ReplaceAll(i.ToKebab(value), "_", "-")
	parts := strings.Split(value, "-")

	for idx, part := range parts {
		if _, ok := i.list[part]; ok {
			parts[idx] = strings.ToUpper(part)
		} else if len(part) > 0 {
			parts[idx] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	if len(parts) == 0 {
		return ""
	}

	parts[0] = strings.ToLower(parts[0])

	return strings.Join(parts, "")
}

// ToKebab convert value to kebab-case.
func (i *Initialism) ToKebab(value string) string {
	var init []string

	rne := []rune(value)

	for key := range i.list {
		if strings.Contains(value, strings.ToUpper(key)) {
			index := strings.Index(value, strings.ToUpper(key))

			pos := index + len(key)

			if (unicode.IsUpper(rne[pos : pos+1][0]) && unicode.IsLower(rne[pos+1 : pos+2][0])) ||
				(!unicode.IsUpper(rne[pos : pos+1][0]) && !unicode.IsLower(rne[pos+1 : pos+2][0])) {
				init = append(init, key)

				continue
			}
		}
	}

	for _, v := range init {
		chars := []rune(v)

		value = strings.Replace(
			value,
			strings.ToUpper(v),
			strings.ToUpper(string(chars[0:1]))+string(chars[1:]),
			1,
		)
	}

	return strcase.ToKebab(value)
}

// ToSnake convert value to snake_case.
func (i *Initialism) ToSnake(value string) string {
	return strings.ToLower(strings.ReplaceAll(i.ToKebab(value), "-", "_"))
}
