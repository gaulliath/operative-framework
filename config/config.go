package config

type Config struct {
	Api        ApiConfig
	Tracker    TrackerConfig
	Database   DataBase
	Common     Common
	Instagram  Network
	Twitter    Network
	PushDriver string
	Gate       ToGate
	Modules    map[string]map[string]string
}

type ToGate struct {
	GateUrl    string
	GateMethod string
	GateTor    string
}

type Network struct {
	Login    string
	Password string
	Api      ApiConfig
}

type DataBase struct {
	Name   string
	User   string
	Pass   string
	Host   string
	Driver string
	Port   string
}

type Common struct {
	HistoryFile       string
	BaseDirectory     string
	ExportDirectory   string
	ConfigurationFile string
	ConfigurationJobs string
}

type ApiConfig struct {
	Host    string
	Port    string
	Key     string
	SKey    string
	Verbose string
}

type TrackerConfig struct {
	Host string
	Port string
}
