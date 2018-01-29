package storage

import "coinmarketscraper/coin"

type CoinStorage interface {
	Save(<-chan coin.Coin) <-chan struct{}
}
