package jup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SwapRequest struct {
	Route         Route  `json:"route"`
	WrapUnwrapSOL bool   `json:"wrapUnwrapSOL,omitempty"`
	FeeAccount    string `json:"feeAccount,omitempty"`
	TokenLedger   string `json:"tokenLedger,omitempty"`
	UserPublicKey string `json:"userPublicKey"`
}

type SwapResponse struct {
	SetupTransaction   string `json:"setupTransaction,omitempty"`
	SwapTransaction    string `json:"swapTransaction"`
	CleanupTransaction string `json:"cleanupTransaction,omitempty"`
}

type Quote struct {
	Routes    []Route `json:"data"`
	TimeTaken float64 `json:"timeTaken"`
}

type Route struct {
	InAmount              float64      `json:"inAmount"`
	OutAmount             float64      `json:"outAmount"`
	OutAmountWithSlippage float64      `json:"outAmountWithSlippage"`
	PriceImpactPct        float64      `json:"priceImpactPct"`
	MarketInfos           []MarketInfo `json:"marketInfos"`
}

type MarketInfo struct {
	ID                 string  `json:"id"`
	Label              string  `json:"label"`
	InputMint          string  `json:"inputMint"`
	OutputMint         string  `json:"outputMint"`
	NotEnoughLiquidity bool    `json:"notEnoughLiquidity"`
	InAmount           float64 `json:"inAmount"`
	OutAmount          float64 `json:"outAmount"`
	PriceImpactPct     float64 `json:"priceImpactPct"`
	LpFee              Fee     `json:"lpFee"`
	PlatformFee        Fee     `json:"platformFee"`
}

type Fee struct {
	Amount float64 `json:"amount"`
	Mint   string  `json:"mint"`
	Pct    float64 `json:"pct"`
}

const swapUrl = "https://quote-api.jup.ag/v1/swap"

func GetSwapTransactions(swap *SwapRequest) (resp SwapResponse, err error) {
	// Get the serialized transaction(s) from Jupiter's Swap API
	var swapJson bytes.Buffer
	err = json.NewEncoder(&swapJson).Encode(&swap)
	if err != nil {
		panic(err)
	}

	r, err := http.Post(swapUrl, "application/json", &swapJson)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	sr := SwapResponse{}
	err = json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", sr)

	return sr, nil
}

type QuoteRequest struct {
	InputMint        string
	OutputMint       string
	Amount           float64
	Slippage         float64
	FeeBps           float64
	OnlyDirectRoutes bool
}

func GetQuote(qr *QuoteRequest) (q Quote, err error) {
	quoteUrl, err := url.Parse("https://quote-api.jup.ag")
	if err != nil {
		panic(err)
	}

	quoteUrl.Path += "/v1/quote"

	amountLamports := qr.Amount * 1000000000
	a := fmt.Sprintf("%f", amountLamports)
	s := fmt.Sprintf("%f", qr.Slippage)
	f := fmt.Sprintf("%f", qr.FeeBps)

	params := url.Values{}
	params.Add("inputMint", qr.InputMint)
	params.Add("outputMint", qr.OutputMint)
	params.Add("amount", a)
	params.Add("slippage", s)
	params.Add("feeBps", f)

	if qr.OnlyDirectRoutes {
		params.Add("onlyDirectRoutes", "true")
	} else {
		params.Add("onlyDirectRoutes", "false")
	}

	quoteUrl.RawQuery = params.Encode()
	fmt.Printf("Encoded URL is %q\n", quoteUrl.String())

	resp, err := http.Get(quoteUrl.String())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	quote := Quote{}
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", quote)

	return quote, nil
}