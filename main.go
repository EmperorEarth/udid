package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	http.HandleFunc("/upload", uploadUDID)
	http.HandleFunc("/", rootHandler)
	m_ := &autocert.Manager{
		Cache:  autocert.DirCache("./certificates"),
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"subdomain.domain.tld",
		),
	}
	go http.ListenAndServe(":http", m_.HTTPHandler(nil))
	s_ := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: m_.GetCertificate,
		},
	}
	s_.ListenAndServeTLS("", "")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		downloadUDID(w, r)
		return
	}
	w.Write([]byte(r.URL.Path[1:]))
}

func downloadUDID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	w.Header().Set("Content-Disposition", "attachment; filename=\"udid.mobileconfig\"")
	file_, err := os.Open("./udid.mobileconfig")
	if err != nil {
		log.Printf("failed to send .mobileconfig: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, file_)
}

func uploadUDID(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read uploaded UDID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/%s",
			string(data)[strings.Index(string(data), "<string>")+8:strings.Index(string(data), "</string>")],
		),
		http.StatusMovedPermanently,
	)
}
