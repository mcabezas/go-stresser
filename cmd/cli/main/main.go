package main

import (
	"flag"
	"net/http"
	"os"

	"golang.org/x/time/rate"
	"stresser/internal/endpoint"
	"stresser/internal/executor"
	"stresser/internal/statistics"
)

func main() {
	cfg := readArgs()
	standardLibClient := &http.Client{}
	endpoints := endpoint.ParseFile(cfg.EndpointsConfigRoute, standardLibClient, rate.NewLimiter(rate.Limit(cfg.MaxPoolSize), cfg.ReqPerSec) )
	ex := executor.NewParallel(
		cfg.Hits,
		endpoints,
		statistics.NewFormatter(',', '\n'),
		os.Stdout)
	_ = ex.Start()
}

type config struct {
	MaxPoolSize          int
	ReqPerSec            int
	Hits                 int
	EndpointsConfigRoute string
}

func readArgs() *config {
	cfg := &config{}
	flag.IntVar(&cfg.MaxPoolSize, "mps", 200, "max pool size")
	flag.IntVar(&cfg.ReqPerSec, "rps", 200, "requests per second")
	flag.IntVar(&cfg.Hits, "hits", 100000, "total requests to be performed")
	flag.StringVar(&cfg.EndpointsConfigRoute, "ecr", "config.json", "endpoint config route")
	flag.Parse()
	return cfg
}
