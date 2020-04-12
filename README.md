[![Build Status](https://github.com/LuckyLukas/whenthengo/workflows/check-build/badge.svg)](https://github.com/LuckyLukas/whenthengo/actions) [![Release Integrationtest Status](https://github.com/LuckyLukas/whenthengo/workflows/release/badge.svg)](https://github.com/LuckyLukas/whenthengo/actions)
# whenthengo

a simple mock http server to use for testing with packages like
[testcontainers-go](https://github.com/testcontainers/testcontainers-go).`

## what?

whenthengo is configured using _whens_ and _thens_, if the app recognizes an http request to match a _when's_ conditions, it will response with the contents present in the matching _then_.

## limitations for 0.0.x

- all headers have to be passed as string values
- simplistic matching, no support for fuzzy matches or best effort responses, _whens_ that depend on whitespace positions will likely fail

### whens

a when is a set of conditions a request must meet to make the server return the matching _then_.

### thens

are linked to one _when_ and describe the expected response.

## parameters

### whens

| property        | type           | desciption  |
| ------------- |-------------| -----|
| method     | string| case insensitive http verb to match|
| url     | string      |   the path to match |
| headers | map[string]string      |    a map of headers the request should include. This checks for "containment" and is case insensitive (eg.: ```"application/json; encoding=UTF-8"``` with match ```"Application/json"```) |
| body| string | the request body to match. this will remove all whitespaces when checking for equiality.|

### thens
| property        | type           | desciption  |
| ------------- |-------------| -----|
| delay     | int| artificial delay for the response in milliseconds |
| status     | int      |   http status code |
| headers | map[string]string      |    a map of headers the response will include |
| body| string | the expected response body|


## input formats

currently you can specify when-thens in ```yaml```
 and ```json``` format

### json

```json
[
  {
    "when": {
      "method": "get",
      "url": "/path/test",
      "headers": {
        "Accept": "application/json",
        "Content-Length": "6"
      },
      "body": "some body"
    },
    "then": {
      "delay": 100,
      "status": 200,
      "headers": {
        "Content-Length": "1"
      },
      "body": "k"
    }
  },
  ...
]
```

### yaml

```
    -
      when:
        method: "get"
        url: "/path/test"
        headers:
          Accept: "application/json"
          Content-Length: 6
        body: |+
          abc
          def
      then:
        delay: 100
        status: 200
        headers:
          Content-Length: 1
        body: "k"
    -
        ...

```

## health

whenthengo provides an endpoint under ```/whenthengoup```
 that just returns an empty ```200 OK```.
 
 This can be used to check if the app is up and ready to handle requests.



## bugs and collaboration

yes and please.
