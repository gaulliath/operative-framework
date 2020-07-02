package export

import (
	"github.com/graniet/operative-framework/session"
)

type Result []map[string]string

type Export struct {
	ModuleName string `json:"module_name"`
	Results    Result `json:"results"`
	Count      int    `json:"count"`
}

func JSON(s *session.Session) *session.Instance {
	return s.CurrentInstance
}
