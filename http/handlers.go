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
		handleRequestWithConfiguration(config.NewInjectorFromQuery)(w, r)
		return

	case "/extract":
		handleRequestWithConfiguration(config.NewExtractorFromQuery)(w, r)
		return

	case "/generate":
		handleRequestWithConfiguration(config.NewGeneratorFromQuery)(w, r)
		return
	}

	http.NotFoundHandler().ServeHTTP(w, r)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(htmlPage))
}

const inputDataMaxSize = 1024 * 1024 * 16

type configCreator func(config.Query) (*config.Config, error)

func handleRequestWithConfiguration(createConfig configCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(inputDataMaxSize); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			writeErrorAsPayload(w, err)
			return
		}
		cfg, err := createConfig(queryArgs{m: r.MultipartForm})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			writeErrorAsPayload(w, err)
			return
		}
		defer cfg.Clear()

		if err = cfg.Execute(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			writeErrorAsPayload(w, err)
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
}

func writeFiles(w io.Writer, files ...string) {
	if len(files) == 1 {
		writeFile(w, files[0])
		return
	}

	tmpFileZip, err := zipFiles(files)
	if err != nil {
		log.Println(err)
		return
	}
	writeFile(w, tmpFileZip)
	if e := os.Remove(tmpFileZip); e != nil {
		log.Println(e)
	}
}

func zipFiles(files []string) (string, error) {
	var tmpFileZip = makeTmpFileName()

	f, err := os.Create(tmpFileZip)
	if err != nil {
		return "", err
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Println(e)
		}
	}()

	z := zip.NewWriter(f)
	for _, fileName := range files {
		fw, err := z.Create(path.Base(fileName))
		if err != nil {
			return "", err
		}
		data, err := os.ReadFile(fileName)
		if err != nil {
			return "", err
		}
		_, err = fw.Write(data)
		if err != nil {
			return "", err
		}
	}
	if err = z.Close(); err != nil {
		return "", err
	}
	return tmpFileZip, nil
}

func writeFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Println(e)
		}
	}()
	if _, err = io.Copy(w, f); err != nil {
		log.Println(err)
	}
	return
}

func makeTmpFileName() string {
	var r [16]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/hideme/%s.zip", hex.EncodeToString(r[:]))
}

type queryArgs struct {
	m *multipart.Form
}

func (a queryArgs) StringVal(key string) string {
	val := a.m.Value[key]
	if len(val) == 0 {
		return ""
	}
	return val[0]
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

func writeErrorAsPayload(w io.Writer, err error) {
	w.Write([]byte("something went wrong: "))
	w.Write([]byte(fmt.Sprintf("<%T> %s", err, err)))
}
