package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/redis.v5"
)

func main() {
	//初始化：加载配置文件
	initParams("config.conf")

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
		rsingle := redis.NewClient(&redis.Options{})
		if rsingle, err = ConnRedisCluster(master); err != nil {
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
			if ScanKeys(rsingle, ParamsMp["IndexPrePattern"], keyChan, keyChanRm) {
				//关闭当前线程的keyChan
				close(keyChan)
				if delMems {
					close(keyChanRm)
				}
			}
			wg.Done()
			fmt.Printf("master节点 " + local + " 已扫描key完成 \n")
			fmt.Printf("================================================== \n")
		}(master)

	}
	wg.Wait()
	fmt.Printf("所有keys数量为: %v \n", allKeys)
	fmt.Printf("遍历keys总耗平均CPU时间: %.2f min \n", (float64(sumTime)/1000/60)/float64(len(masterNodes)))
	fmt.Println("所有节点key已扫描完成")
	fmt.Println("-----------------------------END-----------------------------")

}
