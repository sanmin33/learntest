package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

var fmsgType = flag.String("msg", "hello", "消息内容")
var ftimeSecond = flag.Int("time", 1, "每次发送协程的测试时间长度秒")
var fgoSendNum = flag.Int("s", 1, "同时发送消息的go协程数量,执行的时候是从1个发送协程开始一直运行到最大数量")
var fgoGetNum = flag.Int("g", 1, "同时接收消息的go协程数量")
var fchanNum = flag.Int("c", 1, "传递消息的管道数量，默认为1，如果给不是1的数，则管道数量等于发送消息go协程数量")
var fifcircle = flag.Int("if", 0, "是否循环输入所有生产者、信道、消费者数量组合")

var msg string
var timeSecond int
var goSendNum int
var goGetNum int
var chanNum int

func main() {
	flag.Parse()
	msg = *fmsgType
	timeSecond = *ftimeSecond
	goGetNum = *fgoGetNum
	goSendNum = *fgoSendNum
	chanNum = *fchanNum
	if *fifcircle == 0 {
		creatSendMsg(goSendNum, chanNum, goGetNum)
		return
	}
	for s := 1; s <= goSendNum; s++ {
		for c := 1; c <= chanNum; c++ {
			for g := 1; g <= goGetNum; g++ {
				creatSendMsg(s, c, g)
			}
		}
	}
	//	time.Sleep(time.Second * 1) //没这句主程序瞬间退出，sendMsg没有时间完成打印输出
}

func creatSendMsg(goSendNum int, chanNum int, goGetNum int) {

	counts := make([]int, goSendNum)
	stopCh := make(chan int)
	msgChan := make([]chan interface{}, chanNum, 100)
	for i := 0; i < chanNum; i++ {
		msgChan[i] = make(chan interface{}, 100)
	}
	for i := 0; i < goSendNum; i++ {
		go sendMsg(msgChan, msg, stopCh, &counts[i]) //创建一个到n个发送协程，本函数执行完成后才交回控制权给main函数
	}
	for i := 0; i < goGetNum; i++ {

		go getMsg(msgChan, stopCh)
	}
	time.Sleep(time.Second * time.Duration(timeSecond))

	close(stopCh)
	for runtime.NumGoroutine() > 1 {
		time.Sleep(time.Millisecond * 1) //如果生产者消费者还没有执行完，就等一下。
	}
	totalMsgNum := 0
	for i := 0; i < goSendNum; i++ {
		totalMsgNum = totalMsgNum + counts[i]
	}
	fmt.Println("发送协程数:,", goSendNum, ",消息管道数:,", chanNum, ", 接收协程数：,", goGetNum, ",测试时长:,", timeSecond, ",发送消息总数据量：,", totalMsgNum)
}

func sendMsg(msgChan []chan interface{}, msg string, stopflag chan int, msgNum *int) {
	*msgNum = 0
	for {
		for i := 0; i < len(msgChan); i++ {
			select {
			case msgChan[i] <- msg:
				*msgNum++
			default:
				continue
			}
		}
		select {
		case <-stopflag:
			return
		default:
			continue
		}
	}
}
func getMsg(msgChan []chan interface{}, stopCh chan int) {
	for {
		for i := 0; i < len(msgChan); i++ {
			select {
			case <-stopCh:
				return
			case <-msgChan[i]:
			default:
				continue
			}
		}
	}
}
