package thirdpart

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
)

func FetchJSONFromURI(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 返回状态码错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %v", err)
	}

	return result, nil
}

func GetMetaInfoFromGoPlus(ca string) (GoplusApiResponse, error) {
	url := fmt.Sprintf("https://api.gopluslabs.io/api/v1/solana/token_security?contract_addresses=%s", ca)

	var result GoplusApiResponse

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "*/*")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("GOPLUS: Error making request:", err)
		return result, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("GOPLUS: Error reading response body:", err)
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error("GOPLUS: Error unmarshaling JSON:", err, string(body))
		return result, err
	}

	if result.Code != 1 {
		return result, fmt.Errorf("%v", result.Message)
	}

	return result, nil
}

func CopyMetaIntoObject(metaInfo GoplusTokenInfo) (existMeta model.TokenMeta, pairs []model.TokenPair) {
	createTime := time.Now()

	existMeta.Description = metaInfo.Metadata.Description
	existMeta.Name = metaInfo.Metadata.Name
	existMeta.Symbol = metaInfo.Metadata.Symbol

	if len(metaInfo.Metadata.URI) > 0 {
		existMeta.MetaUri = metaInfo.Metadata.URI
		uriResult, err := FetchJSONFromURI(metaInfo.Metadata.URI)
		if err == nil {
			jsonBytes, err := json.Marshal(uriResult)
			if err != nil {
				log.Error("Error marshalling JSON:", err)
				return
			} else {
				existMeta.MetaContent = string(jsonBytes)
			}
		} else {
			log.Error("[copyMetaIntoObject] Error unable to get meta info ", metaInfo.Metadata.URI)
		}
	}

	dexinfo := metaInfo.Dex
	pairs = make([]model.TokenPair, 0)
	for _, v := range dexinfo {
		pairs = append(pairs, model.TokenPair{
			Pair:       v.ID,
			Dex:        v.DexName,
			OpenTime:   v.OpenTime,
			InitTvl:    v.Tvl,
			LpAmount:   v.LpAmount,
			CreateTime: createTime,
		})
	}

	if metaInfo.Closable.Status == "0" || metaInfo.Closable.Status == "" {
		existMeta.TClosable = "U"
	} else {
		existMeta.TClosable = "C"
	}

	if metaInfo.Freezable.Status == "0" || metaInfo.Freezable.Status == "" {
		existMeta.TFreezable = "U"
	} else {
		existMeta.TFreezable = "C"
	}

	if metaInfo.MetadataMutable.Status == "0" || metaInfo.MetadataMutable.Status == "" {
		existMeta.TMetaChangable = "U"
	} else {
		existMeta.TMetaChangable = "C"
	}

	if metaInfo.Mintable.Status == "0" {
		existMeta.TMintable = "U"
	} else {
		existMeta.TMintable = "C"
	}

	if metaInfo.NonTransferable == "0" {
		existMeta.TTransferable = "C"
	} else {
		existMeta.TTransferable = "U"
	}

	existMeta.TotalSupply, _ = decimal.NewFromString(metaInfo.TotalSupply)

	if metaInfo.DefaultAccountStateUpgradable.Status == "0" {
		existMeta.TUpgradable = "U"
	} else {
		existMeta.TUpgradable = "C"
	}

	existMeta.CreateTime = createTime

	return existMeta, pairs
}

/*
func CopyMetaIntoObject(metaInfo map[string]interface{}) (existMeta model.TokenMeta) {
	if fee, ok := metaInfo["transfer_fee"].(string); ok {
		existMeta.CTransferFee, _ = decimal.NewFromString(fee)
	} else if fee, ok := metaInfo["transfer_fee"].(float64); ok {
		existMeta.CTransferFee = decimal.NewFromFloat(fee)
	} else {
		log.Error("Error: transfer_fee is not a valid decimal value")
	}

	if meta, ok := metaInfo["metadata"].(map[string]interface{}); ok {
		if desc, valid := meta["description"].(string); valid {
			existMeta.Description = desc
		} else {
			log.Error("Error: description is not a valid string")
		}

		if desc, valid := meta["name"].(string); valid {
			existMeta.Name = desc
		} else {
			log.Error("Error: name is not a valid string")
		}

		if desc, valid := meta["symbol"].(string); valid {
			existMeta.Symbol = desc
		} else {
			log.Error("Error: symbol is not a valid string")
		}

		if uri, valid := meta["uri"].(string); valid {
			uriResult, err := FetchJSONFromURI(uri)
			if err == nil {
				existMeta.MetaUri = uri

				jsonBytes, err := json.Marshal(uriResult)
				if err != nil {
					log.Error("Error marshalling JSON:", err)
					return
				} else {
					existMeta.MetaContent = string(jsonBytes)
				}

			}
		}
	} else {
		log.Error("Error unable to get meta info")
	}

	if dex, ok := metaInfo["dex"].([]interface{}); ok {
		if len(dex) > 0 {
			dexmap, _ := dex[0].(map[string]interface{})
			openTimeStr := fmt.Sprintf("%v", dexmap["open_time"])
			timestamp, err := strconv.ParseInt(openTimeStr, 10, 64)
			if err != nil {
				log.Error("Error parsing timestamp:", err)
				return
			} else {
				existMeta.OpenTime = time.Unix(timestamp, 0)
			}
		}
	} else {
		log.Error("Error unable to get dex info")
	}

	if closable, ok := metaInfo["closable"].(map[string]interface{}); ok {
		status, _ := closable["status"].(string)
		if len(status) > 0 && status == "0" {
			existMeta.TClosable = "U"
		} else {
			existMeta.TClosable = "C"
		}
	} else {
		log.Error("Error unable to get closable info")
	}

	if freezable, ok := metaInfo["freezable"].(map[string]interface{}); ok {
		status, _ := freezable["status"].(string)
		if len(status) > 0 && status == "0" {
			existMeta.TFreezable = "U"
		} else {
			existMeta.TFreezable = "C"
		}
	} else {
		log.Error("Error unable to get freezable info")
	}

	if metadata_mutable, ok := metaInfo["metadata_mutable"].(map[string]interface{}); ok {
		status, _ := metadata_mutable["status"].(string)
		if len(status) > 0 && status == "0" {
			existMeta.TMetaChangable = "U"
		} else {
			existMeta.TMetaChangable = "C"
		}
	} else {
		log.Error("Error unable to get metadata_mutable info")
	}

	if mintable, ok := metaInfo["mintable"].(map[string]interface{}); ok {
		status, _ := mintable["status"].(string)
		if len(status) > 0 && status == "0" {
			existMeta.TMintable = "U"
		} else {
			existMeta.TMintable = "C"
		}
	} else {
		log.Error("Error unable to get metadata_mutable info")
	}

	if transferable, ok := metaInfo["non_transferable"].(string); ok {
		if transferable == "0" {
			existMeta.TTransferable = "C"
		} else {
			existMeta.TTransferable = "U"
		}
	}

	if total_supply, ok := metaInfo["total_supply"].(string); ok {
		existMeta.TotalSupply, _ = decimal.NewFromString(total_supply)
	}

	if default_account_state_upgradable, ok := metaInfo["default_account_state_upgradable"].(map[string]interface{}); ok {
		status, _ := default_account_state_upgradable["status"].(string)
		if len(status) > 0 && status == "0" {
			existMeta.TUpgradable = "U"
		} else {
			existMeta.TUpgradable = "C"
		}
	} else {
		log.Error("Error unable to get metadata_mutable info")
	}

	return existMeta
}
*/
