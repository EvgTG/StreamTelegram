package minidb

//YouTube

// true  - есть в базе, идем дальше
// false - новый стрим, добавляем в базу
func (mini *MiniDB) CheckVideo(id string) (bool, error) {
	for _, a := range mini.videoIDs {
		if id == a {
			return true, nil
		}
	}

	if len(mini.videoIDs) >= 50 {
		mini.videoIDs = append(mini.videoIDs[1:50], id)
	} else {
		mini.videoIDs = append(mini.videoIDs, id)
	}

	return false, mini.SetVideoIDs()
}

func (mini *MiniDB) SetVideoIDs() error {
	return mini.write("vids", mini.videoIDs)
}

func (mini *MiniDB) GetVideoIDs() error {
	return mini.read("vids", &mini.videoIDs)
}

// Twitch

// true  - есть в базе, идем дальше
// false - новый стрим, добавляем в базу
func (mini *MiniDB) CheckTwitchVideo(id string) (bool, error) {
	for _, a := range mini.tvideoIDs {
		if id == a {
			return true, nil
		}
	}

	if len(mini.tvideoIDs) >= 50 {
		mini.tvideoIDs = append(mini.tvideoIDs[1:50], id)
	} else {
		mini.tvideoIDs = append(mini.tvideoIDs, id)
	}

	return false, mini.SetTwitchVideoIDs()
}

func (mini *MiniDB) SetTwitchVideoIDs() error {
	return mini.write("tvids", mini.tvideoIDs)
}

func (mini *MiniDB) GetTwitchVideoIDs() error {
	return mini.read("tvids", &mini.tvideoIDs)
}
