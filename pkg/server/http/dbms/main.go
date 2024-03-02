package dbms

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type collectionMetadata struct {
	fileInfo       *os.FileInfo
	fileDescriptor *os.File
	fileSectorSize float64
	fileSize       int
	buff           []byte
	fileName       string
	defaultObject  interface{}
}

// DbMetadata metadata for db
type DbMetadata struct {
	collections       map[string]*collectionMetadata
	name              string
	fullName          string
	path              string
	rootStructureSize uintptr
	rootStructure     interface{}
	MuxHandler        *http.ServeMux
	handlers          map[string]func(w http.ResponseWriter, r *http.Request)
}

// File type
type File struct {
	ID   uint
	Name string
	Mime string
	Size uint
}

var files = make(map[string]*collectionMetadata)

func (db *DbMetadata) saveFile(part *multipart.Part) (err error) {
	var dst *os.File
	dst, err = os.Create("./" + db.name + "/files/" + part.FileName())
	defer dst.Close()
	if err != nil {
		return
	}

	if _, err = io.Copy(dst, part); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// CreateDb by
func CreateDb(rootStructure interface{}) (db *DbMetadata, err error) {
	db = new(DbMetadata)
	db.rootStructure = rootStructure
	db.setDbMetadata()
	db.setCollectionsMetadata()
	db.createCollections()
	db.setHandlers()
	db.createServer()

	return db, nil
}

func (db *DbMetadata) createServer() {
	db.MuxHandler = http.NewServeMux()
	for k := range db.collections {
		db.MuxHandler.HandleFunc("/"+db.name+"/"+k+"/", db.handlers["defaultHandler"])
	}
}

func (db *DbMetadata) setDbMetadata() {
	db.name = strings.ToLower(strings.Split(reflect.TypeOf(db.rootStructure).String(), ".")[1])
	db.rootStructureSize = reflect.TypeOf(db.rootStructure).Size()
	db.handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal()
	}
	if db.path == "" {
		db.path = dir + "/"
	}

	db.fullName = db.path + db.name + "/"
}

func (db *DbMetadata) setCollectionsMetadata() {
	NumFields := reflect.ValueOf(db.rootStructure).Elem().NumField()
	db.collections = make(map[string]*collectionMetadata)
	for i := 0; i < NumFields; i++ {
		field := reflect.ValueOf(db.rootStructure).Elem().Field(i)
		colName := strings.ToLower(field.Type().Name())
		db.collections[colName] = &collectionMetadata{
			fileSectorSize: float64(1000),
			buff:           make([]byte, 1000),
			fileName:       db.fullName + colName,
			defaultObject:  reflect.New(field.Type()),
		}
	}
}

func (db *DbMetadata) createCollections() (created bool, err error) {
	var collectionFullName string
	if !dirExists(db.fullName) {
		mkdirErr := os.Mkdir(db.fullName, 0755)
		os.MkdirAll(db.fullName+"files", 0755)
		if mkdirErr != nil {
			log.Fatal("create db folder: ", mkdirErr)
		}
	}

	for k := range db.collections {
		collectionFullName = db.fullName + k
		if fileExists(collectionFullName) {
			fmt.Println("Example file exists")
			fileDescriptor, err := os.OpenFile(collectionFullName, os.O_RDWR, 0755)
			if err != nil {
				log.Fatal("from open db", err)
			} else {
				db.collections[k].fileDescriptor = fileDescriptor
				if fi, err := fileDescriptor.Stat(); err == nil {
					db.collections[k].fileInfo = &fi
				}
			}
		} else {
			fileDescriptor, err := os.Create(collectionFullName)
			if err != nil {
				log.Fatal("from create db", err)
			} else {
				db.collections[k].fileDescriptor = fileDescriptor
				if fi, err := fileDescriptor.Stat(); err == nil {
					db.collections[k].fileInfo = &fi
				}
			}

		}
	}

	return created, err
}

