package session

import (
	"strconv"
	"strings"
)

func (s *Session) BooleanToString(element bool) string{
	if element == true{
		return "true"
	} else {
		return "false"
	}
}

func (s *Session) StringToBoolean(element string) bool{
	if strings.TrimSpace(element) == "true"{
		return true
	} else{
		return false
	}
}

func (s *Session) IntegerToString(element int) string{
	value := strconv.Itoa(element)
	return value
}

func (s *Session) StringToInteger(element string) int{
	converted, _ := strconv.Atoi(element)
	return converted
}
