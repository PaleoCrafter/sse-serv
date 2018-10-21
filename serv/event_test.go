// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"testing"
)

func TestEventFactory_CreateEvent(t *testing.T) {
	samples := []struct{ given, expect string }{
		{given: "", expect: "\n"},
		{given: "           ", expect: "\n"},
		{given: " f\noo \n\n", expect: "f\noo\n"},
		{given: " :data     ", expect: ":data\n"},
		{given: "X\r\nX\n\r ", expect: "X\nX\n"},
	}

	c := NewEventCreator()

	for _, sample := range samples {
		t.Run("", func(t *testing.T) {
			result := c.CreateEvent([]byte(sample.given))

			if result.Data() != sample.expect {
				t.Errorf("expected %s, got %s", sample.expect, result.Data())
			}
		})
	}
}

func TestEvent_Id(t *testing.T) {
	if NewEventCreator().CreateEvent(nil).Id() != "" {
		t.Errorf("expected empty id")
	}
}

func TestEvent_Event(t *testing.T) {
	if NewEventCreator().CreateEvent(nil).Event() != "" {
		t.Errorf("expected empty event type")
	}
}

func TestEvent_Retry(t *testing.T) {
	if NewEventCreator().CreateEvent(nil).Retry() != 0 {
		t.Errorf("expected retry to be 0")
	}
}

func TestEvent_String(t *testing.T) {
	samples := []struct{ given, expect string }{
		{given: "", expect: "data: \n"},
		{given: "x", expect: "data: x\n"},
		{given: "event: y\ndata: x", expect: "data: event: y\ndata: data: x\n"},
	}

	c := NewEventCreator()

	for _, sample := range samples {
		t.Run("", func(t *testing.T) {
			result := c.CreateEvent([]byte(sample.given))

			if result.String() != sample.expect {
				t.Errorf("expected %s, got %s", sample.expect, result.Data())
			}
		})
	}
}
