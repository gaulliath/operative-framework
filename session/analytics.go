package session

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Entities []Entity
type EntitiesLinks []Link

const (
	OPERATIVE_URL = "http://localhost/"

	S_INSTAGRAM_ACCOUNT = "INSTAGRAM_ACCOUNT"
	S_TWITTER_ACCOUNT   = "TWITTER_ACCOUNT"
	S_TWITTER_HASHTAG   = "TWITTER_HASHTAG"
)

type Analytics struct {
	Link        string
	isPublic    bool
	SessionType string
}

type LoadAnalytics struct {
	Data struct {
		Analytics string `json:"analytics"`
		URL       string `json:"url"`
		Object    struct {
			IPAddress    string        `json:"ip_address"`
			RandomID     string        `json:"random_id"`
			Data         []interface{} `json:"data"`
			DataFriends  []interface{} `json:"data_friends"`
			Links        []interface{} `json:"links"`
			LinksFriends []interface{} `json:"links_friends"`
			UpdatedAt    string        `json:"updated_at"`
			CreatedAt    string        `json:"created_at"`
		} `json:"object"`
	} `json:"data"`
	Success bool `json:"success"`
	Error   bool `json:"error"`
}

type Reader struct {
	Targets       []*Target
	Results       []*TargetResults
	LinkedResults []string
	WhiteListed   map[string]string
}

type Entity struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Code      string `json:"code"`
	CreatedAt string `json:"created_at"`
	Size      string `json:"size"`
	Color     string `json:"color"`
}

type Link struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Color string `json:"color"`
	Size  string `json:"size"`
}

func (a *EntitiesLinks) HasLinked(t1 string, t2 string) bool {
	for _, link := range *a {
		if link.To == t1 && link.From == t2 {
			return true
		} else if link.To == t2 && link.From == t1 {
			return true
		}
	}
	return false
}

func (a *EntitiesLinks) HasFound(from string, to string) bool {
	for _, link := range *a {
		if link.From == from && link.To == to {
			return true
		}
	}
	return false
}

func (r *Reader) HasWhiteList(value string) bool {

	if _, ok := r.WhiteListed[value]; ok {
		return true
	}
	return false
}

