package gdc

import(
	"fmt"
	"os"
	"strings"

	"github.com/codeskyblue/go-sh"
)

func shellCommand(session *sh.Session, cmd string) *sh.Session {
  tokens := strings.Split(cmd, " ")
	tokensInt := []interface{}{} // doesn't seem to be any direct way to cast []string to []interface{}
	for _, t := range(tokens[1:]) {
		tokensInt = append(tokensInt, t)
	}
  return session.Command(tokens[0], tokensInt...)
}

//RunCommand with a command string and arguments for interpolation using Sprintf
func RunCommand(cmd string, args... interface{}) {
	writeDcFile()
	fullCommand := cmd
	if (len(args) > 0) {
		fullCommand = fmt.Sprintf(cmd, args...)
	}
	fmt.Printf("-> %s\n", fullCommand)
	command := shellCommand(sh.InteractiveSession(), fullCommand)
	command.SetEnv("KAFKA_ADV_HOST", os.Getenv("KAFKA_ADV_HOST"))
	command.SetStdin(os.Stdin)
	err := command.Run()
	if (err != nil) {
		Exit("Error running command! %s", fullCommand)
	}
}

// RunCommands run a list of commands to be piped into each other
func RunCommands(commands... string) {
	writeDcFile()
	for i, cmd := range(commands) {
     if (i == 0) {
			 fmt.Printf("-> %s", cmd)
		 } else {
			 fmt.Printf(" | %s", cmd)
		 }
	}
	fmt.Println()
	
	session := sh.InteractiveSession()
	session.PipeStdErrors = true
	session.PipeFail = true
	session.SetStdin(os.Stdin)
	for _, cmd := range(commands) {
		session = shellCommand(session, cmd)
	}
	err := session.Run()
	if (err != nil) {
		Exit("Error running command! %v", commands)
	}
}

