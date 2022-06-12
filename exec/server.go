package exec

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServeConfig interface {
	GetPort() int
}

func Serve(config ServeConfig, handler http.HandlerFunc) error {
	var srv = http.Server{
		Addr:              fmt.Sprintf(":%d", config.GetPort()),
		Handler:           handler,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 1,
		WriteTimeout:      time.Second * 15,
	}
	go func() {
		var ch = make(chan os.Signal)
		signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-ch
		if err := srv.Close(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("listeining on %d port\n", config.GetPort())
	return srv.ListenAndServe()
}
