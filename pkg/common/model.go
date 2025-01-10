package common

import (
	"encoding/xml"
	"time"
)

// AnetApiRequestType defines common properties associated with all API method requests.
type AnetApiRequestType struct {
	MerchantAuthentication MerchantAuthenticationType `xml:"merchantAuthentication" validation:"required"`
	ClientId               string                     `xml:"clientId,omitempty" validation:"max=30"`
	RefId                  string                     `xml:"refId,omitempty" validation:"max=20"`
}

type ImpersonationAuthenticationType struct {
	PartnerLoginId        string `xml:"partnerLoginId" validation:"required,max=25"`
	PartnerTransactionKey string `xml:"partnerTransactionKey" validation:"required,max=16"`
}

type FingerPrintType struct {
	HashValue    string `xml:"hashValue" validation:"required"`
	Sequence     string `xml:"sequence,omitempty"`
	Timestamp    string `xml:"timestamp" validation:"required"`
	CurrencyCode string `xml:"currencyCode"`
	Amount       string `xml:"amount"`
}

// MerchantAuthenticationType containers merchant authentication information.
//
// The following fields are mutually exclusive and one is required for request: TransactionKey, SessionToken, Password,
// ImpersonationAuthentication, FingerPrint, ClientKey, and AccessToken. Validation will check for one of these being
// set and will result in a validation error if more than one or none are set.
type MerchantAuthenticationType struct {
	Name                        string                           `xml:"name,omitempty" validation:"max=30"`
	TransactionKey              string                           `xml:"transactionKey,omitempty" validation:"max=16"`
	SessionToken                string                           `xml:"sessionToken,omitempty"`
	Password                    string                           `xml:"password,omitempty" validation:"max=40"`
	ImpersonationAuthentication *ImpersonationAuthenticationType `xml:"impersonationAuthentication,omitempty"`
	FingerPrint                 *FingerPrintType                 `xml:"fingerPrint,omitempty"`
	ClientKey                   string                           `xml:"clientKey,omitempty"`
	AccessToken                 string                           `xml:"accessToken,omitempty"`
	MobileDeviceId              string                           `xml:"mobileDeviceId,omitempty" validation:"max=60"`
}

// PaymentType indicates the payment type/method for the transactions.
//
// The following fields are mutually exclusive and one is required for request: CreditCard, BankAccount, TrackData,
// EncryptedTrackData, PayPal, OpaqueData, and Emv. Validation will check for one of these being set and will result in
// a validation error if more than one or none are set.
type PaymentType struct {
	CreditCard         *CreditCardType         `xml:"creditCard,omitempty"`
	BankAccount        *BankAccountType        `xml:"bankAccount,omitempty"`
	TrackData          *CreditCardTrackType    `xml:"trackData,omitempty"`
	EncryptedTrackData *EncryptedTrackDataType `xml:"encryptedTrackData,omitempty"`
	PayPal             *PayPalType             `xml:"payPal,omitempty"`
	OpaqueData         *OpaqueDataType         `xml:"opaqueData,omitempty"`
	Emv                *PaymentEmvType         `xml:"emv,omitempty"`
	DataSource         string                  `xml:"dataSource,omitempty"`
}

// CreditCardSimpleType defines common properties for a credit card.
type CreditCardSimpleType struct {
	// CardNumber format should be numeric string or four X's followed by the last four digits.
	CardNumber string `xml:"cardNumber" validation:"required,min=4,max=16"`
	// ExpirationDate format should be gYearMonth (such as 2001-10) or four X's.
	ExpirationDate string `xml:"expirationDate" validation:"required,min=4,max=7"`
}

type CreditCardType struct {
	CreditCardSimpleType
	// CardCode may be passed in for validation, but it will not be stored.
	CardCode string `xml:"cardCode,omitempty" validation:"numeric,min=3,max=4"`
	// IsPaymentToken is to identify whether the CardNumber passed in is a PaymentToken or a real creditCardNumber.
	IsPaymentToken *bool `xml:"isPaymentToken,omitempty"`
	// Cryptogram is needed for one-off payments if the CardNumber passed in is a paymentToken
	Cryptogram string `xml:"cryptogram,omitempty"`
	// TokenRequestorName is only needed for chase pay.
	TokenRequestorName string `xml:"tokenRequestorName,omitempty"`
	// TokenRequestorId is only needed for chase pay.
	TokenRequestorId string `xml:"tokenRequestorId,omitempty"`
	// TokenRequestorEci is only needed for chase pay.
	TokenRequestorEci string `xml:"tokenRequestorEci,omitempty"`
}

