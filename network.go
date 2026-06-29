package phone

import (
	"encoding/json"
	"sync"

	"github.com/oarkflow/pkg/str"
)

type Network struct {
	Plmn        string      `json:"plmn"`
	NibbledPlmn string      `json:"nibbledPlmn"`
	Mcc         string      `json:"mcc"`
	Mnc         string      `json:"mnc"`
	Region      string      `json:"region"`
	Type        string      `json:"type"`
	CountryName string      `json:"countryName"`
	CountryCode string      `json:"countryCode"`
	Lat         string      `json:"lat"`
	Long        string      `json:"long"`
	Brand       string      `json:"brand"`
	Operator    string      `json:"operator"`
	Status      string      `json:"status"`
	Bands       string      `json:"bands"`
	Notes       interface{} `json:"notes"`
}

var CountryNetwork = map[string][]Network{}

var (
	networksOnce sync.Once
	networksErr  error
)

func LoadNetworks() error {
	networksOnce.Do(func() {
		data, err := str.DecodeBinaryString(networkMap)
		if err != nil {
			networksErr = err
			return
		}

		var items []Network
		if err := json.Unmarshal(data, &items); err != nil {
			networksErr = err
			return
		}

		byCountry := make(map[string][]Network)
		for _, item := range items {
			byCountry[item.CountryCode] = append(byCountry[item.CountryCode], item)
		}
		CountryNetwork = byCountry
	})
	return networksErr
}
