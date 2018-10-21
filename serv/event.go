// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import "strings"

// EventCreator ...
type EventCreator interface {
	CreateEvent(data []byte) Event
}

type creator struct{}

// NewEventCreator ...
func NewEventCreator() EventCreator {
	return &creator{}
}

func (e *creator) CreateEvent(data []byte) Event {
	return newEvent("", "", normalize(data), 0)
}

// Event ...
type Event interface {
	Id() string
	Event() string
	Data() string
	Retry() int
	String() string
}

type event struct {
	id    string
	event string
	data  string
	retry int
}

func newEvent(id, kind, data string, retry int) Event {
	return &event{id: id, event: kind, data: data, retry: retry}
}

func (e *event) Id() string {
	return e.id
}

func (e *event) Event() string {
	return e.event
}

func (e *event) Data() string {
	return e.data
}

func (e *event) Retry() int {
	return e.retry
}

func (e *event) String() string {
	s := strings.Trim(e.data, "\n")
	return "data: " + strings.Replace(s, "\n", "\ndata: ", -1) + "\n"
}

func normalize(s []byte) string {
	return strings.Trim(dataClean.Replace(string(s)), " \n") + "\n"
}

var dataClean = strings.NewReplacer("\r\n", "\n", "\r", "")
