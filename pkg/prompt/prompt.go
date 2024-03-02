package prompt

import (
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gofsd/fsd/pkg/cfg"
	"github.com/gofsd/fsd/pkg/util"
	"github.com/spf13/cobra"
)

type currentCommands struct {
	CompiledCommandsTree    []prompt.Suggest
	CommandsTreeHistory     map[string][]prompt.Suggest
	CurrentLine             string
	PrevLine                string
	CmdPrefix               string
	CurrentFullValidCmdPath string
	CurrentCmdPath          string
	CurrentCmdPathArr       []string
	SubCommandsTree         map[string][]prompt.Suggest
}

func (cC *currentCommands) SetSubCommandsTree() {
	if cC.SubCommandsTree == nil {
		cC.SubCommandsTree = make(map[string][]prompt.Suggest)
	}
	for _, i := range cC.CompiledCommandsTree {
		cmdArr := strings.Split(i.Text, " ")
		for idx, _ := range cmdArr {
			var cmdArrSlice []string
			for j := 0; j < idx; j++ {
				cmdArrSlice = append(cmdArrSlice, cmdArr[j])
			}
			itemStr := strings.Join(cmdArrSlice, " ")
			subStrCmd := strings.Replace(i.Text, itemStr, "", 1)
			subStrCmd = strings.Trim(subStrCmd, " ")
			cC.SubCommandsTree[itemStr] = append(cC.SubCommandsTree[itemStr], prompt.Suggest{Text: subStrCmd, Description: i.Description})
		}
	}
}

func (cC *currentCommands) SetCurrentLine(s string) {
	if cC.CurrentLine == s {
		return
	} else if s != "" {
		cC.PrevLine = cC.CurrentLine
		cC.CurrentLine = s
		cC.SetCurrentCmdPath()
		cC.SetValidCmdPath()
		cC.ExecControlSeq()
	}
}
func (cC *currentCommands) ExecControlSeq() {
	if strings.Contains(cC.CurrentLine, "..") {
		cC.CurrentFullValidCmdPath = ""
	}
}
func (cC *currentCommands) SetValidCmdPath() {
	var validCmdPathAr []string
	var validCmdPath string

	for _, c := range cC.CurrentCmdPathArr {
		validCmdPathAr = append(validCmdPathAr, c)
		current := strings.Join(validCmdPathAr, " ")
		if value, ok := cC.SubCommandsTree[current]; ok {
			if len(value) > 0 {
				validCmdPath += current
			}
		}
	}
	cC.CurrentFullValidCmdPath = validCmdPath
}

func (cC *currentCommands) SetCurrentCmdPath() {
	preparedCurrLine := strings.Trim(cC.CurrentLine, " ")
	preparedCmdPrefix := strings.Trim(cC.CmdPrefix, " ")
	if preparedCurrLine != "" && preparedCmdPrefix != "" {
		cC.CurrentCmdPath = preparedCmdPrefix + " " + preparedCurrLine
	} else if preparedCurrLine != "" {
		cC.CurrentCmdPath = preparedCurrLine
	} else {
		cC.CurrentCmdPath = preparedCmdPrefix
	}
	cC.CurrentCmdPathArr = strings.Split(cC.CurrentCmdPath, " ")
}

func (cC *currentCommands) ShowCompiledCommandsTree() bool {
	if len(cC.CurrentCmdPathArr) > 1 {
		return false
	}
	return true
}

func (cC *currentCommands) GetCurrentCommandsTree() []prompt.Suggest {
	return cC.SubCommandsTree[cC.CurrentFullValidCmdPath]
}

func (cC *currentCommands) GetCurrentWord() string {
	currentWord := strings.Replace(cC.CurrentCmdPath, cC.CurrentFullValidCmdPath, "", 1)
	currentWord = strings.Trim(currentWord, " ")

	return ""
}

func (cC *currentCommands) DeleteFromHistory(s string) {
	delete(cC.CommandsTreeHistory, s)
}

func (cCT *currentCommands) rmRootCmdString(root []string) {
	log.Println("Test: ", root)

	for _, s := range root {
		for idx, v := range cCT.CommandsTreeHistory[""] {
			s = strings.Trim(s, " ")
			cCT.CommandsTreeHistory[""][idx].Text = strings.Replace(v.Text, s+" ", "", 1)
		}
	}
}

