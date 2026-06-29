package phone

// NormalizeDiallableCharsOnly strips characters that cannot be dialled on a
// phone keypad, including non-ASCII digits.
func NormalizeDiallableCharsOnly(number string) string {
	return normalizeDiallableCharsOnly(number)
}

// IsPossibleNumberFromRegion parses number as dialled from regionDialingFrom
// and reports whether its length is possible. It is intentionally more
// lenient than IsValidNumber.
func IsPossibleNumberFromRegion(number, regionDialingFrom string) bool {
	parsed, err := Parse(number, regionDialingFrom)
	return err == nil && IsPossibleNumber(parsed)
}

// IsNumberGeographical reports whether number has a geographical association.
func IsNumberGeographical(number *PhoneNumber) bool {
	if number == nil {
		return false
	}
	return IsNumberGeographicalForType(GetNumberType(number), int(number.GetCountryCode()))
}

// IsNumberGeographicalForType is the cheaper form of IsNumberGeographical
// when the caller already knows the number type and calling code.
func IsNumberGeographicalForType(numberType PhoneNumberType, countryCallingCode int) bool {
	return numberType == FIXED_LINE ||
		numberType == FIXED_LINE_OR_MOBILE ||
		(GEO_MOBILE_COUNTRIES[int32(countryCallingCode)] && numberType == MOBILE)
}

// CanBeInternationallyDialled reports whether number may be dialled from
// outside its region. It does not validate the number first.
func CanBeInternationallyDialled(number *PhoneNumber) bool {
	if number == nil {
		return false
	}
	return canBeInternationallyDialled(number)
}

// GetCarrierForValidNumber performs carrier lookup without repeating number
// type validation. Callers must supply a valid mobile-capable number.
func GetCarrierForValidNumber(number *PhoneNumber, lang string) (string, error) {
	carrier, _, err := GetCarrierWithPrefixForNumber(number, lang)
	return carrier, err
}

// GetTimeZonesForNumber is the upstream-compatible spelling of
// GetTimezonesForNumber.
func GetTimeZonesForNumber(number *PhoneNumber) ([]string, error) {
	return GetTimezonesForNumber(number)
}

// GetTimeZonesForGeographicalNumber returns time zones for a parsed
// geographical number.
func GetTimeZonesForGeographicalNumber(number *PhoneNumber) ([]string, error) {
	if number == nil {
		return nil, ErrNotANumber
	}
	return GetTimezonesForPrefix(Format(number, E164))
}
