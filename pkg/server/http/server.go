package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"path/filepath"

	"github.com/gofsd/fsd/pkg/server/http/dbms"
	"github.com/robfig/cron/v3"

	"runtime"
	"syscall"

	"github.com/gofsd/fsd/pkg/pipe"

	"github.com/boltdb/bolt"
	validator "github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const SECOND = 1000000000

var CLIENT_APP_VERSION, STATIC_PATH, RootPath, CmdPath string
var backupTimestamp int64
var collections map[string]interface{} = make(map[string]interface{})
var lastSecondIps map[string]int64 = make(map[string]int64)
var tags []Tag
var db *bolt.DB
var err error

// Handler for testing
var Handler *http.ServeMux

var validate *validator.Validate = validator.New()

var updState = map[int]bool{}

var currentIotState = map[int]bool{}

var state = map[int]bool{
	0:  false,
	2:  false,
	4:  false,
	5:  false,
	12: false,
	13: false,
	14: false,
	16: false,
}
var ESP8266HOST = "192.168.1.6"

type User struct {
	ID       int    `validate:"required`
	UUID     string `validate:"required,len=36"`
	RefUUID  string `validate:"required,len=36"`
	Created  int64  `validate:"required"`
	LastUsed int64  `validate:"required"`
}

type TypeOfStat struct {
	Correct   bool
	Timestamp int64
}

type Statistic struct {
	ID                  int          `validate:"required"`
	TagID               int64        `validate:"required"`
	Tag                 string       `validate:"required,max=20"`
	Listening           []TypeOfStat `validate:"required,len=200"`
	Speaking            []TypeOfStat `validate:"required,len=200"`
	Translating         []TypeOfStat `validate:"required,len=200"`
	MaxStepListening    int          `validate:"required"`
	MaxStepSpeaking     int          `validate:"required"`
	MaxStepTranslating  int          `validate:"required"`
	CurrStepListening   int          `validate:"required"`
	CurrStepSpeaking    int          `validate:"required"`
	CurrStepTranslating int          `validate:"required"`
	RightListening      int          `validate:"required"`
	WrongListening      int          `validate:"reqiured"`
	RightSpeaking       int          `validate:"reqiured"`
	WrongSpeaking       int          `validate:"required"`
	RightTranslating    int          `validate:"required"`
	WrongTranslating    int          `validate:"required"`
	OveralPeaks         int          `validate:"required"`
	NextTimestamp       int64        `validate:"required"`
	Updated             int64        `validate:"required"`
	Lang                string       `validate:"required,len=10"`
}

// {"ID_UUID": "xxxx-xxxx-xxxx", "ExerciseType": "listening", "TagID": 1, "Correct": true}
type StatUpdRequest struct {
	ID           int    `validate:"required"`
	UUID         string `validate:"required,len=36"`
	ExerciseType string `validate:"required,max=15"`
	TagID        int64  `validate:"required,min=0"`
	Tag          string `validate:"required,max=20,min=1"`
	Correct      bool
	Lang         string `validate:"required,max=5"`
}

type Tag struct {
	TagID        int64    `validate:"required"`
	Tag          string   `validate:"required,min=1,max=45"`
	Translations []string `validate:"required,min=1,max=30,dive,max=140"`
	Links        []int64  `validate:"max=100"`
	Disabled     bool
	Lang         string
}

type NextWordResp struct {
	StatRequest StatUpdRequest
	Tag         Tag
}

type NextWordRequest struct {
	ID   int    `validate:"required"`
	UUID string `validate:"required,len=36"`
	Lang string `validate:"required,max=5"`
}

var Listening, Speaking, Translating string = "listening", "speaking", "translating"

var Intervals [8]int64 = [8]int64{120, 600, 3600, 3600 * 5, 3600 * 24, 3600 * 24 * 5, 3600 * 24 * 25, 3600 * 24 * 125}

//var Intervals [8]int64 = [8]int64{15, 30, 50, 70, 90, 110, 130, 150}

