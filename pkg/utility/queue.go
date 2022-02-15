package utility

type Queue struct {
	slots []interface{}
}

//将元素插入到队列尾部。注意，这里的参数不能设置为*Item类型
//因为，空的interface的动态类型是在赋值的时候隐式确定的。所以，这里
//必须是值类型，从而才能确定Item的动态类型。
func (q *Queue) Push(i interface{}) {
	q.slots = append(q.slots, i)
}

//首元素出队
func (q *Queue) Pop() *interface{} {
	res := &q.slots[0]
	q.slots = q.slots[1:]
	return res
}

//判断队列是否为空
func (q *Queue) Empty() bool {
	n := len(q.slots)
	if n > 0 {
		return false
	}

	return true
}

//返回队列的首元素
func (q *Queue) Front() *interface{} {
	return &q.slots[0]
}

//返回队列的尾元素
func (q *Queue) Back() *interface{} {
	return &q.slots[len(q.slots)-1]
}

//返回队列中元素个数
func (q *Queue) Size() int {
	return len(q.slots)
}

//CreateQueue 创建一个空的队列
func CreateQueue() *Queue {
	emptyQueue := new(Queue)
	return emptyQueue
}
