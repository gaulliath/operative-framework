package engine

import "github.com/graniet/operative-framework/session"

func CommandBase(line string, s *session.Session) bool{
	if line == "info"{
		s.ViewInformation()
		return true
	}
	return false
}
