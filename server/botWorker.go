package server

import (
	"github.com/erocheleau/uabot/autobot"
	"github.com/erocheleau/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/satori/go.uuid"
	"math/rand"
)

type BotWorker struct {
	Worker
	bot     *autobot.Autobot
	id      uuid.UUID
	channel chan bool
}

type Worker interface {
	DoWork(goRoutine int)
}

func (worker BotWorker) DoWork(goRoutine int) {
	scenariolib.Info.Printf("Bot starting on worker: %v\n", goRoutine)
	err := worker.bot.Run(worker.channel)
	if err != nil {
		scenariolib.Error.Println(err)
	}
}

func NewWorker(config *explorerlib.Config, quitChannel chan bool, random *rand.Rand, id uuid.UUID) Worker {
	return Worker(WorkWrapper{
		realWorker: &BotWorker{
			bot:     autobot.NewAutobot(config, random),
			id:      id,
			channel: quitChannel,
		},
		workPool: workPool,
	})
}
