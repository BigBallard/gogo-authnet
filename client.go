package authnet

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type AuthNetClient struct {
	config     Config
	apiUrl     string
	httpClient http.Client
}

func NewAuthNetClient(config Config) AuthNetClient {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	}

	host := config.AuthnetHost
	if !strings.HasSuffix(host, "/") {
		host += "/"
	}

	return AuthNetClient{
		config: config,
		apiUrl: host + "xml/v1/request.api",
		httpClient: http.Client{
			Transport: transport,
		},
	}
}

func (c *AuthNetClient) CreateMerchantAuthenticationType() MerchantAuthenticationType {
	return MerchantAuthenticationType{
		Name:           c.config.Auth.ApiLoginId,
		TransactionKey: c.config.Auth.TransactionKey,
	}
}

func (c *AuthNetClient) AuthenticateTest() (*AuthenticateTestResponse, error) {
	testRequest := AuthenticateTestRequest{
		MerchantAuthentication: MerchantAuthentication{
			Name:           c.config.Auth.ApiLoginId,
			TransactionKey: c.config.Auth.TransactionKey,
		},
	}
	var testResponse AuthenticateTestResponse
	rErr := c.SendRequest(testRequest, &testResponse)
	return &testResponse, rErr
}

// RequestError contains the common.ErrorResponse or errors from some other cause. Either could be populated or one of
// them. This fulfills the error interface and provides the Error() function.
type RequestError struct {
	Response *ErrorResponse
	Err      error
}

func (e *RequestError) Error() string {
	if e.Response != nil {
		messageErr := errors.New(e.Response.Messages.Message[0].Text)
		if e.Err != nil {
			messageErr = errors.Join(messageErr, e.Err)
		}
		return messageErr.Error()
	} else if e.Err != nil {
		return e.Err.Error()
	} else {
		return ""
	}
}

// SendRequest takes a request type instance and a response type instance. req can be passed either by reference or by
// value. The res however, is required to be a reference due to the unmarshalling phase of the request.
func (c *AuthNetClient) SendRequest(req any, res any) *RequestError {
	var requestError RequestError
	bodyBytes, mErr := xml.Marshal(req)
	if mErr != nil {
		requestError.Err = errors.Join(errors.New("unable to marshal request body"), mErr)
		return &requestError
	}
	response, reqErr := c.httpClient.Post(c.apiUrl, "text/xml", bytes.NewReader(bodyBytes))
	if reqErr != nil {
		requestError.Err = errors.Join(errors.New("unable to make http request"), reqErr)
		return &requestError
	}
	defer response.Body.Close()
	resBytes := make([]byte, response.ContentLength)
	nRead, readErr := response.Body.Read(resBytes)
	if nRead != int(response.ContentLength) && readErr == io.EOF {
		requestError.Err = errors.Join(errors.New("unable to read response body"), reqErr)
		return &requestError
	}
	if uErr := xml.Unmarshal(resBytes, res); uErr != nil {
		// check if response is ErrorResponse
		var errResponse ErrorResponse
		if ueErr := xml.Unmarshal(resBytes, &errResponse); ueErr != nil {
			requestError.Err = errors.Join(errors.New("unable to unmarshal response body"), reqErr, ueErr)
		} else {
			requestError.Err = reqErr
			requestError.Response = &errResponse
		}
		return &requestError
	}
	return nil
}
