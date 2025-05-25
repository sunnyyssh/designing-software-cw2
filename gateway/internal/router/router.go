package router

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	Locations []Location
}

type Location struct {
	Prefix string
	URL    string
}

type Router struct {
	// Sorted by prefix length
	locs []Location
}

func New(config *Config) *Router {
	return &Router{
		locs: config.Locations,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	log.Printf("got %s %s", req.Method, path)

	var (
		loc       Location
		routePath string
		ok        bool
	)

	for _, l := range r.locs {
		if routePath, ok = strings.CutPrefix(path, l.Prefix); ok {
			loc = l
			break
		}
	}

	if !ok {
		notFound(w)
		return
	}

	routeReq, err := buildReq(req, &loc, routePath)
	if err != nil {
		internalServerError(w)
		return
	}

	resp, err := http.DefaultClient.Do(routeReq)
	if err != nil {
		log.Printf("failed to route request to %s: %s", loc.URL, err)
		badGateway(w)
		return
	}

	if err := write(w, resp); err != nil {
		log.Printf("failed to copy response: %s", err)
		internalServerError(w)
		return
	}

	log.Printf("ret %s %s %d", req.Method, path, resp.StatusCode)
}

func write(w http.ResponseWriter, resp *http.Response) error {
	for key, vals := range resp.Header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func buildReq(baseReq *http.Request, loc *Location, path string) (*http.Request, error) {
	defer baseReq.Body.Close()

	body := &bytes.Buffer{}
	if _, err := io.Copy(body, baseReq.Body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(baseReq.Method, loc.URL+path, body)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = baseReq.URL.RawQuery
	req.Header = baseReq.Header

	return req, nil
}

var notFoundBody = []byte(`{"error":"Not Found","code":404}`)

func notFound(w http.ResponseWriter) error {
	w.WriteHeader(404)
	_, err := w.Write(notFoundBody)
	if err != nil {
		log.Printf("failed to send response: %s", err)
		return err
	}

	return nil
}

var internalServerErrorBody = []byte(`{"error":"Internal server error","code":500}`)

func internalServerError(w http.ResponseWriter) error {
	w.WriteHeader(500)
	_, err := w.Write(internalServerErrorBody)
	if err != nil {
		log.Printf("failed to send response: %s", err)
		return err
	}

	return nil
}

var badGatewayBody = []byte(`{"error":"Bad gateway","code":502}`)

func badGateway(w http.ResponseWriter) error {
	w.WriteHeader(502)
	_, err := w.Write(badGatewayBody)
	if err != nil {
		log.Printf("failed to send response: %s", err)
		return err
	}

	return nil
}
