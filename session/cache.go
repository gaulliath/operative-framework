package session

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type BackupSession struct {
	Id            int               `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName   string            `json:"session_name"`
	Information   Information       `json:"information"`
	Events        Events            `json:"events"`
	SourceFile    string            `json:"source_file"`
	Version       string            `json:"version" sql:"-"`
	TypeLists     []string          `json:"type_lists" sql:"-"`
	ServiceFolder string            `json:"home_folder"`
	Services      []Listener        `json:"services"`
	Alias         map[string]string `json:"-" sql:"-"`
}

type BackupInstances struct {
	Instances []*Instance `json:"instances"`
}

type BackupTargets struct {
	Targets []*Target `json:"targets" sql:"-"`
}

type BackupMonitors struct {
	Monitors Monitors `json:"monitors"`
}

type BackupInterval struct {
	Interval []*Interval `json:"interval"`
}

type BackupWebhooks struct {
	WebHooks []*WebHook `json:"web_hooks"`
}

type Backup struct {
	Id              int               `json:"-" gorm:"primary_key:yes;column:id;AUTO_INCREMENT"`
	SessionName     string            `json:"session_name"`
	Information     Information       `json:"information"`
	Instances       []*Instance       `json:"instances"`
	CurrentInstance *Instance         `json:"current_instance"`
	Events          Events            `json:"events"`
	SourceFile      string            `json:"source_file"`
	Version         string            `json:"version" sql:"-"`
	Targets         []*Target         `json:"targets" sql:"-"`
	Monitors        Monitors          `json:"monitors"`
	TypeLists       []string          `json:"type_lists" sql:"-"`
	ServiceFolder   string            `json:"home_folder"`
	Services        []Listener        `json:"services"`
	Alias           map[string]string `json:"-" sql:"-"`
	Interval        []*Interval       `json:"interval"`
	WebHooks        []*WebHook        `json:"web_hooks"`
}

// Store current session in cache
func (s *Session) ToCache(name string) {

	baseName := strings.Split(name, " ")[0]

	if _, err := os.Stat(s.Config.Common.BaseDirectory + "cache/"); os.IsNotExist(err) {
		s.Stream.Standard("Generate a root cache folder '" + s.Config.Common.BaseDirectory + "cache/'")
		_ = os.Mkdir(s.Config.Common.BaseDirectory+"cache/", os.ModePerm)
	}

	if _, err := os.Stat(s.Config.Common.BaseDirectory + "cache/" + baseName); os.IsNotExist(err) {
		s.Stream.Standard("Generate a cache folder '" + s.Config.Common.BaseDirectory + "cache/" + baseName + "'")
		_ = os.Mkdir(s.Config.Common.BaseDirectory+"cache/"+baseName, os.ModePerm)
	}

	err := s.CreateCacheSession(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	err = s.CreateCacheInstances(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	err = s.CreateCacheTargets(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	err = s.CreateCacheMonitors(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	err = s.CreateCacheIntervals(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	err = s.CreateCacheWebhooks(baseName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	//s.Stream.Success("Session saved at '" + cacheSessionFile + "'")
	return
}

func (s *Session) CreateCacheSession(base string) error {
	var sessionCache BackupSession
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &sessionCache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(sessionCache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/session.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

func (s *Session) CreateCacheInstances(base string) error {
	var instanceCache BackupInstances
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &instanceCache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(instanceCache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/instances.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

func (s *Session) CreateCacheTargets(base string) error {
	var targetsCache BackupTargets
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &targetsCache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(targetsCache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/targets.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

func (s *Session) CreateCacheMonitors(base string) error {
	var cache BackupMonitors
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/monitors.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

func (s *Session) CreateCacheIntervals(base string) error {
	var cache BackupInterval
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/intervals.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

func (s *Session) CreateCacheWebhooks(base string) error {
	var cache BackupWebhooks
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return err
	}

	bytes, err = json.Marshal(cache)
	if err != nil {
		return err
	}

	cacheSessionFile := s.Config.Common.BaseDirectory + "cache/" + base + "/hooks.json"
	file, err := os.OpenFile(cacheSessionFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(string(bytes))
	_ = file.Close()

	return nil
}

// Load session file cache to current session
func (s *Session) LoadCache(name string) {

	var backup Backup
	files := []string{
		"session.json",
		"targets.json",
		"monitors.json",
		"intervals.json",
		"instances.json",
		"hooks.json",
	}

	for _, f := range files {

		path := name + "/" + f
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			s.Stream.Error(err.Error())
			return
		}

		input, err := ioutil.ReadAll(file)
		if err != nil {
			s.Stream.Error(err.Error())
			return
		}

		err = json.Unmarshal(input, &backup)
		if err != nil {
			s.Stream.Error(err.Error())
			return
		}
	}

	// Preload session information
	s.SessionName = backup.SessionName
	s.Information = backup.Information
	s.Events = backup.Events
	s.SourceFile = backup.SourceFile
	s.Alias = backup.Alias

	// Load instances
	for _, instance := range backup.Instances {
		s.Instances = append(s.Instances, instance)
	}

	// Load targets
	for _, target := range backup.Targets {
		target.Sess = s
		s.Targets = append(s.Targets, target)
	}

	// Load monitor
	for _, monitor := range backup.Monitors {
		monitor.Session = s
		s.Monitors = append(s.Monitors, monitor)
	}

	// Load interval
	for _, interval := range backup.Interval {
		interval.S = s
		s.Interval = append(s.Interval, interval)
	}

	// Load hooks
	for _, hook := range backup.WebHooks {
		s.WebHooks = append(s.WebHooks, hook)
	}

	s.Stream.Success("Loaded session: '" + name + "'")
	return
}
