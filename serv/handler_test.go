// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"sse-serv/amqp"
	"sse-serv/mock/amqp"
	"sse-serv/mock/logg"
)

// Create instance
func TestNewResponseHandler(t *testing.T) {
	r := NewResponseHandler(
		nil, http.Request{}, nil, nil, nil, nil, nil, nil,
	)
	if r == nil {
		t.Error("expected constructor to work")
	}
}

// Send 405 when method other than GET
func TestResponseHandler_Handle_1(t *testing.T) {

	unHandled := []string{
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range unHandled {
		recorder := httptest.NewRecorder()

		r := &responseHandler{
			writer:  recorder,
			request: http.Request{Method: method},
		}
		r.Handle()

		if recorder.Code != 405 {
			t.Errorf("expected 405, got %d", recorder.Code)
		}
	}
}

// Send 503 when queue name not determinable
func TestResponseHandler_Handle_2(t *testing.T) {
	recorder := httptest.NewRecorder()

	r := &responseHandler{
		writer:  recorder,
		request: http.Request{Method: "GET"},
		pattern: NewPattern("${cookie:not-here}"),
	}
	r.Handle()

	if recorder.Code != 503 {
		t.Errorf("expected 503, got %d", recorder.Code)
	}
}

// Send 503 when can not connect broker
func TestResponseHandler_Handle_3(t *testing.T) {
	recorder := httptest.NewRecorder()

	r := &responseHandler{
		writer:   recorder,
		request:  http.Request{Method: "GET"},
		pattern:  NewPattern(""),
		provider: amqp.NewProvider(amqp.NewBrokerPool(nil, nil)),
	}
	r.Handle()

	if recorder.Code != 503 {
		t.Errorf("expected 503, got %d", recorder.Code)
	}
}

// Send 503 when can not create consumer
func TestResponseHandler_Handle_4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	providerMock := mock_amqp.NewMockProvider(ctrl)
	providerMock.EXPECT().Consumer(gomock.Any(), gomock.Any()).Return(nil, errors.New(""))

	recorder := httptest.NewRecorder()

	r := &responseHandler{
		writer:   recorder,
		request:  http.Request{Method: "GET"},
		pattern:  NewPattern(""),
		provider: providerMock,
	}
	r.Handle()

	if recorder.Code != 503 {
		t.Errorf("expected 503, got %d", recorder.Code)
	}
}

// Send 503 when consumer can not consume
func TestResponseHandler_Handle_5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumerMock := mock_amqp.NewMockConsumer(ctrl)
	consumerMock.EXPECT().Consume("").Return(nil, errors.New(""))
	consumerMock.EXPECT().Close().Times(1)

	providerMock := mock_amqp.NewMockProvider(ctrl)
	providerMock.EXPECT().Consumer(gomock.Any(), gomock.Any()).Return(consumerMock, nil)

	recorder := httptest.NewRecorder()

	r := &responseHandler{
		writer:   recorder,
		request:  http.Request{Method: "GET"},
		pattern:  NewPattern(""),
		provider: providerMock,
	}
	r.Handle()

	if recorder.Code != 503 {
		t.Errorf("expected 503, got %d", recorder.Code)
	}
}

// Send SSE headers
func TestResponseHandler_Handle_6(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := mock_amqp.NewMockConsumer(ctrl)
	consumer.EXPECT().Consume("").Return(nil, nil)
	consumer.EXPECT().Close().Times(1)

	provider := mock_amqp.NewMockProvider(ctrl)
	provider.EXPECT().Consumer(gomock.Any(), gomock.Any()).Return(consumer, nil)

	logger := mock_logg.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(2)

	recorder := httptest.NewRecorder()

	waitForContextDone := func() <-chan struct{} {
		done := make(chan struct{})
		close(done)
		return done
	}

	key := "Content-Type"
	val := "text/event-stream"

	r := &responseHandler{
		writer:   recorder,
		request:  http.Request{Method: "GET"},
		pattern:  NewPattern(""),
		provider: provider,
		logger:   logger,
		ctxDone:  waitForContextDone,
		headers:  map[string]string{key: val},
	}
	r.Handle()

	if recorder.Code != 200 {
		t.Errorf("expected 200, got %d", recorder.Code)
	}

	if !strings.Contains(recorder.Result().Header.Get(key), val) {
		t.Errorf("expected %s:%s", key, val)
	}

	if !strings.HasPrefix(recorder.Body.String(), ": SSE stream\n\n") {
		t.Error("expected stream banner")
	}
}
