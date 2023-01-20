package push

import (
	"fmt"
	"strings"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
)

// Sender implements message sending
type Sender interface {
	Send(title, msg string)
}

var log = util.NewLogger("push")

// NewMessengerFromConfig creates a new messenger
func NewMessengerFromConfig(typ string, other map[string]interface{}) (res Sender, err error) {
	switch strings.ToLower(typ) {
	case "pushover":
		var cc pushOverConfig
		if err = util.DecodeOther(other, &cc); err == nil {
			res, err = NewPushOverMessenger(cc.App, cc.Recipients)
		}
	case "telegram":
		var cc telegramConfig
		if err = util.DecodeOther(other, &cc); err == nil {
			res, err = NewTelegramMessenger(cc.Token, cc.Chats)
		}
	case "email", "shout":
		var cc shoutrrrConfig
		if err = util.DecodeOther(other, &cc); err == nil {
			res, err = NewShoutrrrMessenger(cc.URI)
		}
	case "script":
		cc := scriptConfig{
			Timeout: request.Timeout,
		}
		if err = util.DecodeOther(other, &cc); err == nil {
			res, err = NewScriptMessenger(cc.CmdLine, cc.Timeout)
		}
	case "ntfy":
		var cc ntfyConfig
		if err = util.DecodeOther(other, &cc); err == nil {
			res, err = NewNtfyMessenger(cc.URI, cc.Priority, cc.Tags)
		}
	default:
		err = fmt.Errorf("unknown messenger type: %s", typ)
	}

	return res, err
}
