// Package queue 提供线程安全的队列实现
// 用于管理扫描任务，支持并发的入队和出队操作
package queue

import (
	"container/list"
	"fmt"
	"sync"
)

// Queue 线程安全的队列结构体
// 基于双向链表实现，使用互斥锁保证并发安全
type Queue struct {
	l    sync.Mutex // 互斥锁，保证并发安全
	data *list.List // 底层双向链表
}

// NewQueue 创建新的队列实例
//
// 返回：
//   - 初始化完成的队列指针
func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

// Push 将元素添加到队列头部
// 线程安全操作
//
// 参数：
//   - v: 要添加的元素（任意类型）
//
// 返回：
//   - 新添加元素的链表节点
func (q *Queue) Push(v interface{}) *list.Element {
	q.l.Lock()
	defer q.l.Unlock()
	return q.data.PushFront(v)
}

// PushBack 将元素添加到队列尾部
// 线程安全操作
//
// 参数：
//   - v: 要添加的元素
//
// 返回：
//   - 新添加元素的链表节点
func (q *Queue) PushBack(v interface{}) *list.Element {
	q.l.Lock()
	defer q.l.Unlock()
	return q.data.PushBack(v)
}

// Pop 从队列尾部取出并移除一个元素
// 实现 FIFO（先进先出）行为
// 线程安全操作
//
// 返回：
//   - 取出的元素，队列为空时返回 nil
func (q *Queue) Pop() interface{} {
	q.l.Lock()
	defer q.l.Unlock()

	// 获取尾部元素
	iter := q.data.Back()
	if nil == iter {
		return nil
	}

	// 移除并返回元素值
	v := iter.Value
	q.data.Remove(iter)
	return v
}

// Pops 批量从队列取出多个元素
// 线程安全操作
//
// 参数：
//   - num: 要取出的元素数量
//
// 返回：
//   - vals: 取出的元素切片
//   - actualLen: 实际取出的元素数量（可能小于请求数量）
func (q *Queue) Pops(num int) ([]interface{}, int) {
	vals := make([]interface{}, num)
	i := 0

	q.l.Lock()
	defer q.l.Unlock()

	for {
		// 已取够数量
		if i >= num {
			break
		}

		// 获取尾部元素
		iter := q.data.Back()
		if iter == nil {
			// 队列已空，返回已取出的元素
			return vals, i
		}

		// 移除元素并保存
		q.data.Remove(iter)
		vals[i] = iter.Value
		i++
	}

	// 如果实际取出数量小于请求数量，截断切片
	if i < num {
		return vals[0:i], i
	}
	return vals, i
}

// Remove 移除指定的链表节点
// 线程安全操作
//
// 参数：
//   - v: 要移除的链表节点
//
// 返回：
//   - 被移除节点的值
func (q *Queue) Remove(v *list.Element) interface{} {
	q.l.Lock()
	defer q.l.Unlock()
	return q.data.Remove(v)
}

// Len 获取队列当前长度
// 注意：此方法不加锁，返回值可能在多线程环境下不准确
//
// 返回：
//   - 队列中的元素数量
func (q *Queue) Len() int {
	return q.data.Len()
}

// Dump 打印队列中的所有元素
// 用于调试目的
func (q *Queue) Dump() {
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item:", iter.Value)
	}
}
