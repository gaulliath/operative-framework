package session

import "github.com/segmentio/ksuid"

type Note struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func (target *Target) AddNote(text string) {
	target.Notes = append(target.Notes, Note{
		Id:   "N_" + ksuid.New().String(),
		Text: text,
	})
	return
}

func (result *OpfResults) AddNoteToResult(text string) {
	result.Notes = append(result.Notes, Note{
		Id:   "N_" + ksuid.New().String(),
		Text: text,
	})
	return
}
