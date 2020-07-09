package daemon

import "fmt"

// Daemon pool
type Pool struct {
	Queue  chan func() error
	Number int
	Total  int

	result         chan error
	finishCallback func()
}

// 初始化
func (g *Pool) Init(number int, total int) {
	g.Queue = make(chan func() error, total)
	g.Number = number
	g.Total = total
	g.result = make(chan error, total)
}

// 开门接客
func (g *Pool) Start() {
	// 开启Number个Daemon
	for i := 0; i < g.Number; i++ {
		go func() {
			for {
				task, ok := <-g.Queue
				if !ok {
					break
				}

				err := task()
				g.result <- err
			}
		}()
	}

	// 获得每个work的执行结果
	for j := 0; j < g.Total; j++ {
		res, ok := <-g.result
		if !ok {
			break
		}

		if res != nil {
			fmt.Println(res)
		}
	}

	// 所有任务都执行完成，回调函数
	if g.finishCallback != nil {
		g.finishCallback()
	}
}

// 关门送客
func (g *Pool) Stop() {
	close(g.Queue)
	close(g.result)
}

// 添加任务
func (g *Pool) AddTask(task func() error) {
	g.Queue <- task
}

// 设置结束回调
func (g *Pool) SetFinishCallback(callback func()) {
	g.finishCallback = callback
}
