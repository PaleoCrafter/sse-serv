// MIT license, dtg [at] lengo [dot] org · 10/2018

package serv

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// QueuePattern ...
type QueuePattern interface {
	QueueName(r *http.Request) (string, error)
}

type queuePattern struct {
	pattern string
	regex   *regexp.Regexp
}

var patternClean = strings.NewReplacer("\n", "", " ", "")
var patternRegex = regexp.MustCompile(`\${([^:]+):([^}]+)}`)

// NewPattern ...
func NewPattern(s string) QueuePattern {
	return &queuePattern{
		pattern: patternClean.Replace(s),
		regex:   patternRegex,
	}
}

func (p *queuePattern) QueueName(r *http.Request) (string, error) {
	var nameErr error

	cb := func(m string) string {
		var cat, key, val string
		if nameErr == nil {
			s := p.regex.ReplaceAllString(m, "$1 $2")
			fmt.Sscanf(s, "%s %s", &cat, &key)
			val, nameErr = resolve(r, cat, key)
		}
		return val
	}

	return p.regex.ReplaceAllStringFunc(p.pattern, cb), nameErr
}

func resolve(r *http.Request, cat string, key string) (string, error) {
	switch cat {
	case "cookie":
		return cookieValue(r, key)
	case "query":
		return queryValue(r, key)
	default:
		return "", errRequestParamMissing
	}
}

func cookieValue(r *http.Request, key string) (string, error) {
	var cookie *http.Cookie
	var err error

	if cookie, err = r.Cookie(key); err != nil {
		return "", err
	}
	if cookie.Value != "" {
		return cookie.Value, nil
	}
	return "", errRequestParamMissing
}

func queryValue(r *http.Request, key string) (string, error) {
	if val := r.URL.Query().Get(key); val != "" {
		return val, nil
	}
	return "", errRequestParamMissing
}

var errRequestParamMissing = errors.New("request parameter(s) missing")
