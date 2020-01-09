package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found: "+r.URL.Path)
		return
	}

	encodings, ok := r.URL.Query()["encoding"]
	if !ok || len(encodings[0]) < 1 {
		logrus.Info("Url Param 'encoding' is missing")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Url Param 'encoding' is missing")
		return
	}

	urls, ok := r.URL.Query()["url"]
	if !ok || len(urls[0]) < 1 {
		logrus.Info("Url Param 'url' is missing")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Url Param 'url' is missing")
		return
	}

	encoding := string(encodings[0])
	url := string(urls[0])
	resp, err := http.Get(url)
	if err != nil {
		logrus.Info("cannot download the contents from URL")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "cannot download the contents from URL")
		return
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("cannot read the contents from URL")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "cannot read the contents from URL")
		return
	}

	if encoding == "sjis" {
		str, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), string(html))
		if err != nil {
			logrus.Info("cannot transform the contents from sjis to UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "cannot transform the contents from sjis to UTF-8")
			return
		}
		fmt.Fprintf(w, str)
		return
	} else if encoding == "euc" {
		str, _, err := transform.String(japanese.EUCJP.NewDecoder(), string(html))
		if err != nil {
			logrus.Info("cannot transform the contents from EUC-JP to UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "cannot transform the contents from EUC-JP to UTF-8")
			return
		}
		fmt.Fprintf(w, str)
		return
	}

}

func main() {
	port := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.ListenAndServe(port, nil)
}
