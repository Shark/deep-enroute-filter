package web

import (
  "encoding/json"
  "fmt"
  "net/http"
  "html/template"
  "sync"
  "time"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
  "github.com/gobuffalo/packr"
  "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var globalQuit chan struct{}
var listeners threadSafeSlice
var processedMessages []webProcessedMessage

// START: https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel
type listener struct {
  source chan types.Event
  quit chan struct{}
}

type threadSafeSlice struct {
    sync.Mutex
    listeners []*listener
}

func (slice *threadSafeSlice) Push(l *listener) {
    slice.Lock()
    defer slice.Unlock()

    slice.listeners = append(slice.listeners, l)
}

func (slice *threadSafeSlice) Iter(routine func(*listener)) {
    slice.Lock()
    defer slice.Unlock()

    for _, listener := range slice.listeners {
        routine(listener)
    }
}
// END: https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type webProcessedMessage struct {
  Timestamp             time.Time
  Destination           string
  Method                string
  Path                  string
  RuleProcessingResults []webRuleProcessingResult
}

type webRuleProcessingResult struct {
  Allowed     bool
  RuleName    string
  RuleMessage *string
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
  listener := &listener{source: make(chan types.Event, 10), quit: globalQuit}
  listeners.Push(listener)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrade failed: %v", err)
		return
	}
	defer c.Close()

  for {
    select {
    case event := <-listener.source:
      var jsonEvent interface{}
      if event.Type() == "ProcessedMessageEvent" {
        jsonEvent = struct {
          Type    string
          Payload interface{}
        }{
          "ProcessedMessages",
          processedMessages,
        }
      } else {
        jsonEvent = struct {
          Type    string
          Payload interface{}
        }{
          event.Type(),
          event.Payload(),
        }
      }

      if json, err := json.Marshal(jsonEvent); err == nil {
        err = c.WriteMessage(websocket.TextMessage, json)
        if err != nil {
          fmt.Printf("Write failed: %v", err)
          break;
        }
      } else {
        fmt.Printf("JSON marshalling failed: %v", err)
      }
    case <-listener.quit:
      return
    }
  }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  box := packr.NewBox("templates")
  t := template.Must(template.New("index").Parse(box.String("index.html.template")))
  t.Execute(w, nil)
}

func ListenAndServe(events <-chan types.Event) {
  go func() {
    for event := range events {
      if event.Type() == "ProcessedMessageEvent" {
        if processedMessage, ok := event.Payload().(types.ProcessedMessage); ok {
          processingResults := make([]webRuleProcessingResult, len(processedMessage.RuleProcessingResults))
          for i, result := range processedMessage.RuleProcessingResults {
            processingResults[i] = webRuleProcessingResult{
              result.Allowed,
              result.Rule.Name(),
              result.RuleMessage,
            }
          }
          processedMessages = append([]webProcessedMessage{webProcessedMessage{
            time.Now(),
            processedMessage.Message.Metadata.DstIP,
            canopus.MethodString(processedMessage.Message.Message.GetCode()),
            processedMessage.Message.Metadata.UriPath,
            processingResults,
          }}, processedMessages...)
        }

        numProcessedMessages := len(processedMessages)
        if numProcessedMessages > 10 {
          processedMessages = processedMessages[0:numProcessedMessages-1]
        }
      }

      listeners.Iter(func(l *listener) {
        select {
        case l.source <- event:
        default:
          fmt.Println("dropped")
        }
      })
    }
  }()

  assetsBox := packr.NewBox("./assets")
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assetsBox)))

  http.HandleFunc("/ws", websocketHandler)
  http.HandleFunc("/",   indexHandler)

  http.ListenAndServe(":8080", nil)
}
