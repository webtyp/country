package country

import "webtyp.com/fmt"

// Normalize writes a phone number in international form, "+" and digits only
// ("+56911111111"). A number that already starts with "+" or "00" keeps its own
// prefix; any other is read as national and gets the default country's prefix.
// The result must be 7 to 15 digits (E.164) and not start with 0.
func Normalize(phone, defaultCode string) (string, error) {
	var digits []byte
	international := false
	for i := 0; i < len(phone); i++ {
		switch ch := phone[i]; {
		case ch == '+' && len(digits) == 0:
			international = true
		case ch >= '0' && ch <= '9':
			digits = append(digits, ch)
		case ch == ' ' || ch == '-' || ch == '(' || ch == ')' || ch == '.':
		default:
			return "", fmt.Err("phone", "has a character that is not a digit")
		}
	}
	if len(digits) >= 2 && digits[0] == '0' && digits[1] == '0' && !international {
		digits, international = digits[2:], true
	}
	if !international {
		c, ok := Find(defaultCode)
		if !ok {
			return "", fmt.Err("phone", "needs a country prefix and no default country is set")
		}
		digits = append([]byte(c.Dial), digits...)
	}
	if len(digits) < 7 || len(digits) > 15 || digits[0] == '0' {
		return "", fmt.Err("phone", "must have 7 to 15 digits with the country prefix")
	}
	return "+" + string(digits), nil
}

// Split separates an international number ("+584121234567") into its country
// and the national digits. It reads the longest dial prefix of the countries a
// prefix belongs to first (+1 is the United States, +7 Russia). ok is false when
// the text is not an international number or no prefix matches.
func Split(international string) (c Country, national string, ok bool) {
	if len(international) < 2 || international[0] != '+' {
		return Country{}, "", false
	}
	digits := international[1:]
	for size := 4; size >= 1; size-- {
		if len(digits) <= size {
			continue
		}
		for _, code := range primary {
			f, found := Find(code)
			if found && len(f.Dial) == size && digits[:size] == f.Dial {
				return f, digits[size:], true
			}
		}
	}
	return Country{}, "", false
}
