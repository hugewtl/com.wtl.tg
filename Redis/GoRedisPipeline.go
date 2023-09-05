package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

var (
	err error
	/*初始化redis客户端*/
	redisdb     = redis.NewClusterClient(&redis.ClusterOptions{})
	masterNodes = make(map[int]string, 23)
	/*存储所有配置文件配置项*/
	ParamsMp = make(map[string]string, 15)
	/*并发任务控制*/
	wg sync.WaitGroup
	/*是否删除zset key中成员*/
	delMems bool
	/*单线程keys数量统计*/
	allKeys int
	/*统计CPU耗时*/
	sumTime int64
)

func main() {
	//初始化：加载配置文件
	initParams("config.conf")
	//初始化数据范围
	count, err := strconv.ParseInt(ParamsMp["opt_count"], 10, 64)
	if err != nil {
		panic(err)
	}
	opt_sec := &redis.ZRangeBy{
		Min:    ParamsMp["opt_sec_min"],
		Max:    ParamsMp["opt_sec_max"],
		Offset: 0,
		Count:  count,
	}
	opt_day := &redis.ZRangeBy{
		Min:    ParamsMp["opt_day_min"],
		Max:    ParamsMp["opt_day_max"],
		Offset: 0,
		Count:  count,
	}
	// fmt.Println(opt_sec.Count)

	//初始化日志文件logScan
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logfile, err := os.OpenFile(ParamsMp["logfile"], os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		logfile.Close()
	}()
	log.SetOutput(logfile)
	/*连接redis集群*/
	if redisdb, err = ConnRedisCluster(ParamsMp["RedisNodesUrl"]); err != nil {
		panic(err)
	}
	defer redisdb.Close()
	/*获取集群所有节点信息*/
	getCMasterNodes(redisdb)
	/*判断是否做删除操作*/
	if ParamsMp["ifRemZsetMems"] == "true" {
		delMems = true
	}
	/*遍历master节点*/
	for _, master := range masterNodes {
		/*连接master节点*/
		if master == "192.168.1.173:6000" {
			rsingle := redis.NewClusterClient(&redis.ClusterOptions{})
			if rsingle, err = ConnRedisCluster(master); err != nil {
				panic(err)
			}
			defer rsingle.Close()
			keyChan := make(chan string, 10000)
			keyChanRm := make(chan string, 10000)

			/*扫描master节点上所有的key,提取zset类型的key,放到keyChan*/
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + " 连接成功,扫描key开始... \n")
				for {
					ok := ScanKeys(rsingle, ParamsMp["IndexPrePattern"], keyChan, keyChanRm)
					if ok {
						break
					}
				}
				wg.Done()

				fmt.Printf("master节点 " + local + " 已扫描key完成 \n")
				fmt.Printf("================================================== \n")
			}(master)

			/*在master节点查询zset类型的key指定score时间范围的数据*/
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + "查询zset key（按时间）开始... \n")
				for {
					ok := ZScanByScoreWITHSCORES(rsingle, keyChan, opt_sec)
					if ok {
						break
					}
				}
				fmt.Printf("master节点 " + local + " 扫描key结束，任务退出！\n")
				wg.Done()
			}(master)
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + "查询zset key（按日期）开始... \n")
				for {
					ok := ZScanByScoreWITHSCORES(rsingle, keyChan, opt_day)
					if ok {
						break
					}
				}
				fmt.Printf("master节点 " + local + " 扫描key结束，任务退出！\n")
				wg.Done()
			}(master)
			/*删除数据*/
			if delMems {
				wg.Add(1)
				go func(local string) {
					fmt.Printf("master节点 " + local + "删除zset key（按时间）元素开始... \n")
					for {
						ok := ZRemTargetKeys(rsingle, keyChanRm, opt_sec.Min, opt_sec.Max)
						if ok {
							break
						}
					}
					fmt.Printf("master节点 " + local + " 指定数据删除" + opt_sec.Min + "->" + opt_sec.Max + "，任务退出！\n")
					wg.Done()
				}(master)
				wg.Add(1)
				go func(local string) {
					fmt.Printf("master节点 " + local + "删除zset key元素（按日期）开始... \n")
					for {
						ok := ZRemTargetKeys(rsingle, keyChanRm, opt_day.Min, opt_day.Max)
						if ok {
							break
						}
					}
					fmt.Printf("master节点 " + local + " 指定数据删除" + opt_day.Min + "->" + opt_day.Max + "，任务退出！\n")
					wg.Done()
				}(master)
			}
		}

	}

	wg.Wait()
	fmt.Printf("所有keys数量为: %v \n", allKeys)
	fmt.Printf("遍历keys总耗平均CPU时间: %.2f min \n", (float64(sumTime)/1000/60)/float64(len(masterNodes)))
	fmt.Println("所有节点key已扫描完成")
	fmt.Println("-----------------------------END-----------------------------")
}
