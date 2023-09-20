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

	"gopkg.in/redis.v5"
)

/*
通过扫描目标key数据，删除zset成员数据
注:SCAN 命令用于迭代当前数据库中的数据库键。
go redis pipeline
*/
func ScanKeys(redisdb *redis.Client, KeyPattern string, keyChan chan string, keyChanRm chan string) bool {
	var keysMap = make(map[string]string, 5000)
	var cursor uint64
	var n int
	var isOver bool
	startTime := time.Now().UnixNano() / 1e6
	pipeline := redisdb.Pipeline()
	for {
		var keys = make([]string, 20)
		keys, cursor, err = redisdb.Scan(cursor, KeyPattern, 1000).Result()

		if err != nil {
			panic(err)
		}
		for _, key := range keys {
			/*将zset类型的key存到map中*/
			n++
			if redisdb.Type(key).Val() == "zset" {
				//处理zset类型key:将key发送给keyChan
				keyChan <- key
				if delMems {
					keyChanRm <- key
				}
				keysMap[key] = key
				log.Printf("Zset类型的key: %v \n", key)
			}
		}

		for cursor != 0 {
			pipeline.Scan(cursor, KeyPattern, 1000)
			cmdScaner, err := pipeline.Exec()
			if err != nil {
				panic(err)
			}
			for _, cmder := range cmdScaner {
				cmd := cmder.(*redis.ScanCmd)
				keys, cursor, err = cmd.Result()
				log.Printf("数据游标：%v", cursor)
				if err != nil {
					panic(err)
				}

				for _, key := range keys {
					/*将zset类型的key存到map中*/
					n++
					if redisdb.Type(key).Val() == "zset" {
						//处理zset类型key:将key发送给keyChan
						keyChan <- key
						if delMems {
							keyChanRm <- key
						}
						keysMap[key] = key
						log.Printf("Zset类型的key1: %v \n", key)
					}
				}
			}
		}

		if cursor == 0 {
			allKeys += n
			fmt.Printf("当前节点遍历Key数量: %v \n", n)
			fmt.Printf("当前节点统计Zset的Key数量: %v \n", len(keysMap))
			sumTime += (time.Now().UnixNano()/1e6 - startTime)
			fmt.Printf("遍历耗时（毫秒）   : %v \n", time.Now().UnixNano()/1e6-startTime)
			isOver = true
			break
		}

	}

	return isOver
}

// func PrintKeys(redisdb *redis.Client, keys []string, keyChan chan string, keyChanRm chan string) {
// 	for _, key := range keys {
// 		/*将zset类型的key存到map中*/
// 		n++
// 		if redisdb.Type(key).Val() == "zset" {
// 			//处理zset类型key:将key发送给keyChan
// 			keyChan <- key
// 			if delMems {
// 				keyChanRm <- key
// 			}
// 			keysMap[key] = key
// 			log.Printf("Zset类型的key: %v \n", key)
// 		}
// 	}
// }

/*连接Redis Connector init*/
func ConnRedisCluster(RClusterUrl string) (*redis.Client, error) {
	//3s超时退出
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	redisdb = redis.NewClient(&redis.Options{
		Addr:     RClusterUrl,
		PoolSize: 200,
	})
	_, err = redisdb.Ping().Result()
	if err != nil {
		fmt.Println("ping redis failed err:", err)
		return nil, err
	}
	return redisdb, err
}

/*依据集群节点角色，提取master节点信息*/
func getCMasterNodes(redisdb *redis.Client) map[int]string {
	rsc := redisdb.ClusterNodes()
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

/*
** 保存修改rdb
 */
func RbgSave(redisdb *redis.Client) {
	pipeline := redisdb.Pipeline()
	oldSave := pipeline.LastSave()
	rSc := pipeline.BgSave()
	newSave := pipeline.LastSave()
	pipeline.Exec()
	if rSc.String() == "bgsave: Background saving started" && newSave.Val() != oldSave.Val() {
		log.Println("Background saving terminated with success!")
	}
}

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
