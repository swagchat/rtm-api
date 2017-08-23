[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Issue Count](https://lima.codeclimate.com/github/fairway-corp/swagchat-api/badges/issue_count.svg)](https://lima.codeclimate.com/github/fairway-corp/swagchat-realtime)
[![Go Report Card](https://goreportcard.com/badge/github.com/fairway-corp/swagchat-api)](https://goreportcard.com/report/github.com/fairway-corp/swagchat-realtime)



# SwagChat Realtime Messaging (Not For Production Use!)

SwagChat is an open source chat components for your webapps.

* **Easy to deploy**
* **Easy to customize**
* **Easy to scale**

## Components

* [RESTful API Server (Go)](https://github.com/fairway-corp/swagchat-api)
* **Realtime Messaging (Go) ---> This repository**
* [Client SDK (TypeScript & JavaScript)](https://github.com/fairway-corp/swagchat-sdk)
* [UIKit (Typescript - React)](https://github.com/fairway-corp/react-swagchat)


## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/architecture-201703011307.png "Architecture")

## Quick start

### Just run the executable binary

You can download binary from [Release page](https://github.com/fairway-corp/swagchat-realtime/releases)

```
# In the case of macOS (Default port is 9100)
./swagchat-realtime_darwin_amd64


# You can also specify the port
./swagchat-realtime_darwin_amd64 -port 9200
```

### docker

```
docker pull swagchat/realtime
docker run swagchat/realtime
```

[Docker repository](https://hub.docker.com/r/swagchat/realtime/)

### heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## License

MIT License.
