package main

import (
	"container/list"
	"time"
)

//暂不考虑并发
//滑动窗口
func test(r request) {
	position := time.Now().Second() % windowSize
	stage := time.Now().Second() / windowSize
	if stage != window[position].stage {
		//新的时间段  提取热点数据 + 清空原阶段数据 + 新增数据
		checkHotTopic(localCount)
		//清空原阶段数据  从localCount中删去window[position]
		clearPartData(localCount, window[position])
		//添加新的数据
		window[position] = &count{stage, make(map[string]int)}
		window[position].m[r.topic]++
		localCount.m[r.topic]++
	} else {
		//同一个1s时间段内的数据
		window[position].m[r.topic]++
		localCount.m[r.topic]++
	}
}

//检查热点数据  5s的滑动窗口内超过100次访问，就识别为热点数据
func checkHotTopic(localCount *count) {
	for k, v := range localCount.m {
		if v > hotCount {
			data := dataFromDB(k)
			for localCache.size+len(data) > localCache.maxSize {
				// 考虑到缓存的上限问题，简单的先进先出缓存淘汰
				localCache.size -= len(localCache.l.Back().Value.(string))
				delete(localCache.m, localCache.l.Back().Value.(string))
				localCache.l.Remove(localCache.l.Back())
			}
			localCache.m[k] = data
			localCache.l.PushFront(k)
		}
	}
}

//发现热点数据之后，从数据库拿出热点数据放到缓存中
func dataFromDB(topic string) string {
	//TODO...
	return " "
}

func clearPartData(c1 *count, c2 *count) {
	for k := range c2.m {
		c1.m[k] -= c2.m[k]
		if c1.m[k] == 0 { //如果清零了，就删除这个key
			delete(c1.m, k)
		}
	}
}

//全局缓存  来的请求先从全局缓存拿
type cache struct {
	maxSize int
	size    int
	l       *list.List
	m       map[string]string
}

var localCache = cache{1024 * 8, 0, list.New(), make(map[string]string)}

//count用于统计topic的出现次数
type count struct {
	stage int
	m     map[string]int //[topic]
}

var localCount = &count{-1, make(map[string]int)} //统计整个滑动窗口内的

//一个标识符记录是否还当前时间段  0-59
var flag = -1

//滑动窗口
const windowSize = 5

var window = [windowSize]*count{}

//出现多少次算是热点数据
var hotCount = 100

type request struct {
	topic string
}

func initWindow() {
	for i := range window {
		window[i] = &count{-1, make(map[string]int)}
	}
}
