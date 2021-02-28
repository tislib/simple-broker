package app

import (
	"backend-processor/model"
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type App struct {
	Addr          string
	CertFile      string
	KeyFile       string
	channel       map[string]chan *model.Record
	makeChanMutex *sync.Mutex
	bufferSize    int
	brokerId      string
}

func (app *App) Run() {
	app.init()
	srv := &http.Server{Addr: app.Addr, Handler: app}
	log.Printf("Serving on " + app.Addr)
	log.Fatal(srv.ListenAndServeTLS(app.CertFile, app.KeyFile))
}

func (app *App) init() {
	log.SetLevel(log.DebugLevel)

	app.channel = make(map[string]chan *model.Record)
	app.bufferSize = 10000
	app.brokerId = uuid.NewString()

	app.makeChanMutex = new(sync.Mutex)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		app.writeMessage(r, w)
		break
	case "GET":
		app.readMessage(r, w)
	}
}

func (app *App) writeMessage(r *http.Request, w http.ResponseWriter) {
	channelName := app.getChannelNameFromRequest(r)

	channel := app.getChannel(channelName)
	recordData := new(model.RecordData)
	record := new(model.Record)
	record.Message = recordData
	record.BrokerId = app.brokerId
	record.MessageId = uuid.NewString()

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		app.handleError(err, w)
		return
	}

	err = json.Unmarshal(bytes, recordData)

	if err != nil {
		app.handleError(err, w)
		return
	}

	channel <- record

	content, err := json.Marshal(record)
	content = append(content, 10)

	if err != nil {
		app.handleError(err, w)
		return
	}

	_, err = w.Write(content)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	if err != nil {
		app.handleError(err, w)
		return
	}
}

func (app *App) handleError(err error, w http.ResponseWriter) {
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Print(err)
	}
}

func (app *App) readMessage(r *http.Request, w http.ResponseWriter) {
	channelName := app.getChannelNameFromRequest(r)
	consumerId := uuid.NewString()

	var record *model.Record

	for {
		record = app.listenMessage(channelName, consumerId)

		if record == nil {
			log.Debugf("no message received for consumer: %s", consumerId)
			continue
		}

		record.ConsumerId = consumerId
		content, err := json.Marshal(record)
		content = append(content, 10)

		if err != nil {
			app.handleError(err, w)
			return
		}

		_, err = w.Write(content)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		if err != nil {
			app.handleError(err, w)
			return
		}
	}
}

func (app *App) listenMessage(name string, consumerId string) *model.Record {
	var cases []reflect.SelectCase

	for channelName := range app.channel {
		if channelName == name || strings.HasPrefix(channelName, name) {
			ch := app.channel[channelName]

			log.Debugf("starting to listen channel: %s / %s", channelName, consumerId)

			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}
	}

	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(time.After(2 * time.Second)),
	})

	_, res, _ := reflect.Select(cases)

	val := res.Interface()

	if reflect.TypeOf(val) == reflect.TypeOf(time.Time{}) {
		return nil
	} else {
		return val.(*model.Record)
	}
}

func (app *App) getChannelNameFromRequest(r *http.Request) string {
	return r.URL.Path
}

func (app *App) getChannel(name string) chan *model.Record {
	if app.channel[name] != nil {
		return app.channel[name]
	}

	app.initChannel(name)

	return app.channel[name]
}

func (app *App) initChannel(name string) {
	defer app.makeChanMutex.Unlock()
	app.makeChanMutex.Lock()

	app.channel[name] = make(chan *model.Record, app.bufferSize)
}
