package mycache

import (
	"fmt"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/system"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

var httpClient = http.Client{Timeout: 10 * time.Second}
var QuoteMap = make(map[string]*QuoteToken)

// Custom Quote Coin (very stablecoin) Add to CustomQuote to increase the price of altcoin: altcoin
// Passively trigger altcoins: Load when the price of the altcoins is obtained. Take the price of the largest pool in tvl and update the maximum pool address regularly? How often does it be updated?
// swap update pool price
var CustomQuote = make(map[string]*QuoteToken)
var quoteKey = "quote"

func init() {

	QuoteMap["So11111111111111111111111111111111111111112"] = &QuoteToken{
		Address:     "So11111111111111111111111111111111111111112",
		PairAddress: "Czfq3xZZDmsdGdUyrNLtRhGc47cXcZtLG4crryfu44zE",
		Symbol:      "SOL",
		ChainCode:   "SOLANA",
		Decimals:    9,
		PairSymbol:  "SOLUSDT",
		Price:       decimal.NewFromFloat(0),
		LestTime:    0,
		Mainnet:     true,
		Logo:        "https://img.apihellodex.lol/BSC/0x570a5d26f7765ecb712c0924e4de545b89fd43df.png",
	}
	QuoteMap["EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"] = &QuoteToken{
		Address:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Symbol:     "USDC",
		PairSymbol: "USDC",
		Decimals:   6,
		Price:      decimal.NewFromFloat(1),
	}
	QuoteMap["Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"] = &QuoteToken{
		Address:    "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
		Symbol:     "USDT",
		PairSymbol: "USDC",
		Decimals:   6,
		Price:      decimal.NewFromFloat(1),
	}
	res := &QuoteToken{}
	err := system.HGet(quoteKey, "SOLANA:So11111111111111111111111111111111111111112", res)
	if err != nil {
		panic(fmt.Sprintf("init quote error: %v", err))
	}
	QuoteMap["So11111111112"] = res
}
func QuotePrice(ca string) *QuoteToken {

	token, exist := QuoteMap[ca]
	if !exist {
		return nil
	}
	if ca == codes.SOL && time.Now().UnixMilli()-token.LestTime > 60*1000 {
		res := &QuoteToken{}
		err := system.HGet(quoteKey, "SOLANA:So11111111111111111111111111111111111111112", res)
		if err != nil {
			log.Errorf("SolPrice error: %v", err)
			return token
		}
		if res.Price.Sign() > 0 {
			token.SetPrice(res.Price)
			return res
		}
	}
	return token
}

type QuoteToken struct {
	Address     string          `json:"address"`
	PairAddress string          `json:"pairAddress"`
	ChainCode   string          `json:"chainCode"`
	Symbol      string          `json:"symbol"`
	Decimals    int8            `json:"decimals"`
	PairSymbol  string          `json:"pairSymbol"`
	Price       decimal.Decimal `json:"price"`
	LestTime    int64           `json:"lastTime"`
	Mainnet     bool            `json:"mainnet"`
	Logo        string          `json:"logo"`
}

// 实现 String 方法
func (q QuoteToken) String() string {
	return fmt.Sprintf(
		"Address: %s, PairAddress: %s, ChainCode: %s, Symbol: %s, Decimals: %d, PairSymbol: %s, Price: %s, LestTime: %d, Mainnet: %t",
		q.Address, q.PairAddress, q.ChainCode, q.Symbol, q.Decimals, q.PairSymbol, q.Price.String(), q.LestTime, q.Mainnet,
	)
}

type PairSymbolPrice struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

// SetPrice sets the price and updates the lestTime to the current timestamp.
func (qt *QuoteToken) SetPrice(price decimal.Decimal) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("QuoteToken SetPrice panic: %v", r)
		}
	}()
	qt.Price = price
	qt.LestTime = time.Now().UnixMilli()
	if qt.PairSymbol != "USDC" {

	}

}

func GetAll() map[string]*QuoteToken {
	return QuoteMap
}