// CreateDB file
func CreateDB() {
	var fileName string = "test.db"
	var fileSectorSize float64 = 100
	if fileExists(fileName) {
		fmt.Println("Example file exists")
		fileDescriptor, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal("from open db", err)
		} else {
			files[fileName] = &collectionMetadata{fileDescriptor: fileDescriptor, fileSectorSize: fileSectorSize}
		}
	} else {
		fileDescriptor, err := os.Create(fileName)
		if err != nil {
			log.Fatal("from create db", err)
		} else {
			files[fileName] = &collectionMetadata{fileDescriptor: fileDescriptor, fileSectorSize: fileSectorSize}
		}

	}
}

func (collection *collectionMetadata) writeAt(data []byte, off int64) (n int, err error) {
	n, err = collection.fileDescriptor.WriteAt(data, off)
	if err != nil {
		log.Fatal("From writeAt: ", err)
	}
	return
}

func (collection *collectionMetadata) readAt(off int64) (n int, err error) {
	n, err = collection.fileDescriptor.ReadAt(collection.buff, off)
	//print("from 139", off, n, err, "coll:", len(collection.buff), "\n")
	if err != nil {
		if err == io.EOF {
			//print("EOF message")
		} else {
			log.Fatal("From readAt: ", err)
		}
	}
	return
}

// Insert to file
func (collection *collectionMetadata) Insert(buf []byte) (n int, err error) {
	len := len(buf)
	len = len + 8
	fi, err := collection.fileDescriptor.Stat()
	if err != nil {
		log.Fatal(err)
	}
	n = int(math.Ceil(float64(float64(fi.Size())/collection.fileSectorSize)) + 1)
	var writeTo int64 = int64(float64(n-1) * collection.fileSectorSize)
	print("from insert", writeTo, fi.Size())

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(len))
	fmt.Println("from insert ", string(buf))
	i := int64(binary.LittleEndian.Uint64(b))
	fmt.Println("from select %V", i)
	buf = append(b, buf...)
	collection.writeAt(buf, writeTo)
	return
}

// Select sector id
func (collection *collectionMetadata) Select(id uint) (data []byte, err error) {
	fmt.Println(collection, "test:test")
	sectorSize := int64(collection.fileSectorSize)
	data = make([]byte, sectorSize)
	collection.readAt(int64(int64(id) * sectorSize))
	totalLen := collection.buff[0:8]
	i := int64(binary.LittleEndian.Uint64(totalLen))
	data = collection.buff[8:i]
	return
}

// Update by sector id
func (collection *collectionMetadata) Update(id int, buf []byte) (n int, err error) {
	if err != nil {
		log.Fatal(err)
	}

	collection.writeAt(buf, int64(id*100))
	return
}

// Delete by sector id
func Delete(id int) (n int, err error) {

	return
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func (db *DbMetadata) setHandlers() {

	db.handlers["defaultHandler"] = func(w http.ResponseWriter, r *http.Request) {
		db.handlers["default"+r.Method](w, r)
	}

	db.handlers["defaultGET"] = func(w http.ResponseWriter, r *http.Request) {
		strs := strings.Split(r.URL.Path, "/")
		id, _ := strconv.ParseUint(strs[3], 10, 64)

		w.WriteHeader(http.StatusOK)
		val, _ := db.collections[strs[2]].Select(uint(id))
		w.Write(val)

	}

	db.handlers["defaultPOST"] = func(w http.ResponseWriter, r *http.Request) {

		if len(db.collections) > 1 {
			db.handlers["multipartFormData"](w, r)
		} else {
			db.handlers["simplePost"](w, r)
		}
	}

	db.handlers["multipartFormData"] = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
			db.saveFile(part)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		fmt.Fprintf(w, `{"foo":"bar"}`)
	}

	db.handlers["simplePost"] = func(w http.ResponseWriter, r *http.Request) {
		strs := strings.Split(r.URL.Path, "/")
		w.WriteHeader(http.StatusOK)
		b := make([]byte, r.ContentLength)
		defer r.Body.Close()
		r.Body.Read(b)
		db.collections[strs[2]].Insert(b)
		w.Write(b)
	}
}
