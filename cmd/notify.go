package cmd

import (
	"log"
	"os"
	"os/exec"

	"time"

	"github.com/fsnotify/fsnotify"
)

func notifyStart() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan int)
	quit := make(chan bool)

	go func() {
		var c *exec.Cmd
		c = nil
		var modif bool
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&event.Op == fsnotify.Write {

					time.Sleep(200 * time.Millisecond)
					port := exec.Command("fuser", "-n", "tcp", "-k", "8080")
					port.Run()
					if c != nil && !modif {
						quit <- true
						c.Process.Kill()
						c = nil
						modif = !modif
					} else {
						modif = !modif
						go func() {

							if c == nil {

								c = exec.Command("dlv", "test", "./backend/")
								c.Stdin = os.Stdin
								c.Stdout = os.Stdout
								c.Stderr = os.Stderr

								c.Run()
							}
							for {
								_, ok := <-quit
								if c != nil && !ok {
									if c.Process != nil {
										c.Process.Kill()
										c = nil

										return
									}
								}
								if c == nil {
									return
								}
							}
						}()
					}

					//done <- 1
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	err = watcher.Add("./backend")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(<-done)
}
