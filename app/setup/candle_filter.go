package setup

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

//購入した分足データから期近のデータだけを抽出する

const TradeDate = 0
const IndexType = 1
const SecurityCode = 2
const SessionID = 3
const IntervalTime = 4
const OpenPrice = 5
const HighPrice = 6
const LowPrice = 7
const ClosePrice = 8
const TradeVolume = 9
const VWAP = 10
const NumberofTrade = 11
const RecordNo = 12
const ContractMonth = 13

func CandleFilter() {

	outCandleFile("/Users/hons/Downloads/nk225.csv", "./out/N225/", "999")
	// outCandleFile("/Users/hons/Downloads/nk225mini.csv", "./out/N225mini/", "999")
	// outCandleFile("/Users/hons/Downloads/jpx400.csv", "./out/JPX400/", "999")
	// outCandleFile("/Users/hons/Downloads/topix.csv", "./out/TOPIX/", "999")
	// outCandleFile("/Users/hons/Downloads/topixmini.csv", "./out/TOPIXmini/", "999")

	outCandleFile("/Users/hons/Downloads/nk225.csv", "./out_evening/N225/", "003")
	// outCandleFile("/Users/hons/Downloads/nk225mini.csv", "./out_evening/N225mini/", "003")
	// outCandleFile("/Users/hons/Downloads/jpx400.csv", "./out_evening/JPX400/", "003")
	// outCandleFile("/Users/hons/Downloads/topix.csv", "./out_evening/TOPIX/", "003")
	// outCandleFile("/Users/hons/Downloads/topixmini.csv", "./out_evening/TOPIXmini/", "003")
}

func outCandleFile(sourceFile, outDir, sessionID string) {

	//読み込みファイル準備
	inFile, err := os.Open(sourceFile)
	if err != nil {
		log.Fatal("Error:", err)
	}
	reader := csv.NewReader(inFile)

	var outFile *os.File
	var writer *csv.Writer
	tradeDate := ""
	contractMonth := ""
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else if len(record[0]) != 8 {
			//ヘッダー行除外
			continue
		} else {
			tradeTime := toDate(record[TradeDate])
			gengetsuTime := toYearMonth(record[ContractMonth])

			if (tradeTime.Month() == 1 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 3)) ||
				(tradeTime.Month() == 2 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 3)) ||
				(tradeTime.Month() == 4 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 6)) ||
				(tradeTime.Month() == 5 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 6)) ||
				(tradeTime.Month() == 7 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 9)) ||
				(tradeTime.Month() == 8 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 9)) ||
				(tradeTime.Month() == 10 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 12)) ||
				(tradeTime.Month() == 11 && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 12)) ||

				((tradeTime.Month() == 3 && !is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 3)) ||
				((tradeTime.Month() == 3 && is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 6)) ||

				((tradeTime.Month() == 6 && !is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 6)) ||
				((tradeTime.Month() == 6 && is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 9)) ||

				((tradeTime.Month() == 9 && !is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 9)) ||
				((tradeTime.Month() == 9 && is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 12)) ||

				((tradeTime.Month() == 12 && !is2ndFridayAfter(tradeTime)) && (tradeTime.Year() == gengetsuTime.Year() && gengetsuTime.Month() == 12)) ||
				((tradeTime.Month() == 12 && is2ndFridayAfter(tradeTime)) && (tradeTime.Year()+1 == gengetsuTime.Year() && gengetsuTime.Month() == 3)) {

				// fmt.Println(record[TradeDate], record[ContractMonth])

				if tradeDate != record[TradeDate] {
					tradeDate = record[TradeDate]
					contractMonth = record[ContractMonth]

					if writer != nil {
						writer.Flush()
					}
					if outFile != nil {
						outFile.Close()
					}

					dirPath := outDir + tradeDate[0:4] + "/" + tradeDate[4:6] + "/"
					_, err = os.Stat(dirPath)
					if err != nil {
						if err = os.MkdirAll(dirPath, 0777); err != nil {
							panic(err)
						}
					}

					//書き込みファイル準備
					outFile, err = os.Create(dirPath + tradeDate[6:8] + ".csv")
					if err != nil {
						log.Fatal("Error:", err)
					}
					writer = csv.NewWriter(outFile) //utf8
					writer.Write([]string{"Date", "Time", "Open", "High", "Low", "Close", "Volume"})
				}
				if contractMonth == record[ContractMonth] && sessionID == record[SessionID] {
					dateTime, _ := time.Parse("200601021504", record[TradeDate]+record[IntervalTime])
					if dateTime.Hour() <= 7 {
						dateTime = dateTime.AddDate(0, 0, 1)
					}
					line := []string{dateTime.Format("20060102"), dateTime.Format("150405"), record[OpenPrice], record[HighPrice], record[LowPrice], record[ClosePrice], record[TradeVolume]}
					// fmt.Println(line)
					writer.Write(line)
				}
			}
		}
	}
	writer.Flush()
	outFile.Close()
	inFile.Close()
}

func toYearMonth(date string) time.Time {

	t, err := time.Parse("200601", date)
	if err != nil {
		panic(err)
	}
	return t
}

func toDate(date string) time.Time {

	t, err := time.Parse("20060102", date)
	if err != nil {
		panic(err)
	}
	return t
}

func is2ndFridayAfter(t time.Time) bool {

	day := t.Day()
	if day < 8 {
		return false
	} else if day > 14 {
		return true
	} else if t.Weekday() < time.Friday {
		return false
	}
	return true
}
