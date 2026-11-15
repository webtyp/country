package country_test

import (
	"testing"

	"webtyp.com/country"
)

func TestFlagIsTheCodeAsRegionalIndicators(t *testing.T) {
	if got := country.Flag("VE"); got != "\U0001F1FB\U0001F1EA" {
		t.Errorf("Flag(VE) = %q", got)
	}
	if country.Flag("ve") != country.Flag("VE") {
		t.Error("the code is case-insensitive")
	}
	for _, bad := range []string{"", "V", "VEN", "1A", "é1"} {
		if country.Flag(bad) != "" {
			t.Errorf("Flag(%q) must be empty", bad)
		}
	}
}

func TestFindAndEveryCountryHasAFlagAndADialPrefix(t *testing.T) {
	c, ok := country.Find("cl")
	if !ok || c.Code != "CL" || c.Dial != "56" {
		t.Fatalf("Find(cl) = %+v, %v", c, ok)
	}
	if _, ok := country.Find("ZZ"); ok {
		t.Error("ZZ is not a country")
	}
	list := country.Catalog{Default: "CL"}.Countries()
	if len(list) < 240 {
		t.Fatalf("expected every country, got %d", len(list))
	}
	for _, f := range list {
		if country.Flag(f.Code) == "" || f.Dial == "" || f.Name == "" {
			t.Errorf("incomplete entry %+v", f)
		}
	}
}

// The default is first, then First in order, no duplicates, then the rest.
func TestCatalogOrder(t *testing.T) {
	cat := country.Catalog{Default: "CL", First: []string{"VE", "PE", "CL", "HT"}}
	if err := cat.Validate(); err != nil {
		t.Fatal(err)
	}
	list := cat.Countries()
	want := []string{"CL", "VE", "PE", "HT"}
	for i, code := range want {
		if list[i].Code != code {
			t.Errorf("position %d = %s, want %s", i, list[i].Code, code)
		}
	}
	seen := map[string]bool{}
	for _, f := range list {
		if seen[f.Code] {
			t.Errorf("%s listed twice", f.Code)
		}
		seen[f.Code] = true
	}
	if opts := cat.PrefixOptions(); opts[0].Key != "CL" || opts[0].Value != country.Flag("CL")+" +56" {
		t.Errorf("first prefix option = %+v", opts[0])
	}
}

func TestCatalogValidate(t *testing.T) {
	if (country.Catalog{}).Validate() == nil {
		t.Error("a catalog without a default country must fail")
	}
	if (country.Catalog{Default: "XX"}).Validate() == nil {
		t.Error("an unknown default must fail")
	}
	if (country.Catalog{Default: "CL", First: []string{"nope"}}).Validate() == nil {
		t.Error("an unknown First code must fail")
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct{ in, def, want string }{
		{"+56 9 1111 1111", "VE", "+56911111111"},
		{"911111111", "CL", "+56911111111"},
		{"0058 412 1234567", "CL", "+584121234567"},
		{"(412) 123-4567", "VE", "+584121234567"},
	}
	for _, c := range cases {
		got, err := country.Normalize(c.in, c.def)
		if err != nil || got != c.want {
			t.Errorf("Normalize(%q,%s) = %q, %v; want %q", c.in, c.def, got, err, c.want)
		}
	}
	for _, bad := range []string{"123", "+0123456789", "abc", "+12345678901234567"} {
		if _, err := country.Normalize(bad, "CL"); err == nil {
			t.Errorf("Normalize(%q) must fail", bad)
		}
	}
	if _, err := country.Normalize("911111111", ""); err == nil {
		t.Error("a national number with no default country must fail")
	}
}

func TestSplit(t *testing.T) {
	cases := map[string]string{"+56911111111": "CL", "+584121234567": "VE", "+12125551234": "US", "+18095551234": "DO", "+79001234567": "RU"}
	for number, code := range cases {
		c, national, ok := country.Split(number)
		if !ok || c.Code != code || "+"+c.Dial+national != number {
			t.Errorf("Split(%s) = %s %q %v, want %s", number, c.Code, national, ok, code)
		}
	}
	if _, _, ok := country.Split("911111111"); ok {
		t.Error("a number without + is not international")
	}
}

// A number has the digits its country's numbers have: Chile's are 9 after +56.
func TestValidateChecksTheCountrysLength(t *testing.T) {
	good := []string{"+56911111111", "+584121234567", "+12125551234", "+5491155551234", "+44207946095"} // GB: unknown rule, only the international limits
	for _, n := range good {
		if err := country.Validate(n); err != nil {
			t.Errorf("Validate(%s) = %v, want nil", n, err)
		}
	}
	bad := []string{"+569851092360", "+5698510923", "+5841212345", "+509371234", "911111111", "+0123456789"}
	for _, n := range bad {
		if country.Validate(n) == nil {
			t.Errorf("Validate(%s) must fail", n)
		}
	}
	if min, max, known := country.NationalLength("cl"); !known || min != 9 || max != 9 {
		t.Errorf("NationalLength(cl) = %d,%d,%v", min, max, known)
	}
	if _, _, known := country.NationalLength("GB"); known {
		t.Error("a country with no certain rule must report unknown")
	}
}
