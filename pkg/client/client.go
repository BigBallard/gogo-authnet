package client

import (
	"authnet/pkg/common"
	"authnet/pkg/config"
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
	config     config.Config
	apiUrl     string
	httpClient http.Client
}

func NewAuthNetClient(config config.Config) AuthNetClient {
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

func (c *AuthNetClient) AuthenticateTest() (*common.AuthenticateTestResponse, error) {
	testRequest := common.AuthenticateTestRequest{
		MerchantAuthentication: common.MerchantAuthentication{
			Name:           c.config.Auth.ApiLoginId,
			TransactionKey: c.config.Auth.TransactionId,
		},
	}
	var testResponse common.AuthenticateTestResponse
	rErr := c.SendRequest(testRequest, &testResponse)
	return &testResponse, rErr
}

// SendRequest takes a request type instance and a response type instance. req can be passed either by reference or by
// value. The res however, is required to be a reference due to the unmarshalling phase of the request.
func (c *AuthNetClient) SendRequest(req any, res any) error {
	bodyBytes, mErr := xml.Marshal(req)
	if mErr != nil {
		return errors.Join(errors.New("unable to marshal request body"), mErr)
	}
	response, reqErr := c.httpClient.Post(c.apiUrl, "text/xml", bytes.NewReader(bodyBytes))
	if reqErr != nil {
		return errors.Join(errors.New("unable to make http request"), reqErr)
	}
	defer response.Body.Close()
	resBytes := make([]byte, response.ContentLength)
	nRead, readErr := response.Body.Read(resBytes)
	if nRead != int(response.ContentLength) && readErr == io.EOF {
		return errors.Join(errors.New("unable to read response body"), reqErr)
	}
	if uErr := xml.Unmarshal(resBytes, res); uErr != nil {
		return errors.Join(errors.New("unable to unmarshal response body"), reqErr)
	}
	return nil
}
