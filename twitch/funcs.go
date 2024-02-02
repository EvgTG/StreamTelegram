package twitch

import (
	"fmt"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/rotisserie/eris"
)

func (t *Twitch) GetStream(nick string) (*helix.Stream, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.ClientOK() || !t.AuthOK() {
		return nil, nil
	}

	resp, err := t.client.GetStreams(&helix.StreamsParams{UserLogins: []string{nick}})
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			time.Sleep(time.Second * 10)

			resp, err = t.client.GetStreams(&helix.StreamsParams{UserLogins: []string{nick}})
			if err != nil {
				return nil, eris.Wrap(err, "t.client.GetStreams()")
			}
		}

		return nil, eris.Wrap(err, "t.client.GetStreams()")
	}

	if resp.StatusCode == 401 {
		err = t.refreshAuth()
		if err != nil {
			return nil, eris.Wrap(err, "t.refreshAuth()")
		}

		resp, err = t.client.GetStreams(&helix.StreamsParams{UserLogins: []string{nick}})
		if err != nil {
			return nil, eris.Wrap(err, "t.client.GetStreams()")
		}
	}
	if resp.StatusCode != 200 {
		return nil, eris.New(fmt.Sprintf("resp: %+v\n", resp))
	}

	if len(resp.Data.Streams) < 1 {
		return nil, nil
	}

	return &resp.Data.Streams[0], nil
}
