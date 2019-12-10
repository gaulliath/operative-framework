package export

import (
	"github.com/graniet/operative-framework/session"
	"strings"
)

type Result []map[string]string

type Export struct {
	ModuleName string `json:"module_name"`
	Results    Result `json:"results"`
	Count      int    `json:"count"`
}

func ExportNow(s *session.Session, module session.Module) Export {

	export := Export{
		ModuleName: module.Name(),
	}

	for _, target := range s.Targets {
		for moduleName, results := range target.Results {
			if moduleName == module.Name() {
				export.Count = len(results)
				if len(results) > 0 {
					resultExport := Result{}
					for _, result := range results {
						resultParsing := make(map[string]string)
						header := strings.Split(result.Header, target.GetSeparator())
						value := strings.Split(result.Value, target.GetSeparator())

						for key, name := range header {
							resultParsing[name] = value[key]
						}
						resultExport = append(resultExport, resultParsing)
					}
					export.Results = resultExport
				}
			}
		}
	}
	return export
}
