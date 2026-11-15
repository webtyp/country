package country

import "webtyp.com/fmt"

// nationalLength is how many digits a number has after the dial prefix, for the
// countries whose rule is certain. It is NOT the whole world: a country that is
// not here is checked only by the international limits (7 to 15 digits in all).
// Add a country here only with a source for its numbering plan.
var nationalLength = []struct {
	code     string
	min, max int
}{
	{"AR", 10, 11}, // 11 when the mobile "9" follows the prefix
	{"BO", 8, 8},
	{"BR", 10, 11},
	{"CA", 10, 10},
	{"CL", 9, 9},
	{"CO", 10, 10},
	{"CR", 8, 8},
	{"EC", 9, 9},
	{"ES", 9, 9},
	{"FR", 9, 9},
	{"GT", 8, 8},
	{"HN", 8, 8},
	{"HT", 8, 8},
	{"MX", 10, 10},
	{"NI", 8, 8},
	{"PA", 7, 8},
	{"PE", 9, 9},
	{"PT", 9, 9},
	{"PY", 9, 9},
	{"SV", 8, 8},
	{"US", 10, 10},
	{"UY", 8, 8},
	{"VE", 10, 10},
}

// NationalLength returns the digits a number of the country has after its dial
// prefix, when the country's rule is known.
func NationalLength(code string) (min, max int, known bool) {
	code = fmt.Convert(code).ToUpper().String()
	for _, n := range nationalLength {
		if n.code == code {
			return n.min, n.max, true
		}
	}
	return 0, 0, false
}

// Validate checks an international number: "+", 7 to 15 digits, not starting
// with 0 and, when the country is known, the digits that country's numbers have.
// It is the single phone rule (the phone input and the module both call it).
func Validate(number string) error {
	norm, err := Normalize(number, "")
	if err != nil || norm != number {
		return fmt.Err("phone", "must be an international number, like +56911111111")
	}
	c, national, ok := Split(number)
	if !ok {
		return nil
	}
	if min, max, known := NationalLength(c.Code); known && (len(national) < min || len(national) > max) {
		return fmt.Err("phone", "has the wrong number of digits for its country")
	}
	return nil
}
