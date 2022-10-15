package common

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const TotalAnalysisYearFrom = 2018
const TotalAnalysisMonthFrom = 10
const TotalAnalysisDayFrom = 1

const TotalAnalysisYearTo = TotalAnalysisYearFrom
const TotalAnalysisMonthTo = TotalAnalysisMonthFrom - 1
const TotalAnalysisDayTo = TotalAnalysisDayFrom

const TradeCostNK225 = 540
const TradeUnitNK225 = 10
const TradeMagnificationNK225 = 1000
const TableNameNK225 = "nk225"
const TableNameNK225Evening = "nk225_evening"

const TradeCostNK225mini = 54
const TradeUnitNK225mini = 5
const TradeMagnificationNK225mini = 100
const TableNameNK225mini = "nk225mini"
const TableNameNK225miniEvening = "nk225mini_evening"

//最大月間負けトレード数
const MaxLoseTrade = 0
const ResetWinTrade = 12

//１月の最大トレード回数
const MaxTradeSize = 25

//ストラテジー保存キー
const StrategyKeyPrefix = "strategy_%d"

//ストラテジー最大テスト数
// const StrategyMaxSize = 25000000
const StrategyMaxSize = 20000000

const (
	CANDLE_TYEP_OPEN  = 1
	CANDLE_TYEP_CLOSE = 2
)

const (
	CSV_INDEX_DATE   = 0
	CSV_INDEX_TIME   = 1
	CSV_INDEX_OPEN   = 2
	CSV_INDEX_HIGH   = 3
	CSV_INDEX_LOW    = 4
	CSV_INDEX_CLOSE  = 5
	CSV_INDEX_VOLUME = 6
)

type OrderInfo struct {
	index int
	price int
	rieki int
	son   int
}

type Candle struct {
	DateTime time.Time
	Open     int
	High     int
	Low      int
	Close    int
	Volume   int
	Type     int
}

type StrategyParameter struct {
	InHour    int
	InMintue  int
	OutHour   int
	OutMintue int
	EndHour   int
	EndMintue int
	Rieki     int
	Son       int
	Hosei     int
	Max       int
	Volume    int
	WinCount  int
	LoseCount int
}

func (m *StrategyParameter) ToIntArray(totalPf, tradeCount, allRieki int) []int {
	return []int{
		m.InHour,
		m.InMintue,
		m.OutHour,
		m.OutMintue,
		m.EndHour,
		m.EndMintue,
		m.Rieki,
		m.Son,
		m.Hosei,
		m.Max,
		m.Volume,
		m.WinCount,
		m.LoseCount,
		totalPf,
		tradeCount,
		allRieki,
	}
}

func (m *StrategyParameter) ToString() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d",
		m.InHour,
		m.InMintue,
		m.OutHour,
		m.OutMintue,
		m.EndHour,
		m.EndMintue,
		m.Rieki,
		m.Son,
		m.Hosei,
		m.Max,
		m.Volume,
		m.WinCount,
		m.LoseCount,
	)
}

func (m *StrategyParameter) Exec(candles []Candle, unit int) (int, int, int) {
	buyOrderIndexs := make([]OrderInfo, 0, MaxTradeSize)
	sellOrderIndexs := make([]OrderInfo, 0, MaxTradeSize)

	var isFirst, isHighTouch, isLowTouch bool
	var pointHigh, pointLow, volume int

	for i, candle := range candles {

		if candle.DateTime.Weekday() == time.Friday {
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
						buyOrderIndexs = append(buyOrderIndexs, OrderInfo{i, pointHigh - m.Hosei, m.Rieki, m.Son})
						isHighTouch = true
					} else if isLowTouch == false && pointLow+m.Hosei >= candle.Low && pointLow+m.Hosei <= candle.High {
						// rieki, son := m.GetRiekiSon(pointHigh, pointLow, unit)
						sellOrderIndexs = append(sellOrderIndexs, OrderInfo{i, pointLow + m.Hosei, m.Rieki, m.Son})
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

func (m *StrategyParameter) IsAnariseFrom(hour, mintue int) bool {

	if hour > m.InHour || (hour == m.InHour && mintue >= m.InMintue) {
		return true
	}
	return false
}

func (m *StrategyParameter) IsAnariseTo(hour, mintue int) bool {

	if hour < m.OutHour || (hour == m.OutHour && mintue < m.OutMintue) {
		return true
	}
	return false
}

func (m *StrategyParameter) GetRiekiSon(high, low, unit int) (int, int) {
	rieki := (high - low) * m.Rieki / 100
	son := (high - low) * m.Son / 100
	rieki -= rieki % unit
	son -= son % unit
	return max(rieki, 10), max(son, 10)
}

func FailOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getTradeHourMinute(tradeTime time.Time) (int, int) {

	return tradeTime.Hour(), tradeTime.Minute()
}

func ToStrategyParameter(values []string) StrategyParameter {

	inHour, _ := strconv.Atoi(values[0])
	inMintue, _ := strconv.Atoi(values[1])
	outHour, _ := strconv.Atoi(values[2])
	outMintue, _ := strconv.Atoi(values[3])
	endHour, _ := strconv.Atoi(values[4])
	endMintue, _ := strconv.Atoi(values[5])
	rieki, _ := strconv.Atoi(values[6])
	son, _ := strconv.Atoi(values[7])
	hosei, _ := strconv.Atoi(values[8])
	max, _ := strconv.Atoi(values[9])
	min, _ := strconv.Atoi(values[10])
	winCount, _ := strconv.Atoi(values[11])
	loseCount, _ := strconv.Atoi(values[12])

	return StrategyParameter{inHour, inMintue, outHour, outMintue, endHour, endMintue, rieki, son, hosei, max, min, winCount, loseCount}
}

// func GetConnetion() *dbr.Session {

// 	user := "hons"
// 	password := "honspw"
// 	host := "127.0.0.1"
// 	port := "3306"
// 	database := "backtest"

// 	conn, err := dbr.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+database+"?parseTime=true", nil)
// 	FailOnError(err)

// 	conn.SetConnMaxLifetime(time.Minute * 5)
// 	return conn.NewSession(nil)
// }

func OutStrategyReportFile(strategys []StrategyParameter, tableName string) {

	fmt.Printf("テスト通過ストラテジー: %d 件\n", len(strategys))

	if len(strategys) > 0 {

		//書き込みファイル準備
		outFile, err := os.Create("./report/strategy_" + tableName + "_" + time.Now().Format("2006-01-02-15-04-05") + ".csv")
		if err != nil {
			log.Fatal("Error:", err)
		}
		defer outFile.Close()

		//redisからデータ読み出し
		for _, strategy := range strategys {
			fmt.Fprintln(outFile, strategy.ToString())
		}
	}
}
