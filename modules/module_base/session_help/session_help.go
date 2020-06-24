package session_help

import (
	"fmt"
	"os"

	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
)

type HelpModule struct {
	session.SessionModule
	sess *session.Session `json:"-"`
}

func PushModuleHelp(s *session.Session) *HelpModule {
	moduleHelp := HelpModule{
		sess: s,
	}

	return &moduleHelp
}

func (module *HelpModule) Name() string {
	return "session_help"
}

func (module *HelpModule) Description() string {
	return "Listing available modules"
}

func (module *HelpModule) Author() string {
	return "Tristan Granier"
}

func (module *HelpModule) GetType() []string {
	return []string{
		session.T_TARGET_SESSION,
	}
}

func (module *HelpModule) GetInformation() session.ModuleInformation {
	information := session.ModuleInformation{
		Name:        module.Name(),
		Description: module.Description(),
		Author:      module.Author(),
		Type:        module.GetType(),
		Parameters:  module.Parameters,
	}
	return information
}

func (module *HelpModule) Start() {
	fmt.Println("ENGINE:")
	t := module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"info session", "Print current session information"},
		{"info api", "Print api rest endpoints information"},
		{"env", "Print environment variable"},
		{"events", "View registered session events"},
		{"help", "Print help information"},
		{"clear", "Clear current screen"},
		{"webhooks", "Print loaded webhooks"},
		{"api <run/stop>", "(Run/Stop) restful API"},
	})
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("INTERVAL:")
	fmt.Println("The operative framework engine natively allows you to perform actions repeatedly via a module called interval")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"interval generate <command>", "Add new interval to session"},
		{"interval list", "Listing of interval(s) available in current session"},
		{"interval set <intervalId> <DELAY> <TIME>", "Set interval delay to command e.g: 10 for 10 minutes"},
		{"interval up <intervalId>", "Run interval command in background every <DELAY>"},
		{"interval down <intervalId>", "Stop interval"},
	})
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("MONITOR:")
	fmt.Println("You have the possibility of monitoring certain keywords in order to be informed of their presence in the results")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"monitor generate <search term>", "Add new monitor to session"},
		{"monitor list", "Listing of monitor(s) available in current session"},
		{"monitor up <monitorId>", "Run monitor"},
		{"monitor down <monitorId>", "Stop monitor"},
	})
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("NOTES:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"note add <id target/result> <text>", "Add new note to target or result"},
		{"note view <id target/result>", "View note linked to target or result "},
	})
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("TARGETS:")
	fmt.Println("Targets are the subjects of your investigations.")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"target add <type> <value>", "Add new target"},
		{"target view result <target_id> <result_id>", "View one result from targets"},
		{"target view results <target_id> <module_name>", "View all targets results from module name"},
		{"target links <target_id>", "View linked targets"},
		{"target update <target_id> <value>", "Update a target"},
		{"target delete <target_id>", "Remove target by ID"},
		{"target list", "List subjects"},
		{"target modules <target_id>", "List modules available with selected target"},
	})
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("FILTERS:")
	fmt.Println("Filters allow you to perform actions at the end of module execution")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Filter Name", "Filter Description"})
	for _, mod := range module.sess.Filters {
		t.AppendRow([]interface{}{mod.Name(), mod.Description()})
	}
	module.sess.Stream.Render(t)
	fmt.Println("")
	fmt.Println("MODULES:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"modules", "List available modules"},
		{"modules <target_type>", "List available modules for target type"},
		{"<module_name> target <target_id>", "Set a target argument"},
		{"<module_name> filter <filter>", "Set a filter argument"},
		{"<module_name> set <argument> <value>", "Set specific argument"},
		{"<module_name> list", "List module arguments"},
		{"<module_name> run", "Run selected module"},
	})
	module.sess.Stream.Render(t)

}
