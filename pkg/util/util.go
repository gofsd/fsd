package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"path/filepath"
	"regexp"
	"runtime/debug"
)

func Exec(s string, timeout int64) (output []byte, err error) {
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", s)
	output, err = cmd.Output()
	return
}

func MapToFlagsArr(flags map[string]string) (s []string) {
	s = make([]string, 0, len(flags))

	for key, value := range flags {
		if len(value) > 0 {
			s = append(s, fmt.Sprintf("--%s %s", key, value))
		} else {
			s = append(s, fmt.Sprintf("--%s", key))
		}
	}
	return
}

// ExecForFilesWithExtension to execute for each file with some extension like *.json, *.go etc. or all files .*
func ExecForFilesWithExtension(dir string, ext string, execFn func(fullName string, name string)) {
	var valid = regexp.MustCompile("\\" + ext)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if FileExists(filepath.Join(path)) {
			if ext == "" {
				execFn(path, info.Name())
			} else if valid.MatchString(info.Name()) {
				execFn(path, info.Name())
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// FileExists check file
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	return !info.IsDir()
}

// FindRootDir in current and parents directories
func FindRootDir(filename string) (wd string) {
	wd, _ = os.Getwd()

	for {
		files, err := ioutil.ReadDir(wd)

		parrent := filepath.Dir(wd)

		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if filename == f.Name() {
				return wd
			}
		}
		if wd == parrent {
			wd, _ = os.Getwd()
			break
		} else {
			wd = parrent
		}
	}
	return
}

// HandleError is func for handle error
func HandleError(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}

func HandlePanic() {}

func FileSize(fileName string) int64 {
	f, err := os.Stat(fileName)
	HandleError(err)
	return f.Size()
}
