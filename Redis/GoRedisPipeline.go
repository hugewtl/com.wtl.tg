package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/redis.v5"
)

func main() {
	if !ifDateValid() {
		fmt.Println("WARNING: 非法的 lic , 请使用有效的 lic ！")
		return
	}
	//初始化：加载配置文件
	dateNow := time.Now().Format("200601021504")
	initParams("config.conf")
	/*
	 * 是否定时任务判断并赋值时间范围
	 */
	GoCron(ParamsMp)
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
	//初始化日志文件logScan
	if ifDirExist(ParamsMp["logdir"]) {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		logfile, err := os.OpenFile(ParamsMp["logdir"]+ParamsMp["logfile"]+"-"+dateNow, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return
		}
		defer func() {
			logfile.Close()
		}()
		log.SetOutput(logfile)
	}

	/*连接redis集群*/
	if redisdb, err = ConnRedisCluster(ParamsMp["RedisNodeUrl"], ParamsMp["Password"]); err != nil {
		log.Printf("RedisNodeUrl refused !")
		/*使用备用节点访问集群*/
		if redisdb, err = ConnRedisCluster(ParamsMp["RedisNodesBackup"], ParamsMp["Password"]); err != nil {
			log.Printf("RedisCluster refused !")
		}
	}
	defer redisdb.Close()
	/*获取集群所有节点信息*/
	getCMasterNodes(redisdb)
	/*监控keyChan_chk1、keyChan_chk1俩队列执行情况，用于关闭队列判定条件*/
	chanChk := make(chan int, len(masterNodes))
	/*判断是否做删除操作*/
	if ParamsMp["ifRemZsetMems"] == "true" {
		delMems = true
	}

	/*遍历master节点*/
	for _, master := range masterNodes {
		/*连接master节点*/
		rsingle := redis.NewClient(&redis.Options{})
		if rsingle, err = ConnRedisCluster(master, ParamsMp["Password"]); err != nil {
			panic(err)
		}
		defer rsingle.Close()
		/*
			为每个master节点定义局部keyChan,存放key数据
		*/
		keyChan := make(chan string, 10000)
		keyChanRm := make(chan string, 10000)

		/*扫描master节点上所有的key,提取zset类型的key,放到keyChan*/
		wg.Add(1)
		go func(local string) {
			fmt.Printf("master节点 " + local + " 连接成功,扫描key开始... \n")
			if ScanKeys(rsingle, ParamsMp["IndexPrePattern"], keyChan, keyChanRm, local) {
				//关闭当前线程的keyChan
				close(keyChan)
				if delMems {
					close(keyChanRm)
				}
			}
			wg.Done()
			fmt.Printf("master节点 " + local + " 已扫描key完成 \n")
			fmt.Printf("============================================================= \n")
		}(master)

		/*删除zset key的数据*/
		if delMems {
			/*
			 * 监控keyChanRm队列完成遍历
			 */
			threads := 2
			jobMonitor := make(chan bool, threads)
			/*
			 * 监控keyChan_chk1、keyChan_chk1完成遍历
			 */
			chkMonitor := make(chan bool, threads)
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + " 删除zset key（按时间）元素开始... \n")
				for {
					ok := ZRemTargetKeys(rsingle, keyChanRm, opt_sec.Min, opt_sec.Max, local, 1)
					if ok {
						break
					}
				}
				log.Printf("master节点 " + local + " 指定数据删除" + opt_sec.Min + "->" + opt_sec.Max + "，任务退出！\n")
				fmt.Printf("master节点 " + local + " 删除zset key元素（按时间）结束！ \n")
				jobMonitor <- true
				wg.Done()
			}(master)
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + " 删除zset key元素（按日期）开始... \n")
				for {
					ok := ZRemTargetKeys(rsingle, keyChanRm, opt_day.Min, opt_day.Max, local, 2)
					if ok {
						break
					}
				}
				log.Printf("master节点 " + local + " 指定数据删除" + opt_day.Min + "->" + opt_day.Max + "，任务退出！\n")

				fmt.Printf("master节点 " + local + " 删除zset key元素（按日期）结束！ \n")
				jobMonitor <- true
				wg.Done()
			}(master)
			/*对漏网数据精确删除*/
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + " 删除zset key（按时间）元素开始... \n")
				for {
					ok := ZRemTargetKeys(rsingle, keyChan_chk1, opt_sec.Min, opt_sec.Max, local, 3)
					if ok {
						break
					}
				}
				log.Printf("master节点 " + local + " 指定数据删除" + opt_sec.Min + "->" + opt_sec.Max + "，任务退出！\n")
				fmt.Printf("master节点 " + local + " 删除zset key元素（按时间）结束！ \n")
				chkMonitor <- true
				wg.Done()
			}(master)
			wg.Add(1)
			go func(local string) {
				fmt.Printf("master节点 " + local + " 删除zset key元素（按日期）开始... \n")
				for {
					ok := ZRemTargetKeys(rsingle, keyChan_chk2, opt_day.Min, opt_day.Max, local, 4)
					if ok {
						break
					}
				}
				log.Printf("master节点 " + local + " 指定数据删除" + opt_day.Min + "->" + opt_day.Max + "，任务退出！\n")
				fmt.Printf("master节点 " + local + " 删除zset key元素（按日期）结束！ \n")
				chkMonitor <- true
				wg.Done()
			}(master)

			/*
				监控队列：keyChan_chk1，keyChan_chk2遍历完成，最后刷盘保存数据到rdb
			*/
			wg.Add(1)
			go func(local string) {
				for {
					if len(chkMonitor) == threads {
						RbgSave(rsingle, local)
						wg.Done()
						break
					}
				}
			}(master)

			/*
			   监控master节点任务状态
			*/
			wg.Add(1)
			go func(local string) {
				for {
					if len(jobMonitor) == threads {
						/*一个master中的keyChanRm执行完了，计数1个master执行任务完成，传递master执行任务完成状态*/
						chanChk <- 1
						wg.Done()
						break
					}
				}
			}(master)

		}

	}
	/*
	 * 全局判断，遍历完所有master节点后，关闭keyChan_chk1，keyChan_chk2
	 */
	if delMems {
		wg.Add(1)
		go func() {
			for {
				if len(chanChk) == len(masterNodes) {
					close(keyChan_chk2)
					close(keyChan_chk1)
					wg.Done()
					break
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("所有keys数量为: %v \n", allKeys)
	fmt.Printf("遍历keys总耗平均CPU时间: %.2f min \n", (float64(sumTime)/1000/60)/float64(len(masterNodes)))
	fmt.Println("所有节点key已扫描完成")
	fmt.Printf("共计删除元素：%d 个\n", sumRemMems)
	fmt.Println("-----------------------------END-----------------------------")
}