type AccountTypeEnum string

const (
	AccountTypeChecking         AccountTypeEnum = "checking"
	AccountTypeSavings                          = "savings"
	AccountTypeBusinessChecking                 = "businessChecking"
)

type BankAccountType struct {
	AccountType AccountTypeEnum `xml:"accountType,omitempty" validation:"oneOf=checking savings businessChecking"`
}

// CreditCardTrackType
//
// The following fields are mutually exclusive and one is required for request: Track1 and Track2. Validation will check
// for one of these being set and will result in a validation error if more than one or none are set.
type CreditCardTrackType struct {
	Track1 string `xml:"track1,omitempty"`
	Track2 string `xml:"track2,omitempty"`
}

type EncodingType string

const (
	EncodingTypeBase64 EncodingType = "Base64"
	EncodingTypeHex                 = "Hex"
)

type EncryptionAlgorithmType string

type KeyManagementScheme struct {
	// TODO Define KeyManagementScheme from scheme
}

const (
	EncryptionAlgorithmTDES EncryptionAlgorithmType = "TDES"
	EncryptionAlgorithmAES                          = "AES"
	EncryptionAlgorithmRSA                          = "RSA"
)

type KeyValue struct {
	Encoding            EncodingType            `xml:"encoding"`
	EncryptionAlgorithm EncryptionAlgorithmType `xml:"encryptionAlgorithm"`
	Scheme              string                  `xml:"scheme"`
}

type KeyBlock struct {
	Value KeyValue `xml:"value"`
}

type EncryptedTrackDataType struct {
	FormatOfPayment KeyBlock `xml:"formatOfPayment"`
}

type PayPalType struct {
	SuccessUrl         string `xml:"successUrl,omitempty" validation:"max=2048"`
	CancelUrl          string `xml:"cancelUrl,omitempty" validation:"max=2048"`
	PayPalLc           string `xml:"payPalLc,omitempty" validation:"max=2"`
	PayPalHdrImg       string `xml:"payPalHdrImg,omitempty" validation:"max=127"`
	PayPalPayflowcolor string `xml:"payPalPayflowcolor,omitempty" validation:"max=6"`
	PayerId            string `xml:"payerId,omitempty" validation:"max=255"`
}

type OpaqueDataType struct {
	DataDescriptor      string     `xml:"dataDescriptor" validation:"required"`
	DataValue           string     `xml:"dataValue" validation:"required"`
	DataKey             string     `xml:"dataKey,omitempty"`
	ExpirationTimestamp *time.Time `xml:"expirationTimestamp,omitempty"`
}

type PaymentEmvType struct {
	EmvData       string `xml:"emvData" validation:"required"`
	EmvDescriptor string `xml:"emvDescriptor" validation:"required"`
	EmvVersion    string `xml:"emvVersion" validation:"required"`
}

type PaymentProfile struct {
	PaymentProfileId string `xml:"paymentProfileId" validation:"required,numeric"`
	CardCode         string `xml:"cardCode,omitempty" validation:"numeric,min=3,max=4"`
}

type CustomerProfilePaymentType struct {
	CreateProfile     *bool           `xml:"createProfile,omitempty"`
	CustomProfileId   string          `xml:"customProfileId,omitempty" validation:"numeric"`
	PaymentProfile    *PaymentProfile `xml:"paymentProfile,omitempty"`
	ShippingProfileId string          `xml:"shippingProfileId,omitempty" validation:"numeric"`
}

type SolutionType struct {
	Id         string `xml:"id" validation:"required"`
	Name       string `xml:"name,omitempty"`
	VendorName string `xml:"vendorName,omitempty"`
}

