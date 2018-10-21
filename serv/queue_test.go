// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"net/http"
	"net/url"
	"testing"
)

func TestPattern_QueueName(t *testing.T) {

	samples := []struct {
		pattern string
		query   string
		cookie  string
		expect  string
	}{
		{
			pattern: "",
			query:   "",
			cookie:  "",
			expect:  "",
		},
		{
			pattern: "${unknown:foo}",
			query:   "",
			cookie:  "",
			expect:  "",
		},
		{
			pattern: "${query:oid}",
			query:   "",
			cookie:  "",
			expect:  "",
		},
		{
			pattern: "${query:oid}",
			query:   "oid=",
			cookie:  "",
			expect:  "",
		},
		{
			pattern: "${query:oid}",
			query:   "oid=123",
			cookie:  "",
			expect:  "123",
		},
		{
			pattern: " ${ query : oid } ",
			query:   "oid=123",
			cookie:  "",
			expect:  "123",
		},
		{
			pattern: " x${ query\n : oid }x ",
			query:   "oid=123",
			cookie:  "",
			expect:  "x123x",
		},
		{
			pattern: "${cookie:sid}",
			query:   "",
			cookie:  "",
			expect:  "",
		},
		{
			pattern: "${cookie:sid}",
			query:   "",
			cookie:  "sid=",
			expect:  "",
		},
		{
			pattern: "${cookie:sid}",
			query:   "",
			cookie:  "sid=123",
			expect:  "123",
		},
		{
			pattern: " ${ cookie : sid } ",
			query:   "",
			cookie:  "sid=123",
			expect:  "123",
		},
		{
			pattern: " x${ cookie\n : sid }x ",
			query:   "",
			cookie:  "sid=123",
			expect:  "x123x",
		},
		{
			pattern: "${query:oid}-${cookie:sid}",
			query:   "",
			cookie:  "",
			expect:  "-",
		},
		{
			pattern: "${query:oid}-${cookie:sid}",
			query:   "oid=",
			cookie:  "sid=",
			expect:  "-",
		},
		{
			pattern: "${query:oid}-${cookie:sid}",
			query:   "oid=ABC",
			cookie:  "sid=123",
			expect:  "ABC-123",
		},
		{
			pattern: " ${ query : oid } ${ cookie : sid } ",
			query:   "oid=ABC",
			cookie:  "sid=123",
			expect:  "ABC123",
		},
		{
			pattern: "x${query :oid\n} x${ cookie\n : sid }x ",
			query:   "oid=ABC",
			cookie:  "sid=123",
			expect:  "xABCx123x",
		},
	}

	for _, sample := range samples {
		pattern := NewPattern(sample.pattern)

		req := &http.Request{
			Header: http.Header{"Cookie": {sample.cookie}},
			URL:    &url.URL{RawQuery: sample.query},
		}

		result, _ := pattern.QueueName(req)

		if result != sample.expect {
			t.Errorf("expected %s, got %s", sample.expect, result)
		}
	}
}
