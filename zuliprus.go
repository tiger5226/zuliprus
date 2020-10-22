// Package zuliprus provides a Slack hook for the logrus loggin package.
package zuliprus

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ifo/gozulipbot"
	"github.com/sirupsen/logrus"
)

// Project version
const (
	VERISON = "0.0.3"
)

// ZuliprusHook is a logrus Hook for dispatching messages to the specified
// channel on Slack.
type ZuliprusHook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If nil, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	APIURL         string
	APIKey         string
	Email          string
	Stream         string
	UserEmails     []string
	Topic          string
	Asynchronous   bool
	Disabled       bool
	FormatFn       FmtFunction
}

type FmtFunction func(e *logrus.Entry) string

var MsgFmtFn = func(e *logrus.Entry) string {
	return fmt.Sprintf("%s>  *%s*  : %s", levelPrefix(e), e.Time.Format("2006-01-02T15:04:05"), e.Message)
}

// Levels sets which levels to sent to slack
func (sh *ZuliprusHook) Levels() []logrus.Level {
	if sh.AcceptedLevels == nil {
		return AllLevels
	}
	return sh.AcceptedLevels
}

// Fire -  Send event to zulip
func (sh *ZuliprusHook) Fire(e *logrus.Entry) error {
	if sh.Disabled {
		return nil
	}

	msg := gozulipbot.Message{
		Stream:  sh.Stream,
		Topic:   sh.Topic,
		Emails:  sh.UserEmails,
		Content: MsgFmtFn(e),
	}
	// If there are fields they will not be send, special handling?
	if sh.Asynchronous {
		go sh.sendMessage(msg)
		return nil
	}
	return sh.sendMessage(msg)
}

func levelPrefix(e *logrus.Entry) string {
	color := ""
	switch e.Level {
	case logrus.DebugLevel:
		color = "```c++\nDEBUG\n```\n"
	case logrus.InfoLevel:
		color = "```bash\nINFO\n```\n"
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = "```terraform\nERROR\n```\n"
	default:
		color = "```apacheconf\nWARN\n```\n"
	}

	return color
}

func (sh *ZuliprusHook) newClient() *gozulipbot.Bot {
	return &gozulipbot.Bot{
		APIKey: sh.APIKey,
		APIURL: sh.APIURL,
		Email:  sh.Email,
		Client: &http.Client{},
	}
}

func (sh *ZuliprusHook) sendMessage(message gozulipbot.Message) error {
	bot := sh.newClient()
	resp, err := bot.Message(message)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("ZulipError: %d %s", resp.StatusCode, t))
	}
	return nil
}
