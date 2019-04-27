package filters

import (
	"github.com/graniet/operative-framework/filters/phone_to_instagram"
	"github.com/graniet/operative-framework/filters/say_hello"
	"github.com/graniet/operative-framework/session"
)

func LoadFilters(s *session.Session){
	s.Filters = append(s.Filters, say_hello.PushSayHelloFilter(s))
	s.Filters = append(s.Filters, phone_to_instagram.PushPhoneToInstagramFilter(s))
}