func (stat *Statistic) UpdListening(upd *StatUpdRequest) {
	if upd.Correct {
		stat.RightListening += 1
		stat.CurrStepListening += 1

		if stat.CurrStepListening > stat.MaxStepListening || stat.CurrStepListening >= 2 {
			if stat.MaxStepListening == 0 {
				stat.OveralPeaks = 0
			}
			stat.MaxStepListening = stat.CurrStepListening

			if stat.OveralPeaks < 7 {
				stat.OveralPeaks = stat.OveralPeaks + 1

			}

		}
		nextTimeStamp := time.Now().Unix() + Intervals[stat.OveralPeaks]
		if stat.NextTimestamp < nextTimeStamp {
			stat.NextTimestamp = nextTimeStamp
		}
		if i := len(stat.Listening); i < 200 {
			stat.Listening = append(stat.Listening, TypeOfStat{
				true,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Listening); i++ {
				stat.Listening[i-1] = stat.Listening[i]
			}
			stat.Listening[199] = TypeOfStat{
				true,
				time.Now().Unix(),
			}
		}
	} else {
		stat.WrongListening += 1
		stat.CurrStepListening = 0
		stat.NextTimestamp = time.Now().Unix() + Intervals[0]
		if i := len(stat.Listening); i < 200 {
			stat.Listening = append(stat.Listening, TypeOfStat{
				false,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Listening); i++ {
				stat.Listening[i-1] = stat.Listening[i]
			}
			stat.Listening[199] = TypeOfStat{
				false,
				time.Now().Unix(),
			}
		}
	}
}

func (stat *Statistic) UpdSpeaking(upd *StatUpdRequest) {
	if upd.Correct {
		stat.RightSpeaking += 1
		stat.CurrStepSpeaking = stat.CurrStepSpeaking + 1
		if stat.CurrStepSpeaking > stat.MaxStepSpeaking || stat.CurrStepSpeaking >= 2 {
			if stat.MaxStepSpeaking == 0 {
				stat.OveralPeaks = 0
			}
			stat.MaxStepSpeaking = stat.CurrStepSpeaking
			if stat.OveralPeaks < 7 {
				stat.OveralPeaks = stat.OveralPeaks + 1
			}

		}
		nextTimeStamp := time.Now().Unix() + Intervals[stat.OveralPeaks]
		if stat.NextTimestamp < nextTimeStamp {
			stat.NextTimestamp = nextTimeStamp
		}
		if i := len(stat.Speaking); i < 200 {
			stat.Speaking = append(stat.Speaking, TypeOfStat{
				true,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Speaking); i++ {
				stat.Speaking[i-1] = stat.Speaking[i]
			}
			stat.Speaking[199] = TypeOfStat{
				true,
				time.Now().Unix(),
			}
		}
	} else {
		stat.WrongSpeaking += 1
		stat.CurrStepSpeaking = 0
		stat.NextTimestamp = time.Now().Unix() + Intervals[0]
		if i := len(stat.Speaking); i < 200 {
			stat.Speaking = append(stat.Speaking, TypeOfStat{
				false,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Speaking); i++ {
				stat.Speaking[i-1] = stat.Speaking[i]
			}
			stat.Speaking[199] = TypeOfStat{
				false,
				time.Now().Unix(),
			}
		}
	}
}

func (stat *Statistic) UpdTranslating(upd *StatUpdRequest) {
	if upd.Correct {
		stat.RightTranslating += 1
		stat.CurrStepTranslating = stat.CurrStepTranslating + 1
		if stat.CurrStepTranslating > stat.MaxStepTranslating || stat.CurrStepSpeaking >= 2 {
			if stat.MaxStepTranslating == 0 {
				stat.OveralPeaks = 0
			}
			stat.MaxStepTranslating = stat.CurrStepTranslating
			if stat.OveralPeaks < 7 {
				stat.OveralPeaks = stat.OveralPeaks + 1
			}

		}
		nextTimeStamp := time.Now().Unix() + Intervals[stat.OveralPeaks]
		if stat.NextTimestamp < nextTimeStamp {
			stat.NextTimestamp = nextTimeStamp
		}
		if i := len(stat.Translating); i < 200 {
			stat.Translating = append(stat.Translating, TypeOfStat{
				true,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Translating); i++ {
				stat.Translating[i-1] = stat.Translating[i]
			}
			stat.Translating[199] = TypeOfStat{
				true,
				time.Now().Unix(),
			}
		}
	} else {
		stat.WrongTranslating += 1
		stat.CurrStepTranslating = 0
		stat.NextTimestamp = time.Now().Unix() + Intervals[0]
		if i := len(stat.Translating); i < 200 {
			stat.Translating = append(stat.Translating, TypeOfStat{
				false,
				time.Now().Unix(),
			})
		} else {
			for i := 1; i < len(stat.Translating); i++ {
				stat.Translating[i-1] = stat.Translating[i]
			}
			stat.Translating[199] = TypeOfStat{
				false,
				time.Now().Unix(),
			}
		}
	}
}

func (stat *Statistic) Update(upd *StatUpdRequest) {
	stat.Updated = time.Now().Unix()
	switch upd.ExerciseType {
	case Listening:
		stat.UpdListening(upd)
		break
	case Speaking:
		stat.UpdSpeaking(upd)
		break
	case Translating:
		stat.UpdTranslating(upd)
		break
	}
}

// Run some comment
func Run(port string) {
	// var updateState = cron.New(cron.WithSeconds())
	// updateState.AddFunc("0,10,20,30,40,50 * * * * *", func() {
	// 	// UpdateAll()
	// 	//print("check\n")
	// 	now := time.Now()
	// 	fmt.Println("Current date and time (RFC3339):", now.Format(time.RFC3339))
	// })
	// updateState.Start()

	// var solutionSupply = cron.New(cron.WithSeconds())
	// solutionSupply.AddFunc("0 7,14,21,28,35,42,49,56 20-23,0-9 * * *", func() {
	// 	SetPin(0, false)

	// 	SetPin(16, true)

	// 	print("solution supply\n")
	// })
	// solutionSupply.Start()

	// var stopSolutionSupply = cron.New(cron.WithSeconds())
	// stopSolutionSupply.AddFunc("20 7,14,21,28,35,42,49,56 20-23,0-9 * * *", func() {
	// 	SetPin(16, false)
	// 	SetPin(0, true)

	// 	print("stop solution supply\n")
	// })
	// stopSolutionSupply.Start()

	// var solutionDrain = cron.New(cron.WithSeconds())
	// solutionDrain.AddFunc("0 0-14,28-44 10-15 * * *", func() {
	// 	SetPin(5, true)

	// 	print("solution drain\n")
	// })
	// solutionDrain.Start()

	// var stopSolutionDrain = cron.New(cron.WithSeconds())
	// stopSolutionDrain.AddFunc("0 15-20,45-50 10-15 * * *", func() {
	// 	SetPin(5, false)

	// 	print("stop solution drain\n")
	// })
	// stopSolutionDrain.Start()

	var light = cron.New(cron.WithSeconds())
	light.AddFunc("40 5 22-23,0-8 * * *", func() {
		SetPin(0, true)
		SetPin(2, true)
		print("light \n")
	})
	light.Start()

	var offLight = cron.New(cron.WithSeconds())
	offLight.AddFunc("0 * 9-21 * * *", func() {
		SetPin(0, false)
		SetPin(2, false)
		print("off light\n")
	})
	offLight.Start()

	var dnat = cron.New(cron.WithSeconds())
	dnat.AddFunc("25 0 23,0-6 * * *", func() {
		SetPin(4, true)
		print("dnat light\n")
	})
	dnat.Start()

	var dnatPeriodcOff = cron.New(cron.WithSeconds())
	dnatPeriodcOff.AddFunc("10 * 7-22 * * *", func() {
		SetPin(4, false)
		print("dnat light off\n")
	})
	dnatPeriodcOff.Start()

	// var offDnat = cron.New(cron.WithSeconds())
	// offDnat.AddFunc("0 0 12-23 * * *", func() {
	// 	SetPin(4, false)
	// 	print("off dnat\n")
	// })
	// offDnat.Start()
	fmt.Sprintf(":%s", port)
	err = http.ListenAndServe(":8881", Handler)

	if os.Getenv("CLIENT_APP_VERSION") == "" {
		os.Setenv("CLIENT_APP_VERSION", "2.0.20")
	}
	CLIENT_APP_VERSION = os.Getenv("CLIENT_APP_VERSION")
	var err error
	if os.Getenv("ENV") == "production" {
		err = http.ListenAndServeTLS(":443", "../fsd_xyz_crt", "../private_key_ssl", Handler)
	} else if os.Getenv("ENV") == "dev" {
		err = http.ListenAndServeTLS(":3000", "../fsd_xyz_crt", "../private_key_ssl", Handler)
	} else {
		///err = http.ListenAndServe(fmt.Sprintf(":%s", port), Handler)
	}

	if err != nil {
		print(err.Error())
	}
}

// Create some comment again
func Create() {
	var tx *bolt.Tx
	print(os.Getenv("STATIC_PATH"))
	if os.Getenv("STATIC_PATH") == "" {

		STATIC_PATH = "/static"

		path, er := os.Getwd()
		if er != nil {
			fmt.Println(er.Error())
		}
		STATIC_PATH = path + STATIC_PATH
		err = os.MkdirAll(STATIC_PATH, 0750)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		STATIC_PATH = os.Getenv("STATIC_PATH")
	}
	if os.Getenv("path_to_db") == "" {
		os.Setenv("path_to_db", "my.db")
	}
	db, err = bolt.Open(os.Getenv("path_to_db"), 0777, nil)
	pipe.SetDB(db)
	if err != nil {
		log.Fatalf("Bolt.open: %s", err)
	}
	//defer db.Close()

	tx, err = db.Begin(true)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists([]byte("user"))
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = tx.CreateBucketIfNotExists([]byte("statistic"))
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = tx.CreateBucketIfNotExists([]byte("tag"))
	if err != nil {
		fmt.Print(err.Error())
	}

	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		fmt.Print(err.Error())
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tag"))
		tag := Tag{}
		c := b.Cursor()
		prefix := []byte(fmt.Sprintf("%s_%s", "1", "en"))
		for k, v := c.First(); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			err = json.Unmarshal(v, &tag)
			if tag.Disabled == true {
				tags = append(tags, tag)
			}
		}
		return err
	})

	go runJobs()
	Handler = http.NewServeMux()
	customDb, _ := dbms.CreateDb(new(dbms.DB))
	Handler = customDb.MuxHandler
	pipe.SetHTTP(Handler)

	Handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}
		Resp(w, "off")
		return
		switch r.Method {
		case http.MethodGet:

			Resp(w, `<!DOCTYPE html>
			<html lang="en">
			  <head>
					<title>File Upload Demo</title>
			  </head>
			  <body>
					<div class="container">
					  <h1>File Upload Demo</h1>
					  <form class="form-signin" method="post" action="/test/" enctype="multipart/form-data">
							  <fieldset>
									<input type="file" name="myfiles" multiple="multiple">
									<input type="file" name="othermyfiles" multiple="multiple">
									<input type="submit" name="submit" value="Submit">
							</fieldset>
					  </form>
					</div>
			  </body>
			</html>`)
			break
		case http.MethodPost:

			reader, err := r.MultipartReader()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				//if part.FileName() is empty, skip this iteration.
				if part.FileName() == "" {
					continue
				}
				dst, err := os.Create("./DB/" + part.FileName())
				defer dst.Close()

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if _, err := io.Copy(dst, part); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			Resp(w, `{"foo":"bar"}`)
			break

		}
	})

	Handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		print(r.URL.Path)
		var p string
		if r.URL.Path == "/" {
			p = STATIC_PATH + "/index.html"
		}

		p = STATIC_PATH + r.URL.Path

		http.ServeFile(w, r, p)
	})

	Handler.HandleFunc("/tag/", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}

		if !checkAuth(r) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		}

		id := r.Header.Get("ID")
		var err error
		var tagBytes []byte
		tag := Tag{}
		if http.MethodPost == r.Method {
			body, _ := ioutil.ReadAll(r.Body)
			err = json.Unmarshal(body, &tag)
			err = validate.Struct(tag)
			if err != nil {
				Resp(w, `{"error": "%s"}`, err.Error())
				return
			}
			tagBytes, err = json.Marshal(tag)
			db.Update(func(tx *bolt.Tx) error {
				prefix := []byte(fmt.Sprintf("%s_%s", id, tag.Lang))
				b := tx.Bucket([]byte("tag"))
				b.Put([]byte(fmt.Sprintf("%s_%s_%d", id, tag.Lang, tag.TagID)), tagBytes)
				if id == "1" {
					AutoLearn(tag, id, tx)
					tags = tags[:0]
					b.Put([]byte(fmt.Sprintf("%s_%s_%d", id, "en", tag.TagID)), tagBytes)
					prefix = []byte(fmt.Sprintf("%s_%s", id, "en"))
					c := b.Cursor()
					for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
						err = json.Unmarshal(v, &tag)
						if tag.Disabled == true {
							tags = append(tags, tag)
						}
					}
				}

				return err
			})

			Resp(w, string(tagBytes))
		} else if http.MethodGet == r.Method {
			pathAr := strings.Split(r.URL.Path, "/")

			tagId := pathAr[2]
			lang := pathAr[3]

			prefix := []byte(fmt.Sprintf("%s_%s_%s", id, lang, tagId))
			db.View(func(tx *bolt.Tx) (err error) {
				b := tx.Bucket([]byte("tag"))
				v := b.Get(prefix)
				if v == nil {
					prefix = []byte(fmt.Sprintf("%s_%s_%s", "1", lang, tagId))
					v = b.Get(prefix)
				}

				if v == nil {
					col, _ := getCollection(lang)
					en, _ := getCollection("en_US")
					length := len(en.([]interface{}))
					lengthCol := len(col.([]interface{}))

					for i := 0; i < length; i++ {
						freak := int64(en.([]interface{})[i].(map[string]interface{})["frequency"].(float64))

						if tagId == strconv.FormatInt(int64(freak), 10) {
							for k := 0; k < lengthCol; k++ {
								if strings.ToLower(en.([]interface{})[i].(map[string]interface{})["name"].(string)) == strings.ToLower(col.([]interface{})[k].(map[string]interface{})["normalizedSource"].(string)) {
									tag.TagID = int64(freak)
									tag.Tag = strings.ToLower(en.([]interface{})[i].(map[string]interface{})["name"].(string))
									countWords := len(col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{}))
									for w := 0; countWords > w; w++ {
										if w > 14 {
											break
										}
										tag.Translations = append(tag.Translations, col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{})[w].(map[string]interface{})["displayTarget"].(string))
									}
								}
							}

						}
					}
					tag.Lang = lang

					tagBytes, err = json.Marshal(tag)
				} else {
					tagBytes = v
				}
				return err
			})

			if err != nil {
				Resp(w, `{"error": "%s"}`, err.Error())
				return
			}

			Resp(w, string(tagBytes))
		}
	})

	Handler.HandleFunc("/next-word", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}

		if !checkAuth(r) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)
		var req StatUpdRequest
		var resp NextWordResp
		err = json.Unmarshal(body, &req)

		if err != nil {
			Resp(w, `{"error": "%s"}`, err.Error())
			return
		}

		err = validate.Struct(req)

		if err != nil {
			Resp(w, `{"error": "%s"}`, err.Error())
			return
		}
		var bk *bolt.Bucket
		tagsPrefix := []byte(fmt.Sprintf("%d_%s", req.ID, req.Lang))
		var ownTags []Tag
		var tag Tag
		var stat Statistic
		var allStat []Statistic
		db.View(func(tx *bolt.Tx) error {
			bk = tx.Bucket([]byte("tag"))
			bkC := bk.Cursor()

			for k, v := bkC.Seek(tagsPrefix); k != nil && bytes.HasPrefix(k, tagsPrefix); k, v = bkC.Next() {
				json.Unmarshal(v, &tag)
				ownTags = append(ownTags, tag)
			}
			b := tx.Bucket([]byte("statistic"))

			c := b.Cursor()

			prefix := []byte(fmt.Sprintf("%d_%s", req.ID, req.Lang))

			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
				json.Unmarshal(v, &stat)
				allStat = append(allStat, stat)
			}

			return nil
		})

		for i := 0; i < len(allStat); i++ {
			stat = allStat[i]
			cont := false
			for i := 0; i < len(tags); i++ {
				if stat.TagID == tags[i].TagID {
					if tags[i].Disabled == true {
						cont = true
						break
					}
				}
			}
			if cont == true {
				cont = false
				continue
			}

			ownCount := false
			for i := 0; i < len(tags); i++ {
				if stat.TagID == tags[i].TagID {
					ownCount = true
				}
			}

			if ownCount == true {
				ownCount = false
				continue
			}

			for i := 0; i < len(ownTags); i++ {
				if stat.TagID == ownTags[i].TagID {
					if ownTags[i].Disabled == true {
						ownCount = true
						break
					}
				}
			}

			if ownCount == true {
				ownCount = false
				continue
			}

			if stat.NextTimestamp < time.Now().Unix() {
				resp.Tag.TagID = stat.TagID
				resp.Tag.Tag = stat.Tag
				resp.Tag.Lang = req.Lang
				resp.StatRequest.ID = req.ID
				resp.StatRequest.UUID = req.UUID
				resp.StatRequest.TagID = stat.TagID
				resp.StatRequest.Tag = stat.Tag
				resp.StatRequest.Lang = req.Lang
				if r.Header.Get("ID") != "1" {
					if stat.MaxStepListening < 2 {
						resp.StatRequest.ExerciseType = Listening
						sendNextWord(w, &req, &resp, allStat)
						return
					} else if stat.MaxStepTranslating < 2 {
						resp.StatRequest.ExerciseType = Translating
						sendNextWord(w, &req, &resp, allStat)
						return
					} else if stat.MaxStepSpeaking < 2 {
						resp.StatRequest.ExerciseType = Speaking
						sendNextWord(w, &req, &resp, allStat)
						return
					}
					return
				} else {
					if stat.MaxStepSpeaking < 2 {
						resp.StatRequest.ExerciseType = Speaking
						sendNextWord(w, &req, &resp, allStat)
						return
					}
				}
			}
		}

		stat = Statistic{}
		resp.Tag.TagID = stat.TagID
		resp.Tag.Lang = req.Lang
		resp.Tag.Tag = ""
		resp.StatRequest.ID = req.ID
		resp.StatRequest.UUID = req.UUID
		resp.StatRequest.TagID = stat.TagID
		resp.StatRequest.Lang = req.Lang
		resp.StatRequest.ExerciseType = Listening

		if r.Header.Get("ID") == "1" {
			resp.StatRequest.ExerciseType = Speaking
		}
		sendNextWord(w, &req, &resp, allStat)
	})

	Handler.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}
		if !checkAuth(r) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		}
		if r.Method == http.MethodPost {
			defer r.Body.Close()
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			var request StatUpdRequest = StatUpdRequest{}
			statItemStr := Statistic{}

			err = json.Unmarshal(b, &request)
			statItemStr.TagID = request.TagID
			statItemStr.ID = request.ID
			statItemStr.Lang = request.Lang
			statItemStr.Tag = request.Tag
			if err != nil {
				Resp(w, err.Error())
				return
			}
			err = validate.Struct(request)
			if err == nil {
				db.Update(func(tx *bolt.Tx) error {
					bucket := tx.Bucket([]byte("statistic"))
					statItem := bucket.Get([]byte(fmt.Sprintf("%d_%s_%d", request.ID, request.Lang, request.TagID)))
					if statItem != nil {
						json.Unmarshal(statItem, &statItemStr)
						statItemStr.ID = request.ID
						statItemStr.TagID = request.TagID
						statItemStr.Lang = request.Lang
						statItemStr.Tag = request.Tag
						statItemStr.Update(&request)
						updatedStat, _ := json.Marshal(statItemStr)
						bucket.Put([]byte(fmt.Sprintf("%d_%s_%d", request.ID, request.Lang, request.TagID)), updatedStat)
						Resp(w, string(updatedStat))
						return err
					} else {
						statItemStr.Update(&request)
						updatedStat, _ := json.Marshal(statItemStr)
						bucket.Put([]byte(fmt.Sprintf("%d_%s_%d", request.ID, request.Lang, request.TagID)), updatedStat)
						Resp(w, string(updatedStat))
						return err
					}
				})
				//fmt.Fprintf(w, "success")
			} else {
				Resp(w, err.Error())
			}
		} else {
			fmt.Print("DONT POST\n" + r.Method)
		}
		//path := strings.Split(r.URL.Path, "/")

	})

	/*Handler.HandleFunc("/statistics/", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}

		fmt.Fprintf(w, `{"error": true, "message": "not enought params"}`)

	})*/

	Handler.HandleFunc("/user1", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}
		if !checkAuth(r) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		} else {
			Resp(w, `{"error": ""}`)
			return
		}
	})

	Handler.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if !allowRequest(w, r) {
			return
		}
		exist := checkAuth(r)
		switch r.Method {
		case http.MethodGet:
			var uid uuid.UUID
			var value []byte

			user := User{}
			uid, _ = uuid.NewV4()
			user.Created = time.Now().UnixNano()
			user.LastUsed = time.Now().UnixNano()
			user.UUID = uid.String()
			uid, _ = uuid.NewV4()
			user.RefUUID = uid.String()
			if exist {
				exId, _ := strconv.ParseInt(r.Header.Get("ID"), 0, 64)
				user.ID = int(exId)
				user.UUID = r.Header.Get("UUID")
			}
			err = validate.Struct(user)
			if err != nil {
				log.Fatal(err)
			}
			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("user"))
				if !exist {
					id, _ := b.NextSequence()
					user.ID = int(id)
				}
				value, _ = json.Marshal(user)
				err := b.Put([]byte(strconv.Itoa(user.ID)), value)
				Resp(w, string(value))
				return err
			})
			break
		case http.MethodPost:
			if exist {
				defer r.Body.Close()
				v, _ := ioutil.ReadAll(r.Body)
				u := User{}
				json.Unmarshal(v, &u)
				ok := getUser(u)
				if ok {
					Resp(w, string(v))
				} else {
					newUser := User{}
					json, _ := json.Marshal(newUser)
					Resp(w, string(json))
				}
				return
			} else {
				print("fail")
				user := User{}
				json, _ := json.Marshal(user)

				Resp(w, string(json))
				return
			}
		}

		fmt.Println(r.URL.Path)
	})

	Handler.HandleFunc("/sync/", func(w http.ResponseWriter, r *http.Request) {
		pathStrings := strings.Split(r.URL.Path, "/")
		timestamp, _ := strconv.ParseInt(pathStrings[2], 0, 64)
		step, _ := strconv.ParseInt(pathStrings[3], 0, 64)

		if !allowRequest(w, r) {
			return
		}
		if !checkAuth(r) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		} else {
			stat := make([]Statistic, 0)

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("statistic"))

				c := b.Cursor()

				prefix := []byte(fmt.Sprintf("%s", r.Header.Get("ID")))
				var statItem Statistic
				for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
					json.Unmarshal(v, &statItem)
					if timestamp < statItem.Updated {
						stat = append(stat, statItem)
					}
				}
				return nil
			})
			sort.Slice(stat, func(i, j int) bool {
				return stat[i].Updated < stat[j].Updated
			})
			var (
				statStr    []byte
				stepLength int64 = 160
			)

			var start, end int64
			start = (step - 1) * stepLength
			end = step * stepLength
			if end > int64(len(stat)) {
				end = int64(len(stat))
			}
			statStr, err = json.Marshal(stat[start:end])
			Resp(w, string(statStr))
			return
		}
	})

	Handler.HandleFunc("/backup", func(w http.ResponseWriter, req *http.Request) {
		if !allowRequest(w, req) {
			return
		}
		if !checkAuth(req) {
			Resp(w, `{"error": "Auth failed"}`)
			return
		}
		if req.Header.Get("ID") == "1" {
			err := db.View(func(tx *bolt.Tx) error {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
				w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
				_, err := tx.WriteTo(w)
				return err
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			Resp(w, `{"error": "Permission denied"}`)
			return
		}
	})

	Handler.HandleFunc("/results/", func(w http.ResponseWriter, req *http.Request) {
		pathAr := strings.Split(req.URL.Path, "/")
		if len(pathAr) != 3 {
			Resp(w, "0.0.1")
			return
		}
		return

	})

	Handler.HandleFunc("/send-error", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			defer r.Body.Close()
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				Resp(w, err.Error())
			}

			f, e := os.OpenFile(STATIC_PATH+"/client_errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
			defer f.Close()
			if e != nil {
				Resp(w, e.Error())
			}
			_, e = f.Write(b)
			if e != nil {
				Resp(w, e.Error())
			}
			Resp(w, string(b))
		}
	})

	Handler.HandleFunc("/getall", func(w http.ResponseWriter, r *http.Request) {
		b, e := json.Marshal(state)
		if e != nil {
			print("error")
		}
		Resp(w, string(b))
	})

	Handler.HandleFunc("/set/", func(w http.ResponseWriter, r *http.Request) {
		pathAr := strings.Split(r.URL.Path, "/")
		// string to int
		i, err := strconv.Atoi(pathAr[2])
		if err != nil {
			// ... handle error
			panic(err)
		}
		if "true" == pathAr[3] {
			state[i] = true
		} else {
			state[i] = false
		}

		requestURL := fmt.Sprintf("http://%s/set?pin=%d&value=%t", ESP8266HOST, i, state[i])
		res, e := http.Get(requestURL)
		if e != nil {
			print("error")
		}

		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})

	Handler.HandleFunc("/setall", func(w http.ResponseWriter, r *http.Request) {

		jsonBody, _ := json.Marshal(state)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%s/%s", ESP8266HOST, "setall")
		req, e := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if e != nil {
			print("error")
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		res, _ := client.Do(req)
		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})

	Handler.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {

		jsonBody, _ := json.Marshal(state)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%s/%s", ESP8266HOST, "setall")
		req, e := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if e != nil {
			print("error")
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		res, _ := client.Do(req)
		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})

	Handler.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {

		jsonBody, _ := json.Marshal(state)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%s/%s", ESP8266HOST, "setall")
		req, e := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if e != nil {
			print("error")
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		res, _ := client.Do(req)
		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})

	Handler.HandleFunc("/signin1", func(w http.ResponseWriter, r *http.Request) {

		jsonBody, _ := json.Marshal(state)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%s/%s", ESP8266HOST, "setall")
		req, e := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if e != nil {
			print("error")
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		res, _ := client.Do(req)
		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})

	Handler.HandleFunc("/signout1", func(w http.ResponseWriter, r *http.Request) {

		jsonBody, _ := json.Marshal(state)
		bodyReader := bytes.NewReader(jsonBody)

		requestURL := fmt.Sprintf("http://%s/%s", ESP8266HOST, "setall")
		req, e := http.NewRequest(http.MethodPost, requestURL, bodyReader)
		if e != nil {
			print("error")
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}
		res, _ := client.Do(req)
		b, _ := io.ReadAll(res.Body)
		Resp(w, string(b))
	})
}

func AutoLearn(tag Tag, id string, tx *bolt.Tx) {

	bucket := tx.Bucket([]byte("statistic"))

	var request StatUpdRequest

	idN, _ := strconv.ParseInt(id, 10, 64)

	request.ID = int(idN)
	request.TagID = tag.TagID
	request.Lang = tag.Lang
	request.Correct = true
	request.ExerciseType = Speaking
	request.Tag = tag.Tag
	statItem := bucket.Get([]byte(fmt.Sprintf("%s_%s_%d", id, request.Lang, request.TagID)))

	var statItemStr Statistic
	json.Unmarshal(statItem, &statItemStr)
	if statItemStr.MaxStepSpeaking < 2 {
		statItemStr.Update(&request)
	}
	if statItemStr.MaxStepSpeaking < 2 {
		statItemStr.Update(&request)
	}
	if statItemStr.MaxStepSpeaking < 2 {
		statItemStr.Update(&request)
	}

	updatedStat, _ := json.Marshal(statItemStr)
	bucket.Put([]byte(fmt.Sprintf("%s_%s_%d", id, request.Lang, request.TagID)), updatedStat)
}

func allowRequest(w http.ResponseWriter, r *http.Request) bool {
	now := time.Now()
	ns := now.UnixNano()
	nsec := ns + SECOND
	if strings.Contains(r.URL.Path, "/sync") {
		nsec = ns + SECOND
	} else if strings.Contains(r.URL.Path, "/statistics") {
		nsec = ns + SECOND/10
	} else if strings.Contains(r.URL.Path, "/next-word") {
		nsec = ns + SECOND/10
	} else if strings.Contains(r.URL.Path, "/backup") {
		nsec = ns + SECOND*100
	} else if strings.Contains(r.URL.Path, "/users") {
		nsec = ns + SECOND*10
	}
	ips := r.Header.Get("X-FORWARDED-FOR")
	var clientIP string
	if ips != "" {
		clientIP = strings.Split(ips, ", ")[0]
	} else {
		clientIP = strings.Split(r.RemoteAddr, ":")[0]
	}

	if _, ok := lastSecondIps[clientIP+r.URL.Path]; ok {
		lastSecondIps[clientIP+r.URL.Path] = nsec
		Resp(w, `{ "error": true, "message":"Too many requests"}`)
		return false
	}
	lastSecondIps[clientIP+r.URL.Path] = nsec
	return true
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func runJobs() {
	ticker := time.NewTicker(10 * time.Second)
	cTicker := time.NewTicker(time.Second / 2)

	//current, err := user.Current()

	f, e := os.OpenFile(STATIC_PATH+"/server_monitor.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	defer f.Close()
	if e != nil {
		log.Fatal(e.Error())
	}
	if true {
		var idle0, total0, idle1, total1 uint64
		for {
			idle0, total0 = getCPUSample()

			select {
			case <-ticker.C:
				idle1, total1 = getCPUSample()
				idleTicks := float64(idle1 - idle0)
				totalTicks := float64(total1 - total0)
				cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				disk := DiskUsage("/")
				fmt.Fprintf(
					f,
					"Time: %s.All: %.2f GB, Used: %.2f GB/Free: %.2f GB. Alloc: %.2f GB, Sys: %.2f GB. CPU usage is %f%% [busy: %f, total: %f]\n",
					time.Now().String(),
					float64(disk.All)/float64(GB),
					float64(disk.Used)/float64(GB),
					float64(disk.Free)/float64(GB),
					float64(m.Alloc)/float64(GB),
					float64(m.Sys)/float64(GB),
					cpuUsage,
					totalTicks-idleTicks,
					totalTicks,
				)
				break
			case <-cTicker.C:
				now := time.Now().UnixNano()
				for k, v := range lastSecondIps {
					if v < now {
						delete(lastSecondIps, k)
					}
				}
				break

			}
		}
	}

	//ticker2 := time.NewTicker(24 * time.Hour)
	if false {
		now := time.Now().UnixNano()
		if now > backupTimestamp {
			var newDb *bolt.DB
			os.Remove("newDb")
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in s", r)
				}
			}()
			var e error
			if os.Getenv("ENV") != "" {
				e = DownloadFile("newDb", "https://test.gofsd.xyz/backup")
			} else {
				e = DownloadFile("newDb", "https://gofsd.xyz/backup")
			}
			if e != nil {
				log.Fatalf("Bolt.open: %s", e)
			}
			newDb, err = bolt.Open("newDb", 0777, nil)
			if err != nil {
				log.Fatalf("Bolt.open: %s", err)
			}
			db.Close()
			newDb.Close()
			os.Remove("my.db")
			os.Rename("newDb", os.Getenv("path_to_db"))
			os.Remove("newDb")
			newDb, err = bolt.Open(os.Getenv("path_to_db"), 0777, nil)
			if err != nil {
				log.Fatalf("Bolt.open: %s", err)
			}
			db = newDb
			newDb = nil
		}
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	/*
			if os.Getenv("ENV") == "" {
			ticker2 = time.NewTicker(2 * time.Minute)
		}
		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()
				for k, v := range lastSecondIps {
					if v < now {
						delete(lastSecondIps, k)
					}
				}
				break
			case <-ticker2.C:
				if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == " " {
					now := time.Now().UnixNano()
					if now > backupTimestamp {
						var newDb *bolt.DB
						os.Remove("newDb")
						defer func() {
							if r := recover(); r != nil {
								fmt.Println("Recovered in s", r)
							}
						}()
						var e error
						if os.Getenv("ENV") == "" {
							e = DownloadFile("newDb", "https://test.gofsd.xyz/backup")
						} else {
							e = DownloadFile("newDb", "https://gofsd.xyz/backup")
						}
						if e != nil {
							log.Fatalf("Bolt.open: %s", e)
						}
						newDb, err = bolt.Open("newDb", 0777, nil)
						if err != nil {
							log.Fatalf("Bolt.open: %s", err)
						}
						db.Close()
						newDb.Close()
						os.Remove("my.db")
						os.Rename("newDb", os.Getenv("path_to_db"))
						os.Remove("newDb")
						newDb, err = bolt.Open(os.Getenv("path_to_db"), 0777, nil)
						if err != nil {
							log.Fatalf("Bolt.open: %s", err)
						}
						db = newDb
						newDb = nil
					}
					break
				}
			}
		}*/

}

func sendNextWord(w http.ResponseWriter, reqData *StatUpdRequest, nextWord *NextWordResp, allStat []Statistic) {
	var col interface{}
	var tag Tag
	col, _ = getCollection(reqData.Lang)
	en, _ := getCollection("en_US")
	length := len(en.([]interface{}))
	lengthCol := len(col.([]interface{}))

	if nextWord.Tag.TagID != 0 {
		db.View(func(tx *bolt.Tx) (err error) {
			var bk *bolt.Bucket
			bk = tx.Bucket([]byte("tag"))
			var tagB []byte = bk.Get([]byte(fmt.Sprintf("%d_%s_%d", nextWord.StatRequest.ID, nextWord.StatRequest.Lang, nextWord.StatRequest.TagID)))

			if tagB == nil {
				tagB = bk.Get([]byte(fmt.Sprintf("%d_%s_%d", 1, nextWord.StatRequest.Lang, nextWord.StatRequest.TagID)))
			}
			if tagB != nil {
				err = json.Unmarshal(tagB, &tag)
			}
			return err
		})
		if tag.TagID == nextWord.StatRequest.TagID {
			nextWord.Tag = tag
			goto END

		}
		for i := 0; i < length; i++ {
			item := en.([]interface{})[i]
			id := int64(item.(map[string]interface{})["frequency"].(float64))
			if id == nextWord.Tag.TagID {
				nextWord.StatRequest.TagID = id
				nextWord.StatRequest.Tag = item.(map[string]interface{})["name"].(string)
				nextWord.Tag.Tag = item.(map[string]interface{})["name"].(string)
				for k := 0; k < lengthCol; k++ {
					if strings.ToLower(nextWord.Tag.Tag) == strings.ToLower(col.([]interface{})[k].(map[string]interface{})["normalizedSource"].(string)) {
						countWords := len(col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{}))
						for w := 0; countWords > w; w++ {
							if w > 14 {
								break
							}
							nextWord.Tag.Translations = append(nextWord.Tag.Translations, col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{})[w].(map[string]interface{})["displayTarget"].(string))
						}
					}
				}
			}
		}
	END:
	} else {

		var newEn []interface{}
		var item interface{}
		for i := 0; i < length; i++ {
			cont := false
			item := en.([]interface{})[i]
			for i := 0; i < len(tags); i++ {
				if int64(item.(map[string]interface{})["frequency"].(float64)) == tags[i].TagID {
					if tags[i].Disabled == true {
						cont = true
					}
					break
				}
			}
			if cont {
				cont = false
				continue
			}
			newEn = append(newEn, item)
		}

		lenStat := len(allStat)
		newLenStat := 0
		canSkip := false
		exitCounter := 0
		print(lenStat, newLenStat)
		for {
			item = newEn[newLenStat]
			for i := 0; i < lenStat; i++ {
				if int64(item.(map[string]interface{})["frequency"].(float64)) == allStat[i].TagID {
					newLenStat++
					canSkip = false
					break
				} else {
					canSkip = true
				}
			}

			if canSkip {
				break
			}
			exitCounter++

			if lenStat == 0 || exitCounter > lenStat {
				break
			}
		}

		id := int64(item.(map[string]interface{})["frequency"].(float64))
		nextWord.StatRequest.TagID = id
		nextWord.StatRequest.Tag = item.(map[string]interface{})["name"].(string)

		db.View(func(tx *bolt.Tx) (err error) {
			var bk *bolt.Bucket
			bk = tx.Bucket([]byte("tag"))
			var tagB []byte = bk.Get([]byte(fmt.Sprintf("%d_%s_%d", nextWord.StatRequest.ID, nextWord.StatRequest.Lang, nextWord.StatRequest.TagID)))

			if tagB == nil {
				tagB = bk.Get([]byte(fmt.Sprintf("%d_%s_%d", 1, nextWord.StatRequest.Lang, nextWord.StatRequest.TagID)))
			}
			if tagB != nil {
				err = json.Unmarshal(tagB, &tag)
			}
			return err
		})

		if tag.TagID == nextWord.StatRequest.TagID {
			nextWord.Tag = tag
		}

		nextWord.Tag.TagID = id
		nextWord.Tag.Tag = item.(map[string]interface{})["name"].(string)
		for k := 0; k < lengthCol; k++ {
			if strings.ToLower(nextWord.Tag.Tag) == strings.ToLower(col.([]interface{})[k].(map[string]interface{})["normalizedSource"].(string)) {
				countWords := len(col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{}))
				for w := 0; countWords > w; w++ {
					if w > 14 {
						break
					}
					nextWord.Tag.Translations = append(nextWord.Tag.Translations, col.([]interface{})[k].(map[string]interface{})["translations"].([]interface{})[w].(map[string]interface{})["displayTarget"].(string))
				}
			}
		}
	}

	resp, _ := json.Marshal(nextWord)
	respStr := string(resp)
	print(fmt.Sprintf("\n%d %s", len(allStat), respStr))

	Resp(w, respStr)
	return
}

