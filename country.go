// Package country is the standard list of countries: ISO 3166 code, telephone
// dial prefix and a translated name, plus the flag, which is not data: it is the
// code written as two regional-indicator characters (VE → 🇻🇪), so it costs no
// file and no request. Windows draws those two characters as letters.
//
// A person never types a country or a prefix: they choose it. An app says which
// countries come first and which is the default (Catalog); the names are texts
// of lang.json, loaded only for the languages the app supports.
package country

import (
	"webtyp.com/fmt"
	"webtyp.com/lang"
)

// Country is one entry of the list.
type Country struct {
	Code string    // ISO 3166-1 alpha-2, upper case: "VE"
	Dial string    // calling code without "+": "58"; a few NANP countries carry their area code ("1809")
	Name lang.Text // English name, translated when shown
}

// Find returns the country with the given code, case-insensitive.
func Find(code string) (Country, bool) {
	code = fmt.Convert(code).ToUpper().String()
	for _, c := range all {
		if c.Code == code {
			return c, true
		}
	}
	return Country{}, false
}

// Flag is the flag emoji of a country code, or "" when the code is not two
// letters A–Z.
func Flag(code string) string {
	code = fmt.Convert(code).ToUpper().String()
	if len(code) != 2 {
		return ""
	}
	for i := 0; i < 2; i++ {
		if code[i] < 'A' || code[i] > 'Z' {
			return ""
		}
	}
	const base = 0x1F1E6 // REGIONAL INDICATOR SYMBOL LETTER A
	return string(rune(base+int(code[0]-'A'))) + string(rune(base+int(code[1]-'A')))
}

// Label is how a country is chosen by name: flag and translated name.
func (c Country) Label() string {
	return Flag(c.Code) + " " + lang.Translate(c.Name).String()
}

// DialLabel is how a country is chosen by prefix: flag and "+58".
func (c Country) DialLabel() string {
	return Flag(c.Code) + " +" + c.Dial
}

// Catalog is what an app decides: the default country and which countries come
// first in every list. The rest of the world follows, by translated name.
type Catalog struct {
	Default string   // ISO code of the default country; required
	First   []string // ISO codes shown right after the default, in this order
}

// Validate fails when the default country is missing or unknown, or a First
// code is unknown: a wrong code must stop the app at start, not hide a country.
func (c Catalog) Validate() error {
	if c.Default == "" {
		return fmt.Err("country: Catalog.Default", "is required")
	}
	if _, ok := Find(c.Default); !ok {
		return fmt.Err("country: Catalog.Default", c.Default, "is not a country code")
	}
	for _, code := range c.First {
		if _, ok := Find(code); !ok {
			return fmt.Err("country: Catalog.First", code, "is not a country code")
		}
	}
	return nil
}

// Countries is the whole list in the order a person sees it: the default, the
// First ones, then every other country by translated name.
func (c Catalog) Countries() []Country {
	out := make([]Country, 0, len(all))
	taken := func(code string) bool {
		for _, o := range out {
			if o.Code == code {
				return true
			}
		}
		return false
	}
	for _, code := range append([]string{c.Default}, c.First...) {
		if f, ok := Find(code); ok && !taken(f.Code) {
			out = append(out, f)
		}
	}
	rest := make([]Country, 0, len(all))
	for _, f := range all {
		if !taken(f.Code) {
			rest = append(rest, f)
		}
	}
	// insertion sort by translated name: the list is ~250 and sorted once per screen
	for i := 1; i < len(rest); i++ {
		for j := i; j > 0 && sortKey(rest[j]) < sortKey(rest[j-1]); j-- {
			rest[j], rest[j-1] = rest[j-1], rest[j]
		}
	}
	return append(out, rest...)
}

// Options are the choices of a select of countries by name: key = ISO code.
func (c Catalog) Options() []fmt.KeyValue {
	list := c.Countries()
	opts := make([]fmt.KeyValue, 0, len(list))
	for _, f := range list {
		opts = append(opts, fmt.KeyValue{Key: f.Code, Value: f.Label()})
	}
	return opts
}

// PrefixOptions are the choices of a select of dial prefixes: key = ISO code,
// text = flag and "+58". The first option is the default country.
func (c Catalog) PrefixOptions() []fmt.KeyValue {
	list := c.Countries()
	opts := make([]fmt.KeyValue, 0, len(list))
	for _, f := range list {
		opts = append(opts, fmt.KeyValue{Key: f.Code, Value: f.DialLabel()})
	}
	return opts
}

// sortKey is the translated name in lower case without the accents of the
// languages the framework ships, so "Álava" does not sort after "Zambia".
func sortKey(c Country) string {
	name := fmt.Convert(lang.Translate(c.Name).String()).ToLower().String()
	out := make([]rune, 0, len(name))
	for _, r := range name {
		out = append(out, fold(r))
	}
	return string(out)
}

func fold(r rune) rune {
	switch r {
	case 'á', 'à', 'â', 'ä', 'ã', 'å':
		return 'a'
	case 'é', 'è', 'ê', 'ë':
		return 'e'
	case 'í', 'ì', 'î', 'ï':
		return 'i'
	case 'ó', 'ò', 'ô', 'ö', 'õ', 'ø':
		return 'o'
	case 'ú', 'ù', 'û', 'ü':
		return 'u'
	case 'ñ':
		return 'n'
	case 'ç':
		return 'c'
	}
	return r
}
