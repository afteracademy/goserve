package micro

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

func RequestNats[S any, R any](client NatsClient, subject string, sData *S) (*R, error) {
	msgJson, err := MsgToJson(sData)
	if err != nil {
		return nil, err
	}

	natsMsg, err := client.GetInstance().Conn.Request(subject, msgJson, client.GetInstance().Timeout)
	if err != nil {
		return nil, err
	}

	rData, err := JsonToMsg[R](natsMsg.Data)
	if err != nil {
		return rData, err
	}

	return rData, err
}

func RequestNatsRaw[S any, R any](client NatsClient, subject string, sData *S) (*R, *nats.Msg, error) {
	jData, err := json.Marshal(sData)
	if err != nil {
		return nil, nil, err
	}

	natsMsg, err := client.GetInstance().Conn.Request(subject, jData, client.GetInstance().Timeout)
	if err != nil {
		return nil, natsMsg, err
	}

	var rData R
	err = json.Unmarshal(natsMsg.Data, &rData)
	if err != nil {
		return &rData, natsMsg, err
	}

	return &rData, natsMsg, err
}
