package main

import (
	"sync"

	"gopkg.in/redis.v5"
)

var (
	err error
	/*初始化redis客户端*/
	redisdb     = redis.NewClient(&redis.Options{})
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
	/*统计删除的数据量*/
	sumRemMems int64
	/*未被删除的数据的key,调用不同score类型函数再次尝试删除*/
	keyChan_chk1 = make(chan string, 10000)
	keyChan_chk2 = make(chan string, 10000)
)
