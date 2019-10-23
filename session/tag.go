package session

import (
	"errors"
	"github.com/segmentio/ksuid"
	"strings"
)

type Tags struct {
	TagId string `json:"-"`
	Text  string `json:"text"`
}

func (target *Target) GetTags() []Tags {
	return target.Tags
}

func (target *Target) HasTag(tag string) bool {
	for _, element := range target.Tags {
		if element.Text == strings.TrimSpace(tag) {
			return true
		}
	}
	return false
}

func (s *Session) AddTag(target *Target, tag string) (bool, error) {
	if target.HasTag(strings.TrimSpace(tag)) {
		return false, errors.New("Tag already exist for target '" + target.GetName() + "'")
	}
	newTag := Tags{
		TagId: "TAG_" + ksuid.New().String(),
		Text:  strings.TrimSpace(tag),
	}

	target.Tags = append(target.Tags, newTag)

	for _, otherTarget := range s.Targets {
		if otherTarget.HasTag(strings.TrimSpace(tag)) {
			if otherTarget.GetId() != target.GetId() {
				target.Link(Linking{
					TargetId:       otherTarget.TargetId,
					TargetResultId: newTag.TagId,
				})
			}
		}
	}
	return true, nil
}
