package services

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

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

	GetPrompt()

	subCommands := make(map[string]*cobra.Command)

	for _, c := range subCmds {
		subCommands[c.Name()] = c
	}

	listOptions := maps.Keys(subCommands)

	sort.Strings(listOptions)

	option, err := prompt.Select(promptTitle, "", listOptions)
	if err != nil {
		os.Exit(1)
	}
	choosen := listOptions[option]
	subCommands[choosen].Run(subCommands[choosen], args)
}

// GetTarget parses the target and all global flags.
// Returns: target(string), flags(GlobalFlags)
func GetTarget(cmd *cobra.Command, args []string, message string) (string, *GlobalFlags) {
	var flags *GlobalFlags

	var target string
	if len(args) < 1 {
		GetPrompt()
		response, err := prompt.Input(message, "")
		if err != nil || len(response) == 0 {
			fmt.Printf("Unable to read input: %v\n", err)
			os.Exit(1)
		}

		input := strings.Split(response, " ")

		flags, err = ParseGlobalFlags(cmd, input)
		if err != nil {
			fmt.Printf("Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		target = input[0]
	} else {
		target = args[0]
		flags = GetGlobalFlags()
	}

	return target, flags
}

// GetOptimalPageSize calculates the page size based on user flag or terminal height
func GetOptimalPageSize(userSize int) int {
	const ApiMax = 100

	// 1. If user specified a flag, use it (clamped to ApiMax)
	if userSize > 0 {
		if userSize > ApiMax {
			return ApiMax
		}
		return userSize
	}

	// 2. Calculate based on Terminal Height
	t := GetTerminal()
	_, height, err := t.Size()
	if err != nil {
		// Fallback if we can't detect terminal size
		return 20
	}

	// 3. Subtract 2 for the prompt/header line
	size := height - 2

	// 4. Safety clamps
	if size > ApiMax {
		return ApiMax
	}
	if size < 1 {
		return 10 // Minimal fallback
	}

	return size
}
