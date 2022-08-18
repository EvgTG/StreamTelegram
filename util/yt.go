package util

import (
	"github.com/rotisserie/eris"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	Err404   = "404"
	Video    = "video"
	Upcoming = "upcoming"
	Live     = "live"
	LiveGo   = "live_go"
	End      = "end"
	End404   = "end404"
)

// 404 video upcoming live end
//TODO добавить wait статус (время стрима пришло, но не началось)
func TypeVideo(url string) (string, *time.Time, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, eris.Wrap(err, "http.Get(url)")
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, eris.Wrap(err, "ioutil.ReadAll()")
	}

	page := string(bs)

	if strings.Contains(page, `{"iconType":"ERROR_OUTLINE"}`) {
		return "404", nil, nil
	}
	if !strings.Contains(page, `<span itemprop="publication" itemscope itemtype="http://schema.org/BroadcastEvent">`) {
		return "video", nil, nil
	}

	tmStrings := []string{`<meta itemprop="startDate" content="`, `">`}
	iTime1 := strings.Index(page, tmStrings[0])
	if iTime1 < 0 {
		return "", nil, eris.New("startDate 404")
	}
	iTime1 += len(tmStrings[0])
	iTime2 := strings.Index(page[iTime1:], tmStrings[1])
	if iTime2 < 0 {
		return "", nil, eris.New("startDate 404")
	}
	iTime2 += iTime1

	tm, err := time.Parse(time.RFC3339, page[iTime1:iTime2])
	if err != nil {
		return "", nil, eris.New("startDate time.Parse()")
	}

	if strings.Contains(page, `"isUpcoming":true,`) {
		return "upcoming", &tm, nil
	}
	if strings.Contains(page, `"isLive":true,`) {
		return "live", &tm, nil
	}

	return "end", &tm, nil
}

func GetChannelIDByUrl(url string) (string, error) {
	ok, err := regexp.MatchString("^https://www.youtube.com/(channel/|c/|user/|)", url)
	if err != nil {
		return "", eris.Wrap(err, "regexp.MatchString()")
	}
	if !ok {
		return "", eris.New("Неверный формат")
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
	search := []string{`browseId":"`, `",`}
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
