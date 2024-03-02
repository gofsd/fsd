package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofsd/fsd/pkg/util"

	"github.com/fsnotify/fsnotify"
)

func execOnChange(dir string, cmds []string) (err error) {
	var process *exec.Cmd
	var canExecute bool = true
	updateTimer := time.NewTicker(300 * time.Millisecond)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		util.HandleError(err)
	}
	defer watcher.Close()

	// Start listening for events.

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !canExecute && event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					canExecute = true
				}
				break
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				util.HandleError(err)
				break
			case <-updateTimer.C:
				if canExecute {
					canExecute = false
					for _, cmdStr := range cmds {
						parts := strings.Split(cmdStr, " ")
						process = exec.Command(parts[0], parts[1:]...)
						process.Run()
					}
					filepath.Walk(config.GetString("project_root")+config.GetString("listenDirToGenerate"), func(p string, info os.FileInfo, err error) error {
						if info.IsDir() {
							watcher.Add(p)
						}
						return nil
					})
				}

			}
		}
	}()
	// Add a path.
	if err != nil {
		util.HandleError(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}