type FsdPrompt struct {
	prompt.Prompt
	prefix          string
	cmdPath         string
	CurrentCommands currentCommands
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func New() *FsdPrompt {
	defer handleExit()
	var p FsdPrompt
	p.SetCompiledCommandsTree(cfg.Config.RootCmd)
	p.CurrentCommands.SetSubCommandsTree()
	p.SortSuggestions()
	p.Prompt = *prompt.New(p.executor, p.completer, prompt.OptionMaxSuggestion(12), prompt.OptionLivePrefix(p.livePrefix), prompt.OptionSetExitCheckerOnInput(p.exit), prompt.OptionBreakLineCallback(p.breakLine))
	return &p
}

func (fp *FsdPrompt) breakLine(doc *prompt.Document) {

}

func (fp *FsdPrompt) exit(in string, breakLine bool) bool {
	if in == "exit" {
		return true
	} else {
		return false
	}
}

func (fp *FsdPrompt) executor(s string) {
	config := cfg.GetCfg()
	args := strings.Split(s, " ")
	for i, a := range args {
		args[i] = strings.Trim(a, " ")
		if args[i] == "" {
			args = append(args[:i], args[i+1:]...)
		}
	}
	if len(fp.CurrentCommands.CurrentCmdPathArr) > 1 && fp.CurrentCommands.CmdPrefix != "" {
		args = fp.CurrentCommands.CurrentCmdPathArr
		s = fp.CurrentCommands.CmdPrefix
	}
	if cmd := config.FindCMD(s); cmd.Name() != cfg.CliName {
		config.CmdPath = cmd.CommandPath()
		config.RootCmd.SetArgs(args)
		config.RootCmd.Execute()
	} else {
		util.Exec(args[0], 2)
	}
	fp.CurrentCommands.CmdPrefix = fp.cmdPath
}

func (fp *FsdPrompt) completer(d prompt.Document) []prompt.Suggest {
	fp.CurrentCommands.SetCurrentLine(d.CurrentLine())
	v := fp.CurrentCommands.GetCurrentCommandsTree()

	return prompt.FilterFuzzy(v, fp.CurrentCommands.GetCurrentWord(), true)
}

func (fp *FsdPrompt) livePrefix() (string, bool) {
	s, _ := os.Getwd()
	fp.cmdPath = cfg.GetCfg().CmdPath
	fp.prefix = strings.Join(strings.Split(fp.cmdPath, " ")[1:], " ")
	fp.CurrentCommands.CmdPrefix = fp.prefix
	fp.prefix = s + ": " + fp.prefix + "> "
	fp.SortWithPrefix()
	//c.addNewSuggestChildOnStart(c.getSuggestionsByPath(c.CmdPath))
	return fp.prefix, true
}

func (fp *FsdPrompt) SortSuggestions() {
	cmdTree := fp.CurrentCommands.GetCurrentCommandsTree()
	sort.Slice(cmdTree, func(i, j int) bool {
		return len(cmdTree[i].Text) < len(cmdTree[j].Text)
	})
}

func (fp *FsdPrompt) addNewSuggestChildOnStart(s []prompt.Suggest) {
	cmdTree := fp.CurrentCommands.GetCurrentCommandsTree()

	cmdTree = append(cmdTree, s...)
}

func (fp *FsdPrompt) GetCompiledCommandsTree(cmd *cobra.Command, sg []prompt.Suggest) []prompt.Suggest {
	for _, cmd := range cmd.Commands() {
		cmdPath := strings.Join(strings.Split(cmd.CommandPath(), " ")[1:], " ")
		sg = append(sg, prompt.Suggest{Text: cmdPath, Description: cmd.Short})
		if len(cmd.Commands()) > 0 {
			sg = fp.GetCompiledCommandsTree(cmd, sg)
		}
	}
	return sg
}

func (fp *FsdPrompt) SetCompiledCommandsTree(cmd *cobra.Command) {
	var sg []prompt.Suggest
	fp.CurrentCommands.CommandsTreeHistory = make(map[string][]prompt.Suggest)
	fp.CurrentCommands.CompiledCommandsTree = fp.GetCompiledCommandsTree(cmd, sg)
}

func (fp *FsdPrompt) SortWithPrefix() {
	cmdTree := fp.CurrentCommands.GetCurrentCommandsTree()
	log.Printf(fp.prefix)
	sort.SliceStable(cmdTree, func(i, j int) bool {
		prefix := cmdTree[i].Text
		prefixJ := cmdTree[j].Text

		if len(prefix) > len(prefix) {
			return false
		}
		if strings.Contains(fp.prefix, prefix) && strings.Contains(fp.prefix, prefixJ) {
			return true
		} else if strings.Contains(fp.prefix, prefix) && !strings.Contains(fp.prefix, prefixJ) {
			return true
		}
		return false

	})
}
