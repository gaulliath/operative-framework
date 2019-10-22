package session

import (
	"fmt"
	"strconv"
	"strings"
)

type Analytics []Entity
type AnalyticsLinks []Link

type Reader struct {
	Targets []*Target
	Results []*TargetResults
}

type Entity struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
	Code string `json:"code"`
}

type Link struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (a *AnalyticsLinks) HasLinked(t1 string, t2 string) bool {
	for _, link := range *a {
		if link.To == t1 && link.From == t2 {
			return true
		} else if link.To == t2 && link.From == t1 {
			return true
		}
	}
	return false
}

func (r *Reader) HasTarget(name string) bool {
	fmt.Println("name", name)
	for _, target := range r.Targets {
		fmt.Println("check:", target.Name)
		if strings.ToLower(strings.TrimSpace(target.Name)) == strings.ToLower(strings.TrimSpace(name)) {
			fmt.Println("check:", target.Name, "ok")
			return true
		}
	}
	return false
}

func (s *Session) GenerateAnalytics() (Analytics, AnalyticsLinks) {

	var analytics Analytics
	analyticsLinks := AnalyticsLinks{}
	reader := Reader{}

	for _, target := range s.Targets {
		reader.Targets = append(reader.Targets, target)
	}

	for _, target := range reader.Targets {
		entity := Entity{
			Key:  target.GetId(),
			Name: target.GetName(),
			Type: target.GetType(),
			Code: target.GetType(),
		}

		analytics = append(analytics, entity)

		if len(target.TargetLinked) > 0 {
			for _, link := range target.TargetLinked {
				if !analyticsLinks.HasLinked(link.TargetId, link.TargetBase) {
					analyticsLinks = append(analyticsLinks, Link{
						From: link.TargetBase,
						To:   link.TargetId,
					})
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
								if !reader.HasTarget(name) {
									entity := Entity{
										Key:  result.ResultId,
										Name: name,
										Type: configuration["RESULT_TYPE"],
										Code: configuration["RESULT_TYPE"],
									}
									analytics = append(analytics, entity)

									analyticsLinks = append(analyticsLinks, Link{
										From: target.GetId(),
										To:   result.ResultId,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return analytics, analyticsLinks
}
