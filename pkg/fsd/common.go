package fsd

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Validate(ctx *gin.Context, value any) (err error) {
	if err = ctx.ShouldBindJSON(value); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return err
	} else {
		return nil
	}
}

func HandlePanic() {
	if err := recover(); err != nil {
		log.Println("panic occurred:", err)
	}
}

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       struct{ Email, Password string }
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		HostURL: HostURL,
		Auth: struct{ Email, Password string }{
			// Username: *username,
			Password: *password,
		},
	}

	if host != nil {
		c.HostURL = *host
	}

	// ar, err := c.SignIn()
	// if err != nil {
	// 	return nil, err
	// }

	// c.Token = ar.Token

	return &c, nil
}

func doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "c.Token")

	// res, err := HTTPClient.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// if res.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	// }

	// return body, err
	return nil, nil
}
