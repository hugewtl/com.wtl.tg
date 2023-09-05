package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

/*删除目标zset中的指定范围【min,max】成员*/
func ZRemTargetKeys(redisdb *redis.ClusterClient, targetKey chan string, min string, max string) bool {
	//上下文变量
	ctx := context.Background()
	pipeline := redisdb.Pipeline()
	for {
		key, ok := <-targetKey
		if !ok {
			break
		}
		if len(key) == 0 {
			break
		}
		pipeline.ZRemRangeByScore(ctx, key, min, max)
		log.Printf("删除key中成员数据: %v", key)

		cmders, err := pipeline.Exec(ctx)
		if err != nil {
			panic(err)
		}
		if len(cmders) == 0 {
			break
		}
	}
	/*保存数据到磁盘，本任务结束，放入计数器*/
	RbgSave(redisdb)
	// threadChan <- true
	return true
}

func ZScanByScoreWITHSCORES(redisdb *redis.ClusterClient, targetKey chan string, opt *redis.ZRangeBy) bool {
	//上下文变量
	ctx := context.Background()
	pipeline := redisdb.Pipeline()
	/*判断管道中要有数据，否则会阻塞*/
	for {
		/*封装pipeline*/
		for i := 0; i < 10; i++ {
			key, ok := <-targetKey
			if !ok {
				break
			}

			if len(key) == 0 {
				break
			}
			pipeline.ZRangeByScoreWithScores(ctx, key, opt)
		}

		cmders, err := pipeline.Exec(ctx)
		if err != nil {
			panic(err)
		}
		if len(cmders) == 0 {
			break
		}
		for _, cmder := range cmders {
			cmd := cmder.(*redis.ZSliceCmd)
			// fmt.Println(cmd)
			members, err := cmd.Result()
			if err != nil {
				panic(err)
			}
			if len(members) == 0 {
				break
			}
			for _, zkMems := range members {
				log.Printf("目标成员数据:  %v\n", fmt.Sprintf("%.f,%v", zkMems.Score, zkMems.Member))
			}
		}
	}
	/*任务线程退出，加入提出队列*/
	// threadChan <- true
	return true
}

func RbgSave(redisdb *redis.ClusterClient) {
	//上下文变量
	ctx := context.Background()
	pipeline := redisdb.Pipeline()
	oldSave := pipeline.LastSave(ctx)
	rSc := pipeline.BgSave(ctx)
	newSave := pipeline.LastSave(ctx)
	pipeline.Exec(ctx)
	/*跟踪后台bgsave保存状态是否成功*/
	// for {
	// fmt.Println(oldSave)
	// fmt.Println(rSc)
	// fmt.Println(newSave)
	if rSc.String() == "bgsave: Background saving started" && newSave.Val() != oldSave.Val() {
		// fmt.Println("Background saving terminated with success!")
		log.Println("Background saving terminated with success!")
		// break
	}
	// }
}

/*
通过扫描目标key数据，删除zset成员数据
注:SCAN 命令用于迭代当前数据库中的数据库键。
go redis pipeline
*/
func ScanKeys(redisdb *redis.ClusterClient, KeyPattern string, keyChan chan string, keyChanRm chan string) bool {
	var keysMap = make(map[string]string, 5000)
	var cursor uint64
	var n int
	startTime := time.Now().UnixNano() / 1e6
	// fmt.Printf("遍历开始：%v \n")
	pipeline := redisdb.Pipeline()
	ctx := context.Background()
	for {
		keys := make([]string, 20)
		pipeline.Scan(ctx, cursor, KeyPattern, 100)
		cmdScaner, err := pipeline.Exec(ctx)
		if err != nil {
			panic(err)
		}

		for _, cmder := range cmdScaner {
			cmd := cmder.(*redis.ScanCmd)
			keys, cursor, err = cmd.Result()
			if err != nil {
				panic(err)
			}
			if len(keys) == 0 {
				break
			}

			for _, key := range keys {
				/*将zset类型的key存到map中*/
				n++
				if GetKeyType(redisdb, key) == "zset" {
					//处理zset类型key:将key发送给keyChan
					keyChan <- key
					if delMems {
						keyChanRm <- key
					}
					keysMap[key] = key
					log.Printf("Zset类型的key: %v \n", key)
				}
			}
		}

		if cursor == 0 {
			allKeys += n
			fmt.Printf("当前节点遍历Key数量: %v \n", n)
			fmt.Printf("当前节点统计Zset的Key数量: %v \n", len(keysMap))
			sumTime += (time.Now().UnixNano()/1e6 - startTime)
			fmt.Printf("遍历耗时（毫秒）   : %v \n", time.Now().UnixNano()/1e6-startTime)
			//关闭当前线程的keyChan
			close(keyChan)
			if delMems {
				close(keyChanRm)
			}
			break
		}
	}
	return true
}

