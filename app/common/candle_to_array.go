package common

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Target struct {
	Path      string
	Table     string
	IsConvert bool
}

var targets = []Target{
	{
		Path:      "./out/N225/",
		Table:     "nk225",
		IsConvert: false,
	},
	{
		Path:      "./out/N225mini/",
		Table:     "nk225mini",
		IsConvert: false,
	},
	// {
	// 	Path:      "./out/TOPIX/",
	// 	Table:     "topix",
	// 	IsConvert: true,
	// },
	// {
	// 	Path:      "./out/TOPIXmini/",
	// 	Table:     "topixmini",
	// 	IsConvert: true,
	// },
	// {
	// 	Path:      "./out/JPX400/",
	// 	Table:     "jpx400",
	// 	IsConvert: false,
	// },
	{
		Path:      "./out_evening/N225/",
		Table:     "nk225_evening",
		IsConvert: false,
	},
	{
		Path:      "./out_evening/N225mini/",
		Table:     "nk225mini_evening",
		IsConvert: false,
	},
	// {
	// 	Path:      "./out_evening/TOPIX/",
	// 	Table:     "topix_evening",
	// 	IsConvert: true,
	// },
	// {
	// 	Path:      "./out_evening/TOPIXmini/",
	// 	Table:     "topixmini_evening",
	// 	IsConvert: true,
	// },
	// {
	// 	Path:      "./out_evening/JPX400/",
	// 	Table:     "jpx400_evening",
	// 	IsConvert: false,
	// },
}

//CSVからSliceへインポート
func CandleToArray(tableName string, fromDateTime, toDateTime time.Time) []Candle {

	fromDate := fromDateTime.Format("20060102")
	toDate := toDateTime.Format("20060102")
	// fmt.Println("fromDate = ", fromDate)

	var target Target
	for _, t := range targets {
		if t.Table == tableName {
			target = t
			break
		}
	}
	prices := []Candle{}

	yearDirs, _ := ioutil.ReadDir(target.Path)
	for _, yearDir := range yearDirs {
		yearDirPath := target.Path + yearDir.Name()
		monthDirs, _ := ioutil.ReadDir(yearDirPath)
		for _, monthDir := range monthDirs {
			monthDirPath := yearDirPath + "/" + monthDir.Name() + "/"
			dayFiles, _ := ioutil.ReadDir(monthDirPath)
			for _, dayFile := range dayFiles {

				inFileDate := yearDir.Name()[0:4] + monthDir.Name()[0:2] + dayFile.Name()[0:2]

				if inFileDate < fromDate || inFileDate >= toDate {
					continue
				}
				//読み込みファイル準備
				dayFilePath := monthDirPath + dayFile.Name()
				dayFile, err := os.Open(dayFilePath)
				FailOnError(err)

				reader := csv.NewReader(dayFile)

				records, err := reader.ReadAll()
				if err != nil {
					FailOnError(err)
				}

				finalLine := len(records) - 1
				for i, record := range records {
					if i == 0 {
						//ヘッダー行は無視
						continue
					}

					candleType := 0
					if i == 1 {
						candleType = CANDLE_TYEP_OPEN
					} else if i == finalLine {
						candleType = CANDLE_TYEP_CLOSE
					}

					priceTime, _ := time.Parse("20060102150405", record[CSV_INDEX_DATE]+record[CSV_INDEX_TIME])
					open := convertPriceInt(record[CSV_INDEX_OPEN], target.IsConvert)
					high := convertPriceInt(record[CSV_INDEX_HIGH], target.IsConvert)
					low := convertPriceInt(record[CSV_INDEX_LOW], target.IsConvert)
					close := convertPriceInt(record[CSV_INDEX_CLOSE], target.IsConvert)
					vol := convertPriceInt(record[CSV_INDEX_VOLUME], target.IsConvert)

					price := Candle{
						DateTime: priceTime,
						Open:     open,
						High:     high,
						Low:      low,
						Close:    close,
						Volume:   vol,
						Type:     candleType,
					}
					prices = append(prices, price)
				}
				// log.Printf("%v Finish !", dayFilePath)
				dayFile.Close()
			}
		}
	}
	// fmt.Println("toDate = ", toDate)
	return prices
}

func convertPriceInt(price string, isConvert bool) int {

	if isConvert {
		syosuSize := 0
		index := strings.Index(price, ".")
		if index >= 0 {
			syosuPrice := price[index+1:]
			syosuSize = len(syosuPrice)
			price = price[0:index] + syosuPrice
		}
		for i := 0; i < 2-syosuSize; i++ {
			price += "0"
		}
	}
	p, _ := strconv.Atoi(price)
	return p
}