type OrderType struct {
	InvoiceNumber                  string     `xml:"invoiceNumber" validation:"max=20"`
	Description                    string     `xml:"description,omitempty" validation:"max=255"`
	DiscountAmount                 *float64   `xml:"discountAmount,omitempty"`
	TaxIsAfterDiscount             *bool      `xml:"taxIsAfterDiscount,omitempty"`
	TotalTaxTypeCode               string     `xml:"totalTaxTypeCode,omitempty" validation:"max=3"`
	PurchaserVATRegistrationNumber string     `xml:"purchaserVATRegistrationNumber,omitempty" validation:"max=21"`
	MerchantVATRegistrationNumber  string     `xml:"merchantVATRegistrationNumber,omitempty" validation:"max=21"`
	VatInvoiceReferenceNumber      string     `xml:"vatInvoiceReferenceNumber,omitempty" validation:"max=15"`
	PurchaserCode                  string     `json:"purchaserCode,omitempty" validation:"max=17"`
	SummaryCommodityCode           string     `xml:"summaryCommodityCode,omitempty" validation:"max=4"`
	PurchaseOrderDateUTC           *time.Time `xml:"purchaseOrderDateUTC,omitempty"`
	SupplierOrderReference         string     `xml:"supplierOrderReference,omitempty" validation:"max=25"`
	AuthorizedContactName          string     `xml:"authorizedContactName,omitempty" validation:"max=36"`
	CardAcceptorRefNumber          string     `xml:"cardAcceptorRefNumber,omitempty" validation:"max=25"`
	AmexDataTAA1                   string     `xml:"amexDataTAA1,omitempty" validation:"max=40"`
	AmexDataTAA2                   string     `xml:"amexDataTAA2,omitempty" validation:"max=40"`
	AmexDataTAA3                   string     `xml:"amexDataTAA3,omitempty" validation:"max=40"`
	AmexDataTAA4                   string     `xml:"amexDataTAA4,omitempty" validation:"max=40"`
}

type ArrayOfLineItem struct {
	LineItem []LineItemType `xml:"lineItem"`
}

type LineItemType struct {
	ItemId                  string   `xml:"itemId" validation:"required,min=1,max=31"`
	Name                    string   `xml:"name" validation:"required,min=1,max=31"`
	Description             string   `xml:"description,omitempty" validation:"max=255"`
	Quantity                float64  `xml:"quantity" validation:"required,min=0.00"`
	UnitPrice               float64  `xml:"unitPrice" validation:"required,min=0.00"`
	Taxable                 *bool    `xml:"taxable,omitempty"`
	UnitOfMeasure           string   `xml:"unitOfMeasure,omitempty" validation:"max=12"`
	TypeOfSupply            string   `xml:"typeOfSupply,omitempty" validation:"max=2"`
	TaxRate                 *float64 `xml:"taxRate,omitempty"`
	TaxAmount               *float64 `xml:"taxAmount,omitempty"`
	NationalTax             *float64 `xml:"nationalTax,omitempty"`
	LocalTax                *float64 `xml:"localTax,omitempty"`
	VatRate                 *float64 `xml:"vatRate,omitempty"`
	AlternateTaxId          string   `xml:"alternateTaxId,omitempty" validation:"max=20"`
	AlternateTaxTypeApplied *bool    `xml:"alternateTaxTypeApplied,omitempty" validation:"max=4"`
	AlternateTaxRate        *float64 `xml:"alternateTaxRate,omitempty"`
	AlternateTaxAmount      *float64 `xml:"alternateTaxAmount,omitempty"`
	TotalAmount             *float64 `xml:"totalAmount,omitempty"`
	CommodityCode           string   `xml:"commodityCode,omitempty" validation:"max=15"`
	ProductCode             string   `xml:"productCode,omitempty" validation:"max=30"`
	ProductSku              string   `xml:"productSku,omitempty" validation:"max=30"`
	DiscountRate            *float64 `xml:"discountRate,omitempty"`
	DiscountAmount          *float64 `xml:"discountAmount,omitempty"`
	TaxIncludedInTotal      *bool    `xml:"taxIncludedInTotal,omitempty"`
	TaxIsAfterDiscount      *bool    `xml:"taxIsAfterDiscount,omitempty"`
}

type ExtendedAmountType struct {
	Amount      float64 `xml:"amount"`
	Name        string  `xml:"name,omitempty" validation:"max=31"`
	Description string  `xml:"description,omitempty" validation:"max=255"`
}

type CustomerTypeEnum string

const (
	CustomerTypeIndividual CustomerTypeEnum = "individual"
	CustomerTypeBusiness   CustomerTypeEnum = "business"
)

