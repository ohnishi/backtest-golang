package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func TotalManualAnalysis() {

	strategys := make([]StrategyParameter, 1, 1)

	strategys[0] = StrategyParameter{9,0,12,45,14,15,200,200,30,260,23000,0,0}
	TotalAnalysis(strategys, TableNameNK225, TradeCostNK225, TradeUnitNK225, TradeMagnificationNK225, false)
}

func TotalManualAnalysisGyakubari() {

	strategys := make([]StrategyParameter, 1, 1)

	strategys[0] = StrategyParameter{12,45,13,30,14,0,180,90,-10,1000,5000,0,0}
	TotalAnalysis(strategys, TableNameNK225, TradeCostNK225, TradeUnitNK225, TradeMagnificationNK225, true)
}

func TotalAutoAnalysis(isGyakubari, isMini bool) {

	reportFolder := "./report/"
	files, _ := ioutil.ReadDir(reportFolder)
	for _, f := range files {

		fileName := f.Name()
		if strings.Index(fileName, "strategy_") != 0 {
			continue
		}
		// fmt.Println(fileName)

		//読み込みファイル準備
		inFile, err := os.Open(reportFolder + "/" + fileName)
		if err != nil {
			log.Fatal("Error:", err)
		}
		reader := csv.NewReader(inFile)

		strategys := []StrategyParameter{}
		for {
			record, err := reader.Read() // 1行読み出す
			if err == io.EOF {
				break
			} else {
				strategys = append(strategys, ToStrategyParameter(record))
			}

		}
		inFile.Close()

		if isMini {
			TotalAnalysis(strategys, TableNameNK225mini, TradeCostNK225mini, TradeUnitNK225mini, TradeMagnificationNK225mini, isGyakubari)
		} else {
			TotalAnalysis(strategys, TableNameNK225, TradeCostNK225, TradeUnitNK225, TradeMagnificationNK225, isGyakubari)
		}

	}
}

func TotalAnalysis(strategys []StrategyParameter, tableName string, cost, unit, magnification int, isGyakubari bool) {

	// sess := GetConnetion()

	toDate := time.Date(TotalAnalysisYearFrom, TotalAnalysisMonthFrom, TotalAnalysisDayFrom, 0, 0, 0, 0, time.UTC)
	fromDate := time.Date(TotalAnalysisYearTo, TotalAnalysisMonthTo, TotalAnalysisDayTo, 0, 0, 0, 0, time.UTC)

	// var candles []Candle
	// _, err := sess.Select("*").From(tableName).Where("date_time >= ? and date_time <= ?", fromDate, toDate).OrderBy("date_time asc").LoadStructs(&candles)
	candles := CandleToArray(tableName, fromDate, toDate)
	// FailOnError(err)

	outResults := [][]int{}
	for i := 0; i < len(strategys); i++ {
		results := totalStrategyExec(&strategys[i], candles, cost, unit, magnification, isGyakubari)
		// fmt.Println(len(results))
		if results != nil {
			outResults = append(outResults, results)
		}
	}

	TotalOutFile(outResults, tableName, isGyakubari)
}

func totalStrategyExec(strategy *StrategyParameter, candles []Candle, cost, unit, magnification int, isGyakubari bool) []int {

	var totalRieki, totalSon, tradeCount int
	if isGyakubari {
		totalRieki, totalSon, tradeCount = strategy.ExecGyakubari(candles, unit)
	} else {
		totalRieki, totalSon, tradeCount = strategy.Exec(candles, unit)
	}

	// fmt.Println("aaaa == ", totalRieki, totalSon, tradeCount)
	tatedamaCnt := 1
	allCost := tradeCount * cost * tatedamaCnt
	allRieki := ((totalRieki-totalSon)*magnification)*tatedamaCnt - allCost
	if totalSon == 0 {
		totalSon = 1
	}
	totalPf := totalRieki * 100 / totalSon
	// if totalPf < 200 {
	// 	return nil
	// }

	return strategy.ToIntArray(totalPf, tradeCount, allRieki)
}

func TotalOutFile(outResults [][]int, tableName string, isGyakubari bool) {

	//書き込みファイル準備
	filePrefix := "./report/total_strategy_"
	if isGyakubari {
		filePrefix = "./report/total_gyakubari_strategy_"
	}
	outFile, err := os.Create(filePrefix + tableName + "_" + time.Now().Format("2006-01-02-15-04-05") + ".csv")
	if err != nil {
		log.Fatal("Error:", err)
	}
	writer := csv.NewWriter(outFile) //utf8

	sort.Sort(ByN{outResults})

	writeResults := [][]string{}
	for _, o := range outResults {
		writeResults = append(writeResults, toStrings(o))
	}
	writer.WriteAll(writeResults)

	writer.Flush()
	outFile.Close()
}

// 基本クラス

type NSslice [][]int

func (n NSslice) Len() int {
	return len(n)
}

func (n NSslice) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// 数字昇順でソートするクラス
type ByN struct {
	NSslice
}

func (n ByN) Less(i, j int) bool {

	if n.NSslice[i][13] == n.NSslice[j][13] {
		if n.NSslice[i][15] == n.NSslice[j][15] {
			if n.NSslice[i][7] == n.NSslice[j][7] {
				if n.NSslice[i][6] == n.NSslice[j][6] {
					return n.NSslice[i][9] < n.NSslice[j][9]
				}
				return n.NSslice[i][6] > n.NSslice[j][6]
			}
			return n.NSslice[i][7] < n.NSslice[j][7]
		}
		return n.NSslice[i][15] > n.NSslice[j][15]
	}
	return n.NSslice[i][13] > n.NSslice[j][13]
}

func toStrings(arrays []int) []string {

	results := []string{}
	for _, a := range arrays {
		results = append(results, fmt.Sprint(a))
	}
	return results
}
