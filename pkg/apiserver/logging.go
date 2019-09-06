package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var loggingCenter *NotificationManager

func init() {
	loggingCenter = NewNotificationManager()

	go func() {
		for {
			b := []byte(time.Now().Format(time.RFC3339))
			if err := loggingCenter.NotifySubscribers(b); err != nil {
				log.Fatal(err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

}

type UnsubscribeFunc func() error

type Subscriber interface {
	Subscribe(c chan []byte) (UnsubscribeFunc, error)
}

func handleSSE(s Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Subscribe
		c := make(chan []byte)
		unsubscribeFn, err := s.Subscribe(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Signal SSE Support
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

	Looping:
		for {
			select {
			case <-r.Context().Done():
				if err := unsubscribeFn(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				break Looping

			default:
				// Find the deployment ID
				id := mux.Vars(r)["id"]
				_ = path.Base(r.URL.Path)
				//r.URL.Path
				b := <-c
				// parse the data from the channel
				// if the correct id then send them the data

				fmt.Fprintf(w, "data: %s %s\n\n", id, b)

				w.(http.Flusher).Flush()
			}
		}
	}
}

type Notifier interface {
	Notify(b []byte) error
}

type NotificationManager struct {
	subscribers   map[chan []byte]struct{}
	subscribersMu *sync.Mutex
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		subscribers:   map[chan []byte]struct{}{},
		subscribersMu: &sync.Mutex{},
	}
}

func (nc *NotificationManager) Subscribe(c chan []byte) (UnsubscribeFunc, error) {
	nc.subscribersMu.Lock()
	nc.subscribers[c] = struct{}{}
	nc.subscribersMu.Unlock()

	unsubscribeFn := func() error {
		nc.subscribersMu.Lock()
		delete(nc.subscribers, c)
		nc.subscribersMu.Unlock()

		return nil
	}

	return unsubscribeFn, nil
}

func (nc *NotificationManager) NotifySubscribers(b []byte) error {
	nc.subscribersMu.Lock()
	defer nc.subscribersMu.Unlock()

	for c := range nc.subscribers {
		select {
		case c <- b:
		default:
		}
	}

	return nil
}
