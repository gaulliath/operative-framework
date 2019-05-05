package session

import (
	"github.com/graniet/go-pretty/table"
	"github.com/fatih/color"
	"fmt"
	"time"
)

type Stream struct{
	Sess *Session
	Verbose bool
	JSON bool
	History map[string]string
}

func (stream *Stream) Standard(text string){
	if stream.Verbose {
		t := time.Now()
		timeText := t.Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] : %s\n", timeText, text)
		stream.Sess.Information.AddEvent()
	}
}

func (stream *Stream) Error(text string){
	if stream.Verbose {
		t := time.Now()
		timeText := t.Format("2006-01-02 15:04:05")
		color.Red("[%s] : %s\n", timeText, text)
		stream.Sess.Information.AddEvent()
	}
}

func (stream *Stream) Success(text string) {
	if stream.Verbose {
		t := time.Now()
		timeText := t.Format("2006-01-02 15:04:05")
		color.Green("[%s] : %s\n", timeText, text)
		stream.Sess.Information.AddEvent()
	}
}

func (stream *Stream) Warning(text string) {
	if stream.Verbose {
		t := time.Now()
		timeText := t.Format("2006-01-02 15:04:05")
		color.Yellow("[%s] : %s\n", timeText, text)
	}
}

func (stream *Stream) WithoutDate(text string){
	if stream.Verbose {
		fmt.Printf("%s\n", text)
		stream.Sess.Information.AddEvent()
	}
}

func (stream *Stream) GenerateTable() table.Writer{
	return table.NewWriter()
}

func (stream *Stream) Render(t table.Writer){
	if stream.Verbose {
		t.Render()
	}
	return
}