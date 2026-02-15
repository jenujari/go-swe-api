package main

import (
	config "github.com/jenujari/go-swe-api/config"
	"github.com/jenujari/go-swe-api/grpc"
	rtc "github.com/jenujari/runtime-context"
)

var masterProcess *rtc.ProcessContext

func init() {
	rtc.InitProcessContext(config.GetLogger())
}

func main() {
	masterProcess = rtc.GetMainProcess()

	masterProcess.Run(grpc.RunGRPCServer)

	masterProcess.WaitForFinish()
}
