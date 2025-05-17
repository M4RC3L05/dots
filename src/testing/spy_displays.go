package testing

import (
	"github.com/m4rc3l05/dots/src/core"
	"github.com/m4rc3l05/dots/src/displays"
	"github.com/onsi/gomega"
)

type SpyDisplaysCalls struct {
	Environment []core.SpyCallNoRt
	Help        []core.SpyCallNoRt
	Version     []core.SpyCallNoRt
}
type SpyDisplaysCallNumber struct {
	Environment int
	Help        int
	Version     int
}

type SpyDisplays struct {
	displays.IDisplays

	Calls SpyDisplaysCalls
}

func (sd *SpyDisplays) Environment(homedir string, dotfilesFilesDir string) {
	sd.Calls.Environment = append(
		sd.Calls.Environment,
		core.SpyCallNoRt{Args: []any{homedir, dotfilesFilesDir}},
	)
}

func (sd *SpyDisplays) Help() {
	sd.Calls.Help = append(sd.Calls.Help, core.SpyCallNoRt{Args: []any{}})
}

func (sd *SpyDisplays) Version(version string) {
	sd.Calls.Version = append(sd.Calls.Version, core.SpyCallNoRt{Args: []any{version}})
}

func MakeSpyDisplays() *SpyDisplays {
	return &SpyDisplays{}
}

func AssertSpyDisplaysCalls(Command SpyDisplays, callNumber *SpyDisplaysCallNumber) {
	var callNumberVal SpyDisplaysCallNumber

	if callNumber != nil {
		callNumberVal = *callNumber
	} else {
		callNumberVal = SpyDisplaysCallNumber{}
	}

	gomega.Expect(Command.Calls.Environment).To(gomega.HaveLen(callNumberVal.Environment))
	gomega.Expect(Command.Calls.Help).To(gomega.HaveLen(callNumberVal.Help))
	gomega.Expect(Command.Calls.Version).To(gomega.HaveLen(callNumberVal.Version))
}
