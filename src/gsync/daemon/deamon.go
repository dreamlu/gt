package daemon

import (
	"github.com/dreamlu/gt/conf"
	"github.com/dreamlu/gt/src/cons"
	"github.com/dreamlu/gt/src/type/time"
	"sync"
	time2 "time"
)

// Daemon progress
type Daemon struct {
	Task  chan *Task // time task
	tasks []*Task    // task queue
	Num   int        // task goroutine pool num
}

type Task struct {
	//ID         uint64      // job id, 0++
	daemonFunc DaemonFunc  // exec func
	Time       *time.CTime // ctime
}

type TimeTask struct {
	Task
	Time time.CTime // ctime
}

type DaemonFunc func()
type Param func(*Task)

func Func(daemonfunc DaemonFunc) Param {
	return func(params *Task) {
		params.daemonFunc = daemonfunc
	}
}

func Time(time time.CTime) Param {
	return func(params *Task) {
		params.Time = &time
	}
}

// ======= Singleton ========
// single daemon
var (
	daemon     *Daemon
	onceDaemon sync.Once
)

func Daemoner() *Daemon {

	onceDaemon.Do(func() {
		daemon = newDaemon()
		go daemon.task()
		go daemon.taskQueue()
	})
	return daemon
}

// new daemon
func newDaemon() *Daemon {
	return &Daemon{
		Num:  conf.Get[int](cons.ConfTaskNum),
		Task: make(chan *Task), // must init via make()
	}
}

// AddTask add daemon task
func (d *Daemon) AddTask(params ...Param) *Daemon {
	task := &Task{}

	for _, p := range params {
		p(task)
	}

	// d.Task <- task
	d.tasks = append(d.tasks, task)
	return d
}

// running task queue
func (d *Daemon) taskQueue() {

	for {
		for k := 0; k < len(d.tasks); k++ {
			task := d.tasks[k]
			if task.Time != nil {
				if time2.Time(*task.Time).After(time2.Now()) {
					continue
				} else {
					d.Task <- task
					d.tasks = append(d.tasks[:k], d.tasks[k+1:]...)
					k--
					continue
				}
			}
			d.Task <- task
		}
	}
}

// running task
// daemon goroutine pool num
func (d *Daemon) task() {

	for i := 0; i < d.Num; i++ {
		//i := i
		go func() {
			for {
				task := <-d.Task
				task.daemonFunc()
			}
		}()
	}
}

func AddTask(params ...Param) *Daemon {
	return Daemoner().AddTask(params...)
}
