package http

import (
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

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
	defer cfg.Clear()

	if err = cfg.Execute(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if files := cfg.Files(); len(files) > 0 {
		fn := path.Base(files[0])
		if len(files) > 1 {
			fn += ".zip"
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+fn+"\"")
		writeFiles(w, files...)
	}
}

func handlerExtract(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(inputDataMaxSize); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg, err := config.NewExtractorFromQuery(queryArgs{
		v: r.URL.Query(),
		m: r.MultipartForm,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer cfg.Clear()

	if err = cfg.Execute(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if files := cfg.Files(); len(files) > 0 {
		fn := path.Base(files[0])
		if len(files) > 1 {
			fn += ".zip"
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+fn+"\"")
		writeFiles(w, files...)
	}
}

func handlerGenerate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(inputDataMaxSize); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg, err := config.NewGeneratorFromQuery(queryArgs{
		v: r.URL.Query(),
		m: r.MultipartForm,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer cfg.Clear()

	if err = cfg.Execute(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if files := cfg.Files(); len(files) > 0 {
		fn := path.Base(files[0])
		if len(files) > 1 {
			fn += ".zip"
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+fn+"\"")
		writeFiles(w, files...)
	}
}

func writeFiles(w io.Writer, files ...string) {
	if len(files) == 1 {
		file, err := os.Open(files[0])
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		if _, err = io.Copy(w, file); err != nil {
			log.Println(err)
		}
		return
	}

	var tmpFileZip = makeTmpFileName()
	f, err := os.Create(tmpFileZip)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		f.Close()
		os.Remove(tmpFileZip)
	}()

	z := zip.NewWriter(f)
	for _, fileName := range files {
		fw, err := z.Create(path.Base(fileName))
		if err != nil {
			log.Println(err)
			return
		}
		data, err := os.ReadFile(fileName)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = fw.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if err = z.Close(); err != nil {
		log.Println(err)
		return
	}
	writeFiles(w, tmpFileZip)
}

func makeTmpFileName() string {
	var r [16]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/hideme/%s.zip", hex.EncodeToString(r[:]))
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
