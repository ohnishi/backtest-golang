package common

import (
	"math"
	"runtime"
	"sync"
	"time"
)

func AnalysisGyakubari(strategys []StrategyParameter, tableName string, tradeCost, unit, magnification int) []StrategyParameter {

	toDate := time.Date(TotalAnalysisYearFrom, TotalAnalysisMonthFrom, 1, 0, 0, 0, 0, time.UTC)
	fromDate := toDate.AddDate(0, -1, 0)
	moveMounthSize := monthDiff(toDate, time.Date(TotalAnalysisYearTo, TotalAnalysisMonthTo, 1, 0, 0, 0, 0, time.UTC))

	moveMounth := 0
	for moveMounth < moveMounthSize {

		addMonth := moveMounth * -1
		candles := CandleToArray(tableName, fromDate.AddDate(0, addMonth, 0), toDate.AddDate(0, addMonth, 0))
		// _, err := sess.Select("*").From(tableName).Where("date_time >= ? and date_time <= ?", fromDate.AddDate(0, addMonth, 0), toDate.AddDate(0, addMonth, 0)).OrderBy("date_time asc").LoadStructs(&candles)
		// FailOnError(err)

		// for _, c := range candles {
		// 	fmt.Println(c)
		// }

		//並列処理
		var wg sync.WaitGroup
		c := make(chan int, runtime.NumCPU())
		for i := 0; i < len(strategys); i++ {
			if strategys[i].LoseCount > MaxLoseTrade {
				continue
			}
			wg.Add(1)
			go func(cur int) {
				c <- 1
				defer func() {
					<-c
					wg.Done()
				}()
				strategyExecGyakubari(&strategys[cur], candles, tradeCost, unit, magnification)
			}(i)
		}
		wg.Wait()
		moveMounth++

		// fmt.Println(candles[0].DateTime.Format("2006-01"))
	}

	return filter(strategys)
}

func strategyExecGyakubari(strategy *StrategyParameter, candles []Candle, tradeCost, unit, magnification int) {

	totalRieki, totalSon, tradeCount := strategy.ExecGyakubari(candles, unit)
	cost := tradeCount * tradeCost
	allRieki := ((totalRieki - totalSon) * magnification) - cost

	if allRieki <= 0 {
		strategy.LoseCount++
	} else {
		strategy.WinCount++
		if strategy.WinCount >= ResetWinTrade {
			strategy.WinCount = 0
			strategy.LoseCount = 0
		}
	}
}

func (m *StrategyParameter) ExecGyakubari(candles []Candle, unit int) (int, int, int) {
	buyOrderIndexs := make([]OrderInfo, 0, MaxTradeSize)
	sellOrderIndexs := make([]OrderInfo, 0, MaxTradeSize)

	var isFirst, isHighTouch, isLowTouch bool
	var pointHigh, pointLow, volume int

	for i, candle := range candles {

		if candle.DateTime.Weekday() != time.Friday {
			continue
		}

		// 始め値の判断もうちょっと上手いやりかたしたい
		if candle.Type == CANDLE_TYEP_OPEN {
			isHighTouch = false
			isLowTouch = false
			isFirst = true
			pointHigh = int(math.MinInt64)
			pointLow = int(math.MaxInt64)
			volume = 0
		}

		if isHighTouch == false || isLowTouch == false {

			hour, minute := getTradeHourMinute(candle.DateTime)

			if m.IsAnariseTo(hour, minute) {

				volume += candle.Volume

				if m.IsAnariseFrom(hour, minute) {
					if pointHigh < candle.High {
						pointHigh = candle.High
						// isHighTouch = false
						// isLowTouch = true
					}
					if pointLow > candle.Low {
						pointLow = candle.Low
						// isHighTouch = true
						// isLowTouch = false
					}
				}
			} else if hour > m.OutHour || (hour == m.OutHour && minute >= m.OutMintue) {

				if isFirst {
					pointHigh -= pointHigh % unit
					amari := (pointLow % unit)
					if amari > 0 {
						pointLow += unit - amari
					}
					isFirst = false
				}

				if (hour < m.EndHour || (hour == m.EndHour && minute < m.EndMintue)) && pointHigh-pointLow <= m.Max && volume >= m.Volume {
					if isHighTouch == false && pointHigh-m.Hosei <= candle.High && pointHigh-m.Hosei >= candle.Low {
						// rieki, son := m.GetRiekiSon(pointHigh, pointLow, unit)
						sellOrderIndexs = append(sellOrderIndexs, OrderInfo{i, pointHigh - m.Hosei, m.Rieki, m.Son})
						isHighTouch = true
					} else if isLowTouch == false && pointLow+m.Hosei >= candle.Low && pointLow+m.Hosei <= candle.High {
						// rieki, son := m.GetRiekiSon(pointHigh, pointLow, unit)
						buyOrderIndexs = append(buyOrderIndexs, OrderInfo{i, pointLow + m.Hosei, m.Rieki, m.Son})
						isLowTouch = true
					}
				} else {
					isHighTouch = true
					isLowTouch = true
				}
			}
		}
	}

	//検証
	kesaiList := make([]int, 0, MaxTradeSize*2)
	for _, buyOrder := range buyOrderIndexs {
		for _, candle := range candles[buyOrder.index+1 : len(candles)] {
			if buyOrder.price+buyOrder.rieki < candle.High {
				kesaiList = append(kesaiList, buyOrder.rieki)
				break
			} else if buyOrder.price-buyOrder.son >= candle.Low {
				kesaiList = append(kesaiList, -buyOrder.son)
				break
			} else if candle.Type == CANDLE_TYEP_CLOSE {
				kesaiList = append(kesaiList, candle.Close-buyOrder.price)
				break
			}
		}
	}
	for _, sellOrder := range sellOrderIndexs {
		for _, candle := range candles[sellOrder.index+1 : len(candles)] {
			if sellOrder.price-sellOrder.rieki > candle.Low {
				kesaiList = append(kesaiList, sellOrder.rieki)
				break
			} else if sellOrder.price+sellOrder.son <= candle.High {
				kesaiList = append(kesaiList, -sellOrder.son)
				break
			} else if candle.Type == CANDLE_TYEP_CLOSE {
				kesaiList = append(kesaiList, sellOrder.price-candle.Close)
				break
			}
		}
	}

	rieki := 0
	son := 0
	for _, saeki := range kesaiList {
		// fmt.Println(saeki * 1000)
		if saeki > 0 {
			rieki += saeki
		} else {
			son -= saeki
		}
	}
	return rieki, son, len(kesaiList)
}
