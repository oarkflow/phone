package phone

import (
	"encoding/json"

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

func LoadNetworks() error {
	data, err := str.DecodeBinaryString(networkMap)
	if err != nil {
		return err
	}
	var items []Network
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, mp := range items {
		CountryNetwork[mp.CountryCode] = append(CountryNetwork[mp.CountryCode], mp)
	}
	return nil
}
