package phone

import (
	"encoding/json"
	"github.com/oarkflow/pkg/str"
)

type Network struct {
	Type        string `json:"type"`
	Mcc         string `json:"mcc"`
	Mnc         string `json:"mnc"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Brand       string `json:"brand"`
	Operator    string `json:"operator"`
	Status      string `json:"status"`
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
