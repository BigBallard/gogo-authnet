package authnet

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
)

var ac *AuthNetClient

func init() {
	conf, loadErr := LoadConfigFromEnv()
	if loadErr != nil {
		panic(loadErr)
	}
	newClient := NewAuthNetClient(*conf)
	ac = &newClient
}

func Test_ChargeCreditCard(t *testing.T) {
	err := retry(3, func() error {
		refId := genRefId()
		request := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeAuthCaptureTransaction,
				Amount:          Float64RefFromInt(5),
				Payment: &PaymentType{
					CreditCard: &CreditCardType{
						CreditCardSimpleType: CreditCardSimpleType{
							CardNumber:     "5424000000000015",
							ExpirationDate: "2025-12",
						},
						CardCode: "999",
					},
				},
				Order: &OrderType{
					InvoiceNumber: "INV-12345",
					Description:   "Product Description",
				},
				LineItems: &ArrayOfLineItem{
					LineItem: []LineItemType{
						{
							ItemId:      "1",
							Name:        "vase",
							Description: "Cannes logo",
							Quantity:    18.0,
							UnitPrice:   45.0,
						},
					},
				},
				Tax: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Tax name",
					Description: "level 2 tax",
				},
				Duty: &ExtendedAmountType{
					Amount:      18.55,
					Name:        "Duty name",
					Description: "duty description",
				},
				Shipping: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Shipping name",
					Description: "shipping description",
				},
				PoNumber: "456654",
				Customer: &CustomerDataType{
					Id: "99999456654",
				},
				BillTo: &CustomerAddressType{
					NameAndAddressType: NameAndAddressType{
						FirstName: "John",
						LastName:  "Doe",
						Company:   "FooBar Inc",
						Address:   "1 A St",
						City:      "FizBuzz",
						State:     "CA",
						Country:   "US",
					},
				},
				ShipTo: &NameAndAddressType{
					FirstName: "Big",
					LastName:  "Bux",
					Company:   "Too Much Money & Sons",
					Address:   "Rich Rd",
					City:      "Richland",
					State:     "NY",
					Country:   "US",
				},
				CustomerIp: "192.168.0.1",
				ProcessingOptions: &ProcessingOptions{
					IsSubsequentAuth: BoolTrueRef(),
				},
				SubsequentAuthInformation: &SubsequentAuthInformation{
					OriginalNetworkTransId: "123456789NNNH",
					OriginalAuthAmount:     Float64RefFromFloat(45.5),
					Reason:                 MerchantInitTransReasonResubmission,
				},
				AuthorizationIndicatorType: &AuthorizationIndicatorType{
					AuthorizationIndicator: AuthIndicatorFinal,
				},
			},
		}
		var response CreateTransactionResponse
		if sendErr := ac.SendRequest(request, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AuthorizeAndCapture(t *testing.T) {
	err := retry(3, func() error {
		refId := genRefId()
		authRequest := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeAuthOnlyTransaction,
				Amount:          Float64RefFromInt(5),
				Payment: &PaymentType{
					CreditCard: &CreditCardType{
						CreditCardSimpleType: CreditCardSimpleType{
							CardNumber:     "5424000000000015",
							ExpirationDate: "2025-12",
						},
						CardCode: "999",
					},
				},
				Order: &OrderType{
					InvoiceNumber: "INV-12345",
					Description:   "Product Description",
				},
				LineItems: &ArrayOfLineItem{
					LineItem: []LineItemType{
						{
							ItemId:      "1",
							Name:        "vase",
							Description: "Cannes logo",
							Quantity:    18.0,
							UnitPrice:   45.0,
						},
					},
				},
				Tax: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Tax name",
					Description: "level 2 tax",
				},
				Duty: &ExtendedAmountType{
					Amount:      18.55,
					Name:        "Duty name",
					Description: "duty description",
				},
				Shipping: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Shipping name",
					Description: "shipping description",
				},
				PoNumber: "456654",
				Customer: &CustomerDataType{
					Id: "99999456654",
				},
				BillTo: &CustomerAddressType{
					NameAndAddressType: NameAndAddressType{
						FirstName: "John",
						LastName:  "Doe",
						Company:   "FooBar Inc",
						Address:   "1 A St",
						City:      "FizBuzz",
						State:     "CA",
						Country:   "US",
					},
				},
				ShipTo: &NameAndAddressType{
					FirstName: "Big",
					LastName:  "Bux",
					Company:   "Too Much Money & Sons",
					Address:   "Rich Rd",
					City:      "Richland",
					State:     "NY",
					Country:   "US",
				},
				CustomerIp: "192.168.0.1",
				ProcessingOptions: &ProcessingOptions{
					IsSubsequentAuth: BoolTrueRef(),
				},
				SubsequentAuthInformation: &SubsequentAuthInformation{
					OriginalNetworkTransId: "123456789NNNH",
					OriginalAuthAmount:     Float64RefFromFloat(45.5),
					Reason:                 MerchantInitTransReasonResubmission,
				},
				AuthorizationIndicatorType: &AuthorizationIndicatorType{
					AuthorizationIndicator: AuthIndicatorFinal,
				},
			},
		}

		var response CreateTransactionResponse
		if sendErr := ac.SendRequest(authRequest, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}
		transId := response.TransactionResponse.TransId

		captRequest := &CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  "2",
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypePriorAuthCaptureTransaction,
				Amount:          Float64RefFromInt(5),
				RefTransId:      transId,
				Order: &OrderType{
					InvoiceNumber: "INV-12345",
					Description:   "Product Description",
				},
			},
		}

		if sendErr := ac.SendRequest(captRequest, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_RefundTransaction(t *testing.T) {
	err := retry(3, func() error {
		refId := genRefId()
		request := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeAuthCaptureTransaction,
				Amount:          Float64RefFromInt(5),
				Payment: &PaymentType{
					CreditCard: &CreditCardType{
						CreditCardSimpleType: CreditCardSimpleType{
							CardNumber:     "5424000000000015",
							ExpirationDate: "2025-12",
						},
						CardCode: "999",
					},
				},
				Order: &OrderType{
					InvoiceNumber: "INV-12345",
					Description:   "Product Description",
				},
				LineItems: &ArrayOfLineItem{
					LineItem: []LineItemType{
						{
							ItemId:      "1",
							Name:        "vase",
							Description: "Cannes logo",
							Quantity:    18.0,
							UnitPrice:   45.0,
						},
					},
				},
				Tax: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Tax name",
					Description: "level 2 tax",
				},
				Duty: &ExtendedAmountType{
					Amount:      18.55,
					Name:        "Duty name",
					Description: "duty description",
				},
				Shipping: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Shipping name",
					Description: "shipping description",
				},
				PoNumber: "456654",
				Customer: &CustomerDataType{
					Id: "99999456654",
				},
				BillTo: &CustomerAddressType{
					NameAndAddressType: NameAndAddressType{
						FirstName: "John",
						LastName:  "Doe",
						Company:   "FooBar Inc",
						Address:   "1 A St",
						City:      "FizBuzz",
						State:     "CA",
						Country:   "US",
					},
				},
				ShipTo: &NameAndAddressType{
					FirstName: "Big",
					LastName:  "Bux",
					Company:   "Too Much Money & Sons",
					Address:   "Rich Rd",
					City:      "Richland",
					State:     "NY",
					Country:   "US",
				},
				CustomerIp: "192.168.0.1",
				ProcessingOptions: &ProcessingOptions{
					IsSubsequentAuth: BoolTrueRef(),
				},
				SubsequentAuthInformation: &SubsequentAuthInformation{
					OriginalNetworkTransId: "123456789NNNH",
					OriginalAuthAmount:     Float64RefFromFloat(45.5),
					Reason:                 MerchantInitTransReasonResubmission,
				},
				AuthorizationIndicatorType: &AuthorizationIndicatorType{
					AuthorizationIndicator: AuthIndicatorFinal,
				},
			},
		}
		var response CreateTransactionResponse
		if sendErr := ac.SendRequest(request, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}

		tranId := response.TransactionResponse.TransId

		refundRequest := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeRefundTransaction,
				Amount:          Float64RefFromInt(5),
				Payment: &PaymentType{
					CreditCard: &CreditCardType{
						CreditCardSimpleType: CreditCardSimpleType{
							CardNumber:     "5424000000000015",
							ExpirationDate: "2025-12",
						},
					},
				},
				RefTransId: tranId,
			},
		}

		if sendErr := ac.SendRequest(refundRequest, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_VoidCompletedTransaction(t *testing.T) {
	err := retry(3, func() error {

		refId := genRefId()
		request := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeAuthCaptureTransaction,
				Amount:          Float64RefFromInt(5),
				Payment: &PaymentType{
					CreditCard: &CreditCardType{
						CreditCardSimpleType: CreditCardSimpleType{
							CardNumber:     "5424000000000015",
							ExpirationDate: "2025-12",
						},
						CardCode: "999",
					},
				},
				Order: &OrderType{
					InvoiceNumber: "INV-12345",
					Description:   "Product Description",
				},
				LineItems: &ArrayOfLineItem{
					LineItem: []LineItemType{
						{
							ItemId:      "1",
							Name:        "vase",
							Description: "Cannes logo",
							Quantity:    18.0,
							UnitPrice:   45.0,
						},
					},
				},
				Tax: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Tax name",
					Description: "level 2 tax",
				},
				Duty: &ExtendedAmountType{
					Amount:      18.55,
					Name:        "Duty name",
					Description: "duty description",
				},
				Shipping: &ExtendedAmountType{
					Amount:      4.26,
					Name:        "Shipping name",
					Description: "shipping description",
				},
				PoNumber: "456654",
				Customer: &CustomerDataType{
					Id: "99999456654",
				},
				BillTo: &CustomerAddressType{
					NameAndAddressType: NameAndAddressType{
						FirstName: "John",
						LastName:  "Doe",
						Company:   "FooBar Inc",
						Address:   "1 A St",
						City:      "FizBuzz",
						State:     "CA",
						Country:   "US",
					},
				},
				ShipTo: &NameAndAddressType{
					FirstName: "Big",
					LastName:  "Bux",
					Company:   "Too Much Money & Sons",
					Address:   "Rich Rd",
					City:      "Richland",
					State:     "NY",
					Country:   "US",
				},
				CustomerIp: "192.168.0.1",
				ProcessingOptions: &ProcessingOptions{
					IsSubsequentAuth: BoolTrueRef(),
				},
				SubsequentAuthInformation: &SubsequentAuthInformation{
					OriginalNetworkTransId: "123456789NNNH",
					OriginalAuthAmount:     Float64RefFromFloat(45.5),
					Reason:                 MerchantInitTransReasonResubmission,
				},
				AuthorizationIndicatorType: &AuthorizationIndicatorType{
					AuthorizationIndicator: AuthIndicatorFinal,
				},
			},
		}
		var response CreateTransactionResponse
		if sendErr := ac.SendRequest(request, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}

		tranId := response.TransactionResponse.TransId

		voidRequest := CreateTransactionRequestType{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
				RefId:                  refId,
			},
			TransactionRequestType: TransactionRequestType{
				TransactionType: TransactionTypeVoidTransaction,
				RefTransId:      tranId,
			},
		}

		if sendErr := ac.SendRequest(voidRequest, &response); sendErr != nil {
			return sendErr
		}
		if response.Messages.ResultCode == MessageTypeError {
			msg := response.Messages.Message[0]
			return errors.New(fmt.Sprintf("%s: %s", msg.Code, msg.Text))
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func genRefId() string {
	return fmt.Sprintf("%d", rand.IntN(99999))
}

func retry(times int, f func() error) error {
	count := 0
	var err error
	for count < times {
		if err = f(); err == nil {
			break
		} else {
			count++
		}
	}
	return err
}
