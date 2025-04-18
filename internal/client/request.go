package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/go-retryablehttp"
)

var NoResErr = errors.New("failed to get response. Check out the Debricked status page: https://status.debricked.com/")
var SupportedFormatsFallbackError = errors.New("get supported formats from the server. Using cached data instead")

func get(uri string, debClient *DebClient, retry bool, format string) (*http.Response, error) {
	request, err := newRequest("GET", *debClient.host+uri, debClient.jwtToken, format, nil)
	if err != nil {
		return nil, err
	}
	res, _ := debClient.httpClient.Do(request)
	req := func() (*http.Response, error) {
		return get(uri, debClient, false, format)
	}

	return interpret(res, req, debClient, retry)
}

func post(uri string, debClient *DebClient, contentType string, body *bytes.Buffer, retry bool) (*http.Response, error) {
	request, err := newRequest("POST", *debClient.host+uri, debClient.jwtToken, "application/json", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)

	res, err := debClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	req := func() (*http.Response, error) {
		return post(uri, debClient, contentType, body, false)
	}

	return interpret(res, req, debClient, retry)
}

func postWithTimeout(uri string, debClient *DebClient, contentType string, body *bytes.Buffer, retry bool, timeout int) (*http.Response, error) {
	request, err := newRequest("POST", *debClient.host+uri, debClient.jwtToken, "application/json", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)

	timeoutDuration := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(request.Context(), timeoutDuration)
	defer cancel()
	request = request.WithContext(ctx)

	res, err := debClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	req := func() (*http.Response, error) {
		return post(uri, debClient, contentType, body, false)
	}

	return interpret(res, req, debClient, retry)
}

// newRequest creates a new HTTP request with necessary headers added
func newRequest(method string, url string, jwtToken string, format string, body io.Reader) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", format)
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	return req, nil
}

// interpret a http response
func interpret(res *http.Response, request func() (*http.Response, error), debClient *DebClient, retry bool) (*http.Response, error) {
	if res == nil {
		return nil, NoResErr
	} else if res.StatusCode == http.StatusForbidden {
		errMsg := `Forbidden. You don't have the necessary access to perform this action. 
		Make sure your access token has proper access https://docs.debricked.com/product/administration/generate-access-token
		For enterprise users: Contact your Debricked company admin or repository admin to request proper access https://docs.debricked.com/product/administration/users/role-based-access-control-enterprise`

		return nil, errors.New(errMsg)
	} else if res.StatusCode == http.StatusUnauthorized {
		errMsg := `Unauthorized. Specify access token. 
Read more on https://docs.debricked.com/product/administration/generate-access-token`
		if retry {
			err := debClient.authenticate()
			if err != nil {
				return nil, errors.New(errMsg)
			}

			return request()
		}

		return nil, errors.New(errMsg)
	}

	return res, nil
}

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (debClient *DebClient) authenticate() error {
	if debClient.accessToken != nil { // To avoid segfault
		if len(*debClient.accessToken) != 0 {
			return debClient.authenticateExplicitToken()
		}
	}

	return debClient.authenticateCachedToken()
}

func (debClient *DebClient) authenticateCachedToken() error {
	token, err := debClient.authenticator.Token()
	if err == nil {
		debClient.jwtToken = token.AccessToken
	}

	return err
}

func (debClient *DebClient) authenticateExplicitToken() error {
	uri := "/api/login_refresh"

	data := map[string]string{"refresh_token": *debClient.accessToken}
	jsonData, _ := json.Marshal(data)
	res, reqErr := debClient.httpClient.Post(
		*debClient.host+uri,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if reqErr != nil {
		return reqErr
	}

	if res != nil {
		defer res.Body.Close()
	}

	var tokenData map[string]string
	body, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(body, &tokenData)
	if err != nil {
		var errMessage errorMessage
		_ = json.Unmarshal(body, &errMessage)

		return fmt.Errorf("%s %s\n", color.RedString("⨯"), errMessage.Message)
	}
	debClient.jwtToken = tokenData["token"]

	return nil
}
