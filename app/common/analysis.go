package common

import (
	"runtime"
	"sync"
	"time"
)

func Analysis(strategys []StrategyParameter, tableName string, tradeCost, unit, magnification int) []StrategyParameter {

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
				strategyExec(&strategys[cur], candles, tradeCost, unit, magnification)
			}(i)
		}
		wg.Wait()
		moveMounth++

		// fmt.Println(candles[0].DateTime.Format("2006-01"))
	}

	return filter(strategys)
}

func strategyExec(strategy *StrategyParameter, candles []Candle, tradeCost, unit, magnification int) {

	totalRieki, totalSon, tradeCount := strategy.Exec(candles, unit)
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

func filter(strategys []StrategyParameter) []StrategyParameter {

	// fmt.Println("strategys === ", strategys)
	// fmt.Println("MaxLoseTrade === ", MaxLoseTrade)
	ret := []StrategyParameter{}
	for i := 0; i < len(strategys); i++ {

		// fmt.Println("strategys[i].LoseCount === ", strategys[i].LoseCount)
		if strategys[i].LoseCount <= MaxLoseTrade {
			ret = append(ret, strategys[i])
		}
	}
	// fmt.Println("ret === ", ret)
	return ret
}

func monthDiff(d1, d2 time.Time) int {
	if d2.After(d1) {
		d1, d2 = d2, d1
	}

	return (d1.Year()-d2.Year())*12 + int(d1.Month()) - int(d2.Month())
}
