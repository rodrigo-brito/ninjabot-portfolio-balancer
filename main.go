package main

import (
	"context"

	"balancer/strategy"

	"github.com/rodrigo-brito/ninjabot"
	"github.com/rodrigo-brito/ninjabot/exchange"
	"github.com/rodrigo-brito/ninjabot/plot"
	"github.com/rodrigo-brito/ninjabot/plot/indicator"
	"github.com/rodrigo-brito/ninjabot/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	// Ninjabot settings
	settings := ninjabot.Settings{
		Pairs: []string{
			"BTCUSDT",
			"ETHUSDT",
			"DOTUSDT",
			"ADAUSDT",
		},
	}

	// Wallet with five slots
	// 20% of BTC
	// 20% of ETH
	// 20% of DOT
	// 20% of ADA
	// 20% of USDT (free cash)

	balancer := strategy.NewBalancer(
		strategy.Weight{
			Pair:   "BTCUSDT",
			Weight: 0.2,
		},
		strategy.Weight{
			Pair:   "ETHUSDT",
			Weight: 0.2,
		},
		strategy.Weight{
			Pair:   "DOTUSDT",
			Weight: 0.2,
		},
		strategy.Weight{
			Pair:   "ADAUSDT",
			Weight: 0.2,
		})

	// Load your CSV with historical data
	csvFeed, err := exchange.NewCSVFeed(
		balancer.Timeframe(),
		exchange.PairFeed{
			Pair:      "BTCUSDT",
			File:      "./data/btc.csv",
			Timeframe: "1d",
		},
		exchange.PairFeed{
			Pair:      "ETHUSDT",
			File:      "./data/eth.csv",
			Timeframe: "1d",
		},
		exchange.PairFeed{
			Pair:      "ADAUSDT",
			File:      "./data/ada.csv",
			Timeframe: "1d",
		},
		exchange.PairFeed{
			Pair:      "DOTUSDT",
			File:      "./data/dot.csv",
			Timeframe: "1d",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a storage in memory
	storage, err := storage.FromMemory()
	if err != nil {
		log.Fatal(err)
	}

	// Create a virtual wallet with 1.000 USDT
	wallet := exchange.NewPaperWallet(
		ctx,
		"USDT",
		exchange.WithPaperAsset("USDT", 1000),
		exchange.WithDataFeed(csvFeed),
	)

	// Initialize a chart to plot trading results
	chart, err := plot.NewChart(plot.WithIndicators(
		indicator.EMA(5, "blue"),
	), plot.WithPaperWallet(wallet))
	if err != nil {
		log.Fatal(err)
	}

	bot, err := ninjabot.NewBot(
		ctx,
		settings,
		wallet,
		balancer,
		ninjabot.WithStorage(storage),
		ninjabot.WithBacktest(wallet),
		ninjabot.WithCandleSubscription(chart),
		ninjabot.WithOrderSubscription(chart),
		ninjabot.WithLogLevel(log.WarnLevel),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Execute backtest
	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print trading results
	bot.Summary()
	err = chart.Start()
	if err != nil {
		log.Fatal(err)
	}
}
