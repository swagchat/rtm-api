[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/swagchat/rtm-api)](https://goreportcard.com/report/github.com/swagchat/rtm-api)

# swagchat Real Time Messaging

swagchat is an open source chat components for your webapps.

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")

##### Related repositories

* [Chat API](https://github.com/swagchat/chat-api)
* [SDK (TypeScript & JavaScript)](https://github.com/swagchat/swagchat-sdk-js)
* [UIKit (A set of React components)](https://github.com/swagchat/react-swagchat)

## Quick start

### Just run the executable binary

You can download binary from [Release page](https://github.com/swagchat/rtm-api/releases)

```
# In the case of macOS (Default port is 9100)
./swagchat-rtm-api_darwin_amd64


# You can also specify the port
./swagchat-rtm-api_darwin_amd64 -port 9200
```

### docker

```
docker run swagchat/rtm-api
```

[Docker repository](https://hub.docker.com/r/swagchat/rtm-api/)

### heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## License

MIT License.
