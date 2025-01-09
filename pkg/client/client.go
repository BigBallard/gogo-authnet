package client

import (
	"authnet/pkg/auth"
	"authnet/pkg/config"
	"bytes"
	"encoding/json"
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

func requestToBody(requestName string, request any) map[string]any {
	return map[string]any{
		requestName: request,
	}
}

func (c *AuthNetClient) AuthenticateTest() (*auth.AuthenticateTestResponse, error) {
	testRequest := auth.AuthenticateTestRequest{
		MerchantAuthentication: auth.MerchantAuthentication{
			Name:           c.config.Auth.ApiLoginId,
			TransactionKey: c.config.Auth.TransactionId,
		},
	}
	requestBody := requestToBody("authenticateTestRequest", testRequest)
	bodyBytes, mErr := json.Marshal(requestBody)
	if mErr != nil {
		return nil, errors.Join(errors.New("unable to marshal request body"), mErr)
	}

	response, reqErr := c.httpClient.Post(c.apiUrl, "application/json", bytes.NewReader(bodyBytes))
	if reqErr != nil {
		return nil, errors.Join(errors.New("unable to make http request"), reqErr)
	}
	defer response.Body.Close()

	resBytes := make([]byte, response.ContentLength)
	nRead, readErr := response.Body.Read(resBytes)
	if nRead != int(response.ContentLength) && readErr == io.EOF {
		return nil, errors.Join(errors.New("unable to read response body"), reqErr)
	}

	var testResponse auth.AuthenticateTestResponse
	resBytes = bytes.TrimPrefix(resBytes, []byte("\xef\xbb\xbf")) // remove byte-order mark (BOM)
	if uErr := json.Unmarshal(resBytes, &testResponse); uErr != nil {
		return nil, errors.Join(errors.New("unable to unmarshal response body"), reqErr)
	}
	return &testResponse, nil
}
