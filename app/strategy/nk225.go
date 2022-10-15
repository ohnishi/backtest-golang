package strategy

import (
	"backtest/app/common"
)

func AnalysisN225() {

	passedStrategys := []common.StrategyParameter{}
	strategys := make([]common.StrategyParameter, 0, common.StrategyMaxSize)
	for _, inHour := range []int{9, 10} {
		for _, inMintue := range []int{0, 15, 30, 45} {
			for _, outHour := range []int{12, 13} {
				for _, outMintue := range []int{0, 15, 30, 45} {
					for _, endHour := range []int{13, 14} {
						for _, endMintue := range []int{0, 15, 30, 45} {
							for rieki := 100; rieki <= 200; rieki += 10 {
								for son := 100; son <= 200; son += 10 {
									for hosei := 10; hosei <= 40; hosei += 10 {
										for max := 100; max <= 300; max += 20 {
											for volume := 20000; volume <= 30000; volume += 1000 {
												if (inHour < outHour || (inHour == outHour && inMintue < outMintue)) &&
													(outHour < endHour || (outHour == endHour && outMintue < endMintue)) &&
													(rieki >= son) &&
													(inHour != 8 || inHour >= 45) &&
													(endHour != 15 || endMintue <= 15) {

													strategys = append(strategys,
														common.StrategyParameter{
															InHour:    inHour,
															InMintue:  inMintue,
															OutHour:   outHour,
															OutMintue: outMintue,
															EndHour:   endHour,
															EndMintue: endMintue,
															Rieki:     rieki,
															Son:       son,
															Hosei:     hosei,
															Max:       max,
															Volume:    volume,
															WinCount:  0,
															LoseCount: 0,
														},
													)
													if len(strategys) >= common.StrategyMaxSize {
														passedStrategys = append(passedStrategys, common.Analysis(strategys, common.TableNameNK225, common.TradeCostNK225, common.TradeUnitNK225, common.TradeMagnificationNK225)...)
														strategys = make([]common.StrategyParameter, 0, common.StrategyMaxSize)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	passedStrategys = append(passedStrategys, common.Analysis(strategys, common.TableNameNK225, common.TradeCostNK225, common.TradeUnitNK225, common.TradeMagnificationNK225)...)
	common.OutStrategyReportFile(passedStrategys, common.TableNameNK225)
}

func AnalysisN225Gyakubari() {

	passedStrategys := []common.StrategyParameter{}
	strategys := make([]common.StrategyParameter, 0, common.StrategyMaxSize)
	for _, inHour := range []int{11, 12, 13} {
		for _, inMintue := range []int{0, 15, 30, 45} {
			for _, outHour := range []int{13, 14} {
				for _, outMintue := range []int{0, 15, 30, 45} {
					for _, endHour := range []int{14, 15} {
						for _, endMintue := range []int{0, 15, 30, 45} {
							for rieki := 50; rieki <= 200; rieki += 10 {
								for son := 50; son <= 200; son += 10 {
									for hosei := -50; hosei <= 0; hosei += 10 {
										for max := 1000; max <= 1000; max += 20 {
											for volume := 5000; volume <= 20000; volume += 1000 {
												if (inHour < outHour || (inHour == outHour && inMintue < outMintue)) &&
													(outHour < endHour || (outHour == endHour && outMintue < endMintue)) &&
													(rieki >= son) &&
													(inHour != 8 || inHour >= 45) &&
													(endHour != 15 || endMintue <= 15) {

													strategys = append(strategys,
														common.StrategyParameter{
															InHour:    inHour,
															InMintue:  inMintue,
															OutHour:   outHour,
															OutMintue: outMintue,
															EndHour:   endHour,
															EndMintue: endMintue,
															Rieki:     rieki,
															Son:       son,
															Hosei:     hosei,
															Max:       max,
															Volume:    volume,
															WinCount:  0,
															LoseCount: 0,
														},
													)
													if len(strategys) >= common.StrategyMaxSize {
														passedStrategys = append(passedStrategys, common.AnalysisGyakubari(strategys, common.TableNameNK225, common.TradeCostNK225, common.TradeUnitNK225, common.TradeMagnificationNK225)...)
														strategys = make([]common.StrategyParameter, 0, common.StrategyMaxSize)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	passedStrategys = append(passedStrategys, common.AnalysisGyakubari(strategys, common.TableNameNK225, common.TradeCostNK225, common.TradeUnitNK225, common.TradeMagnificationNK225)...)
	common.OutStrategyReportFile(passedStrategys, common.TableNameNK225)
}
