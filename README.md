# ☎️ phone
[![Build Status](https://github.com/oarkflow/phone/workflows/CI/badge.svg)](https://github.com/oarkflow/phone/actions?query=workflow%3ACI) 
[![Go Reference](https://pkg.go.dev/badge/github.com/oarkflow/phone.svg)](https://pkg.go.dev/github.com/oarkflow/phone)

Production-oriented Go phone parsing, validation, formatting, number matching,
carrier/geocoding/time-zone lookup, and MCC/MNC enrichment. Its core tracks the
Java reference implementation of Google's [libphonenumber](https://github.com/google/libphonenumber),
through [nyaruka/phonenumbers](https://github.com/nyaruka/phonenumbers).

> [!IMPORTANT]
> The aim of this project is strictly to be a port and match as closely as possible the functionality in libphonenumber. Please don't submit feature requests for functionality that doesn't exist in libphonenumber.

> [!IMPORTANT]
> We use the metadata from libphonenumber so if you encounter unexpected parsing results, please first verify if the problem affects libphonenumber and report there if so. You can use their [online demo](https://libphonenumber.appspot.com) to quickly check parsing results.

## Version Numbers

As we don't want to bump our major semantic version number in step with the upstream library, we use independent version numbers than the Google libphonenumber repo. The release notes will mention what version of the metadata a release was built against.

## Usage

```go
import phone "github.com/oarkflow/phone"

// parse our phone number
num, err := phone.Parse("6502530000", "US")

// format it using national format
formattedNum := phone.Format(num, phone.NATIONAL)

// or use the enriched convenience API
verified := phone.Verify("+1 650 253 0000")
```

## Updating Metadata

The `buildmetadata` command clones an upstream libphonenumber release and rebuilds the
embedded metadata into the gzipped files under `data/`:

 * `data/metadata.xml.gz` - core territory metadata (number formats, validation rules, etc.)
 * `data/shortnumber_metadata.xml.gz` - short-number metadata
 * `data/alternateformats_metadata.xml.gz` - alternate format patterns used when matching
 * `data/countrycode_to_region.xml.gz` - maps a country code to its region(s)
 * `data/prefix_to_carriers/*.gz` - maps a phone number prefix to a carrier
 * `data/prefix_to_geocodings/*.gz` - maps a phone number prefix to a geographic area
 * `data/prefix_to_timezone.xml.gz` - maps a phone number prefix to a timezone

By default it resolves the **latest** upstream release tag, rebuilds `data/`, and records
the release it used in the generated `metadataversion.go` (the exported `MetadataVersion`
constant):

```bash
% go run ./cmd/buildmetadata
```

To rebuild from a specific release instead, pass the tag:

```bash
% go run ./cmd/buildmetadata v9.0.33
```

After syncing, run the tests and update [SYNC.md](SYNC.md) — which records the upstream
version the port is reconciled against and the deliberate divergences from upstream.
