package server

import (
	"github.com/adambbolduc/uabot/autobot"
	"github.com/erocheleau/uabot/scenariolib"
	"github.com/satori/go.uuid"
)

type BotWorker struct {
	bot *autobot.Autobot
	id uuid.UUID
}

func (worker *BotWorker) DoWork(goRoutine int) {
	quitChannel := make(chan bool)
	quitChannels[worker.id] = quitChannel
	scenariolib.Info.Printf("Bot starting on worker: %v\n", goRoutine)
	err := worker.bot.Run(quitChannel)
	if err != nil {
		scenariolib.Error.Println(err)
	}
}

