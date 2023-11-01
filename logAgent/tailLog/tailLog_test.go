package tailLog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"testing"
	"time"
)

func TestTailLog(t *testing.T) {
	logFilePath := "C:\\Users\\19406\\Desktop\\go\\shyCollector\\logAgent\\logtest1.log"
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true}
	tails, err := tail.TailFile(logFilePath, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}

	for {
		msg, ok := <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
