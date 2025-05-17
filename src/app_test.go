package src_test

import (
	"github.com/m4rc3l05/dots/src"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("app()", func() {
	var cmds *testing.SpyCommands
	var displays *testing.SpyDisplays
	var logger *testing.SpyLogger

	_ = BeforeEach(func() {
		logger = testing.MakeSpyLogger()
		cmds = testing.MakeSpyCommands()
		displays = testing.MakeSpyDisplays()
	})

	It("should print help if `help` flag is provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Flags: src.CmdFlagsArgs{
					Help: true,
				},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, nil)
		testing.AssertSpyDisplaysCalls(*displays, &testing.SpyDisplaysCallNumber{Help: 1})
	})

	It("should print environment if `printEnvironment` flag is provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Flags: src.CmdFlagsArgs{
					PrintEnvironment: true,
				},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, nil)
		testing.AssertSpyDisplaysCalls(*displays, &testing.SpyDisplaysCallNumber{Environment: 1})
	})

	It("should print version if `version` flag is provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Flags: src.CmdFlagsArgs{
					Version: true,
				},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, nil)
		testing.AssertSpyDisplaysCalls(*displays, &testing.SpyDisplaysCallNumber{Version: 1})
	})

	It("should run diff if `diff` command provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"diff"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, &testing.SpyCommandsCallNumber{Diff: 1})
		testing.AssertSpyDisplaysCalls(*displays, nil)
	})

	It("should run apply if `apply` command provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"apply"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, &testing.SpyCommandsCallNumber{Apply: 1})
		testing.AssertSpyDisplaysCalls(*displays, nil)
	})

	It("should run apply with first arg if `apply` command provided with an arg", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"apply", "foo"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, &testing.SpyCommandsCallNumber{Apply: 1})
		testing.AssertSpyDisplaysCalls(*displays, nil)
		Expect(cmds.Calls.Apply[0].Args).To(Equal([]any{commands.ApplyArgs{From: "foo"}}))
	})

	It("should run adopt if `adopt` command provided", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"adopt"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, &testing.SpyCommandsCallNumber{Adopt: 1})
		testing.AssertSpyDisplaysCalls(*displays, nil)
	})

	It("should run adopt with first arg if `adopt` command provided with an arg", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"adopt", "foo"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, nil)
		testing.AssertSpyCommandsCalls(*cmds, &testing.SpyCommandsCallNumber{Adopt: 1})
		testing.AssertSpyDisplaysCalls(*displays, nil)
		Expect(cmds.Calls.Adopt[0].Args).To(Equal([]any{commands.AdoptArgs{From: "foo"}}))
	})

	It("should print help if command is not supported", func() {
		src.App(src.Args{
			CmdArgs: src.CmdArgs{
				Rest: []string{"foo"},
			},
			Displays: displays,
			Commands: cmds,
			Logger:   logger,
		})

		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Warnnl: 1})
		testing.AssertSpyCommandsCalls(*cmds, nil)
		testing.AssertSpyDisplaysCalls(*displays, &testing.SpyDisplaysCallNumber{Help: 1})
		Expect(logger.Calls.Warnnl[0].Args).To(Equal([]any{"Command %s not found", "foo"}))
	})
})
