package auth

type MerchantAuthentication struct {
	Name           string `json:"name"`
	TransactionKey string `json:"transactionKey"`
}
type AuthenticateTestRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"`
	RefId                  *string                `json:"refId,omitempty"`
}

type ResultCode string

const (
	ResultCodeOk    ResultCode = "Ok"
	ResultCodeError ResultCode = "Error"
)

type Message struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

type Messages struct {
	ResultCode ResultCode `json:"resultCode"`
	Message    []Message  `json:"message"`
}

type AuthenticateTestResponse struct {
	RefId    *string   `json:"refId,omitempty"`
	Messages *Messages `json:"messages,omitempty"`
}
