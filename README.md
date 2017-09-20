[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Issue Count](https://lima.codeclimate.com/github/fairway-corp/swagchat-api/badges/issue_count.svg)](https://lima.codeclimate.com/github/fairway-corp/swagchat-realtime)
[![Go Report Card](https://goreportcard.com/badge/github.com/fairway-corp/swagchat-api)](https://goreportcard.com/report/github.com/fairway-corp/swagchat-realtime)



# swagchat Real Time Messaging

swagchat is an open source chat components for your webapps.

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")

##### Related repositories

* [Chat API](https://github.com/fairway-corp/swagchat-chat-api)
* [SDK (TypeScript & JavaScript)](https://github.com/swagchat/swagchat-sdk-js)
* [UIKit (A set of React components)](https://github.com/swagchat/react-swagchat)

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
