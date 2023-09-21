package minidb

func (mini *MiniDB) SetChannelID(id string) error {
	return mini.write("id", id)
}

func (mini *MiniDB) GetChannelID() (string, error) {
	var id string
	err := mini.read("id", &id)
	return id, err
}

func (mini *MiniDB) SetCycleDuration(dur int) error {
	return mini.write("dur", dur)
}

func (mini *MiniDB) GetCycleDuration() (int, error) {
	var dur int
	err := mini.read("dur", &dur)
	return dur, err
}

func (mini *MiniDB) SetLocs(locs []string) error {
	return mini.write("locs", locs)
}

func (mini *MiniDB) GetLocs() ([]string, error) {
	var locs []string
	err := mini.read("locs", &locs)
	return locs, err
}

func (mini *MiniDB) SetTimeWithCity(bl bool) error {
	return mini.write("timecity", bl)
}

func (mini *MiniDB) GetTimeWithCity() (bool, error) {
	var bl bool
	err := mini.read("timecity", &bl)
	return bl, err
}

type Channel struct {
	ID          int64
	EndOfStream bool
}

func (mini *MiniDB) SetNotifyList(list []Channel) error {
	return mini.write("notifylist", list)
}

func (mini *MiniDB) GetNotifyList() ([]Channel, error) {
	var list []Channel
	err := mini.read("notifylist", &list)
	return list, err
}
