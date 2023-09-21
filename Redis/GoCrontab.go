package main

import (
	"strconv"
	"time"
)

/*
 * 定时任务赋值时间范围函数
 */
func GoCron(ParamsInit map[string]string) {

	if ParamsInit["ifCrond"] == "true" {
		/*
			依据配置项，是否定时进行取值判断
		*/
		days, err := strconv.ParseInt(ParamsInit["days"], 10, 64)
		if err != nil {
			panic(err)
		}
		interval, err := strconv.ParseInt(ParamsInit["interval"], 10, 64)
		if err != nil {
			panic(err)
		}
		currentTime := time.Now()
		delEndTime := currentTime.AddDate(0, 0, -int(days))

		delStartTime := delEndTime.AddDate(0, 0, -int(interval))
		/*依据系统时间初始化，赋值扫描redis zset成员数据的score范围*/
		ParamsMp["opt_sec_min"] = delStartTime.Format("20060102150405")
		ParamsMp["opt_sec_max"] = delEndTime.Format("20060102150405")
		ParamsMp["opt_day_min"] = delStartTime.Format("20060102")
		ParamsMp["opt_day_max"] = delEndTime.Format("20060102")
		// delEndTime = delStartTime
		// fmt.Println(ParamsMp["opt_sec_min"])
		// fmt.Println(ParamsMp["opt_sec_max"])
		// fmt.Println(ParamsMp["opt_day_min"])
		// fmt.Println(ParamsMp["opt_day_max"])
	}
}
