# country
<img src="docs/img/badges.svg">

Standard list of countries: ISO code, dial prefix and translated name, with a flag derived from the code.

- **The person chooses, never types.** `Catalog{Default: "CL", First: []string{"VE", "PE"}}` is what an app decides (`App.Countries`): the default first, then `First`, then the rest by translated name. `Options()` is a select by name (`🇻🇪 Venezuela`), `PrefixOptions()` a select by prefix (`🇻🇪 +58`).
- **The flag costs nothing:** the code written as two regional-indicator characters; no file, no request. Windows draws them as the two letters.
- **Names are `lang.json` texts** (CLDR names), loaded only for the languages the app supports. The table (code, prefix, English name) is about 10 KB in the binary.
- **Phones:** `Normalize(phone, defaultCode)` returns `+<digits>` (7 to 15, E.164) and `Split(number)` returns the country and the national digits.

The table is generated from CLDR region names and ITU calling codes; shared prefixes (+1, +7, +44…) read as the main country (`primary`).
