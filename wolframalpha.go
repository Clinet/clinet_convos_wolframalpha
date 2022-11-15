package wolframalpha

import (
	"errors"
	"strings"

	"github.com/Clinet/clinet_convos"
	"github.com/Clinet/clinet_storage"
	"github.com/JoshuaDoes/go-wolfram"
)

var WolframAlpha *ClientWolframAlpha

type ClientWolframAlpha struct {
	Client *wolfram.Client
}

func (wa *ClientWolframAlpha) Login() error {
	cfg := &storage.Storage{}
	if err := cfg.LoadFrom("wolframalpha"); err != nil {
		return err
	}
	appID, err := cfg.ExtraGet("cfg", "appID")
	if err != nil {
		return err
	}
	WolframAlpha = &ClientWolframAlpha{
		Client: &wolfram.Client{
			AppID: appID.(string),
		},
	}
	return nil
}

func (wa *ClientWolframAlpha) Query(query *convos.ConversationQuery, lastState *convos.ConversationState) (*convos.ConversationResponse, error) {
	resp := &convos.ConversationResponse{}
	if lastState != nil {
		resp.WolframAlpha = lastState.Response.WolframAlpha
	}

	wolframConvo, err := wa.Client.GetConversationalQuery(query.Text, wolfram.Metric, resp.WolframAlpha)
	if err != nil {
		return nil, err
	}

	if wolframConvo.ErrorMessage != "" {
		return nil, errors.New("wolframalpha: " + wolframConvo.ErrorMessage)
	}

	if wolframConvo.Result == "" {
		return nil, errors.New("wolframalpha: empty result")
	}

	if !strings.HasSuffix(wolframConvo.Result, ".") {
		wolframConvo.Result += "."
	}

	resp.TextSimple = wolframConvo.Result
	resp.WolframAlpha = wolframConvo

	return resp, nil
}