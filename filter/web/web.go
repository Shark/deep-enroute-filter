package web

import (
  "net/http"
  "html/template"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"
)

type webContentType struct {
  ProcessedMessages []types.ProcessedMessage
}

var webContent webContentType

func handler(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("web/index.html")
  t.Execute(w, webContent)
}

func ListenAndServe(processedMessages <-chan types.ProcessedMessage) {
  webContent = webContentType{
    ProcessedMessages: make([]types.ProcessedMessage, 0),
  }

  go func() {
    for message := range processedMessages {
      webContent.ProcessedMessages = append(webContent.ProcessedMessages, message)
    }
  }()

  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
