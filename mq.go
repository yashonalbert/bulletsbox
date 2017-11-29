package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MQ struct
type MQ struct {
	queues map[string][][]byte
}

func (mq *MQ) enQueue(name string, content []byte) {
	if _, ok := mq.queues[name]; ok {
		mq.queues[name] = append(mq.queues[name], content)
	} else {
		mq.queues[name] = [][]byte{content}
	}
	return
}

func (mq *MQ) deQueue(name string) []byte {
	if len(mq.queues[name]) == 0 {
		return nil
	}
	content := mq.queues[name][0]
	mq.queues[name] = mq.queues[name][1:]
	return content
}

func main() {
	var mq MQ
	mq.queues = make(map[string][][]byte)
	http.HandleFunc("/queues", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if name, ok := query["name"]; ok {
			switch r.Method {
			case "GET":
				log.Print(mq.queues)
				fmt.Fprintln(w, mq.deQueue(name[0]))
			case "POST":
				b, err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				mq.enQueue(name[0], b)
			default:
				http.Error(w, "not found", 404)
			}
		} else {
			http.Error(w, "query error", 400)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
