package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type handler struct {
	Action string
	Method string
	URI    string
	Params []string
	Auth   bool
	Before []http.HandlerFunc
	After  []http.HandlerFunc
	Do     http.HandlerFunc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Auth {
		ok := auth(w, r)
		if !ok {
			http.Error(w, "Unauthorized Access", http.StatusUnauthorized)
			return
		}
	}

	for _, m := range h.Before {
		m(w, r)
	}

	// validate request parameters
	if len(h.Params) > 0 {
		err := validate(r, h.Params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	h.Do(w, r)

	for _, m := range h.After {
		m(w, r)
	}
}

func NewServeMux(groups ...[]handler) *http.ServeMux {
	mux := http.NewServeMux()

	// To be able to register different handlers for the same uri with the standard http.ServeMux,
	// we don't register the handler directly, we create a top layer handlerFunc
	// that routes the request to the right handler based on the request method.
	var routes = make(map[string]map[string]handler)
	for _, group := range groups {
		for _, h := range group {
			if routes[h.URI] == nil {
				routes[h.URI] = make(map[string]handler)
			}

			method := strings.ToUpper(h.Method)
			if routes[h.URI][method].Method != "" {
				panic("Duplicate Routes: " + method + " " + h.URI)
			}
			routes[h.URI][strings.ToUpper(h.Method)] = h
		}
	}

	for uri, m := range routes {
		mux.HandleFunc(uri, routeFunc(m))
	}
	return mux
}

func routeFunc(m map[string]handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m[r.Method].Method != "" {
			m[r.Method].ServeHTTP(w, r)
		} else {
			for k := range m {
				m[k].ServeHTTP(w, r)
				return
			}
		}
	}
}

// validate presence of params and checks specified validations
// numeric: should be numerical value
// max=x:   can't be more than x
func validate(r *http.Request, params []string) error {
	var err error
	err = r.ParseForm()
	if err != nil {
		return err
	}

	for _, param := range params {
		p := strings.Split(param, ":")

		if len(p) < 2 { // no rules, just validate presence
			if r.FormValue(p[0]) == "" {
				return errors.New(p[0] + " can't be empty")
			}
			continue
		}

		if r.FormValue(p[0]) == "" && !strings.Contains(p[1], "empty") {
			return errors.New(p[0] + " can't be empty")
		}

		s := strings.Split(p[1], ",")

		for _, c := range s {
			vs := strings.Split(c, "=")
			switch vs[0] {
			case "max":
				max, _ := strconv.Atoi(vs[1])
				prm := strings.TrimSpace(r.FormValue(p[0]))
				if len(prm) > max {
					return errors.New(p[0] + " can't be more than " + vs[1] + " characters")
				}
			case "min":
				min, _ := strconv.Atoi(vs[1])
				prm := strings.TrimSpace(r.FormValue(p[0]))
				if len(prm) < min {
					return errors.New(p[0] + " can't be less than " + vs[1] + " characters")
				}
			case "numeric":
				_, err := strconv.Atoi(r.FormValue(p[0]))
				if err != nil {
					return errors.New(p[0] + " should be an integer string")
				}

			}
		}

	}
	return nil
}

// TODO: more rules
// whitelist
// blacklist
// match
// regexp
// len
// contain
// url
// abs
// bytemin
// bytemax
// nonzero

func auth(w http.ResponseWriter, r *http.Request) bool {
	// authentication logic
}
