package util

import (
	"github.com/rotisserie/eris"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func GetChannelIDByUrl(url string) (string, error) {
	ok, err := regexp.MatchString("^https://www.youtube.com/(channel|c|user)/", url)
	if err != nil {
		return "", eris.Wrap(err, "regexp.MatchString()")
	}
	if !ok {
		return "", eris.New("404")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", eris.Wrap(err, "http.Get(url)")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", eris.Wrap(err, "ioutil.ReadAll()")
	}

	page := string(bytes)
	search := []string{"\"channelId\":\"", "\","}
	i1 := strings.Index(page, search[0])
	if i1 < 0 {
		return "", eris.New("404")
	}
	i1 += len(search[0])
	i2 := strings.Index(page[i1:], search[1])
	if i2 < 0 {
		return "", eris.New("404")
	}
	i2 += i1

	return page[i1:i2], nil
}
