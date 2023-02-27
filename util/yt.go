package util

import (
	"github.com/rotisserie/eris"
	"io"
	"net/http"
	"os"
	"regexp"
	"streamtg/go-log"
	"strings"
	"time"
)

const (
	Err404   = "404"
	Video    = "video"
	Upcoming = "upcoming"
	Wait     = "wait"
	Live     = "live"
	LiveGo   = "live_go"
	End      = "end"
	End404   = "end404"
)

// 404 video upcoming wait live end
func TypeVideo(videoID string, debugSave bool) (string, *time.Time, error) {
	resp, err := http.Get("https://www.youtube.com/watch?v=" + videoID)
	if err != nil {
		return "", nil, eris.Wrap(err, "http.Get(url)")
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, eris.Wrap(err, "io.ReadAll()")
	}
	if debugSave {
		err = os.WriteFile("files/"+time.Now().Format(time.RFC3339)+".txt", bs, os.ModePerm)
		if err != nil {
			log.Error(eris.Wrap(err, "os.WriteFile"))
		}
	}
	page := string(bs)

	if strings.Contains(page, `{"iconType":"ERROR_OUTLINE"}`) {
		return Err404, nil, nil
	}
	if !strings.Contains(page, `<span itemprop="publication" itemscope itemtype="http://schema.org/BroadcastEvent">`) {
		return Video, nil, nil
	}

	getTime := func() (time.Time, error) {
		tmStrings := []string{`<meta itemprop="startDate" content="`, `">`}
		iTime1 := strings.Index(page, tmStrings[0])
		if iTime1 < 0 {
			return time.Unix(322, 0), nil
		}
		iTime1 += len(tmStrings[0])
		iTime2 := strings.Index(page[iTime1:], tmStrings[1])
		if iTime2 < 0 {
			return time.Unix(322, 0), nil
		}
		iTime2 += iTime1

		return time.Parse(time.RFC3339, page[iTime1:iTime2])
	}

	tm, err := getTime()
	if err != nil && tm.Unix() != 322 {
		return "", nil, eris.New("startDate time.Parse()")
	}

	if strings.Contains(page, `"isUpcoming":true,`) {
		if tm.Unix() == 322 {
			tm = time.Now().In(time.UTC)
			return Wait, &tm, nil
		}
		return Upcoming, &tm, nil
	}
	if strings.Contains(page, `"isLive":true,`) {
		return Live, &tm, nil
	}

	return End, &tm, nil
}

func GetChannelIDByUrl(url string) (string, error) {
	ok, err := regexp.MatchString(`^https:\/\/www\.youtube\.com\/(channel\/|c\/|user\/|)`, url)
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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", eris.Wrap(err, "io.ReadAll()")
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
