# Ninjabot Portfolio Balancer

Example of strategy that balances a given portfolio with weights.

**Related Discussion**: https://github.com/rodrigo-brito/ninjabot/discussions/65

### Downloading Data

```bash
ninjabot download --pair BTCUSDT --timeframe 1d --days 365 --output ./data/btc.csv
ninjabot download --pair ETHUSDT --timeframe 1d --days 365 --output ./data/eth.csv
ninjabot download --pair DOTUSDT --timeframe 1d --days 365 --output ./data/dot.csv
ninjabot download --pair ADAUSDT --timeframe 1d --days 365 --output ./data/ada.csv
```

### How to execute
```go
go run main.go
```

### Teorical Wallet

```go
20% of BTC
20% of ETH
20% of DOT
20% of ADA
20% of USDT (free cash)
```

The strategy will buy or sell each week to balance the portfolio.

### Backtest result

```
INFO[2021-11-19 22:45] [SETUP] Using paper wallet                   
INFO[2021-11-19 22:45] [SETUP] Initial Portfolio = 1000.000000 USDT 
+---------+--------+-----+------+---------+--------+---------+----------+
|  PAIR   | TRADES | WIN | LOSS |  % WIN  | PAYOFF | PROFIT  |  VOLUME  |
+---------+--------+-----+------+---------+--------+---------+----------+
| ETHUSDT |     24 |  24 |    0 | 100.0 % |  0.000 |  755.49 |  2747.22 |
| ADAUSDT |     22 |  22 |    0 | 100.0 % |  0.000 | 1293.77 |  3982.54 |
| BTCUSDT |     19 |  15 |    4 | 78.9 %  |  1.185 |  220.34 |  3471.36 |
| DOTUSDT |     24 |  19 |    5 | 79.2 %  |  8.688 |  990.18 |  4027.96 |
+---------+--------+-----+------+---------+--------+---------+----------+
|   TOTAL |     89 |  80 |    9 |  89.9 % |  2.596 | 3259.78 | 14229.09 |
+---------+--------+-----+------+---------+--------+---------+----------+

--------------
WALLET SUMMARY
--------------
619.457127 ADA = 948.361987 USDT
0.019902 BTC = 1819.799873 USDT
26.713543 DOT = 2450.025785 USDT
0.272265 ETH = 3010.440556 USDT

TRADING VOLUME
ADAUSDT        = 3982.54 USDT
BTCUSDT        = 3471.36 USDT
DOTUSDT        = 4027.96 USDT
ETHUSDT        = 2747.22 USDT

1249.341469 USDT
--------------
START PORTFOLIO = 1000.00 USDT
FINAL PORTFOLIO = 4259.78 USDT
GROSS PROFIT    =  3259.782025 USDT (325.98%)
MARKET (B&H)    =  784.12%
MAX DRAWDOWN    =  -37.48 %
VOLUME          =  14229.09 USDT
COSTS (0.001*V) =  14.23 USDT (ESTIMATION) 
--------------
```
![image](https://user-images.githubusercontent.com/7620947/142710022-6890e84b-8b73-4948-8192-6a6990097cc0.png)

