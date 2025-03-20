package gogo_authnet

import (
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

func Test_ChargeCreditCardRequest(t *testing.T) {
	var response CreateTransactionResponse
	request := CreateTransactionRequestType{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
			RefId:                  "12345",
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
						Name:        "Vase",
						Description: "Cannes Logo",
						Quantity:    float64(18),
						UnitPrice:   float64(45),
					},
				},
			},
			Tax: &ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 tax name",
				Description: "level2 tax",
			},
			Duty: &ExtendedAmountType{
				Amount:      8.55,
				Name:        "duty name",
				Description: "duty description",
			},
			Shipping: &ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 shipping name",
				Description: "shipping description",
			},
			PoNumber: "456654",
			Customer: &CustomerDataType{
				Id: "99999456654",
			},
			BillTo: &CustomerAddressType{
				NameAndAddressType: NameAndAddressType{
					FirstName: "China",
					LastName:  "Bayles",
					Company:   "Thyme for Tea",
					Address:   "12 Main Street",
					City:      "Pecan Springs",
					State:     "TX",
					Zip:       "44628",
					Country:   "US",
				},
			},
			CustomerIp: "192.168.1.1",
			UserFields: &UserFields{
				UserField: []UserField{
					{
						Name:  "MerchantDefinedFieldName1",
						Value: "MerchantDefinedFieldValue1",
					},
					{
						Name:  "favorite_color",
						Value: "blue",
					},
				},
			},
			ProcessingOptions: &ProcessingOptions{
				IsSubsequentAuth: BoolFalseRef(),
			},
			SubsequentAuthInformation: &SubsequentAuthInformation{
				OriginalNetworkTransId: "1234567890",
				OriginalAuthAmount:     Float64RefFromInt(45),
				Reason:                 MerchantInitTransReasonResubmission,
			},
			AuthorizationIndicatorType: &AuthorizationIndicatorType{
				AuthorizationIndicator: AuthIndicatorFinal,
			},
		},
	}

	if err := ac.SendRequest(request, &response); err != nil {
		t.Fatal(err)
	}

	if response.TransactionResponse.Messages.ResultCode == MessageTypeError {
		t.Fatal(response.TransactionResponse.Messages.Message[0].Text)
	}

}

func Test_AuthorizeCreditCardRequest(t *testing.T) {
	var response CreateTransactionResponse
	request := CreateTransactionRequestType{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
			RefId:                  "12345",
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
			LineItems: &ArrayOfLineItem{
				LineItem: []LineItemType{
					{
						ItemId:      "1",
						Name:        "Vase",
						Description: "Cannes Logo",
						Quantity:    float64(18),
						UnitPrice:   float64(45),
					},
				},
			},
			Tax: &ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 tax name",
				Description: "level2 tax",
			},
			Duty: &ExtendedAmountType{
				Amount:      8.55,
				Name:        "duty name",
				Description: "duty description",
			},
			Shipping: &ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 shipping name",
				Description: "shipping description",
			},
			PoNumber: "456654",
			Customer: &CustomerDataType{
				Id: "99999456654",
			},
			BillTo: &CustomerAddressType{
				NameAndAddressType: NameAndAddressType{
					FirstName: "Ellen",
					LastName:  "Johnson",
					Company:   "Souveniropolis",
					Address:   "14 Main Street",
					City:      "Pecan Springs",
					State:     "TX",
					Zip:       "44628",
					Country:   "US",
				},
			},
			ShipTo: &NameAndAddressType{
				FirstName: "China",
				LastName:  "Bayles",
				Company:   "Thyme for Tea",
				Address:   "12 Main Street",
				City:      "Pecan Springs",
				State:     "TX",
				Zip:       "44628",
				Country:   "US",
			},
			CustomerIp: "192.168.1.1",
			UserFields: &UserFields{
				UserField: []UserField{
					{
						Name:  "MerchantDefinedFieldName1",
						Value: "MerchantDefinedFieldValue1",
					},
					{
						Name:  "favorite_color",
						Value: "blue",
					},
				},
			},
			ProcessingOptions: &ProcessingOptions{
				IsSubsequentAuth: BoolFalseRef(),
			},
			SubsequentAuthInformation: &SubsequentAuthInformation{
				OriginalNetworkTransId: "123456789NNNH",
				OriginalAuthAmount:     Float64RefFromInt(45),
				Reason:                 MerchantInitTransReasonResubmission,
			},
			AuthorizationIndicatorType: &AuthorizationIndicatorType{
				AuthorizationIndicator: AuthIndicatorPre,
			},
		},
	}

	if err := ac.SendRequest(request, &response); err != nil {
		t.Fatal(err)
	}

	if response.TransactionResponse.Messages.ResultCode == MessageTypeError {
		t.Fatal(response.TransactionResponse.Messages.Message[0].Text)
	}
}

func Test_PreviouslyAuthorizedAmountRequest(t *testing.T) {
	var response CreateTransactionResponse
	request := CreateTransactionRequestType{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
			RefId:                  "12345",
		},
		TransactionRequestType: TransactionRequestType{
			TransactionType: TransactionTypePriorAuthCaptureTransaction,
			Amount:          Float64RefFromInt(5),
			RefTransId:      "1234567890",
		},
	}

	if err := ac.SendRequest(request, &response); err != nil {
		t.Fatal(err)
	}

	if response.TransactionResponse.Messages.ResultCode == MessageTypeError {
		t.Fatal(response.TransactionResponse.Messages.Message[0].Text)
	}
}
