package main

import (
	"backtest/app/common"
	"fmt"
	"time"
)

func main() {

	startTime := time.Now() //処理開始時間
	fmt.Printf("開始時間 = %s\n", startTime.Format("2006-01-02 15:04:05"))

	// setup.CandleFilter()

	// strategy.AnalysisN225()
	// strategy.AnalysisN225Gyakubari()
	// strategy.AnalysisN225Mini()
	// strategy.AnalysisN225MiniGyakubari()

	// common.TotalAutoAnalysis(false, false)
	// common.TotalAutoAnalysis(true, false)
	// common.TotalAutoAnalysis(false, true)
	// common.TotalAutoAnalysis(true, true)

	common.TotalManualAnalysis()
	common.TotalManualAnalysisGyakubari()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("終了時間 = %s\n", endTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("処理時間 = %v\n", duration)
}
