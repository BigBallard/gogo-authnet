package auth

import "encoding/xml"

type MerchantAuthentication struct {
	Name           string `xml:"name"`
	TransactionKey string `xml:"transactionKey"`
}
type AuthenticateTestRequest struct {
	XMLName                xml.Name               `xml:"AnetApi/xml/v1/schema/AnetApiSchema.xsd authenticateTestRequest"`
	MerchantAuthentication MerchantAuthentication `xml:"merchantAuthentication"`
	RefId                  *string                `xml:"refId,omitempty"`
}

type ResultCode string

const (
	ResultCodeOk    ResultCode = "Ok"
	ResultCodeError ResultCode = "Error"
)

type Message struct {
	Code string `xml:"code"`
	Text string `xml:"text"`
}

type Messages struct {
	ResultCode ResultCode `xml:"resultCode"`
	Message    []Message  `xml:"message"`
}

type AuthenticateTestResponse struct {
	RefId    *string   `xml:"refId,omitempty"`
	Messages *Messages `xml:"messages,omitempty"`
}