func getCollection(lang string) (col interface{}, err error) {
	var data []byte
	col = collections[lang]

	if col != nil {
		return
	}
	path, _ := filepath.Abs(fmt.Sprintf("lookup/%s.json", lang))
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	} else {
		json.Unmarshal(data, &col)
		if lang == "en_US" {
			sort.Slice(col, func(i int, j int) bool {
				return col.([]interface{})[i].(map[string]interface{})["frequency"].(float64) < col.([]interface{})[j].(map[string]interface{})["frequency"].(float64)
			})
		}
		collections[lang] = col
		return
	}
}

func checkAuth(r *http.Request) bool {

	id := r.Header.Get("ID")
	uuid := r.Header.Get("UUID")
	ok := false
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(k) == id {
				user := User{}
				json.Unmarshal(v, &user)
				if user.UUID == uuid {
					ok = true
				}
			}
		}
		return nil
	})
	return ok
}

func getUser(u User) bool {
	id := u.ID
	uuid := u.UUID
	ok := false
	var user User
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			if string(k) == fmt.Sprint(id) {

				user = User{}
				json.Unmarshal(v, &user)
				if user.UUID == uuid {
					ok = true
				}
			}
		}
		return nil
	})
	return ok
}

func DownloadFile(filepath string, url string) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in d", r)
		}
	}()
	// Get the data
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("ID", "1")
	req.Header.Set("UUID", "d0b35f00-7372-4dec-873d-65df4063d3c2")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func Resp(w http.ResponseWriter, pattern string, args ...interface{}) {
	w.Header().Set("CLIENT_APP_VERSION", CLIENT_APP_VERSION)
	w.Header().Set("ADS_PERIOD", "60")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,DELETE")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, pattern)

}

func SetPin(pin int, value bool) {
	requestURL := fmt.Sprintf("http://%s/set?pin=%d&value=%t", ESP8266HOST, pin, value)
	http.Get(requestURL)
	//state[pin] = value
}

func UpdateAll() {
	requestURL := fmt.Sprintf("http://%s/getall", ESP8266HOST)
	r, _ := http.Get(requestURL)
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, currentIotState)
}
