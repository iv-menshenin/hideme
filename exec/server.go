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

func Serve(config ServeConfig) error {
	var srv = http.Server{
		Addr:              fmt.Sprintf(":%d", config.GetPort()),
		Handler:           http.HandlerFunc(handler),
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
	return srv.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s: %s (%s)", r.Method, r.URL.Path, r.URL.RawQuery))
	switch r.URL.Path {

	case "/":
		handlerRoot(w, r)
		return

	case "/inject":
		handlerInject(w, r)
		return

	case "/extract":
		handlerExtract(w, r)
		return

	case "/generate":
		handlerGenerate(w, r)
		return
	}

	http.NotFoundHandler().ServeHTTP(w, r)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {

}

func handlerInject(w http.ResponseWriter, r *http.Request) {

}

func handlerExtract(w http.ResponseWriter, r *http.Request) {

}

func handlerGenerate(w http.ResponseWriter, r *http.Request) {

}
