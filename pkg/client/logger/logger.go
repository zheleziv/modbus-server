package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	Lock           sync.Mutex
	ParentNodeName string
	ParentNodeIp   string
	IsLogOutput    bool
}

const (
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
)

func (thisLogger *Logger) WriteWithTag(level string, nodestate string, tag string, msg string) {

	text := fmt.Sprintf("%s >> Node: name - %s, IP - %s, state - %s. Tag: name - %s. Message [%s]: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		thisLogger.ParentNodeName,
		thisLogger.ParentNodeIp,
		nodestate,
		tag,
		level,
		msg)

	fmt.Print(text)
	if thisLogger.IsLogOutput {
		thisLogger.Lock.Lock()
		defer thisLogger.Lock.Unlock()
		save(text)
	}
}

func (thisLogger *Logger) Write(level string, nodestate string, msg string) {
	text := fmt.Sprintf("%s >> Node: name - %s, IP - %s, state - %s. Message [%s]: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		thisLogger.ParentNodeName,
		thisLogger.ParentNodeIp,
		nodestate,
		level,
		msg)

	fmt.Print(text)
	if thisLogger.IsLogOutput {
		thisLogger.Lock.Lock()
		defer thisLogger.Lock.Unlock()
		save(text)
	}
}

func save(text string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		fmt.Println(err)
	}
}
