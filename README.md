AimMatic SDKs are easiest and best supported way for most developers to use AimMatic APIs.

### placeNext SDK ###

placeNext SDK is a client library to help developers quickly deploy applications with programmatic connections to placeNext Rest APIs.

### Getting Started ###

```
go get github.com/aimmatic/aimmatic-go-sdk-placenext
```

If you're using dep as dependency management

```
dep ensure -add github.com/aimmatic/aimmatic-go-sdk-placenext
```

#### Paired Key Connections ####

This SDK requires an API Key and Secret Key to establish a secure connection to placeNext.
To set-up your application's API Key and Secret Key.

**Using a variable environment**

Set the variable environment PLACENEXT_APIKEY and PLACENEXT_SECRETKEY then
create core rest api to access with our api.

```go
restApi := core.NewRestApi(rest.DefaultClient())
restApi.V1().GetNSS()
```

**Setup ApiKey and SecretKey at runtime**

```go
config, _ := rest.NewConfig("Your Api Key", "Your Secret Key")
restApi := core.NewRestApi(rest.NewRestClient(config))
restApi.V1().GetNSS()
```