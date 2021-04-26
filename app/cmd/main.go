package main

import (
	"flag"
	"time"
	"syscall"
	"os"
	"os/signal"

	"Asura/conf"
	log "Asura/src/logger"
	"Asura/app/http"
	"Asura/app/service"
)

var (
	s *service.Service
)

func main() {
	flag.Parse()
	conf.Init()
	if conf.Conf.Debug {
		log.Init(&log.Config{
			Dir: conf.Conf.Log.Dir,
		})
		defer log.Close()
	}
	s = service.New(conf.Load())
	log.Info("dts-cache-center [version: s%] start", conf.Load().Version)
	// rpc.Init(&conf.Conf, s)
	http.Init(conf.Load(), s)
	signalHandler()
	return
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("get a signal %s, stop the dts-cache-center process", si.String())
			s.Close()
			s.Wait()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
