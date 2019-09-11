package apiserver

import (
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/gorilla/mux"
)

// MVP of a streaming logging provider

// The notificationCenter is in charge of handling the various notification managers, whihc
// in turn will notify all of their subscribers
var notificationCenter map[string]*notificationManager

// Notification is what will be sent to subscribers of a manager
type Notification struct {
	ID      string
	RawData []byte
}

// RegisterNotificationManager will create a manager and an endpoint
func RegisterNotificationManager(managerName, endpoint string) error {
	// Register the new Manager to the Notification Center
	notificationCenter[managerName] = newNotificationManager()
	AddDynamicEndpoint(endpoint,
		endpoint,
		fmt.Sprintf("Automatically generated notification endpoint for [%s]", managerName),
		managerName,
		http.MethodGet,
		handleSubscribers(notificationCenter[managerName]))
	return nil
}

// NotifyManager - This will Notify a Manager that there is a new notification that needs to go to subscribers
func NotifyManager(managerName string, n Notification) error {
	manager := notificationCenter[managerName]
	if manager == nil {
		return fmt.Errorf("Notification Manager [%s], hasn't been registered", managerName)
	}
	manager.notifySubscribers(n)
	return nil
}

//   --------------  Notication MAGIC below --------------

func init() {
	// Initialise the notificationCenter map

	notificationCenter = make(map[string]*notificationManager)

}

type unsubscribeFunc func() error

type subscriber interface {
	subscribe(n chan Notification) (unsubscribeFunc, error)
}

func handleSubscribers(s subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Subscribe
		n := make(chan Notification)
		unsubscribeFn, err := s.subscribe(n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set environment for streaming events
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
				newNotification := <-n
				if newNotification.ID == id {
					// parse the data from the channel
					// if the correct id then send them the data

					fmt.Fprintf(w, "%s\n", newNotification.RawData)
				}

				w.(http.Flusher).Flush()
			}
		}
	}
}

type notifier interface {
	Notify(n Notification) error
}

type notificationManager struct {
	subscribers   map[chan Notification]struct{}
	subscribersMu *sync.Mutex
}

func newNotificationManager() *notificationManager {
	return &notificationManager{
		subscribers:   map[chan Notification]struct{}{},
		subscribersMu: &sync.Mutex{},
	}
}

func (nc *notificationManager) subscribe(n chan Notification) (unsubscribeFunc, error) {
	nc.subscribersMu.Lock()
	nc.subscribers[n] = struct{}{}
	nc.subscribersMu.Unlock()

	unsubscribeFn := func() error {
		nc.subscribersMu.Lock()
		delete(nc.subscribers, n)
		nc.subscribersMu.Unlock()

		return nil
	}

	return unsubscribeFn, nil
}

func (nc *notificationManager) notifySubscribers(n Notification) error {
	// Lock them until updates are complete
	nc.subscribersMu.Lock()
	defer nc.subscribersMu.Unlock()

	for c := range nc.subscribers {
		select {
		case c <- n:
		default:
		}
	}

	return nil
}
