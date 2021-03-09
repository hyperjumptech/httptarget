package server

import (
	"context"
	"fmt"
	"github.com/hyperjumptech/httptarget/model"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	endPoints = &model.EndPoints{}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func Start(host string, port int, initEndpoint *model.EndPoint) error {
	err := endPoints.Add(initEndpoint)
	if err != nil {
		return err
	}
	listen := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr: listen,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:      5 * time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       5 * time.Second,
		IdleTimeout:       2 * time.Second,
		Handler:           &HttpTargetHandler{},
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		//log.Infof("This RETTER instance will forwards GET request...")
		//log.Infof("  From : %s/*", l)
		//log.Infof("  To   : %s/*", Config.GetString(BackendURL))
		//log.Infof("URL Query Detect       : %s", Config.GetString(CacheDetectQuery))
		//log.Infof("URL Session Detect     : %s", Config.GetString(CacheDetectSession))
		//log.Infof("RETTER is listening on : [%s]", l)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	os.Exit(0)
	return nil
}

type HttpTargetHandler struct {
}

func (h *HttpTargetHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ep := endPoints.GetByPath(req.URL.Path)
	if ep != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Not Found"))
	}

	randDelay := ep.DelayMinMs + rand.Intn(ep.DelayMaxMs-ep.DelayMinMs)
	time.Sleep(time.Duration(randDelay) * time.Millisecond)

	if ep.ReturnHeaders != nil {
		for k, v := range ep.ReturnHeaders {
			res.Header()[k] = v
		}
	}
	res.WriteHeader(ep.ReturnCode)
	if len(ep.ReturnBody) > 0 {
		res.Write([]byte(ep.ReturnBody))
	}
}
