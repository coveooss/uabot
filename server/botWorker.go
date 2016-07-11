package server

import (
	"github.com/adambbolduc/uabot/autobot"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/satori/go.uuid"
	"math/rand"
)

type BotWorker struct {
	Worker
	bot *autobot.Autobot
	id  uuid.UUID
}

type Worker interface {
	DoWork(goRoutine int)
}

func (worker BotWorker) DoWork(goRoutine int) {
	quitChannel := make(chan bool)
	quitChannels[worker.id] = quitChannel
	scenariolib.Info.Printf("Bot starting on worker: %v\n", goRoutine)
	err := worker.bot.Run(quitChannel)
	if err != nil {
		scenariolib.Error.Println(err)
	}
}

func NewWorker(config *explorerlib.Config, random *rand.Rand, id uuid.UUID) Worker {
	return Worker(WorkWrapper{
		realWorker: &BotWorker{
			bot: autobot.NewAutobot(config, random),
			id:  id,
		},
		workPool: workPool,
	})
}
