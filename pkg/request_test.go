package pkg

import (
	"authnet/pkg/client"
	"authnet/pkg/common"
	"authnet/pkg/config"
	"authnet/pkg/util"
	"testing"
)

var ac *client.AuthNetClient

func init() {
	conf, loadErr := config.LoadConfigFromEnv(false)
	if loadErr != nil {
		panic(loadErr)
	}
	newClient := client.NewAuthNetClient(*conf)
	ac = &newClient
}

func Test_CreateTransactionRequest(t *testing.T) {
	var response common.CreateTransactionResponse
	request := common.CreateTransactionRequestType{
		ANetApiRequest: common.ANetApiRequest{
			MerchantAuthentication: ac.CreateMerchantAuthenticationType(),
			RefId:                  "12345",
		},
		TransactionRequestType: common.TransactionRequestType{
			TransactionType: "authCaptureTransaction",
			Amount:          util.Float64RefFromInt(5),
			Payment: &common.PaymentType{
				CreditCard: &common.CreditCardType{
					CreditCardSimpleType: common.CreditCardSimpleType{
						CardNumber:     "5424000000000015",
						ExpirationDate: "2025-12",
					},
					CardCode: "999",
				},
			},
			Order: &common.OrderType{
				InvoiceNumber: "INV-12345",
				Description:   "Product Description",
			},
			LineItems: &common.ArrayOfLineItem{
				LineItem: []common.LineItemType{
					{
						ItemId:      "1",
						Name:        "Vase",
						Description: "Cannes Logo",
						Quantity:    float64(18),
						UnitPrice:   float64(45),
					},
				},
			},
			Tax: &common.ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 tax name",
				Description: "level2 tax",
			},
			Duty: &common.ExtendedAmountType{
				Amount:      8.55,
				Name:        "duty name",
				Description: "duty description",
			},
			Shipping: &common.ExtendedAmountType{
				Amount:      4.26,
				Name:        "level2 shipping name",
				Description: "shipping description",
			},
			PoNumber: "456654",
			Customer: &common.CustomerDataType{
				Id: "99999456654",
			},
			BillTo: &common.CustomerAddressType{
				NameAndAddressType: common.NameAndAddressType{
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
			UserFields: &common.UserFields{
				UserField: []common.UserField{
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
			ProcessingOptions: &common.ProcessingOptions{
				IsSubsequentAuth: util.BoolFalseRef(),
			},
			SubsequentAuthInformation: &common.SubsequentAuthInformation{
				OriginalNetworkTransId: "123456789NNNH",
				OriginalAuthAmount:     util.Float64RefFromInt(45),
				Reason:                 common.MerchantInitTransReasonResubmission,
			},
			AuthorizationIndicatorType: &common.AuthorizationIndicatorType{
				AuthorizationIndicator: common.AuthIndicatorFinal,
			},
		},
	}

	if err := ac.SendRequest(request, &response); err != nil {
		t.Fatal(err)
	}

	if response.TransactionResponse.Messages.ResultCode == common.MessageTypeError {
		t.Fatal(response.TransactionResponse.Messages.Message[0].Text)
	}

}
