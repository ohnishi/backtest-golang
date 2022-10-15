package strategy

import (
	"backtest/app/common"
)

func AnalysisN225Mini() {

	passedStrategys := []common.StrategyParameter{}
	strategys := make([]common.StrategyParameter, 0, common.StrategyMaxSize)
	for _, inHour := range []int{9, 10} {
		for _, inMintue := range []int{0, 15, 30, 45} {
			for _, outHour := range []int{12, 13} {
				for _, outMintue := range []int{0, 15, 30, 45} {
					for _, endHour := range []int{13, 14} {
						for _, endMintue := range []int{0, 15, 30, 45} {
							for rieki := 70; rieki <= 250; rieki += 10 {
								for son := 70; son <= 250; son += 10 {
									for hosei := 0; hosei <= 50; hosei += 10 {
										for max := 100; max <= 300; max += 20 {
											for volume := 100000; volume <= 300000; volume += 20000 {
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
														passedStrategys = append(passedStrategys, common.Analysis(strategys, common.TableNameNK225mini, common.TradeCostNK225mini, common.TradeUnitNK225mini, common.TradeMagnificationNK225mini)...)
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
	passedStrategys = append(passedStrategys, common.Analysis(strategys, common.TableNameNK225mini, common.TradeCostNK225mini, common.TradeUnitNK225mini, common.TradeMagnificationNK225mini)...)
	common.OutStrategyReportFile(passedStrategys, common.TableNameNK225mini)
}

func AnalysisN225MiniGyakubari() {

	passedStrategys := []common.StrategyParameter{}
	strategys := make([]common.StrategyParameter, 0, common.StrategyMaxSize)
	for _, inHour := range []int{11, 12} {
		for _, inMintue := range []int{0, 15, 30, 45} {
			for _, outHour := range []int{13, 14} {
				for _, outMintue := range []int{0, 15, 30, 45} {
					for _, endHour := range []int{14, 15} {
						for _, endMintue := range []int{0, 15, 30, 45} {
							for rieki := 100; rieki <= 250; rieki += 10 {
								for son := 100; son <= 250; son += 10 {
									for hosei := -30; hosei <= 30; hosei += 10 {
										for max := 500; max <= 500; max += 20 {
											for volume := 1; volume <= 1; volume += 1000 {
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
														passedStrategys = append(passedStrategys, common.AnalysisGyakubari(strategys, common.TableNameNK225mini, common.TradeCostNK225mini, common.TradeUnitNK225mini, common.TradeMagnificationNK225mini)...)
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
	passedStrategys = append(passedStrategys, common.AnalysisGyakubari(strategys, common.TableNameNK225mini, common.TradeCostNK225mini, common.TradeUnitNK225mini, common.TradeMagnificationNK225mini)...)
	common.OutStrategyReportFile(passedStrategys, common.TableNameNK225mini)
}
