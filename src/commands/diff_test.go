package commands_test

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("diff()", func() {
	var workingDir string
	var logger *core.SpyLogger

	_ = BeforeEach(func() {
		workingDir, _ = os.MkdirTemp(os.TempDir(), "wd-*")
		logger = core.MakeSpyLogger()
		color.NoColor = true
	})

	_ = AfterEach(func() {
		os.RemoveAll(workingDir)
	})

	It("should not diff anything if no files exists on `fromDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should not diff anything if no readable file exists on `fromDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "*")

		os.Chmod(from.Name(), 0o300)

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, &core.SpyLoggerCallNumber{Warnnl: 1})
		Expect(logger.Calls.Warnnl[0].Args).To(
			Equal(

				[]any{
					"File %s does not exists or is not a file or is not readable, skipping...",
					from.Name(),
				},
			),
		)
	})

	It("should not diff anything if no matching file exists on `toDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "*")
		defer from.Close()
		to := strings.Replace(from.Name(), fromDir, toDir, 1)

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, &core.SpyLoggerCallNumber{Warnnl: 1})
		Expect(logger.Calls.Warnnl[0].Args).To(
			Equal(
				[]any{
					"File %s does not exists or is not a file or is not readable, skipping...",
					to,
				},
			),
		)
	})

	It("should not diff anything if matching file is not readable on `toDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "*")
		defer from.Close()
		to := strings.Replace(from.Name(), fromDir, toDir, 1)

		fsutil.CopyFile(from.Name(), to)
		os.Chmod(to, 0o300)

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, &core.SpyLoggerCallNumber{Warnnl: 1})
		Expect(logger.Calls.Warnnl[0].Args).To(
			Equal(
				[]any{
					"File %s does not exists or is not a file or is not readable, skipping...",
					to,
				},
			),
		)
	})

	It("should diff and not show differences if files are the same", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "*")
		defer from.Close()
		to := strings.Replace(from.Name(), fromDir, toDir, 1)

		from.WriteString("foobar")
		fsutil.CopyFile(from.Name(), to)

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, &core.SpyLoggerCallNumber{Lognl: 1, Log: 1})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Diffing %s against %s ...", from.Name(), to}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
	})

	It("should diff and show diferences if files are not the same", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "1-*")
		defer from.Close()
		from2, _ := os.CreateTemp(fromDir, "2-*")
		defer from2.Close()
		to := strings.Replace(from.Name(), fromDir, toDir, 1)
		to2 := strings.Replace(from2.Name(), fromDir, toDir, 1)

		from.WriteString("foobar")
		fsutil.CopyFile(from.Name(), to)
		from2.WriteString("foobuz")
		fsutil.CopyFile(from2.Name(), to2)
		from2.Truncate(0)
		from2.Seek(0, 0)
		from2.WriteString("foobiz")

		result, err := commands.DiffCmd(commands.DiffCmdArgs{
			FromDir: fromDir,
			ToDir:   toDir,
			Extra: commands.DiffCmdArgsExtra{
				Logger: logger,
			},
		})

		Expect(result).To(BeFalse())
		Expect(err).To(BeNil())
		core.AssertSpyLoggerCalls(*logger, &core.SpyLoggerCallNumber{Lognl: 9, Log: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Diffing %s against %s ...", from.Name(), to}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Diffing %s against %s ...", from2.Name(), to2}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✕"}))
		Expect(
			logger.Calls.Lognl[2].Args,
		).To(Equal([]any{strings.Join([]string{"---", from2.Name()}, " ")}))
		Expect(logger.Calls.Lognl[3].Args).To(Equal([]any{strings.Join([]string{"+++", to2}, " ")}))
		Expect(logger.Calls.Lognl[4].Args).To(Equal([]any{"@@ -1 +1 @@"}))
		Expect(logger.Calls.Lognl[5].Args).To(Equal([]any{"-foobiz"}))
		Expect(logger.Calls.Lognl[6].Args).To(Equal([]any{"\\ No newline at end of file"}))
		Expect(logger.Calls.Lognl[7].Args).To(Equal([]any{"+foobuz"}))
		Expect(logger.Calls.Lognl[8].Args).To(Equal([]any{"\\ No newline at end of file"}))
	})
})
