package dummy

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Username  string
	Password  string
	Debug     bool
	UserAgent string
	BaseUrl   string
	UrlPath   string
}

type Client struct {
	config     Config
	HttpClient *http.Client
	Token      string
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func encodeToken(config Config) string {
	credentials := config.Username + ":" + config.Password
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	token := "Basic " + encoded
	return token
}

func NewClient(config Config) *Client {
	token := encodeToken(config)
	return &Client{
		config:     config,
		HttpClient: http.DefaultClient,
		Token:      token,
	}
}

// doJsonRequest is the simplest type of request: a method on a URI that returns
// some JSON result which we unmarshal into the passed interface.
func (client *Client) doJsonRequest(method, api string,
	reqBody map[string]interface{}) ([]byte, error) {
	req, err := client.createRequest(method, api, reqBody)
	if err != nil {
		return nil, err
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var resp *http.Response
	resp, err = client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("API error %s: %s", resp.Status, body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// If we got no body, by default let's just make an empty JSON dict. This
	// saves us some work in other parts of the code.
	if len(body) == 0 {
		body = []byte{'{', '}'}
	}

	// Try to parse common response fields to check whether there's an error reported in a response.
	var common *Response
	err = json.Unmarshal(body, &common)
	if err != nil {
		// UnmarshalTypeError errors are ignored, because in some cases API response is an array that cannot
		// unmarshal into a struct.
		_, ok := err.(*json.UnmarshalTypeError)
		if !ok {
			return nil, err
		}
	}
	if common != nil && common.Status == "error" {
		return nil, fmt.Errorf("API returned error: %s", common.Error)
	}

	return body, nil
}
func (client *Client) createRequest(method, api string, reqBody map[string]interface{}) (*http.Request, error) {
	// Handle the body if they gave us one.
	var bodyReader io.Reader
	if method == "POST" && reqBody != nil {
		bJson, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bJson)
	}

	apiUrlStr, err := url.Parse(client.config.BaseUrl + client.config.UrlPath + api)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, apiUrlStr.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	if bodyReader != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", client.Token)
		req.Header.Add("User-Agent", client.config.UserAgent)
	}
	return req, nil
}

func (client *Client) redactError(err error) error {
	if err == nil {
		return nil
	}
	errString := err.Error()

	if len(client.Token) > 0 {
		errString = strings.Replace(errString, client.Token, "redacted", -1)
	}
	if len(client.config.UserAgent) > 0 {
		errString = strings.Replace(errString, client.config.UserAgent, "redacted", -1)
	}

	// Return original error if no replacements were made to keep the original,
	// probably more useful error type information.
	if errString == err.Error() {
		return err
	}
	return fmt.Errorf("%s", errString)
}

func (client *Client) PostExample() (string, error) {

	jsonBody := map[string]interface{}{
		"params": map[string]interface{}{
			"SOME": "THING",
		},
	}

	resp, err := client.doJsonRequest("POST", "/ExampleEndPoint",
		jsonBody)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}
