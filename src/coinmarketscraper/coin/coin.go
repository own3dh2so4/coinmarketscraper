package coin

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