type DriversLicenseType struct {
	Number      string `xml:"number" validation:"required,min=5,max=20"`
	State       string `xml:"state" validation:"required,min=1,max=2"`
	DateOfBirth string `xml:"dateOfBirth" validation:"required,min=8,max=10"`
}

type CustomerDataType struct {
	Type           *CustomerTypeEnum   `xml:"type,omitempty"`
	Id             string              `xml:"id,omitempty" validation:"max20"`
	Email          string              `xml:"email,omitempty" validation:"max=255"`
	DriversLicense *DriversLicenseType `xml:"driversLicense,omitempty"`
	TaxId          string              `xml:"taxId,omitempty" validation:"min=8,max=9"`
}

type CustomerAddressType struct {
	NameAndAddressType
	PhoneNumber string `xml:"phoneNumber,omitempty" validation:"max=25"`
	FaxNumber   string `xml:"faxNumber,omitempty" validation:"max=25"`
	Email       string `xml:"email,omitempty"`
}

type NameAndAddressType struct {
	FirstName string `xml:"firstName,omitempty" validation:"max=50"`
	LastName  string `xml:"lastName,omitempty" validation:"max=50"`
	Company   string `xml:"company,omitempty" validation:"max=50"`
	Address   string `xml:"address,omitempty" validation:"max=60"`
	City      string `xml:"city,omitempty" validation:"max=40"`
	State     string `xml:"state,omitempty" validation:"max=40"`
	Zip       string `xml:"zip,omitempty" validation:"max=20"`
	Country   string `xml:"country,omitempty" validation:"max=60"`
}

type CCAuthenticationType struct {
	AuthenticationIndicator       string `xml:"authenticationIndicator" validation:"required"`
	CardholderAuthenticationValue string `xml:"cardholderAuthenticationValue" validation:"required"`
}

type TransRetailInfoType struct {
	MarketType        string `xml:"marketType,omitempty" validation:"required"` // default 2
	DeviceType        string `xml:"deviceType,omitempty"`
	CustomerSignature string `xml:"customerSignature,omitempty"`
	TerminalNumber    string `xml:"terminalNumber,omitempty"`
}

type ArrayOfSetting struct {
	Setting []SettingType `xml:"setting,omitempty"`
}

type SettingType struct {
	SettingName  string `xml:"settingName,omitempty"`
	SettingValue string `xml:"settingValue,omitempty"`
}

type UserFields struct {
	UserField []UserField `xml:"userField,omitempty" validation:"min=0,max=20"`
}