func GetKeyType(redisdb *redis.ClusterClient, KeyName string) string {
	//上下文变量
	ctx := context.Background()
	pipeline := redisdb.Pipeline()
	pipeline.Type(ctx, KeyName)
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, cmder := range cmders {
		cmd := cmder.(*redis.StatusCmd)
		return cmd.Val()
	}
	return ""
}

/*连接Redis Connector init*/
func ConnRedisCluster(RClusterUrl string) (*redis.ClusterClient, error) {
	//3s超时退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	redisdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{RClusterUrl},
		PoolSize: 200,
	})
	_, err = redisdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("ping redis failed err:", err)
		return nil, err
	}
	return redisdb, err
}

/*依据集群节点角色，提取master节点信息*/
func getCMasterNodes(redisdb *redis.ClusterClient) map[int]string {
	rsc := redisdb.ClusterNodes(redisdb.Context())
	//空格切割所有字符串
	slicInfo := strings.Fields(rsc.String())
	for i, substr := range slicInfo {
		/*按数据排布规律，角色名称的上一个下标就是角色对应的IP:PORT*/
		if strings.Contains(substr, "master") {
			masterNodes[i] = slicInfo[i-1]
		}
	}
	/*
		for _, m := range masterNodes {
			fmt.Println(m)
		}
	*/
	return masterNodes
}

func initLogFile(filename string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		logfile.Close()
	}()
	// multiWriter := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(logfile)
}

// 将参数存储到map中
// var ParamsMp = make(map[string]string, 15)

/*
读取配置文件，将变量->值提取到map中备用
*/
func initParams(filepath string) map[string]string {
	/*
		读取配置文件，获取配置参数
	*/
	file, err := os.Open(filepath) //try to open file
	if err != nil {
		log.Printf("读取配置文件失败：%v\n", err)
		panic(err)
	}
	/*
		函数调用结束时关闭文件
	*/
	defer file.Close()
	/*创建文件读取对象 */
	r := bufio.NewReader(file)
	for {
		//按行读取
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		//去掉已读取的行首位空格
		str := strings.TrimSpace(string(b))
		//如果以注释符"#"打头，不做处理
		if strings.Index(str, "#") == 0 {
			continue
		} else if len(str) != 0 {
			/**参数配置的"="两边值，映射字段*/
			/*
				获取参数字段名称:k
			*/
			k := strings.TrimSpace(strings.Split(str, "=")[0])
			/*
				获取对应的参数值：v
			*/
			v := strings.TrimSpace(strings.Split(str, "=")[1])
			/*
				将参数赋值存储到map:ParamsMp
			*/
			ParamsMp[k] = v
		}
	}
	return ParamsMp
}

// func ZScanByScore(rdb *redis.ClusterClient, keyName string, opt *redis.ZRangeBy) {
// 	//上下文变量
// 	ctx := context.Background()
// 	pipeline := redisdb.Pipeline()
// 	pipeline.Pipeline().ZRangeByScore(ctx, keyName, opt)

// 	cmders, err := pipeline.Exec(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, cmder := range cmders {
// 		cmd := cmder.(*redis.StringSliceCmd)
// 		// fmt.Println(cmd)
// 		keys, err := cmd.Result()
// 		for _, key := range keys {
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Println(key)
// 		}

// 	}
// }
