package image_reverse_search

import (
	"bytes"
	"errors"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type ImageReverseModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

const (
	baseurl = "https://www.google.com"
)

// requestParams : Parameters for fetchURL
type requestParams struct {
	Method      string
	URL         string
	Contenttype string
	Data        io.Reader
	Client      *http.Client
}

// Imgdata : Image URL
type Imgdata struct {
	OU      string `json:"ou"`
	WebPage bool
}

// DefImg : Initialize imagdata.
func DefImg(webpages bool) *Imgdata {
	return &Imgdata{
		WebPage: webpages,
	}
}

// fetchURL : Fetch method
func (r *requestParams) fetchURL() *http.Response {
	req, err := http.NewRequest(
		r.Method,
		r.URL,
		r.Data,
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	if len(r.Contenttype) > 0 {
		req.Header.Set("Content-Type", r.Contenttype)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/26.0")
	res, _ := r.Client.Do(req)
	return res
}

// getWebPages : Retrieve web pages with matching images on Google top page. When this is not used, images are retrieved.
func getWebPages(doc *goquery.Document) []string {
	var ar []string
	doc.Find("h3.r").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			ar = append(ar, url)
		})
	})
	return ar
}

// ImgFromFile : Search images from an image file
func (im *Imgdata) ImgFromFile(file string) []string {
	var url string
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fs, err := os.Open(file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		return []string{}
	}
	defer fs.Close()
	data, err := w.CreateFormFile("encoded_image", file)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		return []string{}
	}
	if _, err = io.Copy(data, fs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		return []string{}
	}
	_ = w.Close()
	r := &requestParams{
		Method: "POST",
		URL:    baseurl + "/searchbyimage/upload",
		Data:   &b,
		Client: &http.Client{
			Timeout:       time.Duration(10) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("Redirect") },
		},
		Contenttype: w.FormDataContentType(),
	}
	var res *http.Response
	for {
		res = r.fetchURL()
		if res.StatusCode == 200 {
			break
		}
		reurl, _ := res.Location()
		r.URL = reurl.String()
		r.Method = "GET"
		r.Data = nil
		r.Contenttype = ""
	}
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	var ar []string
	if im.WebPage {
		ar = getWebPages(doc)
	} else {
		doc.Find(".iu-card-header").Each(func(_ int, s *goquery.Selection) {
			url, _ = s.Attr("href")
		})
		r.URL = baseurl + url
		r.Client = &http.Client{
			Timeout: time.Duration(10) * time.Second,
		}
		res = r.fetchURL()
		doc, _ = goquery.NewDocumentFromReader(res.Body)
		doc.Find(".rg_meta").Each(func(_ int, s *goquery.Selection) {
			_ = json.Unmarshal([]byte(s.Text()), &im)
			ar = append(ar, im.OU)
		})
	}
	return ar
}

func PushImageReverseModule(s *session.Session) *ImageReverseModule{
	mod := ImageReverseModule{
		sess: s,
		Stream: &s.Stream,
	}

	mod.CreateNewParam("TARGET", "Image path", "", true, session.STRING)
	mod.CreateNewParam("limit", "Limit search", "10", false, session.STRING)
	return &mod
}

func (module *ImageReverseModule) Name() string{
	return "image_search"
}

func (module *ImageReverseModule) Description() string{
	return "Search possible image(s) connection on the internet"
}

func (module *ImageReverseModule) Author() string{
	return "Tristan Granier"
}

func (module *ImageReverseModule) GetType() string{
	return "file"
}

func (module *ImageReverseModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *ImageReverseModule) Start(){
	trg, err := module.GetParameter("TARGET")
	if err != nil{
		module.sess.Stream.Error(err.Error())
		return
	}
	target, err2 := module.sess.GetTarget(trg.Value)
	if err2 != nil{
		module.sess.Stream.Error(err2.Error())
		return
	}
	results := DefImg(false).ImgFromFile(target.GetName())
	if len(results) > 0{
		t := module.sess.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"URL",
		})
		t.SetAllowedColumnLengths([]int{90,})
		for _, url := range results{
			t.AppendRow(table.Row{
				url,
			})
			res := session.TargetResults{
				Header: "URL" + target.GetSeparator(),
				Value: url + target.GetSeparator(),
			}
			module.Results = append(module.Results, url)
			target.Save(module, res)
		}
		module.sess.Stream.Render(t)
	} else{
		module.Stream.Warning("No result found.")
	}
}
