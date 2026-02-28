package main

import (
	"sync"
	"time"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
	_ "github.com/bamgoo/builtin"
	"github.com/bamgoo/log"
	_ "github.com/bamgoo/log-file"
	_ "github.com/bamgoo/log-greptime"
)

func main() {
	bamgoo.Go()
}

func init() {
	bamgoo.Register(bamgoo.START, bamgoo.Trigger{
		Name: "Log Demo",
		Desc: "emit many logs to test async batch and overflow strategy",
		Action: func(ctx *bamgoo.Context) {
			const workers = 16
			const perWorker = 5000

			start := time.Now()
			var wg sync.WaitGroup
			wg.Add(workers)

			for w := 0; w < workers; w++ {
				worker := w
				go func() {
					defer wg.Done()
					for i := 0; i < perWorker; i++ {
						if i%1000 == 0 {
							log.Warningw("burst", Map{"worker": worker, "seq": i})
						} else {
							log.Infof("burst worker=%d seq=%d", worker, i)
						}
					}
				}()
			}

			wg.Wait()
			took := time.Since(start)
			log.Noticew("log-demo done", Map{"workers": workers, "per_worker": perWorker, "total": workers * perWorker, "cost": took.String()})
		},
	})
}
