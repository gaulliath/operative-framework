package compiler

import (
	"bufio"
	"fmt"
	"github.com/graniet/go-pretty/table"
	"github.com/graniet/operative-framework/session"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
	"regexp"
	"time"
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

func (s *Scripts) GetVariable(name string) ([]string, error){
	for _, variable := range s.Vars{
		if variable.Name == name{
			return variable.Value, nil
		}
	}
	return []string{}, errors.New("Variable not found.")
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

func (s *Scripts) ResetForeach(){
	s.Foreach.Statements = []string{}
	s.Foreach.Context = []string{}
	s.Foreach.InForeach = false
}

func Settings(sess *session.Session, line string, s *Scripts) bool{
	if string(line[0]) == "@"{
		if line == "@verbose off"{
			sess.ParseCommand("session_stream set VERBOSE false")
			sess.ParseCommand("session_stream run")
			return true
		} else if line == "@history"{
			s.Config.History = true
			return true
		} else if line == "@verbose on"{
			sess.ParseCommand("session_stream set VERBOSE true")
			sess.ParseCommand("session_stream run")
			return true
		} else if line == "@debug"{
			s.Config.Debug = true
			return true
		}
	}
	return false
}

func Process(line string, script *Scripts, sess *session.Session){
	if strings.HasPrefix(line, "@"){
		Settings(sess, line, script)
		if script.Config.History == true{
			script.AddHistory(line)
		}
		if script.Config.Debug{
			log.Println("set setting : " + line)
		}
	} else if strings.HasPrefix(line, "{") && strings.Contains(line, "="){
		variable := strings.TrimSpace(strings.Split(line, "=")[0])
		command  := strings.TrimSpace(strings.Split(line, "=")[1])
		if strings.Contains(command, "{") && strings.Contains(command, "}"){
			r, _ := regexp.Compile("\\{(.*?)\\}")
			v := strings.TrimSpace(r.FindString(command))
			get, err := script.GetVariable(v)
			if err != nil{
				return
			}
			command = strings.Replace(command, v, get[0], 1)
			sess.ParseCommand(line)
			if script.Config.Debug{
				log.Println("assignation : " + command)
			}
		} else {
			ret := sess.ParseCommand(command)
			script.AddVar(variable, false, ret)
			if script.Config.Debug{
				log.Println("assignation : " + command)
			}
		}
		if script.Config.History == true{
			script.AddHistory(command)
		}
		return
	} else if strings.Contains(line, "{") && strings.Contains(line, "}"){
		r, _ := regexp.Compile("\\{(.*?)\\}")
		v := strings.TrimSpace(r.FindString(line))
		get, err := script.GetVariable(v)
		if err != nil{
			return
		}
		line = strings.Replace(line, v, get[0], -1)
		sess.ParseCommand(line)
		if script.Config.Debug{
			log.Println("execute : " + line)
		}
		if script.Config.History == true{
			script.AddHistory(line)
		}
	}
	return
}

func ProcessForeach(script *Scripts, sess *session.Session){
	if script.Foreach.InForeach == true{
		if script.Config.Debug{
			log.Println("Foreach execution...")
		}
		for _, context := range script.Foreach.Context{
			for _, statement := range script.Foreach.Statements{
				r, _ := regexp.Compile("\\{F(.*?)\\}")
				v := strings.TrimSpace(r.FindString(statement))
				if v != ""{
					statement = strings.Replace(statement, v, context, -1)
				}

				if strings.HasPrefix(statement, "@"){
					Settings(sess, statement, script)
					if script.Config.Debug{
						log.Println("set setting : " + statement)
					}
				} else if strings.HasPrefix(statement, "{") && strings.Contains(statement, "="){
					variable := strings.TrimSpace(strings.Split(statement, "=")[0])
					command  := strings.TrimSpace(strings.Split(statement, "=")[1])
					if strings.Contains(command, "{") && strings.Contains(command, "}"){
						r, _ := regexp.Compile("\\{(.*?)\\}")
						v := strings.TrimSpace(r.FindString(command))
						get, err := script.GetVariable(v)
						if err != nil{
							continue
						}
						command = strings.Replace(command, v, get[0], 1)
						sess.ParseCommand(statement)

						if script.Config.Debug{
							log.Println("assignation : " + statement)
						}
					}
					ret := sess.ParseCommand(command)
					script.AddVar(variable, false, ret)
				} else if strings.Contains(statement, "{") && strings.Contains(statement, "}"){
					r, _ := regexp.Compile("\\{(.*?)\\}")
					v := strings.TrimSpace(r.FindString(statement))
					get, err := script.GetVariable(v)
					if err != nil{
						continue
					}
					if len(get) > 0 {
						statement = strings.Replace(statement, v, get[0], -1)
						sess.ParseCommand(statement)
						if script.Config.Debug{
							log.Println("execute : " + statement)
						}
					}
				} else{
					sess.ParseCommand(statement)
					if script.Config.Debug{
						log.Println("execute : " + statement)
					}
				}

				if script.Config.History == true{
					script.AddHistory(statement)
				}
			}
		}
	}
	return
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
		if line != ""{
			if !strings.HasPrefix(line, "//"){
				if strings.HasPrefix(line, "foreach:") && strings.Contains(line, "=>"){
					newScript.Foreach.InForeach = true
					r, _ := regexp.Compile("\\{(.*?)\\}")
					v := strings.TrimSpace(r.FindString(line))
					results, err := newScript.GetVariable(v)
					if err != nil{
						newScript.Foreach.InForeach = false
						continue
					}
					newScript.AddContext(results)
				} else if newScript.Foreach.InForeach == true{
					if !strings.HasPrefix(line, ":endforeach") {
						newScript.AddStatement(line)
					} else{
						ProcessForeach(&newScript, sess)
						newScript.ResetForeach()
					}
				} else{
					Process(line, &newScript, sess)
				}
			} else{
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

		if newScript.Config.Debug{
			t := time.Now()
			timeText := t.Format("15:04:05")
			Settings(sess, "@verbose on", &newScript)
			sess.Stream.Standard("executed at " + timeText)
		}
	}
}

func Run2(sess *session.Session, script string){
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
			if !strings.HasPrefix(line, "//") {
				if string(line[0]) == "@" {
					Settings(sess, line, &newScript)
				} else if newScript.Foreach.InForeach{
					if newScript.Config.Debug{
						log.Println("Foreach execution...")
					}
					if strings.Contains(line,"<=)"){
						for range newScript.Foreach.Context{
							for _, stmts := range newScript.Foreach.Statements{
								if string(stmts[0]) == "@"{
									Settings(sess, stmts, &newScript)
								} else if strings.Contains(stmts, "{") && strings.Contains(stmts, "}") && !strings.Contains(stmts,"=>"){
									r, _ := regexp.Compile("\\{F(.*?)\\}")
									v := strings.TrimSpace(r.FindString(stmts))
									if string(stmts[0]) == "{"{
										variable := strings.TrimSpace(strings.Split(v, "=")[0])
										commands := strings.TrimSpace(strings.Split(v, "=")[1])
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
						if newScript.Config.Debug{
							log.Println("Foreach executed.")
						}
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
					if newScript.Config.Debug{
						log.Println("Assignation", variable)
					}
				} else if strings.Contains(line, "{") && strings.Contains(line, "}") && !strings.Contains(line,"=>") && !strings.Contains(line, "="){
					r, _ := regexp.Compile("\\{(.*?)\\}")
					v := strings.TrimSpace(r.FindString(line))
					for _, vars := range newScript.Vars{
						if vars.Name == v{
							for _, value := range vars.Value{
								newCommand := strings.Replace(line, v, value, -1)
								sess.ParseCommand(newCommand)
								newScript.AddHistory(newCommand)
								if newScript.Config.Debug{
									log.Println("Command execution", newCommand)
								}
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

		if newScript.Config.Debug{
			t := time.Now()
			timeText := t.Format("15:04:05")
			Settings(sess, "@verbose on", &newScript)
			sess.Stream.Standard("executed at " + timeText)
		}
	}
}
