package http

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/iv-menshenin/hideme/config"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(htmlPage))
}

const inputDataMaxSize = 1024 * 1024 * 16

func handlerInject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(inputDataMaxSize); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg, err := config.NewInjectorFromQuery(queryArgs{
		v: r.URL.Query(),
		m: r.MultipartForm,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = cfg.Execute(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handlerExtract(w http.ResponseWriter, r *http.Request) {

}

func handlerGenerate(w http.ResponseWriter, r *http.Request) {

}

type queryArgs struct {
	v url.Values
	m *multipart.Form
}

func (a queryArgs) StringVal(key string) string {
	return a.v.Get(key)
}

func (a queryArgs) ByteVal(key string) (r io.ReadCloser, name string, err error) {
	f := a.m.File[key]
	if len(f) == 0 {
		return
	}
	name = f[0].Filename
	r, err = f[0].Open()
	return
}