func (r *Reader) HasTarget(name string) bool {
	for _, target := range r.Targets {
		if strings.ToLower(strings.TrimSpace(target.Name)) == strings.ToLower(strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}

func (r *Reader) HasLinkedResults(id string) bool {
	for _, result := range r.LinkedResults {
		if result == id {
			return true
		}
	}
	return false
}

func (s *Session) AnalyticsUp() {

	if s.Analytics.Link != "" {
		s.Stream.Success("Analytics as been opened at '" + s.Analytics.Link + "'")
		return
	}

	client := http.Client{}
	request, err := http.NewRequest("POST", OPERATIVE_URL, nil)

	if err != nil {
		s.NewEvent(ERROR_ANALYTICS, err.Error())
		return
	}

	res, err := client.Do(request)
	if err != nil {
		s.NewEvent(ERROR_ANALYTICS, err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		s.NewEvent(ERROR_ANALYTICS, err.Error())
		return
	}

	responseAnalytics := LoadAnalytics{}
	_ = json.Unmarshal(body, &responseAnalytics)

	if responseAnalytics.Success {
		s.Stream.Success("Analytics as been opened at '" + responseAnalytics.Data.URL + "'")
		s.Analytics.Link = responseAnalytics.Data.URL
		s.Analytics.isPublic = false
		s.Analytics.SessionType = S_INSTAGRAM_ACCOUNT
		return
	}

	s.Stream.Error("A error as been occurred.")
	return
}

func (s *Session) PutAnalytics() (Entities, EntitiesLinks, Entities, EntitiesLinks, Entities, EntitiesLinks) {

	var analytics Entities
	var analyticsFriends Entities
	var analyticsSpot Entities

	analyticsLinks := EntitiesLinks{}
	analyticsLinksFriends := EntitiesLinks{}
	analyticsLinksSpots := EntitiesLinks{}

	reader := Reader{
		WhiteListed: make(map[string]string),
	}

	for _, target := range s.Targets {
		reader.Targets = append(reader.Targets, target)
	}

	for _, target := range reader.Targets {

		var size int
		for _, results := range target.Results {
			size = size + len(results)
		}

		entity := Entity{
			Key:       target.GetId(),
			Name:      target.GetName(),
			Type:      target.GetType(),
			Code:      target.GetType(),
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			Size:      strconv.Itoa(size) + " " + strconv.Itoa(size),
			Color:     "#F2711C",
		}

		analytics = append(analytics, entity)
		analyticsFriends = append(analyticsFriends, entity)
		analyticsSpot = append(analyticsSpot, entity)

		if len(target.TargetLinked) > 0 {
			for _, link := range target.TargetLinked {
				if !analyticsLinks.HasLinked(link.TargetId, link.TargetBase) {

					if !analyticsLinks.HasFound(link.TargetBase, link.TargetResultId) {
						links := Link{
							From: link.TargetBase,
							To:   link.TargetResultId,
						}
						analyticsLinks = append(analyticsLinks, links)
						analyticsLinksFriends = append(analyticsLinksFriends, links)
					}

					if !analyticsLinks.HasFound(link.TargetId, link.TargetResultId) {
						links := Link{
							From:  link.TargetId,
							To:    link.TargetResultId,
							Color: "#F2711C",
							Size:  "3",
						}
						analyticsLinks = append(analyticsLinks, links)
						analyticsLinksFriends = append(analyticsLinksFriends, links)
					}
					reader.LinkedResults = append(reader.LinkedResults, link.TargetResultId)
				}
			}
		}

		if len(target.Results) > 0 {
			for moduleName, results := range target.Results {
				module, err := s.SearchModule(moduleName)
				if err == nil {
					configuration, err := s.LoadModuleConfiguration(module.Name())
					if err != nil {
						continue
					}
					if len(results) > 0 {
						for _, result := range results {

							var name string
							var primaryResult, err = strconv.Atoi(configuration["PRIMARY_RESULT"])
							if err != nil {
								continue
							}

							if strings.Contains(result.Value, target.GetSeparator()) {
								name = strings.Split(result.Value, target.GetSeparator())[primaryResult]
							}

							if name != "" {
								if !reader.HasLinkedResults(result.ResultId) {
									continue
								}
								if !reader.HasTarget(name) {

									if !reader.HasWhiteList(name) {
										entity := Entity{
											Key:       result.ResultId,
											Name:      name,
											Type:      configuration["RESULT_TYPE"],
											Code:      configuration["RESULT_TYPE"],
											CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
											Color:     "lightgray",
										}
										analytics = append(analytics, entity)
										analyticsFriends = append(analyticsFriends, entity)

										reader.WhiteListed[name] = result.ResultId
									} else {
										result.ResultId = reader.WhiteListed[name]
									}

									if !analyticsLinks.HasFound(target.GetId(), result.ResultId) {
										links := Link{
											From:  target.GetId(),
											To:    result.ResultId,
											Color: "#F2711C",
											Size:  "3",
										}
										analyticsLinks = append(analyticsLinks, links)
										analyticsLinksFriends = append(analyticsLinksFriends, links)
									}

									//reader.LinkedResults = append(reader.LinkedResults, result.ResultId)
								}
							}
						}

						for _, result := range results {

							var name string
							var primaryResult, err = strconv.Atoi(configuration["PRIMARY_RESULT"])
							if err != nil {
								continue
							}

							if strings.Contains(result.Value, target.GetSeparator()) {
								name = strings.Split(result.Value, target.GetSeparator())[primaryResult]
							}

							if name != "" {

								if reader.HasLinkedResults(result.ResultId) {
									continue
								}

								if !reader.HasTarget(name) {

									if !reader.HasWhiteList(name) {
										entity := Entity{
											Key:       result.ResultId,
											Name:      name,
											Type:      configuration["RESULT_TYPE"],
											Code:      configuration["RESULT_TYPE"] + "_NL",
											CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
										}
										analytics = append(analytics, entity)
										//analyticsSpot = append(analyticsSpot, entity)

										if !analyticsLinks.HasFound(target.GetId(), result.ResultId) {
											links := Link{
												From:  target.GetId(),
												To:    result.ResultId,
												Color: "black",
												Size:  "1",
											}
											analyticsLinks = append(analyticsLinks, links)
											analyticsLinksSpots = append(analyticsLinksSpots, links)
										}
										reader.LinkedResults = append(reader.LinkedResults, result.ResultId)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return analytics, analyticsLinks, analyticsFriends, analyticsLinksFriends, analyticsSpot, analyticsLinksSpots
}

func (s *Session) WaitAnalytics() {
	for {
		time.Sleep(3 * time.Second)
		if s.Analytics.Link != "" {
			analytics, analyticsLinks, analyticsFriends, analyticsLinksFriends, analyticsSpots, analyticsLinksSpots := s.PutAnalytics()

			toJsonAnalytics, err := json.Marshal(analytics)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			toJsonAnalyticsLinks, err := json.Marshal(analyticsLinks)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			toJsonAnalyticsFriends, err := json.Marshal(analyticsFriends)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			toJsonAnalyticsLinksFriends, err := json.Marshal(analyticsLinksFriends)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			toJsonAnalyticsSpots, err := json.Marshal(analyticsSpots)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			toJsonAnalyticsLinksSpots, err := json.Marshal(analyticsLinksSpots)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			if s.LastAnalyticsModel == string(toJsonAnalytics) &&
				s.LastAnalyticsLinks == string(toJsonAnalyticsLinks) {
				continue
			}

			client := http.Client{}
			data := url.Values{}
			data.Add("data", string(toJsonAnalytics))
			data.Add("links", string(toJsonAnalyticsLinks))
			data.Add("data_friends", string(toJsonAnalyticsFriends))
			data.Add("links_friends", string(toJsonAnalyticsLinksFriends))
			data.Add("data_spots", string(toJsonAnalyticsSpots))
			data.Add("links_spots", string(toJsonAnalyticsLinksSpots))
			data.Add("session_type", s.Analytics.SessionType)

			req, err := http.NewRequest(http.MethodPut, s.Analytics.Link, strings.NewReader(data.Encode()))
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
			res, err := client.Do(req)
			if err != nil {
				s.NewEvent(ERROR_ANALYTICS, err.Error())
				continue
			}

			_, _ = ioutil.ReadAll(res.Body)

			s.LastAnalyticsModel = string(toJsonAnalytics)
			s.LastAnalyticsLinks = string(toJsonAnalyticsLinks)
		}
	}
}
