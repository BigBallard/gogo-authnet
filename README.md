# GoGo Authnet

A library that interfaces with [Authorize.net](https://authorize.net) which is a transaction authorization solution.

***

> While a Go library already exists for this, the aim of GoGo Authnet is to provide a more complete and updated solution.
> This API will give the user more freedom to tailor requests to their business model.

## Getting Started

Start by importing the package:

```
go get github.com/BigBallard/gogo-authnet
```

```go
import "github.com/BigBallard/gogo-authnet
```

### Configuration
The GoGo Authnet client has multiple configuration options:

- Configuration File
- CLI Arguments
- Environment Variables.

The order of these options are in order of least precedence to most, meaning that a value set in the configuration file
can be overwritten by the value set in the environment variables. Refer to the [configuration documentation](docs/CONFIG.md)
for available settings.

### Basic Example

In this trivial example, we load the configuration and create a new `AuthNetClient`. It is wise to always ensure that
the intended credentials loaded correctly by running `AuthenticateTest`. An error returns if the credentials don't work.

After testing authentication, we then create a zero value variable `GetCustomerProfileResponse`. This will be used shortly.
We next create a `GetCustomerProfileRequest`. All **Request**'s and **Response**'s types embed the `ANetApiRequest` and
`ANetApiResponse` types respectfully.

> Requests require the `MerchantAuthentication` which the client provides a
> convenience function `CreateMerchantAuthenticationType` which will have the values populated with the **API Login ID**
> and **Transaction Key**.

We then finish creating our response and send the request using the client. We pass the request, which can be a reference
or value type, and the response which must be a reference type.

If no error is returned from the `SendRequest` call, then the response should be populated with the appropriate values.

```go
import (
    "github.com/BigBallard/gogo-authnet/config"
    "github.com/BigBallard/gogo-authnet/client"
    "github.com/BigBallard/gogo-authnet/common"
    "github.com/BigBallard/gogo-authnet/util"
)

func main() {
    conf, loadErr := config.LoadConfigFromEnv()
    if loadErr != nil {
		panic(loadErr)
    }
	
    client := client.NewAuthNetClient(*conf)
	
    if _, authErr := client.AuthenticateTest(); authErr != nil {
        pani(authErr)
    }

    var response common.GetCustomerProfileResponse
	
    request := common.GetCustomerProfileRequest{
        ANetApiRequest: common.ANetApiRequest{
            MerchantAuthentication: client.CreateMerchantAuthenticationType(),
        },
        CustomerProfileId: "10000",
        IncludeIssuerInfo: util.BoolTrueRef(),
    }

    if rErr := client.SendRequest(request, &response); rErr != nil {
        panic(rErr)
    }
    ...
}
```

The `SendRequest` function returns a possible `RequestError` which holds either an `error`, `ErrorResponse`, or both.
The `ErrorRepsonse` is an Authorize.net type that provides more detained information about the reason for the failed
request.

## Development

To run the tests and develop gogo-authnet, you will first need to acquire sandbox credentials which you can get by 
creating a sandbox account [here](https://developer.authorize.net/hello_world/sandbox.html).

Once you have credentials, place them in the configuration file prior to running any tests or using your preferred
configuration method.