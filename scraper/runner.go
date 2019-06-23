package scraper

import (
	"errors"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/mgutz/logxi/v1"
	"github.com/osiloke/grapple/scraper"
)

var logger = log.New("jumler.scraper")

func init() {
	scraper.AddCustomType("title", func(prop *scraper.Schema, sel *goquery.Selection) interface{} {
		// sel.Get
		return nil
	})
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	logger.Debug(fmt.Sprintf("%s took %s", name, elapsed))
}

//ScrapeCategory scrapes a jumia category
func ScrapeCategory(url, category string, pages int) (chan map[string]interface{}, error) {
	var schema = scraper.SchemaFromString(categoryItemSchema)
	client, err := scraper.NewDefaultClient(&scraper.DefaultClient{Retry: 3}) //, Socks5Proxy: "127.0.0.1:9050"})
	if err != nil {
		logger.Info(err.Error())
		return nil, errors.New("failed making client")
	}

	job := &scraper.Job{
		Name:      category,
		URL:       url,
		JobSchema: schema,
		Con:       client,
	}
	allResults := make(chan map[string]interface{})
	scrapeResult := func(job *scraper.Job) error {
		results, err := job.ScrapeStream()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		logger.Info("handling results for " + job.URL)
		for r := range results {
			// logger.Info("retrieved result", "result", r)
			allResults <- r
		}
		return nil
	}
	if pages > 0 {
		go func() {
			defer close(allResults)
			curPage := 1
			job.URL = fmt.Sprintf("%s/%s?page=%d", url, category, curPage)
			for pages > 0 {
				if err := scrapeResult(job); err != nil {
					logger.Warn("unable to scrape result", err)
					break
				}
				pages = pages - 1
				curPage++
				job = &scraper.Job{
					Name:      category,
					URL:       fmt.Sprintf("%s/%s?page=%d", url, category, curPage),
					JobSchema: schema,
					Con:       client,
				}
			}
		}()
	} else {
		go func() {
			defer close(allResults)
			curPage := 1
			for {
				job = &scraper.Job{
					Name:      category,
					URL:       fmt.Sprintf("%s/page-%d", url, curPage),
					JobSchema: schema,
					Con:       client,
				}
				if err := scrapeResult(job); err != nil {
					logger.Warn(err.Error())
					break
				}
				curPage++
			}
		}()
	}
	return allResults, nil
}
