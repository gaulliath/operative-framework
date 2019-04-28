package diagram

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/graniet/operative-framework/session"
)

type DiagramModule struct{
	session.SessionModule
	sess *session.Session
	Stream *session.Stream
}

func PushDiagramModuleRetrieval(s *session.Session) *DiagramModule{
	mod := DiagramModule{
		sess: s,
		Stream: &s.Stream,
	}
	return &mod
}

func (module *DiagramModule) Name() string{
	return "diagram"
}

func (module *DiagramModule) Author() string{
	return "Tristan Granier"
}

func (module *DiagramModule) Description() string{
	return "Generate a diagram of project"
}

func (module *DiagramModule) GetType() string{
	return ""
}

func (module *DiagramModule) GetInformation() session.ModuleInformation{
	information := session.ModuleInformation{
		Name: module.Name(),
		Description: module.Description(),
		Author: module.Author(),
		Type: module.GetType(),
		Parameters: module.Parameters,
	}
	return information
}

func (module *DiagramModule) Start(){
	graphAst, _ := gographviz.ParseString(`digraph Project {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	_ = graph.AddNode("Project", "a", nil)
	_ = graph.AddNode("Project", "b", nil)
	_ = graph.AddEdge("a", "b", true, nil)
	output := graph.String()
	fmt.Println(output)
}
