package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/francoispqt/gojay"
	"github.com/gammazero/workerpool"
	"github.com/osiloke/dostow-contrib/api"
)

// Our structure which will be pushed to our stream
type Item struct {
	Brand string
	Image string
	Name  string
	Sku   string
	Link  string
}

func (u *Item) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("image", u.Image)
	enc.StringKey("brand", u.Brand)
	enc.StringKey("name", u.Name)
	enc.StringKey("sku", u.Sku)
	enc.StringKey("link", u.Link)
}
func (u *Item) IsNil() bool {
	return u == nil
}

// Our MarshalerStream implementation
type StreamChan chan *Item

func (s StreamChan) MarshalStream(enc *gojay.StreamEncoder) {
	select {
	case <-enc.Done():
		return
	case o := <-s:
		fmt.Println(o)
		enc.Object(o)
	}
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// StreamItemResultsToIO takes a scraper response and sends it to a jsonlines `io.writer`
func StreamItemResultsToIO(ws io.Writer, res chan map[string]interface{}) {
	enc := gojay.Stream.BorrowEncoder(ws).NConsumer(10).LineDelimited()
	defer enc.Release()
	// s := StreamChan(make(chan *Item))
	// go enc.EncodeStream(s)
	for r := range res {
		if r["sku"] == nil {
			continue
		}
		sku := toString(r["sku"])
		brand := toString(r["brand"])
		image := toString(r["image"])
		name := toString(r["name"])
		link := toString(r["link"])
		item := &Item{
			Brand: brand,
			Image: image,
			Name:  name,
			Sku:   sku,
			Link:  link,
		}
		enc.Encode(item)

		// ws.Write(enc.)
	}
	// Wait
	// <-enc.Done()
}

type SearchQuery struct {
	Name string `json:"name,omitempty"`
	Sku  string `json:"sku"`
}

func getImage(url string) (io.ReadCloser, error) {
	response, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return response.Body, nil
}

type Result struct {
	Data []map[string]interface{}
}

// DostowWriter stores results in a store
func DostowWriter(apiURL, apiKey, storeName, category string, res chan map[string]interface{}) {
	dostow := api.NewClient(apiURL, apiKey)
	store := dostow.Store
	upload := dostow.File
	wp := workerpool.New(6)
	for r := range res {
		// if processed > 0 {
		// 	return
		// }

		if r["sku"] == nil {
			continue
		}
		r := r
		wp.Submit(func() {
			sku := toString(r["sku"])
			r["category"] = category
			image := toString(r["image"])
			fmt.Println(r)
			delete(r, "image")
			delete(r, "link")
			fmt.Printf("saving %+v", r)
			raw, err := store.Search(storeName, store.Query(SearchQuery{Sku: sku}))
			var result Result
			if err == nil {
				if err := json.Unmarshal(*raw, &result); err == nil {
					if result.Data[0]["photo"] == nil {
						// upload image

						if file, err := getImage(image); err == nil {
							defer file.Close()
							if _, _, err := upload.Create(storeName, result.Data[0]["id"].(string), "photo", "im.jpg", file); err != nil {
								fmt.Println(err.Error())
								panic(err)
							}
						}
					}
				}
				fmt.Printf("error %v", err)
				return
			}
			raw, err = store.Create(storeName, r)
			if err != nil {
				fmt.Println(err.Error())
				if e, ok := err.(api.APIError); ok && e.Status != "400" {
					return
				}
			} else {
				var m map[string]interface{}
				if err := json.Unmarshal(*raw, &m); err == nil {
					// upload image
					if file, err := getImage(image); err == nil {
						defer file.Close()
						if _, _, err := upload.Create(storeName, m["id"].(string), "photo", "im.jpg", file); err != nil {
							fmt.Println(err.Error())
						}
					}
				}
			}
		})

	}
	wp.StopWait()
	// right first item with header
	// first := <- res
	// err = gocsv.MarshalFile(records, analyticsFile)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// for r := range res {
	// 	fmt.Println(r)
	// 	// r["forumName"] = hash
	// 	// r["forumLink"] = url
	// 	// fmt.Printf("saving %+v", r)
	// 	// _, err := store.Create("threads", r)
	// 	// if err != nil {
	// 	// 	fmt.Println(err.Error())
	// 	// 	if e, ok := err.(api.APIError); ok && e.Status != "400" {
	// 	// 		break
	// 	// 	}
	// 	// }

	// 	// err = gocsv.MarshalCSVWithoutHeaders(records, gocsv.DefaultCSVWriter(categoryFile))
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// }
}
