package web

import (
  "fmt"
  "net/http"
  "html/template"
  "time"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
  "github.com/gobuffalo/packr"
)

type webProcessedMessage struct {
  Timestamp time.Time
  ShortHash string
  Method string
  UriPath string
  RuleProcessingResults []webRuleProcessingResult
}

func (m *webProcessedMessage) RelativeTimestamp() string {
  seconds := int(time.Since(m.Timestamp).Seconds())
  return fmt.Sprintf("%ds ago", seconds)
}

type webContentType struct {
  ProcessedMessages []webProcessedMessage
}

type webRuleProcessingResult struct {
  Allowed bool
  RuleName string
  RuleMessage string
}

var webContent webContentType

func handler(w http.ResponseWriter, r *http.Request) {
  box := packr.NewBox("templates")
  t := template.Must(template.New("index").Parse(box.String("index.html.template")))
  t.Execute(w, webContent)
}

func ListenAndServe(processedMessages <-chan types.ProcessedMessage) {
  webContent = webContentType{
    ProcessedMessages: make([]webProcessedMessage, 0),
  }

  go func() {
    for message := range processedMessages {
      webProcessingResults := make([]webRuleProcessingResult, 0)
      for _, processingResult := range message.RuleProcessingResults {
        webProcessingResult := webRuleProcessingResult{
          Allowed: processingResult.Allowed,
          RuleName: processingResult.Rule.Name(),
        }
        if(processingResult.RuleMessage != nil) {
          webProcessingResult.RuleMessage = *processingResult.RuleMessage
        } else {
          webProcessingResult.RuleMessage = ""
        }
        webProcessingResults = append(webProcessingResults, webProcessingResult)
      }
      webMessage := webProcessedMessage{
        time.Now(),
        message.Message.Metadata.Hash()[0:6],
        canopus.MethodString(message.Message.Message.GetCode()),
        message.Message.Message.GetURIPath(),
        webProcessingResults,
      }

      numProcessedMessages := len(webContent.ProcessedMessages)
      if(numProcessedMessages > 10) {
        webContent.ProcessedMessages = webContent.ProcessedMessages[1:numProcessedMessages]
      }
      webContent.ProcessedMessages = append(webContent.ProcessedMessages, webMessage)
    }
  }()

  assetsBox := packr.NewBox("./assets")
  http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assetsBox)))

  http.HandleFunc("/", handler)

  http.ListenAndServe(":8080", nil)
}
