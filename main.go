package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

type WebHookJson struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Ref    string `json:"ref"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Info : リクエストを受け付けました")

	//Validate request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error : POST以外受け付けないのでエラー")
		return
	}

	// header
	method := r.Method
	fmt.Println("[method] " + method)
	for k, v := range r.Header {
		log.Print("[header] " + k + ": " + strings.Join(v, ","))
	}

	// body
	buf, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("bodyErr ", bodyErr.Error())
		http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Printf("BODY: %q", rdr1)
	r.Body = rdr2

	var webHookJson WebHookJson

	// paramaterを取得 & パース
	// JSONがURLエンコードされているのでデコード
	var payload = r.FormValue("payload")
	payload, _ = url.QueryUnescape(payload)

	log.Print("Decode body: " + payload)

	// JSONデコード
	if err := json.Unmarshal([]byte(payload), &webHookJson); err != nil {
		log.Println("Error : JSONデコード時にエラー")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flag.Parse()
	// 申し訳程度に実行時引数のチェック
	size := len(flag.Args())
	if size != 2 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error : 実行時引数のエラー(シェルのパス、pullしたいgitリポジトリのパスであること)")
		return
	}

	// シェル実行
	out, err := exec.Command("sh", flag.Arg(0), flag.Arg(1), webHookJson.After, webHookJson.Ref).Output()
	if err != nil {
		log.Println("シェル実行時にエラー")
		log.Println(err.Error())
		return
	}
	fmt.Println(string(out))

	w.WriteHeader(http.StatusOK)
	log.Println("Info : 処理完了")
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Listen Server ....")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
