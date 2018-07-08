package core

import (
  "regexp"
)

func parseDefinition(definition string) []string {
  r := regexp.MustCompile(`(?:<(?P<path>[^>]+)>)`)
  rSubexpNames := r.SubexpNames()

  matches := r.FindAllStringSubmatch(definition, -1)

  if matches == nil {
    return []string{}
  }

  var paths []string

  for _, match := range matches {
    for i, value := range match {
      if(rSubexpNames[i] == "path") {
        paths = append(paths, value)
      }
    }
  }

  return paths
}
