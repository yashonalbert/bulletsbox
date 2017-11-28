package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type MQ struct {
	queue []string
}

func (mq *MQ) inQueue(msg string) {
	mq.queue = append(mq.queue, msg)
	return
}

func (mq *MQ) outQueue() string {
	if len(mq.queue) == 0 {
		return ""
	}
	msg := mq.queue[0]
	mq.queue = mq.queue[1:]
	return msg
}

func main() {
	var mq MQ
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		mq.inQueue(string(b))
		log.Print(mq.outQueue())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
