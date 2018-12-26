[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/swagchat/rtm-api)](https://goreportcard.com/report/github.com/swagchat/rtm-api)

# swagchat rtm-api

swagchat is an open source chat components for your webapps.

rtm-api is designed to be easy to introduce to your microservices as well.

**Currently developing for version 1**

## Architecture

![Architecture](https://client.fairway.ne.jp/swagchat/img/swagchat-start-guide-20170920.png "Architecture")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fswagchat%2Frtm-api.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fswagchat%2Frtm-api?ref=badge_shield)

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

## Configuration

### Specify the setting file (yaml format)

To override the default configuration options, make a copy of `defaultConfig.yaml` and then specify that file name in runtime parameter `config` and execute.

```
./rtm-api -config myConfig.yaml
```

### Specify environment variables

You can overwrite it with environment variable.

```
export HTTP_PORT=80 && ./rtm-api
```

### Specify runtime parameters

You can overwrite it with runtime parameters.

```
./rtm-api -httpPort 80
```

You can check the variables that can be set with the help command of the executable binary.

```
./rtm-api -h
```

## License

MIT License.


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fswagchat%2Frtm-api.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fswagchat%2Frtm-api?ref=badge_large)