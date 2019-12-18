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
var ifQuit bool = false

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
    for c := 1; c <= chanNum; c++ {
        for s := c; s <= goSendNum; s++ {
            for g := c; g <= goGetNum; g++ {
                creatSendMsg(s, c, g)
            }
        }
    }
    //  time.Sleep(time.Second * 1) //没这句主程序瞬间退出，sendMsg没有时间完成打印输出
}

func creatSendMsg(goSendNum int, chanNum int, goGetNum int) {
    ifQuit = false
    counts := make([]int, goSendNum)
    stopCh := make(chan int)
    msgChan := make([]chan interface{}, chanNum, 100)
    for i := 0; i < chanNum; i++ {
        msgChan[i] = make(chan interface{}, 100)
    }
    go sendMsg(goSendNum, msgChan, msg, stopCh, counts)

    go getMsg(goGetNum, msgChan, stopCh)
    time.Sleep(time.Second * time.Duration(timeSecond))

    maxgos := runtime.NumGoroutine()
    //  fmt.Println("最大协程数量：", maxgos)
    close(stopCh)
    ifQuit = true
	
    //等待发送协程全部退出
    for maxgos-runtime.NumGoroutine() < goSendNum {
        time.Sleep(time.Millisecond * 1) //如果生产者消费者还没有执行完，就等一下。
    }

    //发送协程全部退出后，关闭信道，让接收协程也退出
    for i := 0; i < len(msgChan); i++ {
        close(msgChan[i])
    }

    //等退接收协程全部退出
    for runtime.NumGoroutine() > 1 {
        time.Sleep(time.Millisecond)
    }

    //全部发送协程和接收协程都退出后，开始统计传输的消息数量
    //  fmt.Println("当前协程数量：", runtime.NumGoroutine())
    totalMsgNum := 0
    for i := 0; i < goSendNum; i++ {
        totalMsgNum = totalMsgNum + counts[i]
    }
    fmt.Println("发送协程数:,", goSendNum, ",消息管道数:,", chanNum, ", 接收协程数：,", goGetNum, ",测试时长:,", timeSecond, ",发送消息总数据量：,", totalMsgNum)
}

func sendMsg(s int, msgChan []chan interface{}, msg string, stop chan int, msgNums []int) {
    for i := 0; i < s; i++ {
        msgNums[i] = 0
    }
    c := 0 //信道数量是小于等于发送协程的，如果信道已经用光了，那么要从第0个信道开始重新使用
    for i := 0; i < s; i++ {
        go func(ch chan interface{}, num *int, stop chan int) {
            for {
                /*  select {
                    case <-stop:
                        return
                    default:
                        ch <- msg
                        *num++
                    }
                */
                if !ifQuit {
                    ch <- msg
                    *num++
                } else {
                    return
                }
            }
        }(msgChan[c], &msgNums[i], stop)
        c++
        //信道数量是小于等于发送协程的，如果信道已经用光了，那么要从第0个信道开始重新使用
        if c == len(msgChan) {
            c = 0
        }
    }
}
func getMsg(g int, msgChan []chan interface{}, stopCh chan int) {
    c := 0
    for i := 0; i < g; i++ {
        go func(ch chan interface{}) {
            for _ = range ch {
            }
        }(msgChan[c])
        c++
        if c == len(msgChan) {
            c = 0
        }
    }
}
