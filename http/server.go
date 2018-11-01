package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ManageServer starts srvr and wait for SIGINT or SIGQUIT or SIGKILL or SIGTERM
// to stop the server gracefully. If a second signal found or the gracePeriod
// has expired it stops the server immedietly. It returns any error that is returned
// by the srvr
func ManageServer(srvr *http.Server, gracePeriod time.Duration) error {
	sigCh := make(chan os.Signal, 0)
	errCh := make(chan error, 0)
	go func() {
		log.Println("Starting web server on", srvr.Addr)
		errCh <- srvr.ListenAndServe()
	}()

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt)

	for i := 0; i < 2; i++ {
		select {
		case err := <-errCh:
			return err
		case <-sigCh:
			if i == 0 {
				d := gracePeriod
				log.Println("Suttingdown server gracefully with in", d)
				log.Println("To shutdown immedietly press again")
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), d)
					defer cancel()

					if err := srvr.Shutdown(ctx); err != nil {
						errCh <- err
						return
					}
					errCh <- nil
				}()
				continue
			}

			log.Println("Suttingdown web server forcefully")
			if err := srvr.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
