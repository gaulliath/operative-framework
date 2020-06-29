package session

type Information struct {
	ApiStatus      bool `json:"api_status"`
	TrackerStatus  bool `json:"tracker_status"`
	ModuleLaunched int  `json:"module_launched"`
	Event          int  `json:"event"`
}

func (i *Information) AddEvent() {
	i.Event = i.Event + 1
	return
}

func (i *Information) AddModule() {
	i.ModuleLaunched = i.ModuleLaunched + 1
	return
}

func (i *Information) SetApi(s bool) {
	i.ApiStatus = s
	return
}

func (i *Information) SetTracker(s bool) {
	i.TrackerStatus = s
	return
}
