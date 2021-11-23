package main


import (
	"errors"
	"fmt"
)

type Queue interface {
	// 将 msg 添加到 队尾
	Enqueue(msg string)
	// 返回队首的i, 如果队空，则返回error
	Dequeue() (msg string, err error)
}

type queue struct {
	slice []string
}

func (q *queue) Enqueue(msg string) {
	q.slice = append(q.slice, msg)
}

func (q *queue) Dequeue() (msg string, err error) {
	if len(q.slice) == 0 {
		return "",errors.New("no int in queue")
	}
	// 获取队首
	msg = q.slice[0]
	// 删除队首
	q.slice = q.slice[1:]
	return
}

func (q *queue) String() string {
	return fmt.Sprintf("%#v", q.slice)
}
