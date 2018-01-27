package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
)

type Coin struct {
	ID                string
	Name              string
	Acronym           string
	RankPosition      int
	MarketCap         float64
	Price             float64
	Volume24h         float64
	CirculatingSupply float64
	Change24h         float64
	IsMineable        bool
}

func main() {
	resp, _ := http.Get("https://coinmarketcap.com/2")
	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isTbody := t.Data == "tbody"
			if isTbody {

				coinsScraper(z)
			}

		}
	}

}

func coinsScraper(token *html.Tokenizer) {
	//fmt.Println("Tabla!")
	for {
		tt := token.Next()
		switch {
		case tt == html.StartTagToken:
			t := token.Token()
			if t.Data == "tr" {
				fmt.Println(fmt.Sprintf("%#v", coinScraper2(token, t.Attr[0].Val)))
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "tbody" {
				//fmt.Println("End body")
				return
			}
		}
	}
}
func coinScraper(token *html.Tokenizer) {
	fmt.Println("Moneda!")
	for {
		tt := token.Next()
		switch {
		case tt == html.TextToken:
			t := token.Token()
			dataClean := strings.Replace(strings.Replace(t.Data, " ", "", -1), "\n", "", -1)
			if dataClean != "" {
				fmt.Println(dataClean)
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "tr" {
				//fmt.Println("End Coin")
				return
			}
		}
	}
}

func coinScraper2(token *html.Tokenizer, id string) Coin {
	coinInfo := [10]string{}
	count := 0
	for {
		tt := token.Next()
		switch {
		case tt == html.TextToken:
			t := token.Token()
			dataClean := strings.Replace(strings.Replace(t.Data, " ", "", -1), "\n", "", -1)
			if dataClean != "" {
				coinInfo[count] = dataClean
				count++
			}
		case tt == html.EndTagToken:
			t := token.Token()
			if t.Data == "tr" {
				return buildCoin(id, coinInfo)
			}
		}
	}
}

func buildCoin(coinID string, coinInfo [10]string) Coin {
	coin := Coin{
		ID:         coinID,
		Name:       coinInfo[2],
		Acronym:    coinInfo[1],
		IsMineable: false,
	}
	if pos, err := strconv.Atoi(coinInfo[0]); err == nil {
		coin.RankPosition = pos
	}
	if mc, err := parseDollar(coinInfo[3]); err == nil {
		coin.MarketCap = mc
	}
	if price, err := parseDollar(coinInfo[4]); err == nil {
		coin.Price = price
	}
	if vol, err := parseDollar(coinInfo[5]); err == nil {
		coin.Volume24h = vol
	}
	if cir, err := parseDollar(coinInfo[6]); err == nil {
		coin.CirculatingSupply = cir
	}
	if coinInfo[8] != "*" {
		coin.IsMineable = true
		if change, err := strconv.ParseFloat(coinInfo[8][0:len(coinInfo[8])-1], 64); err == nil {
			coin.Change24h = change
		}
	} else {
		coin.IsMineable = false
		if change, err := strconv.ParseFloat(coinInfo[9][0:len(coinInfo[9])-1], 64); err == nil {
			coin.Change24h = change
		}
	}

	return coin
}

func parseDollar(dollarValue string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(dollarValue[1:], ",", "", -1), 64)
}
