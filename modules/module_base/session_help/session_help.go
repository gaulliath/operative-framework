package session_help

import (
	"github.com/graniet/operative-framework/session"
	"github.com/graniet/go-pretty/table"
	"os"
	"fmt"
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
	fmt.Println("Targets:")
	t := module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"target add <type> <value>", "Add new target"},
		{"target links <target_id>", "View linked targets"},
		{"target update <target_id> <value>", "Update a target"},
		{"target delete <target_id>", "Remove target by ID"},
		{"target list", "List subjects"},
		{"target modules <target_id>", "List modules available with selected target"},
	})
	module.sess.Stream.Render(t)

	fmt.Println("")
	fmt.Println("Modules:")
	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"command", "description"})
	t.AppendRows([]table.Row{
		{"<module> target <target_id>", "Set a target argument"},
		{"<module> set <argument> <value>", "Set specific argument"},
		{"<module> list", "List module arguments"},
		{"<module> run", "Run selected module"},
	})
	module.sess.Stream.Render(t)

	t = module.sess.Stream.GenerateTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Module Name", "Module Description"})
	for _, mod := range module.sess.Modules{
		t.AppendRow([]interface{}{mod.Name(), mod.Description()})
	}
	module.sess.Stream.Render(t)
}
