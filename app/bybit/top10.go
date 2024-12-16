package bybit

import (
	"encoding/json"

	"github.com/marcioaso/consult/utils"
)

type BybitTop10Data struct {
	AlgoType  string `json:"algoType"`
	BaseCoin  string `json:"baseCoin"`
	BuType    string `json:"buType"`
	QuoteCoin string `json:"quoteCoin"`
	SymbolId  string `json:"symbolId"`
}

func GetTop10() ([]BybitTop10Data, error) {
	var requestData = &struct {
		Result struct {
			Data []BybitTop10Data `json:"data"`
		} `json:"result"`
	}{}

	result, err := utils.Request(
		getUrl("/de/cht/api/searchSymbolRecommend/list"),
		DefaultHeaders,
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &requestData)
	if err != nil {
		return nil, err
	}

	var response = make([]BybitTop10Data, 0)

	response = append(response, requestData.Result.Data...)
	return response, nil
}
