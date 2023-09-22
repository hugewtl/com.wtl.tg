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

/*删除目标zset中的指定范围【min,max】成员*/
func ZRemTargetKeys(redisdb *redis.Client, targetKey chan string, min string, max string, master string, funcFlag int) bool {
	/*
		对目标数据进行删除操作
	*/
	for {
		key, ok := <-targetKey
		if !ok {
			break
		}
		sumRemMems_rem := redisdb.ZRemRangeByScore(key, min, max)

		if sumRemMems_rem.Val() > 0 {
			/*
				删除元素计数器
			*/
			sumRemMems += sumRemMems_rem.Val()
			time.Sleep(100 * time.Millisecond)
		} else {
			if funcFlag == 1 {
				keyChan_chk2 <- key
			}
			if funcFlag == 2 {
				keyChan_chk1 <- key
			}
		}

		if err != nil {
			panic(err)
		} else {
			log.Printf("%v 已删除 [ %v ] 中成员目标【withscores: %v--->%v】数据.", master, key, min, max)
		}
	}

	return true
}

/*
通过扫描目标key数据，删除zset成员数据
注:SCAN 命令用于迭代当前数据库中的数据库键。
go redis pipeline
*/
func ScanKeys(redisdb *redis.Client, KeyPattern string, keyChan chan string, keyChanRm chan string, master string) bool {
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
					/*需要发送给两种score类型的队列*/
					keyChanRm <- key
				}
				keysMap[key] = key
				log.Printf("%v Zset类型的key from client: %v \n", master, key)
			}
		}
		/*
		   使用pipeline减少网络交互次数
		*/
		for cursor != 0 {
			// log.Printf("%v 当前cursor1: %v", master, cursor)
			pipeline.Scan(cursor, KeyPattern, 1000)
			cmdScaner, err := pipeline.Exec()
			if err != nil {
				panic(err)
			}
			// log.Printf("%v 当前cursor2: %v", master, cursor)
			for _, cmder := range cmdScaner {
				cmd := cmder.(*redis.ScanCmd)
				keys, cursor, err = cmd.Result()
				// log.Printf("%v 当前cursor3: %v", master, cursor)
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
						log.Printf("%v Zset类型的key from pipeline: %v \n", master, key)
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

/*连接Redis Connector init*/
func ConnRedisCluster(RClusterUrl string, password string) (*redis.Client, error) {
	//3s超时退出
	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if strings.Contains(RClusterUrl, ",") {
		/*如果传值是备用节点列表，需要遍历备用节点并验证可用性，返回可用节点*/
		for _, redisNodeBak := range strings.Split(RClusterUrl, ",") {
			redisdb = redis.NewClient(&redis.Options{
				Addr:     redisNodeBak,
				PoolSize: 200,
				Password: password,
			})
			if _, err = redisdb.Ping().Result(); err == nil {
				/*找到可用节点就退出遍历*/
				break
			}
		}
		return redisdb, err
	}
	redisdb = redis.NewClient(&redis.Options{
		Addr:     RClusterUrl,
		PoolSize: 200,
		Password: password,
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
		/*按数据排布规律，角色名称的上一个下标就是角色对应的IP:PORT,剔除fail的master节点，因为已转移*/
		if strings.Contains(substr, "master") && !strings.Contains(substr, "master,fail") {
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
保存修改rdb
*/
func RbgSave(redisdb *redis.Client, master string) {
	// oldSave := redisdb.LastSave()
	rSc := redisdb.BgSave()
	if rSc.String() == "bgsave: Background saving started" {
		log.Printf("%v --> %v --> %v \n", master, redisdb.LastSave().Val(), rSc)
		log.Printf("%v Background saving terminated with success! \n", master)
	}
}

func ifDirExist(path string) bool {
	// 判断路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，创建它
		err := os.MkdirAll(path, 0700)
		if err != nil {
			fmt.Println("创建路径失败:", err)
			return false
		}
		fmt.Println("路径已创建:", path)
		return true
	} else {
		fmt.Println("路径已存在:", path)
		return true
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
