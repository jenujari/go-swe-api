// Package router have all the routes that the server handle
package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jenujari/go-swe-api/config"
	rtc "github.com/jenujari/runtime-context"
)

var (
	server *http.Server
	router *http.ServeMux
)

func init() {
	cfg := config.GetConfig()

	server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	router = http.NewServeMux()

	SetRoutes(router)

	server.Handler = router
	server.SetKeepAlivesEnabled(true)
	config.GetLogger().Println("server initialization complete.")
}

func RunServer() {
	pc := rtc.GetMainProcess()

	go func(cmdx *rtc.ProcessContext) {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			cmdx.FatalErrorChan <- fmt.Errorf("ListenAndServe(): %v", err)
		}
	}(pc)

	<-pc.CTX.Done()
	config.GetLogger().Println("shutting down server...")
	if err := server.Shutdown(pc.CTX); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	config.GetLogger().Println("server shutdown complete...")
}

func GetServer() *http.Server {
	return server
}
