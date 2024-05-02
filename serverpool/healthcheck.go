package serverpool

import (
	"log"
	"sync"
	"time"

	"github.com/rimo10/load_balancer/backend"
)

func Healthcheck(s Serverpool) {
	var wg sync.WaitGroup
	for _, b := range s.GetBackends() {
		wg.Add(1)
		b := b
		go func(bac backend.Backend) {
			defer wg.Done()
			status := "up"
			alive := bac.IsAlive()
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
			log.Printf("%s [%s]", bac.GetUrl(), status)
		}(b)
	}
	wg.Wait()
}

func LaunchHealthCheck(sp Serverpool) {
	t := time.NewTicker(time.Minute * 2)
	log.Printf("Starting Health check...")
	for {
		select {
		case <-t.C:
			go Healthcheck(sp)
			log.Printf("Health check completed")
		}
	}
}
