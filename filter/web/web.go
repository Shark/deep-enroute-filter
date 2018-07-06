package web

import (
  "net/http"
  "html/template"

  "gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/types"

  "github.com/zubairhamed/canopus"
)

type webProcessedMessage struct {
  ShortHash string
  Method string
  UriPath string
  RuleProcessingResults []webRuleProcessingResult
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
  t, _ := template.ParseFiles("web/index.html")
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
        message.Message.Metadata.Hash()[0:6],
        canopus.MethodString(message.Message.Message.GetCode()),
        message.Message.Message.GetURIPath(),
        webProcessingResults,
      }
      webContent.ProcessedMessages = append(webContent.ProcessedMessages, webMessage)
    }
  }()

  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
