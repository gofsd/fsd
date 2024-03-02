package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	typs "github.com/gofsd/fsd-types"
	"github.com/gofsd/fsd/pkg/log"
	"github.com/gofsd/fsd/pkg/store"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	ROOT = iota
	EQUAL
	AND
	OR
	ELSE
)

type KV struct {
	Key   string `json:"key" validate:"omitempty,min=0,max=9"`
	Value string `json:"value" validate:"omitempty,min=0,max=256"`
}

type Command struct {
	typs.Command
}

func (cmd *Command) Json() []byte {
	b, _ := json.Marshal(cmd)
	return b
}
func (cmd *Command) Gob() []byte {
	b, _ := json.Marshal(cmd)
	return b
}

func (cmd *Command) String() (s string) {
	json.Marshal(cmd)
	return s
}

type Result []byte

// type Commands map[string]Cmd

// var commands = make(Commands)

// func (c *Commands) Add(cmdName string, cmd Cmd) *Commands {
// 	if _, ok := (*c)[cmdName]; ok {
// 		panic("Cmd duplication")
// 	}
// 	(*c)[cmdName] = cmd
// 	return c
// }

type Helper struct {
	canExecute, andCondition, orCondition, executed bool
	cmd                                             *cobra.Command
	prev                                            int
	E                                               error
}

func Log(message string) {
	l := log.Get("")
	l.Info(message)
}

func Set(cmd *cobra.Command) *Helper {
	var h Helper
	h.cmd = cmd
	h.prev = ROOT
	return &h
}

func (h *Helper) Error() *Helper {
	if !h.canExecute && !h.executed && h.E == nil {
		// h.cmd.PrintErr("handler not found")
		h.E = fmt.Errorf("handler not found for command: %s. With flags: %s.", h.cmd.CommandPath(), h.FlagsToString())
		h.cmd.PrintErr(h.E.Error())
	}
	h.executed = false

	return h
}

func (h *Helper) FlagsToString() (s string) {
	flagSet := h.cmd.Flags()
	flagSet.VisitAll(func(f *pflag.Flag) {
		s = fmt.Sprintf("%s  %s=%s", s, f.Name, f.Value)
	})
	return
}

func (h *Helper) Throw(message string) *Helper {
	if h.canExec() {
		if len(message) > 0 {
			h.setExecuted(errors.New(message))
		}
	}

	return h
}

func (h *Helper) Equal(a, b any) *Helper {
	if h.prev == AND && !h.canExecute {
		h.prev = EQUAL
		return h
	} else if h.prev == OR && h.canExecute {
		h.prev = EQUAL
		return h
	}
	h.canExecute = reflect.DeepEqual(a, b)
	h.prev = EQUAL
	// h.executed = false

	return h
}

func (h *Helper) NotEqual(a, b any) *Helper {
	if h.prev == AND && !h.canExecute {
		h.prev = EQUAL
		return h
	} else if h.prev == OR && h.canExecute {
		h.prev = EQUAL
		return h
	}
	h.canExecute = !reflect.DeepEqual(a, b)

	h.prev = EQUAL

	return h
}

func (h *Helper) AND() *Helper {
	h.prev = AND
	return h
}

func (h *Helper) OR() *Helper {
	h.prev = OR
	return h
}

func (h *Helper) Else() *Helper {
	if !h.executed {
		h.canExecute = true
	}
	h.prev = ELSE
	return h
}

func (h *Helper) canExec() bool {
	return h.E == nil && h.canExecute && !h.executed
}

func (h *Helper) setExecuted(e error) {
	h.executed = true
	h.canExecute = false
	h.E = e
}

func (h *Helper) JustFn(handler func()) *Helper {
	if !h.canExec() {
		return h
	}
	handler()
	return h
}

func (h *Helper) JustB(handler func() []byte) *Helper {
	if !h.canExec() {
		return h
	}
	b := handler()
	h.pickHandler(b, nil)
	h.setExecuted(nil)
	return h
}

func (h *Helper) JustStr(handler func() string) *Helper {
	if !h.canExec() {
		return h
	}
	s := handler()
	h.pickHandler(s, nil)
	h.setExecuted(nil)
	return h
}

func (h *Helper) HandleB(handler func() ([]byte, error)) *Helper {
	if !h.canExec() {
		return h
	}
	b, e := handler()
	h.pickHandler(b, e)
	h.setExecuted(e)
	return h
}

func (h *Helper) HandleError(t *Helper) *Helper {

	return h
}

func (h *Helper) HandleCRUD(handler func(typs.ICrud) error, args typs.ICrud) *Helper {
	if !h.canExec() {
		return h
	}
	e := handler(args)
	h.setExecuted(e)
	return h
}

func (h *Helper) HandleStr(handler func() (string, error)) *Helper {
	if !h.canExec() {
		return h
	}
	b, e := handler()
	h.pickHandler(b, e)
	h.setExecuted(e)
	return h
}

func (h *Helper) pickHandler(b any, e error) {
	switch v := b.(type) {
	case []byte:
		h.handleBytes(v, e)
	case string:
		h.handleString(v, e)
	}
}

func (h *Helper) handleBytes(b []byte, e error) *Helper {
	if !h.canExecute {
		return h
	}

	if e != nil {
		h.cmd.PrintErr(e)
	}
	h.cmd.OutOrStdout().Write(b)
	return h

}

func (h *Helper) HandleJson(b []byte, e error) *Helper {
	if !h.canExecute {
		return h
	}
	if e != nil {
		h.cmd.PrintErr(e)
	}
	h.cmd.OutOrStdout().Write(b)
	return h
}

func (h *Helper) handleString(s string, e error) *Helper {
	if !h.canExecute {
		return h
	}
	if e != nil {
		h.cmd.PrintErr(e)
	}
	h.cmd.Print(s)
	return h
}

type CommandStore struct {
	DB *store.Store
}

func (str *CommandStore) Save(command typs.Command, cmdOutput []byte) (cmd typs.CommandResponse) {
	if len(cmdOutput) < 256 {
		cmd.Command = command
		json.Unmarshal(cmdOutput, &cmd.Result)
		str.DB.JustSet(&cmd)
	}

	return cmd
}

func Store() CommandStore {
	var str CommandStore
	str.DB = store.New(store.SetFullDbName("commands"))
	return str
}

func Test(c *cobra.Command, params typs.Command) (out []byte, er error) {
	cmd, args, _ := c.Find(params.Name)

	for _, kv := range params.Flags {
		cmd.Flags().Set(kv.K, kv.V)
	}
	var o, e *bytes.Buffer
	o, e = bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	cmd.SetOut(o)
	cmd.SetErr(e)

	er = cmd.RunE(cmd, args)
	out, _ = io.ReadAll(o)
	return out, er

}

func If[T any](b bool, first, second T) T {
	if b == true {
		return first
	} else {
		return second
	}
}
