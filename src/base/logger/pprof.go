// -------------------------------------------
// @file      : pprof.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午1:40
// -------------------------------------------

package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// Interval 间隔保活时间
const Interval = 2 * time.Minute

type PprofTask struct {
	sync.Once
	LastUpdate  time.Time
	Timer       *time.Timer
	triggerBack func()
}

var pprofTask *PprofTask

func init() {
	pprofTask = &PprofTask{}
	pprofTask.LastUpdate = time.Now()
}

func (task *PprofTask) loop() {
	defer func() {
		if err := recover(); err != nil {
			Errorw(fmt.Sprintf("pprof task panic. err:%s", err), zap.Stack("stack"))
		}
		if task.Timer != nil {
			task.Timer.Stop()
		}
	}()
	Debugf("PprofTask loop start")
	for {
		select {
		case <-task.Timer.C:
			if time.Since(task.LastUpdate) >= Interval {
				Debugf("PprofTask dump start")
				// 如果在间隔时间内没有调用keepalive就会进入这里打印一次pprof文件并退出pprof task
				task.DumpGoroutineInfo()
				task.DumpCPUInfo()
				task.DumpHeapInfo()
				if task.triggerBack != nil {
					task.triggerBack()
				}
				return
			}
			task.Timer.Reset(Interval)
		}
	}
}

func (task *PprofTask) GetFilename(category string) string {
	filename := global.fileName
	filename = strings.Replace(filename, ".log", "", 1)
	return fmt.Sprintf("%s-%s-%s.prof", filename, category, time.Now().UTC().Format(time.RFC3339))
}

func (task *PprofTask) DumpGoroutineInfo() {
	p := pprof.Lookup("goroutine")
	filename := task.GetFilename("goroutine")
	file, err := os.Create(filename)
	if err != nil {
		Errorf("DumpGoroutineInfo create file(%s) err:%s", filename, err)
		return
	}
	Debugf("DumpGoroutineInfo create file(%s)", filename)
	err = p.WriteTo(file, 0)
	if err != nil {
		Errorf("DumpGoroutineInfo write to %s err:%s", filename, err)
	}
	Debugf("DumpGoroutineInfo dumped to file(%s)", filename)
}

func (task *PprofTask) DumpHeapInfo() {
	p := pprof.Lookup("heap")
	filename := task.GetFilename("heap")
	file, err := os.Create(filename)
	if err != nil {
		Errorf("DumpHeapInfo create file(%s) err:%s", filename, err)
		return
	}
	Debugf("DumpHeapInfo create file(%s)", filename)
	err = p.WriteTo(file, 0)
	if err != nil {
		Errorf("DumpHeapInfo write to %s err:%s", filename, err)
	}
	Debugf("DumpHeapInfo dumped to file(%s)", filename)
}

func (task *PprofTask) DumpCPUInfo() {
	go func() {
		filename := task.GetFilename("cpu")
		file, err := os.Create(filename)
		if err != nil {
			Errorf("DumpCPUInfo create file(%s) err:%s", filename, err)
			return
		}
		Debugf("DumpCPUInfo create file(%s)", filename)
		duration := 1 * time.Minute
		tmpTimer := time.NewTimer(duration)
		defer func() {
			pprof.StopCPUProfile()
			err := file.Close()
			if err != nil {
				Errorf("DumpCPUInfo close file(%s) err:%s", filename, err)
			}
			tmpTimer.Stop()
		}()
		if err := pprof.StartCPUProfile(file); err != nil {
			Errorf("DumpCPUInfo StopCPUProfile err:%s", err)
		}
		for {
			select {
			case <-tmpTimer.C:
				Debugf("DumpCPUInfo dumped to file(%s)", filename)
				return
			}
		}
	}()
}

// StartPprofTask 启动pprof task
func StartPprofTask() {
	pprofTask.Do(func() {
		pprofTask.Timer = time.NewTimer(Interval)
		go pprofTask.loop()
	})
}

// PprofKeepAlive 保活pprof task timer，当timer到期时会 打印pprof信息并退出task
func PprofKeepAlive(callback func()) {
	if pprofTask == nil || pprofTask.Timer == nil {
		return
	}
	pprofTask.LastUpdate = time.Now()
	pprofTask.triggerBack = callback
}
