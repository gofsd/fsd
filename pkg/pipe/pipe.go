package pipe

import (
	"bytes"
	"fmt"
	i "io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"sync"

	"io/ioutil"

	"github.com/boltdb/bolt"
	"github.com/gofsd/fsd/pkg/io"
	"github.com/gofsd/fsd/pkg/types"
	"github.com/gofsd/fsd/pkg/util"
)

var (
	db      *bolt.DB
	handler *http.ServeMux
)

func MoveTop20kWordListToDB(s, dbPath string) {
	var db *bolt.DB
	db, _ = bolt.Open(dbPath, 0777, nil)
	var v = io.S[*types.TopWordsByGoogle, *types.TopWordsByGoogle]{
		Path: "top_words_by_google",
		DB:   db,
	}
	v.ReadFile(s).
		Split("\n").
		ForEachFromSplit(func(it string, idx int) *types.TopWordsByGoogle {
			return &types.TopWordsByGoogle{
				ID:   idx,
				Word: it,
			}
		}).
		FilterMap(func(str *types.TopWordsByGoogle, idx int) (*types.TopWordsByGoogle, bool) {
			return str, true
		}).
		Put()
}

func SelectTopWordsByGoogle(dbPath string) {
	var db *bolt.DB
	db, _ = bolt.Open(dbPath, 0777, nil)
	v := io.S[*types.TopWordsByGoogle, *types.TopWordsByGoogle]{
		T:    &types.TopWordsByGoogle{},
		R:    &types.TopWordsByGoogle{},
		Path: "top_words_by_google",
		DB:   db,
	}
	v.Select(func(item *types.TopWordsByGoogle, idx int) (*types.TopWordsByGoogle, bool) {
		str, e := item.Marshal()
		if e != nil {
			util.HandleError(e)
		}
		print(string(idx) + string(str) + item.Word + "\n")
		return item, false
	})
}

func Buckets(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println(string(name))
			b.Stats()
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SetDB(db *bolt.DB) {

}

// adapt HTTP connection to ReadWriteCloser
type HttpConn struct {
	in  i.Reader
	out i.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

type SC struct {
	hc    *HttpConn
	codec rpc.ServerCodec
}

func SetHTTP(h *http.ServeMux) {

	arith := new(Arith)
	s := rpc.NewServer()
	err := s.Register(arith)
	if err != nil {
		log.Fatalf("Format of service Arith isn't correct. %s", err)
	}

	var pool = sync.Pool{
		New: func() interface{} {
			var s SC
			s.hc = &HttpConn{}
			s.codec = jsonrpc.NewServerCodec(s.hc)
			return s
		},
	}

	h.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		sc := pool.Get().(SC)
		defer pool.Put(sc)
		sc.hc.in = r.Body
		sc.hc.out = w
		err := s.ServeRequest(sc.codec)
		if err != nil {
			log.Printf("Error while serving JSON request: %v", err)
			http.Error(w, "Error while serving JSON request, details have been logged.", 500)
			return
		}

	})
	handler = h

}

// Holds arguments to be passed to service Arith in RPC call
type Args struct {
	A, B int
}

// Representss service Arith with method Multiply
type Arith int

// Result of RPC call is of this type
type Result int

// This procedure is invoked by rpc and calls rpcexample.Multiply which stores product of args.A and args.B in result pointer
func (t *Arith) Multiply(args Args, result *Result) error {
	return Multiply(args, result)
}

// stores product of args.A and args.B in result pointer
func Multiply(args Args, result *Result) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}
func serveJSONRPC(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		http.Error(w, "method must be connect", 405)
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}
	defer conn.Close()
	i.WriteString(conn, "HTTP/1.0 Connected\r\n\r\n")
	jsonrpc.ServeConn(conn)
}
func Test(s string) {
	print(s + "from test")
	resp, err := http.Post("http://localhost:3000/rpc", "application/json", bytes.NewBufferString(
		`{"jsonrpc":"2.0","id":1,"method":"Arith.Multiply","params": [{"A": 1, "B": 2}]}, {"jsonrpc":"2.0","id":2,"method":"Arith.Multiply","params": [{"A": 2, "B": 3}]}`,
	))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("returned JSON: %s\n", string(b))
}
