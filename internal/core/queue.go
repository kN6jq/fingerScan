package core

import (
	"container/list"
	"sync"
)

// Queue 线程安全的队列实现
type Queue struct {
	mutex sync.Mutex
	data  *list.List
}

// NewQueue 创建新队列
func NewQueue() *Queue {
	return &Queue{
		data: list.New(),
	}
}

// Push 添加元素到队首
func (q *Queue) Push(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.data.PushFront(v)
}

// Pop 从队尾取出元素
func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if elem := q.data.Back(); elem != nil {
		return q.data.Remove(elem)
	}
	return nil
}

// Len 返回队列长度
func (q *Queue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.data.Len()
}

// IsEmpty 检查队列是否为空
func (q *Queue) IsEmpty() bool {
	return q.Len() == 0
}

// Clear 清空队列
func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.data.Init()
}

// Peek 查看队尾元素但不移除
func (q *Queue) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if elem := q.data.Back(); elem != nil {
		return elem.Value
	}
	return nil
}
