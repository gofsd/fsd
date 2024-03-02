package bolt

import (
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gofsd/fsd/pkg/util"
)

var (
	ErrorHandler func(error)
)

// Bolt structure with custom fields and methods
type Bolt struct {
	db *bolt.DB
}

const (
	EXTENSION   string      = ".bolt"
	RW_FOR_ALL  os.FileMode = 0666
	RW_FOR_USER os.FileMode = 0600
)

// Create or open bolt db file by absolute file name (absolute path + file name)
func CreateOrOpenDB(fileName string) *Bolt {
	var (
		b *Bolt
		e error
	)
	b = &Bolt{}
	b.db, e = bolt.Open(fileName+EXTENSION, RW_FOR_ALL, &bolt.Options{Timeout: time.Second})
	util.HandleError(e)
	return b
}

// Create if not exist and return bucket
func (b *Bolt) GetBucket(bucketFullName []byte) (*bolt.Bucket, *bolt.Tx) {
	var (
		tx     *bolt.Tx
		bucket *bolt.Bucket
		e      error
	)
	tx, e = b.db.Begin(true)
	util.HandleError(e)
	for idx, bBucketName := range bucketFullName {
		if idx == 0 {
			bucket, e = tx.CreateBucketIfNotExists([]byte{bBucketName})
		} else {
			bucket, e = bucket.CreateBucketIfNotExists([]byte{bBucketName})
		}
	}
	util.HandleError(e)

	return bucket, tx
}

func (b *Bolt) Set(bucketFullName []byte, key, value []byte) *bolt.Bucket {
	var (
		bucket *bolt.Bucket
		e      error
		tx     *bolt.Tx
	)
	bucket, tx = b.GetBucket(bucketFullName)
	util.HandleError(e)
	defer tx.Rollback()

	e = bucket.Put(key, value)
	util.HandleError(e)
	e = tx.Commit()
	util.HandleError(e)
	return bucket
}

func (b *Bolt) Get(bucketFullName []byte, key []byte) []byte {
	var (
		bucket *bolt.Bucket
		e      error
		data   []byte
		tx     *bolt.Tx
	)
	bucket, tx = b.GetBucket(bucketFullName)
	data = bucket.Get(key)
	util.HandleError(e)
	defer tx.Rollback()
	tx.Commit()
	return data
}

func (b *Bolt) IsOpen() bool {
	return b.db.AllocSize > 0
}

func (b *Bolt) Path() string {
	return b.db.Path()
}

func (b *Bolt) Close() {
	b.db.Close()
}

func (b *Bolt) GetPageSize() int {
	return b.db.Info().PageSize
}

func (b *Bolt) GetSize() int64 {
	defer func() {
		b.db.Close()
		Buckets("/home/madi/.bolt/my.db.bolt")

	}()
	return util.FileSize(b.Path())
}

func Buckets(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		print(err)
		return
	}

	db, err := bolt.Open(path, RW_FOR_ALL, nil)
	if err != nil {
		print(err)
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			print(string(name))
			return nil
		})
	})
	if err != nil {
		print(err)
		return
	}
}
