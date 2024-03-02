package http

import (
	"bytes"
	"encoding/json"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

//HttpTest
type HttpTest struct {
	Host string
	TestCases []EndpointTestData
	T *testing.T
}

func (t HttpTest) Run(){
	if t.T == nil {
		return
	}
	for _, tt := range t.TestCases {
		tt.Host = t.Host
		tt.T = t.T
		t.T.Run(tt.Name, func(t *testing.T) {
			tt.Execute()
		})
	}
}

//Request test params for making request
type Request struct {
	URL         string
	QueryParams map[string]string
	Headers     map[string][]string
	MethodType  string
	Body []byte
}
//Response
type Response struct {
	Headers map[string][]string
	Body []byte
}
//ExpectedResponse
type ExpectedResponse struct {
	Headers map[string][]string
	Body []byte
}
//EndpointTestData all data what you need to make http request and make test
type EndpointTestData struct {
	Host string
	Name          string
	Request
	Response
	ExpectedResponse
	Err           error
	T *testing.T
}

//Execute execute test
func (testItem *EndpointTestData) Execute() {
	testItem.MakeRequest()
	if !testItem.JSONBytesEqual() {
		testItem.T.Errorf("(%s): expected %s, actual %s", testItem.Response.Body, testItem.ExpectedResponse.Body, testItem.Response.Body)
	}
	testItem.CompareHeaders()
}

func (testItem *EndpointTestData) CompareHeaders() {
	for key, _ := range testItem.ExpectedResponse.Headers{
		if testItem.Response.Headers[key] != nil {
			expected := len(testItem.ExpectedResponse.Headers[key])
			actual := len(funk.IntersectString(testItem.Response.Headers[key], testItem.ExpectedResponse.Headers[key]))
			if expected != actual {
				testItem.T.Errorf("(%d): expected intersect length %d, actual intersect length %d", actual, expected, actual)
			}
		} else {
			testItem.T.Errorf("(%s): not found",key)
		}
	}
}

//MakeRequest make req
func (testItem *EndpointTestData) MakeRequest() {

	switch testItem.Request.MethodType {
	case http.MethodGet:
		testItem.GetRequest()
		break
	case http.MethodPost:
		testItem.PostRequest()
		break
	case http.MethodPatch:
		testItem.PatchRequest()
		break
	case http.MethodPut:
		testItem.PutRequest()
		break
	case http.MethodDelete:
		testItem.DeleteRequest()
		break
	}
}

func (testItem *EndpointTestData) GetRequest() {
	var (
		err  error
		resp *http.Response
	)

	resp, testItem.Err = http.Get(testItem.Host + testItem.Request.URL)
	if testItem.Err != nil {
		log.Fatal("test data: ", err)
	}
	defer resp.Body.Close()
	testItem.Response.Headers = resp.Header

	testItem.Response.Body, testItem.Err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func (testItem *EndpointTestData) PostRequest() {
	var (
		err  error
		resp *http.Response
	)
	resp, err = http.Post(testItem.Host+testItem.Request.URL, "application/json", bytes.NewBuffer(testItem.Request.Body))
	if testItem.Err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	testItem.Response.Body, testItem.Err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
func (testItem *EndpointTestData) PutRequest() {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		client http.Client = http.Client{}
	)
	req, err = http.NewRequest(testItem.Request.MethodType, testItem.Host+testItem.Request.URL, bytes.NewBuffer(testItem.Request.Body))
	req.Header.Set("Content-Type", "application/json")
	if testItem.Err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	testItem.Response.Body, testItem.Err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
func (testItem *EndpointTestData) PatchRequest() {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		client http.Client = http.Client{}
	)
	req, err = http.NewRequest(testItem.Request.MethodType, testItem.Host+testItem.Request.URL, bytes.NewBuffer(testItem.Request.Body))
	req.Header.Set("Content-Type", "application/json")

	if testItem.Err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	testItem.Response.Body, testItem.Err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func (testItem *EndpointTestData) DeleteRequest() {
	var (
		err    error
		req    *http.Request
		resp   *http.Response
		client http.Client = http.Client{}
	)
	req, err = http.NewRequest(testItem.Request.MethodType, testItem.Host+testItem.Request.URL, bytes.NewBuffer(testItem.Request.Body))
	if testItem.Err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	resp, err = client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	testItem.Response.Body, testItem.Err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func (testItem *EndpointTestData) JSONBytesEqual() bool {
	var j, j2 interface{}
	if err := json.Unmarshal(testItem.Response.Body, &j); err != nil {
		return false
	}
	if err := json.Unmarshal(testItem.ExpectedResponse.Body, &j2); err != nil {
		return false
	}
	return reflect.DeepEqual(j2, j)
}
