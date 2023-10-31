// package tailLog
//
// import (
//
//	"fmt"
//	"github.com/hpcloud/tail"
//	"log/slog"
//	"sync"
//	"time"
//
// )
//
// const (
//
//	StatusNormal = 1
//	StatusDelete = 2
//
// )
//
// var (
//
//	tailObjMgr *TailObjMgr
//
// )
//
//	type CollectConf struct {
//		LogPath string `json:"log_path"`
//		Topic   string `json:"topic"`
//	}
//
//	type TailObj struct {
//		tail     *tail.Tail
//		conf     CollectConf
//		status   int
//		exitChan chan int
//	}
//
//	type TextMsg struct {
//		Msg   string
//		Topic string
//	}
//
// // TailObjMgr 管理系统所有tail对象
//
//	type TailObjMgr struct {
//		tailsObjs []*TailObj
//		msgChan   chan *TextMsg
//		lock      sync.Locker
//	}
//
//	func GetOneLine() (msg *TextMsg) {
//		msg = <-tailObjMgr.msgChan
//		return
//	}
//
//	func UpdateConfig(configs []CollectConf) (err error) {
//		for _, oneConf := range configs {
//			// 对于已经运行的所有实例, 路径是否一样
//			isRunning := false
//			for _, obj := range tailObjMgr.tailsObjs {
//				// 路径一样则证明是同一实例
//				if oneConf.LogPath == obj.conf.LogPath {
//					isRunning = true
//					obj.status = StatusNormal
//					break
//				}
//			}
//
//			// 如果不存在该配置项 新建一个tail task任务
//			if !isRunning {
//				createNewTask(oneConf)
//			}
//
//		}
//		// 遍历所有查看是否存在删除操作
//		var tailObjs []*TailObj
//		for _, obj := range tailObjMgr.tailsObjs {
//			obj.status = StatusDelete
//			for _, oneConf := range configs {
//				if oneConf.LogPath == obj.conf.LogPath {
//					obj.status = StatusNormal
//					break
//				}
//			}
//			// 如果status为删除, 则将exitChan置为1
//			if obj.status == StatusDelete {
//				obj.exitChan <- 1
//			}
//			// 将obj存入临时的数组中
//			tailObjs = append(tailObjs, obj)
//		}
//		// 将临时数组传入tailsObjs中
//		tailObjMgr.tailsObjs = tailObjs
//		return
//	}
//
//	func createNewTask(conf CollectConf) {
//		// 初始化TailFile实例
//		tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
//			ReOpen:    true,
//			Follow:    true,
//			Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
//			MustExist: false,
//			Poll:      true,
//		})
//
//		if errTail != nil {
//			slog.Error(fmt.Sprintf("收集文件[%s]错误: %v", conf.LogPath, errTail))
//			return
//		}
//		// 导入配置项
//		obj := &TailObj{
//			conf:     conf,
//			exitChan: make(chan int, 1),
//		}
//
//		obj.tail = tails
//		tailObjMgr.tailsObjs = append(tailObjMgr.tailsObjs, obj)
//
//		go readFromTail(obj)
//	}
//
// // InitTail 初始化tail
// func InitTail(conf []CollectConf, chanSize int) (err error) {
//
//		tailObjMgr = &TailObjMgr{
//			msgChan: make(chan *TextMsg, chanSize), // 定义Chan管道
//		}
//
//		// 加载配置项
//		if len(conf) == 0 {
//			slog.Error(fmt.Sprintf("无效的日志collect配置: %v\n", conf))
//		}
//
//		// 循环导入
//		for _, v := range conf {
//			createNewTask(v)
//		}
//
//		return
//	}
//
// // 读入日志数据
//
//	func readFromTail(tailObj *TailObj) {
//		for true {
//			select {
//
//			case msg, ok := <-tailObj.tail.Lines:
//				if !ok {
//					slog.Warn(fmt.Sprintf("Tail file close reopen, filename:%s\n", tailObj.tail.Filename))
//					time.Sleep(100 * time.Millisecond)
//					continue
//				}
//				textMsg := &TextMsg{
//					Msg:   msg.Text,
//					Topic: tailObj.conf.Topic,
//				}
//				// 放入chan里
//				tailObjMgr.msgChan <- textMsg
//
//			// 如果exitChan为1, 则删除对应配置项
//			case <-tailObj.exitChan:
//				slog.Warn(fmt.Sprintf("tail obj 退出, 配置项为conf:%v", tailObj.conf))
//				return
//			}
//		}
//	}
package tailLog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"shyCollector/logAgent/kafka"
	"time"
)

var (
	tailObj *tail.Tail
	LogChan chan string
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	tailObj.Init()
	return tailObj
}

func (t *TailTask) Init() {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true}
	tailObj, err := tail.TailFile(t.path, config)
	if err != nil {
		panic(fmt.Sprintf("tailTask init failed: %v\n", err))
	}
	t.instance = tailObj

	go t.run()
}

func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines:
			kafka.SendToKafka(t.topic, line.Text)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
