package minidb

func (mini *MiniDB) SetTwitchNick(nick string) error {
	return mini.write("twitchnick", nick)
}

func (mini *MiniDB) GetTwitchNick() (string, error) {
	var nick string
	err := mini.read("twitchnick", &nick)
	return nick, err
}

func (mini *MiniDB) SetCliID(clientID string) error {
	return mini.write("clientID", clientID)
}

func (mini *MiniDB) GetCliID() (string, error) {
	var clientID string
	err := mini.read("clientID", &clientID)
	return clientID, err
}

func (mini *MiniDB) SetCliSecret(clientSecret string) error {
	return mini.write("clientSecret", clientSecret)
}

func (mini *MiniDB) GetCliSecret() (string, error) {
	var clientSecret string
	err := mini.read("clientSecret", &clientSecret)
	return clientSecret, err
}

func (mini *MiniDB) SetAccessToken(accessToken string) error {
	return mini.write("accessToken", accessToken)
}

func (mini *MiniDB) GetAccessToken() (string, error) {
	var accessToken string
	err := mini.read("accessToken", &accessToken)
	return accessToken, err
}

func (mini *MiniDB) SetRefreshToken(refreshToken string) error {
	return mini.write("refreshToken", refreshToken)
}

func (mini *MiniDB) GetRefreshToken() (string, error) {
	var refreshToken string
	err := mini.read("refreshToken", &refreshToken)
	return refreshToken, err
}
