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

func NewClient(a Auth) Client {
	client := &http.Client{}
	auth := Auth{
		Username: a.Username,
		Password: a.Password,
		TGT:      a.TGT,
	}
	return Client{
		Auth:       auth,
		httpClient: client,
	}
}

func (c Client) request(method string, url string, body string, headers map[string]string) ([]byte, error) {
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
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, err
}

func (c Client) login(username string, password string) error {
	c.Auth.TGT = ""
	c.Auth.Username = username
	c.Auth.Password = password

	return c.refreshTGT()
}

func (c Client) refreshTGT() error {
	url := "https://cas.apiit.edu.my/cas/v1/tickets"
	body := fmt.Sprintf("tgt=%s", c.Auth.TGT)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp, err := c.request("POST", url, body, headers)

	if err == nil {
		c.Auth.TGT = string(resp)
		setAuth(c.Auth)
		return err
	}

	if strings.Contains(err.Error(), "401") {
		return fmt.Errorf("Incorrect username or password")
	}

	return fmt.Errorf("Error logging in")
}

func (c Client) getTicket(service string) (string, error) {
	url := fmt.Sprintf("https://cas.apiit.edu.my/cas/v1/tickets/%s", c.Auth.TGT)
	body := fmt.Sprintf("service=%s", service)
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp, err := c.request("POST", url, body, headers)

	return string(resp), err
}

func (c Client) authenticatedRequest(method, url, body string, headers map[string]string, service string) ([]byte, error) {
	ticket, err := c.getTicket(service)

	if err != nil {
		if strings.Contains(err.Error(), "401") {

			err = c.refreshTGT()
			if err != nil {
				return nil, err
			}

			ticket, err = c.getTicket(service)
			if err != nil {
				return nil, err
			}
		}
	}

	maps.Copy(headers, map[string]string{
		"ticket": ticket,
	},
	)
	return c.request(method, url, body, headers)
}

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors"`
}

func (c Client) submitAttendance(code string) error {
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

	// TODO test graphql response
	if len(result.Errors) > 0 {
		return fmt.Errorf("GraphQL Error: %s", result.Errors[0].Message)
	}
	return err
}
