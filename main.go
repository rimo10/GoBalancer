package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/rimo10/load_balancer/backend"
	"github.com/rimo10/load_balancer/lb"
	"github.com/rimo10/load_balancer/serverpool"
	"github.com/rimo10/load_balancer/utils"
)

func main() {
	config, err := utils.GetLBConfig()
	if err != nil {
		log.Fatal(err)
	}
	sp, err := serverpool.NewServerPool()
	if err != nil {
		log.Fatalf(err.Error())
	}
	// ctx := context.Background()
	loadbalancer := lb.NewLoadBalancer(sp)
	for _, b := range config.Backends {
		serverUrl, err := url.Parse(b)
		if err != nil {
			log.Fatal(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[%s] %s \n", serverUrl.Host, err.Error())
			retries := loadbalancer.GetRetryFromContext(r)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(r.Context(), lb.Retry, retries+1)
					proxy.ServeHTTP(w, r.WithContext(ctx))
				}
				return
			}
			sp.MarkBackendStatus(serverUrl, false)
		}
		newBackend := backend.NewBackend(serverUrl, proxy)
		sp.AddBackend(newBackend)
		log.Printf("Configured Server : %s\n", serverUrl)
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: http.HandlerFunc(loadbalancer.Serve),
	}
	go serverpool.LaunchHealthCheck(sp)
	log.Printf("Load Balancer started at : %d\n", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
