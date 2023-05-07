package client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sharovik/toggl-jira/src/log"
)

// BaseHTTPClientInterface base interface for all http clients
type BaseHTTPClientInterface interface {
	//Configuration methods
	SetOauthToken(token string)
	SetBaseURL(baseURL string)
	BasicAuth(username string, password string) string
	GetClientID() string
	GetClientSecret() string
	GetOAuthToken() string

	//Http methods
	Request(string, string, interface{}, map[string]string) ([]byte, int, error)
	Post(string, interface{}, map[string]string) ([]byte, int, error)
	Get(string, map[string]string, map[string]string) ([]byte, int, error)
	Put(string, interface{}, map[string]string) ([]byte, int, error)
}

// HTTPClient main http client
type HTTPClient struct {
	Client *http.Client

	//Configuration of client
	OAuthToken   string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

// SetOauthToken method sets the oauth token and retrieves its self
func (client *HTTPClient) SetOauthToken(token string) {
	client.OAuthToken = token
}

// GetClientID method retrieves the clientID
func (client *HTTPClient) GetClientID() string {
	return client.ClientID
}

// GetClientSecret method retrieves the clientSecret
func (client *HTTPClient) GetClientSecret() string {
	return client.ClientSecret
}

// GetOAuthToken method retrieves the oauth token
func (client *HTTPClient) GetOAuthToken() string {
	return client.OAuthToken
}

// SetBaseURL method sets the base url and retrieves its self
func (client *HTTPClient) SetBaseURL(baseURL string) {
	client.BaseURL = baseURL
}

// BasicAuth retrieves the encode string for basic auth
func (client *HTTPClient) BasicAuth(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
}

// Request method for API requests
//
// This method accepts parameters:
// method - the method of request. Ex: POST, GET, PUT, DELETE and etc
// endpoint - endpoint to which we should do a request
// body - it's a request body. Accepted types of body: string, url.Values(for form_data requests), byte
// headers - request headers
func (client *HTTPClient) Request(method string, endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {

	log.Logger().StartMessage("Http request")

	var (
		resp    *http.Response
		request *http.Request
		err     error
	)

	switch b := body.(type) {
	case string:
		log.Logger().Debug().
			Str("endpoint", endpoint).
			Str("method", method).
			Str("body", b).
			Msg("Endpoint call")
		request, err = http.NewRequest(method, endpoint, strings.NewReader(b))
		if err != nil {
			log.Logger().AddError(err).Msg("Error during the request generation")
			log.Logger().FinishMessage("Http request")
			return nil, 0, err
		}

		request.Header.Set("Content-Type", "application/json")
	case url.Values:
		log.Logger().Debug().
			Str("endpoint", endpoint).
			Str("method", method).
			Interface("body", b).
			Msg("Endpoint call")

		request, err = http.NewRequest(method, endpoint, strings.NewReader(b.Encode()))
		if err != nil {
			log.Logger().AddError(err).Msg("Error during the request generation")
			log.Logger().FinishMessage("Http request")
			return nil, 0, err
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
		log.Logger().Debug().
			Str("endpoint", endpoint).
			Str("method", method).
			Str("body", string(b.([]byte))).
			Msg("Endpoint call")
		request, err = http.NewRequest(method, endpoint, bytes.NewReader(b.([]byte)))
		if err != nil {
			log.Logger().AddError(err).Msg("Error during the request generation")
			log.Logger().FinishMessage("Http request")
			return nil, 0, err
		}
		request.Header.Set("Content-Type", "application/json")
	}

	for attribute, value := range headers {
		request.Header.Set(attribute, value)
	}

	if client.OAuthToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.OAuthToken))
	}

	resp, errorResponse := client.Client.Do(request)

	if resp == nil {
		err = errors.New("response cannot be null")
		errMsg := err.Error()
		if errorResponse != nil {
			errMsg = errorResponse.Error()
		}
		log.Logger().AddError(errorResponse).
			Str("response_error", errMsg).
			Msg("Error during response body parse")

		log.Logger().FinishMessage("Http request")
		return nil, 0, err
	}

	defer resp.Body.Close()
	byteResp, errorConversion := io.ReadAll(resp.Body)
	if errorConversion != nil {
		log.Logger().AddError(errorConversion).
			Err(errorConversion).
			Msg("Error during response body parse")
		log.Logger().FinishMessage("Http request")
		return byteResp, 0, errorConversion
	}

	var response []byte
	if string(byteResp) == "" {
		response = []byte(`{}`)
	} else {
		response = byteResp
	}

	log.Logger().FinishMessage("Http request")
	return response, resp.StatusCode, nil
}

// Post method for POST http requests
func (client *HTTPClient) Post(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPost, client.generateAPIUrl(endpoint), body, headers)
}

// Put method for PUT http requests
func (client *HTTPClient) Put(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPut, client.generateAPIUrl(endpoint), body, headers)
}

// Get method for GET http requests
func (client *HTTPClient) Get(endpoint string, query map[string]string, headers map[string]string) ([]byte, int, error) {
	if client.OAuthToken != "" {
		query["access_token"] = client.OAuthToken
	}

	var queryString = ""
	for fieldName, value := range query {
		if queryString == "" {
			queryString += "?"
		} else {
			queryString += "&"
		}

		queryString += fmt.Sprintf("%s=%s", fieldName, value)
	}

	return client.Request(http.MethodGet, client.generateAPIUrl(endpoint)+queryString, []byte(``), headers)
}

func (client *HTTPClient) generateAPIUrl(endpoint string) string {
	log.Logger().Debug().
		Str("base_url", client.BaseURL).
		Str("endpoint", endpoint).
		Msg("Generate API url")

	return client.BaseURL + endpoint
}

func GetHttpClient() http.Client {
	netTransport := &http.Transport{
		TLSHandshakeTimeout: 7 * time.Second,
	}

	return http.Client{
		Timeout:   time.Duration(15) * time.Second,
		Transport: netTransport,
	}
}
