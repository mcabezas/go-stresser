package endpoint

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"stresser/internal/statistics"
)

type EndpointsConfig struct {
	Default struct {
		Headers []struct {
			Key   string
			Value string
		}
		URL    string
		Method string
	}
	Endpoints []*endpoint
}

type endpoint struct {
	Name    string
	Headers []struct {
		Key   string
		Value string
	}
	URL    string
	Method string
	Body   string
	limiter *rate.Limiter
	client *http.Client
}

func ParseFile(path string, client *http.Client, limiter *rate.Limiter) []EndPoint {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg := EndpointsConfig{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	var result []EndPoint
	for _, e := range cfg.Endpoints {
		if len(e.Headers) == 0 {
			e.Headers = cfg.Default.Headers
		}
		if e.URL == "" {
			e.URL = cfg.Default.URL
		}
		if e.Method == "" {
			e.Method = cfg.Default.Method
		}
		e.client = client
		e.limiter = limiter
		result = append(result, e)
	}
	return result
}

func (e *endpoint) Execute(formatter statistics.Formatter) statistics.Statistics {
	body := strings.NewReader(e.Body)
	request, err := http.NewRequest(e.Method, e.URL, body)
	if err != nil {
		return statistics.NewErrStatistics(e.Name, time.Now(), time.Now(), err, formatter)
	}
	headers := http.Header{}
	for _, h := range e.Headers {
		headers.Add(h.Key, h.Value)
	}
	request.Header = headers
	if err := e.limiter.Wait(context.Background()); err != nil {
		return statistics.NewErrStatistics(e.Name, time.Now(), time.Now(), err, formatter)
	}
	start := time.Now()
	response, err := e.client.Do(request)
	if err != nil {
		return statistics.NewErrStatistics(e.Name, start, time.Now(), err, formatter)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	respBody, _ := ioutil.ReadAll(response.Body)
	return statistics.NewStatistics(e.Name, start, time.Now(), response.StatusCode, respBody, formatter)
}
