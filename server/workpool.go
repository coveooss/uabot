package server

import (
	"github.com/goinggo/workpool"
	"github.com/satori/go.uuid"
)

const routine = "bot"

type WorkPool struct {
	goinggoWorkpool         *workpool.WorkPool
	workerInfo              []map[string]interface{}
	NumberConcurrentRoutine int
	QueueLength             int32
}

func NewWorkPool(numConcurrentRoutine int, queueLength int32) *WorkPool {
	return &WorkPool{
		goinggoWorkpool:         workpool.New(numConcurrentRoutine, queueLength),
		workerInfo:              make([]map[string]interface{}, numConcurrentRoutine),
		NumberConcurrentRoutine: numConcurrentRoutine,
		QueueLength:             queueLength,
	}
}

func (workPool *WorkPool) PostWork(worker *WorkWrapper) error {
	return workPool.goinggoWorkpool.PostWork(routine, worker)
}

type WorkWrapper struct {
	realWorker *BotWorker
	workPool   *WorkPool
	workerID   uuid.UUID
}

func (_workWrapper *WorkWrapper) DoWork(goRoutine int) {
	info := _workWrapper.realWorker.bot.GetInfo()
	info["workerId"] = _workWrapper.realWorker.id.String()
	_workWrapper.workPool.workerInfo[goRoutine] = info
	_workWrapper.realWorker.DoWork(goRoutine)
	_workWrapper.workPool.workerInfo[goRoutine] = nil
}

func (workPool *WorkPool) getInfo() []map[string]interface{} {
	filteredInfo := make([]map[string]interface{}, 0)
	for _, info := range workPool.workerInfo {
		if info != nil {
			filteredInfo = append(filteredInfo, info)
		}
	}
	return filteredInfo
}

func (workPool *WorkPool) ActiveRoutines() int32 {
	return workPool.goinggoWorkpool.ActiveRoutines()
}

func (workPool *WorkPool) QueuedWork() int32 {
	return workPool.goinggoWorkpool.QueuedWork()
}
