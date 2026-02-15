package main

import (
	config "github.com/jenujari/go-swe-api/config"
	router "github.com/jenujari/go-swe-api/router"
	rtc "github.com/jenujari/runtime-context"
)

var masterProcess *rtc.ProcessContext

func init() {
	rtc.InitProcessContext(config.GetLogger())
}

func main() {
	masterProcess = rtc.GetMainProcess()

	srv := router.GetServer()
	masterProcess.Run(router.RunServer)
	config.GetLogger().Println("Server is running at", srv.Addr)

	masterProcess.WaitForFinish()
}
