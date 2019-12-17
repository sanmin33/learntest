package main

import (
	"flag"
	"fmt"
	"time"
)

var fmsgType = flag.String("msg", "hello", "消息内容")
var ftimeSecond = flag.Int("time", 1, "每次发送协程的测试时间长度秒")
var fgoSendNum = flag.Int("gosendnum", 1, "同时发送消息的go协程数量,执行的时候是从1个发送协程开始一直运行到最大数量")
var fgoGetNum = flag.Int("gogetnum", 1, "同时接收消息的go协程数量")
var fchanNum = flag.Int("channum", 1, "传递消息的管道数量，默认为1，如果给不是1的数，则管道数量等于发送消息go协程数量")

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

	for i := 1; i <= goSendNum; i++ {
		creatSendMsg(i)
	}
	time.Sleep(time.Second * 1) //没这句主程序瞬间退出，sendMsg没有时间完成打印输出
}

func creatSendMsg(n int) {

	if chanNum == 1 {
		counts := make([]int, n)
		stopCh := make(chan int)
		msgChan := make(chan interface{}, 100)

		for i := 0; i < n; i++ {
			go sendMsg(msgChan, msg, stopCh, &counts[i]) //创建一个到n个发送协程，本函数执行完成后才交回控制权给main函数
		}
		for i := 0; i < goGetNum; i++ {
			go getMsg(msgChan)
		}
		time.Sleep(time.Second * time.Duration(timeSecond))
		close(stopCh) //发消息给发送协程可以停止工作了
		totalMsgNum := 0
		for i := 0; i < n; i++ {
			totalMsgNum = totalMsgNum + counts[i]
		}
		fmt.Println("发送协程数:", n, " 接收协程数：", goGetNum, "消息管道数:", 1, "测试时长:", timeSecond, "发送消息总数量：", totalMsgNum)
	} else {
		counts := make([]int, n)
		stopCh := make(chan int)
		msgChan := make([]chan interface{}, n, 100)
		for i := 0; i < n; i++ {
			msgChan[i] = make(chan interface{}, 100)
		}
		//对于多对多的情况，这里简化为一对一对情况，即一个生产者对一个管道，对应一个消费者。但同时有多个这种组合在运行
		for i := 0; i < n; i++ {
			go sendMsg(msgChan[i], msg, stopCh, &counts[i]) //创建一个到n个发送协程，本函数执行完成后才交回控制权给main函数
			go getMsg(msgChan[i])
		}
		time.Sleep(time.Second * time.Duration(timeSecond))

		close(stopCh)
		totalMsgNum := 0
		for i := 0; i < n; i++ {
			totalMsgNum = totalMsgNum + counts[i]
		}
		fmt.Println("发送协程数:", n, " 接收协程数：", n, "消息管道数:", n, "测试时长:", timeSecond, "发送消息总数据量：", totalMsgNum)

	}

}

func sendMsg(msgChan chan interface{}, msg string, stopflag chan int, msgNum *int) {
	*msgNum = 0
	for {
		msgChan <- msg
		select {
		case <-stopflag:
			close(msgChan)
			return
		default:
			*msgNum++
		}
	}
}
func getMsg(msg chan interface{}) {
	for _ = range msg {
		//	time.Sleep(time.Microsecond * 1)
	}
}
