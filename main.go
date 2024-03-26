package main

import (
	"github.com/ipoluianov/bdata/app"
	"github.com/ipoluianov/bdata/application"

	bitcoinclient "github.com/ipoluianov/bdata/bitcoin_client"
	"github.com/ipoluianov/bdata/logger"
)

func main() {
	//bybit.Start()

	bitcoinclient.ParseRawFile()

	//bitcoinclient.CheckConnect()
	return

	application.Name = "bdata"
	application.ServiceName = "bdata"
	application.ServiceDisplayName = "bdata"
	application.ServiceDescription = "bdata"
	application.ServiceRunFunc = app.RunAsService
	application.ServiceStopFunc = app.StopService

	logger.Init(logger.CurrentExePath() + "/logs")

	if !application.TryService() {
		app.RunDesktop()
	}

}
