package services

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/cli/go-gh/v2/pkg/tableprinter"
	"github.com/cli/go-gh/v2/pkg/term"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var (
	prompt       *prompter.Prompter
	terminal     *term.Term
	tablePrinter *tableprinter.TablePrinter
)

func GetTerminal() *term.Term {
	if terminal == nil {
		t := term.FromEnv()
		terminal = &t
	}

	return terminal
}

func GetPrompt() *prompter.Prompter {

	if prompt == nil {

		terminal := GetTerminal()

		in, ok := terminal.In().(*os.File)
		if !ok {
			log.Fatal("error casting to file")
		}
		out, ok := terminal.Out().(*os.File)
		if !ok {
			log.Fatal("error casting to file")
		}
		errOut, ok := terminal.ErrOut().(*os.File)
		if !ok {
			log.Fatal("error casting to file")
		}

		prompt = prompter.New(in, out, errOut)
	}

	return prompt
}

func getTablePrinter() (tableprinter.TablePrinter, error) {
	if tablePrinter == nil {
		t := GetTerminal()
		w, _, err := t.Size()
		if err != nil {
			return nil, err
		}
		tb := tableprinter.New(t.Out(), t.IsTerminalOutput(), w)
		tablePrinter = &tb
	}
	return *tablePrinter, nil
}

func enabledOrDisabled(b bool) string {
	if b {
		return "Enabled"
	}
	return "Disabled"
}

func ChooseSubCommand(subCmds []*cobra.Command, args []string, promptTitle string) {

	sort.Sort(ByCommandName(subCmds))

	GetPrompt()

	subCommands := make(map[string]*cobra.Command)

	for _, c := range subCmds {
		subCommands[c.Name()] = c
	}

	listOptions := maps.Keys(subCommands)
	fmt.Println(listOptions)

	option, err := prompt.Select(promptTitle, "", listOptions)
	if err != nil {
		os.Exit(1)
	}
	choosen := listOptions[option]
	subCommands[choosen].Run(subCommands[choosen], args)

}

func GetTarget(args []string, message string) []string {

	if len(args) < 1 {
		name, err := prompt.Input(message, "")
		if err != nil || len(name) == 0 {
			fmt.Printf("Unable to %s.", message)
			os.Exit(1)
		}
		args = append(args, name)
	}

	return args
}
