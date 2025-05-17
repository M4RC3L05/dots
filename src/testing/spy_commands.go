package testing

import (
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
	"github.com/onsi/gomega"
)

type SpyCommandsCalls struct {
	Diff  []core.SpyCallNoRt
	Adopt []core.SpyCallNoRt
	Apply []core.SpyCallNoRt
}
type SpyCommandsCallNumber struct {
	Diff  int
	Adopt int
	Apply int
}

type SpyCommands struct {
	commands.ICommands

	Calls SpyCommandsCalls
}

func (sl *SpyCommands) Diff(args commands.DiffArgs) (bool, error) {
	sl.Calls.Diff = append(sl.Calls.Diff, core.SpyCallNoRt{Args: []any{args}})

	return true, nil
}

func (sl *SpyCommands) Adopt(args commands.AdoptArgs) (bool, error) {
	sl.Calls.Adopt = append(sl.Calls.Adopt, core.SpyCallNoRt{Args: []any{args}})

	return true, nil
}

func (sl *SpyCommands) Apply(args commands.ApplyArgs) (bool, error) {
	sl.Calls.Apply = append(sl.Calls.Apply, core.SpyCallNoRt{Args: []any{args}})

	return true, nil
}

func MakeSpyCommands() *SpyCommands {
	return &SpyCommands{}
}

func AssertSpyCommandsCalls(Command SpyCommands, callNumber *SpyCommandsCallNumber) {
	var callNumberVal SpyCommandsCallNumber

	if callNumber != nil {
		callNumberVal = *callNumber
	} else {
		callNumberVal = SpyCommandsCallNumber{}
	}

	gomega.Expect(Command.Calls.Diff).To(gomega.HaveLen(callNumberVal.Diff))
	gomega.Expect(Command.Calls.Adopt).To(gomega.HaveLen(callNumberVal.Adopt))
	gomega.Expect(Command.Calls.Apply).To(gomega.HaveLen(callNumberVal.Apply))
}
