package session_help

import (
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"os"
)

type HelpModule struct{
	session.SessionModule
	sess *session.Session
}

func PushModuleHelp(s *session.Session) *HelpModule{
	moduleHelp := HelpModule{
		sess: s,
	}

	return &moduleHelp
}

func (module *HelpModule) Name() string{
	return "session_help"
}

func (module *HelpModule) Description() string{
	return "Listing available modules"
}

func (module *HelpModule) Author() string{
	return "Tristan Granier"
}

func (module *HelpModule) GetType() string{
	return "session"
}


func (module *HelpModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *HelpModule) Start(){
	fmt.Println("ENGINE:")
	t := module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"info session", "Print current session information"},
		{"info api", "Print api rest endpoints information"},
		{"env", "Print environment variable"},
		{"help", "Print help information"},
		{"clear", "Clear current screen"},
		{"api <run/stop>", "(Run/Stop) restful API"},
	})
	module.sess.Stream.Render(t)
	fmt.Println("NOTES:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"note add <id target/result> <text>", "Add new note to target or result"},
		{"note view <id target/result>", "View note linked to target or result "},
	})
	module.sess.Stream.Render(t)
	fmt.Println("TARGETS:")
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
	fmt.Println("FILTERS:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Filter Name", "Filter Description"})
	for _, mod := range module.sess.Filters{
		t.AppendRow([]interface{}{mod.Name(), mod.Description()})
	}
	module.sess.Stream.Render(t)
	fmt.Println("MODULES:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"<module> target <target_id>", "Set a target argument"},
		{"<module> filter <filter>", "Set a filter argument"},
		{"<module> set <argument> <value>", "Set specific argument"},
		{"<module> list", "List module arguments"},
		{"<module> run", "Run selected module"},
	})
	module.sess.Stream.Render(t)

	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Module Name", "Module Description", "Target"})
	for _, mod := range module.sess.Modules{
		if mod.GetType() == ""{
			t.AppendRow([]interface{}{mod.Name(), mod.Description(), "<blank>"})
		} else{
			t.AppendRow([]interface{}{mod.Name(), mod.Description(), mod.GetType()})
		}
	}
	module.sess.Stream.Render(t)

}
