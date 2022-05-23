package log

var errsNotif ErrsNotifications

type ErrsNotifications struct {
	i int
}

func GetErrN() int {
	x := errsNotif.i
	errsNotif.i = 0
	return x
}
