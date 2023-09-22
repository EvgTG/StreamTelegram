package twitch

import (
	"fmt"
	"github.com/nicklaw5/helix/v2"
	"github.com/rotisserie/eris"
	"streamtg/minidb"
	"sync"
)

type Twitch struct {
	clientID, clientSecret    string
	accessToken, refreshToken string

	client *helix.Client

	db *minidb.MiniDB

	mutex sync.Mutex
}

func NewTwitch(db *minidb.MiniDB) (*Twitch, error) {
	var err error
	t := Twitch{}
	t.db = db

	err = t.reloadClient_1()
	if err != nil {
		return nil, eris.Wrap(err, "t.reloadClient_1()")
	}
	err = t.reloadClient_2()
	if err != nil {
		return nil, eris.Wrap(err, "t.reloadClient_2()")
	}

	return &t, nil
}

func (t *Twitch) reloadClient_1() error {
	var err error

	if t.clientID == "" {
		t.clientID, err = t.db.GetCliID()
		if err != nil {
			return eris.Wrap(err, "GetCliID")
		}
	}
	if t.clientSecret == "" {
		t.clientSecret, err = t.db.GetCliSecret()
		if err != nil {
			return eris.Wrap(err, "GetCliSecret")
		}
	}

	if !t.ClientOK() {
		return nil
	}

	t.client, err = helix.NewClient(&helix.Options{
		ClientID:     t.clientID,
		ClientSecret: t.clientSecret,
		RedirectURI:  "http://localhost",
	})
	if err != nil {
		return eris.Wrap(err, "helix.NewClient()")
	}

	return nil
}

func (t *Twitch) reloadClient_2() error {
	var err error

	if t.accessToken == "" {
		t.accessToken, err = t.db.GetAccessToken()
		if err != nil {
			return eris.Wrap(err, "GetAccessToken")
		}
	}
	if t.refreshToken == "" {
		t.refreshToken, err = t.db.GetRefreshToken()
		if err != nil {
			return eris.Wrap(err, "GetRefreshToken")
		}
	}

	if !t.AuthOK() {
		return nil
	}

	t.client.SetUserAccessToken(t.accessToken)

	return nil
}

func (t *Twitch) ClientOK() bool {
	return t.clientID != "" && t.clientSecret != ""
}

func (t *Twitch) AuthOK() bool {
	return t.accessToken != "" && t.refreshToken != ""
}

func (t *Twitch) SetClient(id, secret string) error {
	var err error

	t.clientID = id
	t.clientSecret = secret
	t.accessToken = ""
	t.refreshToken = ""

	err = t.db.SetCliID(id)
	if err != nil {
		return eris.Wrap(err, "t.db.SetCliID()")
	}
	err = t.db.SetCliSecret(secret)
	if err != nil {
		return eris.Wrap(err, "t.db.SetCliSecret()")
	}
	err = t.db.SetAccessToken("")
	if err != nil {
		return eris.Wrap(err, "t.db.SetAccessToken()")
	}
	err = t.db.SetRefreshToken("")
	if err != nil {
		return eris.Wrap(err, "t.db.SetRefreshToken()")
	}

	err = t.reloadClient_1()
	if err != nil {
		return eris.Wrap(err, "t.reloadClient_1()")
	}

	return nil
}

func (t *Twitch) GetAuthURL() (string, error) {
	if t.clientID == "" || t.clientSecret == "" {
		return "", eris.New("clientID or clientSecret nil")
	}

	return t.client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       []string{"user:read:email"},
		State:        "some-state",
		ForceVerify:  false,
	}), nil
}

func (t *Twitch) SetCode(code string) error {
	var err error

	resp, err := t.client.RequestUserAccessToken(code)
	if err != nil {
		return eris.Wrap(err, "t.client.RequestUserAccessToken(code)")
	}

	if resp.Data.AccessToken == "" || resp.Data.RefreshToken == "" || resp.StatusCode != 200 || resp.Error != "" || resp.ErrorMessage != "" {
		return eris.New(fmt.Sprintf("resp: %+v\n", resp))
	}

	t.accessToken = resp.Data.AccessToken
	t.refreshToken = resp.Data.RefreshToken

	err = t.db.SetAccessToken(resp.Data.AccessToken)
	if err != nil {
		return eris.Wrap(err, "t.db.SetAccessToken()")
	}
	err = t.db.SetRefreshToken(resp.Data.RefreshToken)
	if err != nil {
		return eris.Wrap(err, "t.db.SetRefreshToken()")
	}

	t.client.SetUserAccessToken(t.accessToken)

	return nil
}

func (t *Twitch) refreshAuth() error {
	var err error

	resp, err := t.client.RefreshUserAccessToken(t.refreshToken)
	if err != nil {
		return eris.Wrap(err, "t.client.RefreshUserAccessToken()")
	}

	if resp.Data.AccessToken == "" || resp.Data.RefreshToken == "" || resp.StatusCode != 200 || resp.Error != "" || resp.ErrorMessage != "" {
		return eris.New(fmt.Sprintf("resp: %+v\n", resp))
	}

	t.accessToken = resp.Data.AccessToken
	t.refreshToken = resp.Data.RefreshToken

	err = t.db.SetAccessToken(resp.Data.AccessToken)
	if err != nil {
		return eris.Wrap(err, "t.db.SetAccessToken()")
	}
	err = t.db.SetRefreshToken(resp.Data.RefreshToken)
	if err != nil {
		return eris.Wrap(err, "t.db.SetRefreshToken()")
	}

	return nil
}

/*
func (t *Twitch) a()  {

}
*/
