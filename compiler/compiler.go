package compiler

import (
	"bufio"
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"log"
	"os"
	"strings"
	"regexp"
)

type Scripts struct{
	Lines []string
	Vars []ScriptVar
	Comments []string
	Foreach Foreach
	History []string
	Config Config
}

type Foreach struct{
	InForeach bool
	Statements []string
	Context []string
}

type ScriptVar struct{
	Name string
	Foreach bool
	Value []string
}

type Config struct{
	Debug bool
	History bool
}

func (s *Scripts) AddComment(line string){
	s.Comments = append(s.Comments, line)
}

func (s *Scripts) AddHistory(line string){
	s.History = append(s.History, line)
}

func (s *Scripts) AddStatement(line string){
	s.Foreach.Statements = append(s.Foreach.Statements, line)
}

func (s *Scripts) AddContext(lines []string){
	s.Foreach.Context = lines
}

func (s *Scripts) AddVar(name string, foreach bool, value []string){
	for  k, v := range s.Vars{
		if v.Name == name{
			s.Vars[k].Value = value
			return
		}
	}
	s.Vars = append(s.Vars, ScriptVar{
		Name: name,
		Foreach: foreach,
		Value: value,
	})
}

func Settings(sess *session.Session, line string, s *Scripts){
	if string(line[0]) == "@"{
		if line == "@debug off"{
			s.Config.Debug = false
			sess.ParseCommand("session_stream set VERBOSE false")
			sess.ParseCommand("session_stream run")
		} else if line == "@history"{
			s.Config.History = true
		} else if line == "@debug on"{
			s.Config.Debug = true
			sess.ParseCommand("session_stream set VERBOSE true")
			sess.ParseCommand("session_stream run")
			return
		}
	}
}

func Run(sess *session.Session, script string){
	file, err := os.Open(script)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	var newScript Scripts

	fScanner := bufio.NewScanner(file)
	for fScanner.Scan() {
		line := strings.TrimSpace(fScanner.Text())
		if line != "" {
			if !strings.Contains(line, "//") {
				Settings(sess, line, &newScript)
				if newScript.Foreach.InForeach{
					if strings.Contains(line,"<=)"){
						for _, element := range newScript.Foreach.Context{
							for _, stmts := range newScript.Foreach.Statements{
								if string(stmts[0]) == "@"{
									Settings(sess, stmts, &newScript)
								} else if strings.Contains(stmts, "{") && strings.Contains(stmts, "}") && !strings.Contains(stmts,"=>"){
									r, _ := regexp.Compile("\\{F(.*?)\\}")
									v := strings.TrimSpace(r.FindString(stmts))
									newCommand := strings.Replace(stmts, v, element, 1)
									if string(stmts[0]) == "{"{
										variable := strings.TrimSpace(strings.Split(newCommand, "=")[0])
										commands := strings.TrimSpace(strings.Split(newCommand, "=")[1])
										ret := sess.ParseCommand(commands)
										newScript.AddHistory(commands)
										if ret != nil {
											newScript.AddVar(variable, false, ret)
										}
									} else{
										r, _ := regexp.Compile("\\{(.*?)\\}")
										v := strings.TrimSpace(r.FindString(stmts))
										for _, vars := range newScript.Vars{
											if vars.Name == v{
												for _, value := range vars.Value{
													newCommand := strings.Replace(stmts, v, value, -1)
													sess.ParseCommand(newCommand)
													newScript.AddHistory(newCommand)
													break
												}
												break
											}
										}
									}
								} else{
									sess.ParseCommand(stmts)
									newScript.AddHistory(stmts)
								}
							}
						}
						newScript.Foreach.InForeach = false
						continue
					}
					newScript.AddStatement(line)
				} else if strings.Contains(line, "=") && !strings.Contains(line, "=>"){
					variable := strings.TrimSpace(strings.Split(line, "=")[0])
					commands := strings.TrimSpace(strings.Split(line, "=")[1])
					ret := sess.ParseCommand(commands)
					newScript.AddHistory(commands)
					if ret != nil {
						newScript.AddVar(variable, false, ret)
					}
				} else if strings.Contains(line, "{") && strings.Contains(line, "}") && !strings.Contains(line,"=>"){
					r, _ := regexp.Compile("\\{(.*?)\\}")
					v := strings.TrimSpace(r.FindString(line))
					for _, vars := range newScript.Vars{
						if vars.Name == v{
							for _, value := range vars.Value{
								newCommand := strings.Replace(line, v, value, -1)
								sess.ParseCommand(newCommand)
								newScript.AddHistory(newCommand)
							}
						}
					}

				} else if strings.Contains(line, "({") && strings.Contains(line, "=>"){
					newScript.Foreach.InForeach = true
					r, _ := regexp.Compile("\\{(.*?)\\}")
					v := r.FindAllString(line, -1)
					if len(v) < 2{
						log.Println("Failed to compile at : '" + line)
						return
					} else{
						for _, vars := range newScript.Vars{
							if vars.Name == v[0]{
								newScript.AddContext(vars.Value)
							}
						}
					}
				} else{
					sess.ParseCommand(line)
					newScript.AddHistory(line)
				}
			} else {
				newScript.AddComment(line)
			}
		}
	}

	if newScript.Config.History {
		t := sess.Stream.GenerateTable()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"LINE",
			"CONTENT",
		})
		for k, history := range newScript.History {
			t.AppendRow(table.Row{
				k,
				history,
			})
		}
		t.Render()
	}
}
