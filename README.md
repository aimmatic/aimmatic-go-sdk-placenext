AimMatic SDKs are easiest and best supported way for most developers to use AimMatic APIs.

### placeNext SDK ###

placeNext SDK is a client library to help developers quickly deploy applications with programmatic connections to placeNext Rest APIs.

### Paired Key Connections ###

This SDK requires an API Key and Secret Key to establish a secure connection to placeNext.
There are two ways to set-up your application's API Key and Secret Key.

1. Use a variable environment PLACENEXT_APIKEY and PLACENEXT_SECRETKEY

**Using a variable environment**

```go
restApi := core.NewRespApi(rest.DefaultClient())
restApi.V1().IngestGeometry(....)
```

2. Use the below code in your application

**Setup ApiKey and SecretKey for runtime Globally**

```go
config, err := rest.NewConfig("Your Api Key", "Your Secret Key")
rest.SetConfig(config)
restApi := core.NewRespApi(rest.DefaultClient())
restApi.V1().IngestGeometry(...)
```
