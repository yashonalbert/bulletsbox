package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Item struct
type Item struct {
	score   int
	content []byte
}

// MQ struct
type MQ struct {
	queues map[string][]*Item
}

func (mq *MQ) enQueue(name string, score int, content []byte) {
	var item Item
	item.score = score
	item.content = content
	if _, ok := mq.queues[name]; ok {
		for i := len(mq.queues[name]) - 1; i >= 0; i-- {
			if item.score == mq.queues[name][i].score {
				mq.queues[name] = append(append(mq.queues[name][:i], &item), mq.queues[name][i:]...)
				return
			}
			if item.score < mq.queues[name][i].score {
				mq.queues[name] = append(append(mq.queues[name][:i-1], &item), mq.queues[name][i-1:]...)
				return
			}
		}
		mq.queues[name] = append([]*Item{&item}, mq.queues[name]...)
	} else {
		mq.queues[name] = []*Item{&item}
		return
	}
}

func (mq *MQ) deQueue(name string) []byte {
	if len(mq.queues[name]) == 0 {
		return nil
	}
	content := mq.queues[name][0].content
	mq.queues[name] = mq.queues[name][1:]
	return content
}

func main() {
	var mq MQ
	mq.queues = make(map[string][]*Item)
	http.HandleFunc("/queues", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if name, ok := query["name"]; ok {
			switch r.Method {
			case "GET":
				fmt.Fprintln(w, mq.deQueue(name[0]))
			case "POST":
				b, err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				if score, ok := query["score"]; ok {
					scoreInt, err := strconv.ParseInt(score[0], 10, 32)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					mq.enQueue(name[0], int(scoreInt), b)
				} else {
					mq.enQueue(name[0], 1024, b)
				}
			default:
				http.Error(w, "not found", 404)
			}
		} else {
			http.Error(w, "query error", 400)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
