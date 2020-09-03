package main

import "context"

func main() {
	app := New()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())
}

/*
TODO:
tg команды - время проверки (через сколько проверять rss), ch id
перенести loc и текста в конфиг
перебрать ошибки, + errors.Wrap, перебрать имена пм
*/

// тест на лимиты запросов.
/*
import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gilliek/go-opml/opml"
	"github.com/mmcdole/gofeed"
	"sync"
	"time"
)

func main() {
	ids := []string{}
	doc, _ := opml.NewOPMLFromFile("subscription_manager")

	for {
		chid := make(chan string, 200)
		idsCache := []string{}
		var errs, itemsN int = 0, 0
		var errsstr []string

		bar := pb.StartNew(len(doc.Body.Outlines[0].Outlines))
		s := sync.WaitGroup{}
		s1 := sync.WaitGroup{}
		s1.Add(1)
		tm := time.Now()

		//переделать на sync lock
		go func() {
			for {
				id, ass := <-chid
				if ass {
					idsCache = append(idsCache, id)
				} else {
					s1.Done()
					return
				}
			}
		}()

		s.Add(len(doc.Body.Outlines[0].Outlines))
		for _, v := range doc.Body.Outlines[0].Outlines {
			go func(v opml.Outline) {
				feed, err := gofeed.NewParser().ParseURL(v.XMLURL)
				if err != nil {
					errsstr = append(errsstr, err.Error()+" url: "+v.XMLURL)
					errs++
				} else {
					itemsN += len(feed.Items)
					for _, itm := range feed.Items {
						if !stringInSlice(itm.Link, &ids) {
							chid <- itm.Link
						}
					}
				}
				bar.Increment()
				s.Done()
			}(v)
		}

		s.Wait()
		close(chid)
		s1.Wait()
		tm2 := time.Since(tm)
		bar.Finish()

		if len(ids) != 0 && len(idsCache) != 0 {
			fmt.Printf("%v\n", idsCache)
		}
		ids = append(ids, idsCache...)

		fmt.Printf("%v idsCache-%v items-%v errs-%v errstr-%v\n", tm2, len(idsCache), itemsN, errs, errsstr)

		time.Sleep(time.Minute * 2)
	}
}

func stringInSlice(a string, list *[]string) bool {
	for _, b := range *list {
		if b == a {
			return true
		}
	}
	return false
}

*/
