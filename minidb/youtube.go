package minidb

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
