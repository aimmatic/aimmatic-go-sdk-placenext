### placeNext SDK ###

placeNext SDK a library that help developer to quick deploy an application
that connected to Placenext Rest API.

### Usage ###

The SDK is required api key and secret key to establish connection to placeNext Api.
There are 2 way to setup api key and secret key.

1. Using variable environment PLACENEXT_APIKEY and PLACENEXT_SECRETKEY

2. Using the below code in your application

**Using variable environment**

```go
restApi := core.NewRespApi(rest.DefaultClient())
restApi.V1().IngestGeometry(....)
```

**Setup ApiKey and SecretKey at runtime Globally**

```go
config, err := rest.NewConfig("Your Api Key", "Your Secret Key")
rest.SetConfig(config)
restApi := core.NewRespApi(rest.DefaultClient())
restApi.V1().IngestGeometry(...)
```
