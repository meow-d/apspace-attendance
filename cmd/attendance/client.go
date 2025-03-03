package main

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
)

type Auth struct {
	Username string
	Password string
	TGT      string
}

type Client struct {
	Auth       Auth
	httpClient *http.Client
}

// get auth info and creates a new client class
func NewClient() *Client {
	a := getAuth()

	client := &http.Client{}
	auth := Auth{
		Username: a.Username,
		Password: a.Password,
		TGT:      a.TGT,
	}

	return &Client{
		Auth:       auth,
		httpClient: client,
	}
}

// sends a http request. doesn't close the response body
func (c *Client) request(method string, url string, body string, headers map[string]string) (*http.Response, error) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// send a http request and returns the response body
func (c *Client) requestBody(method string, url string, body string, headers map[string]string) ([]byte, error) {
	resp, err := c.request(method, url, body, headers)

	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return respBody, err
}

// overrides auth info, but doesn't do the actual login as that's in refreshTGT()
func (c *Client) login(username string, password string) {
	c.Auth.TGT = ""
	c.Auth.Username = username
	c.Auth.Password = password
	setAuth(c.Auth)
}

func (c *Client) removeTGT() {
	c.Auth.TGT = ""
	setAuth(c.Auth)
}

func (c *Client) logout() {
	c.Auth.Username = ""
	c.Auth.Password = ""
	c.Auth.TGT = ""
	setAuth(c.Auth)
}

// logs in and gets the tgt
// ideally password shouldn't be stored at all but, but the tgt expires relatively often. it's how APSpace does stuff
func (c *Client) refreshTGT() error {
	// remove tgt from auth first
	c.removeTGT()

	url := "https://cas.apiit.edu.my/cas/v1/tickets"
	body := fmt.Sprintf("username=%s&password=%s", c.Auth.Username, c.Auth.Password)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp, err := c.request("POST", url, body, headers)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("error loggin in. Incorrect username or password?")
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return fmt.Errorf("error loggin in. Incorrect username or password?")
	}

	parts := strings.Split(location, "/")
	tgt := parts[len(parts)-1]
	if !strings.HasPrefix(tgt, "TGT-") {
		return fmt.Errorf("error loggin in. Incorrect username or password?")
	}

	c.Auth.TGT = tgt
	err = setAuth(c.Auth)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getTicket(service string) (string, error) {
	url := fmt.Sprintf("https://cas.apiit.edu.my/cas/v1/tickets/%s", c.Auth.TGT)
	body := fmt.Sprintf("service=%s", service)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp, err := c.requestBody("POST", url, body, headers)

	if err != nil {
		return "", err
	}
	respStr := string(resp)

	if !strings.HasPrefix(respStr, "ST-") {
		return "", fmt.Errorf("error getting ticket")
	}

	return respStr, err
}

// sends an authanticated request
// requests a ticket, refreshes tgt if ticket request fails
func (c *Client) authenticatedRequest(method, url, body string, headers map[string]string, service string) ([]byte, error) {
	// login if tgt not present
	if c.Auth.TGT == "" {
		err := c.refreshTGT()
		if err != nil {
			return nil, err
		}
	}

	// get ticket
	ticket, err := c.getTicket(service)

	// refresh tgt if ticket request fails
	if err != nil {
		if err.Error() == "Error getting ticket" {
			err = c.refreshTGT()
			if err != nil {
				return nil, err
			}
			ticket, err = c.getTicket(service)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	// send request
	maps.Copy(headers, map[string]string{
		"ticket": ticket,
	},
	)
	return c.requestBody(method, url, body, headers)
}

// attendance stuff

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

func (c *Client) submitAttendance(code string) error {
	service := "https://api.apiit.edu.my/attendix"
	url := "https://attendix.apu.edu.my/graphql"
	body := fmt.Sprintf(`{"operationName":"updateAttendance","variables":{"otp":"%s"},"query":"mutation updateAttendance($otp: String!) {\n  updateAttendance(otp: $otp) {\n    id\n    attendance\n    classcode\n    date\n    startTime\n    endTime\n    classType\n    __typename\n  }\n}\n"}`, code)
	headers := map[string]string{
		"X-Api-Key": "da2-u4ksf3gspnhyjcokxzugo3mqr4", // idk why but the api key is hardcoded
	}

	resp, err := c.authenticatedRequest("POST", url, body, headers, service)
	if err != nil {
		return err
	}

	var result GraphQLResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		return fmt.Errorf("error: %s", result.Errors[0].Message)
	}
	return err
}
