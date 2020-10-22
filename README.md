zuliprus
========

Zulip hook for [Logrus](https://github.com/sirupsen/logrus). 

## Use

```go
package main

import (
	"os"

	logrus "github.com/sirupsen/logrus"
	"github.com/tiger5226/zuliprus"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stderr)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.AddHook(&zuliprus.ZuliprusHook{
		APIURL:         "https://zulip.mycompany.com/api/v1/",
		APIKey:         "aP8vzq5gwfZHBd4V6ztcYzO4Jugczgt6",
		Email:          "my-bot@zulip.mycompany.com",
		AcceptedLevels: zuliprus.LevelThreshold(logrus.DebugLevel),
		Stream:         "mystream",
		Topic:          "that-topic",
	})

	logrus.Debug("this is a debug level message")
	logrus.Info("this is an info level message")
	logrus.Error("this is an error level message")
	logrus.Warning("this is a warning level message")

}

```

## Parameters

#### Required
  * APIURL
  * APIKey
  * Email
  * Stream
  
#### Optional
  * AcceptedLevels
  * UserEmails
  * Topic
  * Asynchronous
  * Disabled
  * FormatFn
## Installation

    go get github.com/tiger5226/zuliprus

## Credits

Based on slackrus hook by [johntdyer](https://github.com/johntdyer/zuliprus)
