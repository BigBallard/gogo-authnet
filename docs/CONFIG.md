# Configuration Properties

The following are a list of properties and their environment/CLI equivalents.

> Config properties that have a colon `:` denote objects and properties within the object. (ie object:prop equates to 
> {"object":{"property": ""}})

### Authorize.net Host
Config: **authnet-host**

Env/CLI: **AUTHNET_HOST**

Either `prod` or `dev`.
- `prod`: Requests will be made to https://api.authorize.net
- `dev`: Requests will be made to https://apitest.authorize.net

### API Login ID
Config: **auth:api-login-id**

Env/CLI: **AUTH_API_LOGIN_ID**

The login ID provided by Authorize.net for making calls to their API. Used in conjunction with 
`Transaction Key`.Will require an account even for development and testing.

### Transaction Key
Config: **auth:transaction-key**

Env/CLI: **AUTH_TRANSACTION_KEY**

The authorization key provided by Authorize.net for making calls to their API. Used in 
conjunction with `API Login ID`.
