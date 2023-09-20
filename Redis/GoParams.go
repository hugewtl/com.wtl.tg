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
)