type UserField struct {
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type SubMerchantType struct {
	Identifier                 string `xml:"identifier" validation:"required,max=40"`
	DoingBusinessAs            string `xml:"doingBusinessAs,omitempty" validation:"max=50"`
	PaymentServiceProviderName string `xml:"paymentServiceProviderName,omitempty" validation:"max=40"`
	PaymentServiceFacilitator  string `xml:"paymentServiceFacilitator,omitempty" validation:"max=20"`
	StreetAddress              string `xml:"streetAddress,omitempty" validation:"max=40"`
	Phone                      string `xml:"phone,omitempty" validation:"max=20"`
	Email                      string `xml:"email,omitempty" validation:"max=40"`
	PostalCode                 string `xml:"postalCode,omitempty" validation:"max=20"`
	City                       string `xml:"city,omitempty" validation:"max=30"`
	RegionCode                 string `xml:"regionCode,omitempty" validation:"max=10"`
	CountryCode                string `xml:"countryCode,omitempty" validation:"max=10"`
}

type ProcessingOptions struct {
	IsFirstRecurringPayment *bool `xml:"isFirstRecurringPayment,omitempty"`
	IsFirstSubsequentAuth   *bool `xml:"isFirstSubsequentAuth,omitempty"`
	IsSubsequentAuth        *bool `xml:"isSubsequentAuth,omitempty"`
	IsStoredCredentials     *bool `xml:"isStoredCredentials,omitempty"`
}

type MerchantInitTransReasonEnum string

const (
	MerchantInitTransReasonResubmission    MerchantInitTransReasonEnum = "resubmission"
	MerchantInitTransReasonDelayedCharge                               = "delayedCharge"
	MerchantInitTransReasonReauthorization                             = "reauthorization"
	MerchantInitTransReasonNoShow                                      = "noShow"
)

type SubsequentAuthInformation struct {
	// TODO Customer validator for alpha numeric space string
	OriginalNetworkTransId string                       `xml:"originalNetworkTransId,omitempty" validation:"max=255"`
	OriginalAuthAmount     *float64                     `xml:"originalAuthAmount,omitempty"`
	Reason                 *MerchantInitTransReasonEnum `xml:"reason,omitempty"`
}

type OtherTaxType struct {
	NationalTaxAmount  *float64 `xml:"nationalTaxAmount,omitempty"`
	LocalTaxAmount     *float64 `xml:"localTaxAmount,omitempty"`
	AlternateTaxAmount *float64 `xml:"alternateTaxAmount,omitempty"`
	AlternateTaxId     string   `xml:"alternateTaxId,omitempty" validation:"max=15"`
	VatTaxRate         *float64 `xml:"vatTaxRate,omitempty"`
	VatTaxAmount       *float64 `xml:"vatTaxAmount,omitempty"`
}

type AuthIndicatorEnum string

const (
	AuthIndicatorPre   AuthIndicatorEnum = "pre"
	AuthIndicatorFinal                   = "final"
)

type AuthorizationIndicatorType struct {
	AuthorizationIndicator *AuthIndicatorEnum `xml:"authorizationIndicator,omitempty"`
}

type TransactionRequestType struct {
	AnetApiRequestType
	TransactionType            string                      `xml:"transactionType,omitempty"`
	Amount                     *float64                    `xml:"amount,omitempty"`
	CurrencyCode               string                      `xml:"currencyCode,omitempty"`
	Payment                    *PaymentType                `xml:"payment,omitempty"`
	Profile                    *CustomerProfilePaymentType `xml:"profile,omitempty"`
	Solution                   *SolutionType               `xml:"solution,omitempty"`
	CallId                     string                      `xml:"callId,omitempty"`
	TerminalNumber             string                      `xml:"terminalNumber,omitempty"`
	AuthCode                   string                      `xml:"authCode,omitempty"`
	RefTransId                 string                      `xml:"refTransId,omitempty"`
	SplitTenderId              string                      `xml:"splitTenderId,omitempty"`
	Order                      *OrderType                  `xml:"order,omitempty"`
	LineItems                  *ArrayOfLineItem            `xml:"lineItems,omitempty"`
	Tax                        *ExtendedAmountType         `xml:"tax,omitempty"`
	Duty                       *ExtendedAmountType         `xml:"duty,omitempty"`
	Shipping                   *ExtendedAmountType         `xml:"shipping,omitempty"`
	TaxExempt                  *bool                       `xml:"taxExempt,omitempty"`
	PoNumber                   string                      `xml:"poNumber,omitempty"`
	Customer                   *CustomerDataType           `xml:"customer,omitempty"`
	BillTo                     *CustomerAddressType        `xml:"billTo,omitempty"`
	ShipTo                     *NameAndAddressType         `xml:"shipTo,omitempty"`
	CustomerIp                 string                      `xml:"customerIp,omitempty"`
	CardHolderAuthentication   *CCAuthenticationType       `xml:"cardHolderAuthentication,omitempty"`
	Retail                     *TransRetailInfoType        `xml:"retail,omitempty"`
	EmployeeId                 string                      `xml:"employeeId,omitempty"`
	TransactionSettings        *ArrayOfSetting             `xml:"transactionSettings,omitempty"`
	UserFields                 *UserFields                 `xml:"userFields,omitempty"`
	Surcharge                  *ExtendedAmountType         `xml:"surcharge,omitempty"`
	MerchantDescriptor         string                      `xml:"merchantDescriptor,omitempty"`
	SubMerchant                *SubMerchantType            `xml:"subMerchant,omitempty"`
	Tip                        *ExtendedAmountType         `xml:"tip,omitempty"`
	ProcessingOptions          *ProcessingOptions          `xml:"processingOptions,omitempty"`
	SubsequentAuthInformation  *SubsequentAuthInformation  `xml:"subsequentAuthInformation,omitempty"`
	OtherTax                   *OtherTaxType               `xml:"otherTax,omitempty"`
	ShipFrom                   *NameAndAddressType         `xml:"shipFrom,omitempty"`
	AuthorizationIndicatorType *AuthorizationIndicatorType `xml:"authorizationIndicatorType,omitempty"`
}

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
