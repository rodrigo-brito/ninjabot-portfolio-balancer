package strategy

import (
	"math"

	"github.com/rodrigo-brito/ninjabot"
	"github.com/rodrigo-brito/ninjabot/model"
	"github.com/rodrigo-brito/ninjabot/service"
	log "github.com/sirupsen/logrus"
)

type Balancer struct {
	AssetWeight map[string]float64
	LastClose   map[string]float64
}

type Weight struct {
	Pair   string
	Weight float64
}

func NewBalancer(weights ...Weight) *Balancer {
	s := &Balancer{
		AssetWeight: make(map[string]float64),
		LastClose:   make(map[string]float64),
	}

	for _, w := range weights {
		s.AssetWeight[w.Pair] = w.Weight
	}

	return s
}

func (b Balancer) Timeframe() string {
	return "1w"
}

func (b Balancer) WarmupPeriod() int {
	return 1
}

func (b Balancer) Indicators(df *model.Dataframe) {
	b.LastClose[df.Pair] = df.Close.Last(0)
}

func (b Balancer) CalculatePositionAdjustment(df *ninjabot.Dataframe, broker service.Broker) (expect, diff float64, err error) {
	totalEquity := 0.0

	for p, _ := range b.AssetWeight {
		asset, _, err := broker.Position(p)
		if err != nil {
			return 0, 0, err
		}

		totalEquity += asset * b.LastClose[p]
	}

	asset, quote, err := broker.Position(df.Pair)
	if err != nil {
		return 0, 0, err
	}

	totalEquity += quote // include free cash to calculate the total equity

	targetSize := b.AssetWeight[df.Pair] * totalEquity
	return targetSize, asset*b.LastClose[df.Pair] - targetSize, nil
}

func (b Balancer) OnCandle(df *model.Dataframe, broker service.Broker) {
	_, quotePosition, err := broker.Position(df.Pair)
	if err != nil {
		log.Error(err)
		return
	}

	expected, diff, err := b.CalculatePositionAdjustment(df, broker)
	if err != nil {
		log.Error(err)
		return
	}

	// avoid small operations
	if math.Abs(diff)/expected < 0.01 || math.Abs(diff) < 10 {
		return
	}

	if diff > 0 {
		// Sell excess of coins
		_, err = broker.CreateOrderMarketQuote(ninjabot.SideTypeSell, df.Pair, diff)
	} else {
		if diff > quotePosition {
			log.Errorf("free cash not enough, DIFF = %.2f USDT, CASH = %.2f USDT", diff, quotePosition)
			return
		}

		// Buy more coins
		_, err = broker.CreateOrderMarketQuote(ninjabot.SideTypeBuy, df.Pair, -diff)
	}
	if err != nil {
		log.Error(err)
	}
}
