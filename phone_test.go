package phone

import (
	"strings"
	"testing"
)

func TestNumberVerifyRejectsInvalidInputWithoutPanicking(t *testing.T) {
	inputs := []string{"", "not a phone", "+999999999999999999999", "1"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			number := Verify(input, "US")
			if !number.Invalid {
				t.Fatalf("Verify(%q) unexpectedly returned valid: %+v", input, number)
			}
		})
	}
}

func TestNumberVerifyCanBeReusedWithoutStaleEnrichment(t *testing.T) {
	number := Number{Phone: "+14155552671"}
	number.Verify(false)
	if number.Invalid || number.Phone != "+14155552671" || number.CountryCode != "US" {
		t.Fatalf("unexpected valid result: %+v", number)
	}

	number.Phone = "invalid"
	number.Verify(false)
	if !number.Invalid {
		t.Fatalf("reused number unexpectedly valid: %+v", number)
	}
	if number.CountryCode != "" || number.DialCode != 0 || number.Timezone != "" {
		t.Fatalf("reused number retained stale enrichment: %+v", number)
	}
}

func TestNumberVerifySupportsInternationalAccessPrefix(t *testing.T) {
	number := Verify("0044 20 7031 3000")
	if number.Invalid || number.Phone != "+442070313000" || number.CountryCode != "GB" {
		t.Fatalf("unexpected result: %+v", number)
	}
}

func TestFindNumbersAndAsYouTypeFormatter(t *testing.T) {
	var matches []string
	for match := range FindNumbers("Call +1 415-555-2671 or (212) 555-0198.", "US") {
		matches = append(matches, Format(match.Number(), E164))
	}
	if len(matches) != 2 || matches[0] != "+14155552671" || matches[1] != "+12125550198" {
		t.Fatalf("unexpected matches: %v", matches)
	}

	formatter := GetAsYouTypeFormatter("US")
	var formatted string
	for _, digit := range strings.ReplaceAll("4155552671", " ", "") {
		formatted = formatter.InputDigit(digit)
	}
	if formatted != "(415) 555-2671" {
		t.Fatalf("unexpected as-you-type output: %q", formatted)
	}
}
