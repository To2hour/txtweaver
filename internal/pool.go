package internal

import (
	"sync"
)

type Pool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func NewPool(maxCount int) *Pool {
	p := &Pool{
		tasks: make(chan func(), 10000),
	}

	// 提前启动固定数量的协程（Worker）
	for i := 0; i < maxCount; i++ {
		go p.worker()
	}
	return p
}

// 工作协程
func (p *Pool) worker() {
	for task := range p.tasks {
		task()
		p.wg.Done()
	}
}

func (p *Pool) Execute(task func()) {
	p.wg.Add(1) // 放入任务前，计数器 +1
	p.tasks <- task
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	close(p.tasks)
}
