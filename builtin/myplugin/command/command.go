package command

import (
	"github.com/DavidGamba/go-getoptions/option"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	plugincore "github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Subcommand interface {
	CommandInfo() (*component.CommandInfo, error)
}

type CommandConfig struct {
}

// Command is the Command implementation for myplugin.
type Command struct {
	config CommandConfig
}

func (c *Command) ConfigSet(v interface{}) error {
	return nil
}

func (c *Command) CommandFunc() interface{} {
	return nil
}

func (c *Command) Config() (interface{}, error) {
	return &c.config, nil
}

func (c *Command) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&CommandConfig{}))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ExecuteFunc implements component.Command
func (c *Command) ExecuteFunc(cliArgs []string) interface{} {
	if len(cliArgs) < 2 {
		return c.Execute
	}
	n := cliArgs[1]
	switch n {
	case "info":
		return c.ExecuteInfo
	case "dothing":
		return c.ExecuteDoThing
	case "interactive":
		return c.ExecuteInteractive
	}

	return c.Execute
}

func (c *Command) ExecuteInfo(trm terminal.UI, env plugincore.Project) int64 {
	return (&Info{Command: c}).Execute(trm, env)
}

func (c *Command) ExecuteDoThing(trm terminal.UI, flags map[string]interface{}) int64 {
	return (&DoThing{Command: c}).Execute(trm, flags)
}

func (c *Command) ExecuteInteractive(trm terminal.UI, flags map[string]interface{}) int64 {
	return (&Interactive{Command: c}).Execute(trm)
}

// CommandInfoFunc implements component.Command
func (c *Command) CommandInfoFunc() interface{} {
	return c.CommandInfo
}

func (c *Command) CommandInfo() *component.CommandInfo {
	return &component.CommandInfo{
		Name:        "myplugin",
		Help:        c.Help(),
		Synopsis:    c.Synopsis(),
		Flags:       c.Flags(),
		Subcommands: c.subcommandsInfo(),
	}
}

func (c *Command) Synopsis() string {
	return "I don't do much, just hanging around"
}

func (c *Command) Help() string {
	return "I'm here for testing, try running some subcommands"
}

func (c *Command) Flags() []*option.Option {
	stringflag := option.New("hehe", option.StringType)
	stringflag.Description = "a test flag for strings"
	stringflag.DefaultStr = "message"
	stringflag.Aliases = append(stringflag.Aliases, "hh")

	return []*option.Option{stringflag}
}

func (c *Command) Execute(trm terminal.UI, flags map[string]interface{}) int64 {
	trm.Output("You gave me the flag: " + flags["hehe"].(string))

	trm.Output(c.Help())
	trm.Output("My subcommands are: ")
	for _, cmd := range c.subcommandsInfo() {
		trm.Output("    " + cmd.Name)
	}
	return 0
}

func (c *Command) subcommandsInfo() (r []*component.CommandInfo) {
	for _, cmd := range c.subcommands() {
		v, _ := cmd.CommandInfo()
		r = append(r, v)
	}
	return
}

func (c *Command) subcommands() map[string]Subcommand {
	return map[string]Subcommand{
		"info":        &Info{Command: c},
		"dothing":     &DoThing{Command: c},
		"interactive": &Interactive{Command: c},
	}
}

var (
	_ component.Command = (*Command)(nil)
)
