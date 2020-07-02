package session_import

import (
	"github.com/graniet/operative-framework/session"
	"strconv"
	"strings"
)

type ImportModule struct {
	session.SessionModule
	sess *session.Session `json:"-"`
}

func PushModuleImport(s *session.Session) *ImportModule {
	mod := ImportModule{
		sess: s,
	}
	mod.CreateNewParam("TARGET", "Source file path (csv)", "", true, session.STRING)
	mod.CreateNewParam("DELIMITER", "CSV file separator", ";", true, session.STRING)
	mod.CreateNewParam("PRIMARY", "Set a key of primary element", "0", true, session.INT)
	mod.CreateNewParam("LINKING", "Key of linked value. e.g: 0,5,3", "0", true, session.STRING)
	mod.CreateNewParam("VERBOSE", "Display imported line output", "false", true, session.BOOL)
	return &mod
}

func (module *ImportModule) Name() string {
	return "import.csv"
}

func (module *ImportModule) Description() string {
	return "Import data from CSV"
}

func (module *ImportModule) Author() string {
	return "Tristan Granier"
}

func (module *ImportModule) GetType() []string {
	return []string{
		session.T_TARGET_IMPORT,
	}
}

func (module *ImportModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *ImportModule) Start() {
	targetId, err := module.GetParameter("TARGET")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	target, err := module.sess.GetTarget(targetId.Value)
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	delimiter, err := module.GetParameter("DELIMITER")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	primary, err := module.GetParameter("PRIMARY")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	verbose, err := module.GetParameter("VERBOSE")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	main, err := strconv.Atoi(primary.Value)
	if err != nil {
		module.sess.Stream.Error("Primary parameter value is incorrect please use 'INTEGER'")
		return
	}

	linking, err := module.GetParameter("LINKING")
	if err != nil {
		module.sess.Stream.Error(err.Error())
		return
	}

	links := []int{}
	values := strings.Split(linking.Value, ",")
	for _, val := range values {
		nb, err := strconv.Atoi(val)
		if err != nil {
			continue
		}

		links = append(links, nb)
	}

	module.sess.Stream.Backgound("Import as started")
	go module.sess.ImportFromCsv(target.GetName(), delimiter.Value, main, module.sess.StringToBoolean(verbose.Value), links)

}
