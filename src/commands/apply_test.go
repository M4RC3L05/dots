package commands_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/core"
	"github.com/m4rc3l05/dots/src/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("adopt()", func() {
	var workingDir string
	var logger *testing.SpyLogger
	var cmd commands.Commands

	_ = BeforeEach(func() {
		workingDir, _ = os.MkdirTemp(os.TempDir(), "wd-*")
		logger = testing.MakeSpyLogger()
		color.NoColor = true
		cmd = commands.Commands{Logger: logger}
	})

	_ = AfterEach(func() {
		os.RemoveAll(workingDir)
	})

	It("should normalize `from` path", func() {
		result, err := cmd.Apply(commands.ApplyArgs{
			From: "/foo/bar/..",
			Extra: commands.ApplyArgsExtra{
				Homedir:          workingDir,
				DotfilesFilesDir: workingDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not readable", "/foo")))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `from` does not exists", func() {
		result, err := cmd.Apply(commands.ApplyArgs{
			From: workingDir + "/foo",
			Extra: commands.ApplyArgsExtra{
				Homedir:          workingDir,
				DotfilesFilesDir: workingDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not readable", workingDir+"/foo")))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `from` is not readable", func() {
		from, _ := os.MkdirTemp(workingDir, "*")

		os.Chmod(from, 0o300)

		result, err := cmd.Apply(commands.ApplyArgs{
			From: from,
			Extra: commands.ApplyArgsExtra{
				Homedir:          workingDir,
				DotfilesFilesDir: workingDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(err).To(MatchError(fmt.Sprintf("path %s does not exists or is not readable", from)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `from` is not a subdirectory of `dotfilesFilesDir`", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.MkdirTemp(workingDir, "*")

		result, err := cmd.Apply(commands.ApplyArgs{
			From: from,
			Extra: commands.ApplyArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s is not a subpath of %s", from, dotfilesFilesDir)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should apply a single file", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(dotfilesFilesDir, "*")
		defer from.Close()

		result, err := cmd.Apply(commands.ApplyArgs{
			From: from.Name(),
			Extra: commands.ApplyArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 1, Lognl: 1})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Applying %s to %s ...", from.Name(), strings.Replace(from.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
	})

	It("should return an error if something appens while applying file", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(dotfilesFilesDir, "*")
		defer from.Close()
		to, _ := os.Create(strings.Replace(from.Name(), dotfilesFilesDir, homedir, 1))
		defer to.Close()

		os.Chmod(to.Name(), 0o400)

		result, err := cmd.Apply(commands.ApplyArgs{
			From: from.Name(),
			Extra: commands.ApplyArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		unwrapErrors, _ := err.(interface{ Unwrap() []error })
		Expect(
			unwrapErrors.Unwrap()[0],
		).To(MatchError(fmt.Sprintf("error applying %s to %s", from.Name(), to.Name())))
		Expect(
			unwrapErrors.Unwrap()[1],
		).To(MatchError(fmt.Sprintf("open %s: permission denied", to.Name())))
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 1, Lognl: 1})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Applying %s to %s ...", from.Name(), strings.Replace(from.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✕"}))
	})

	It("should apply a directory", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		f1, _ := os.CreateTemp(dotfilesFilesDir, "1-*")
		defer f1.Close()
		f2, _ := os.CreateTemp(dotfilesFilesDir, "2-*")
		defer f2.Close()

		result, err := cmd.Apply(commands.ApplyArgs{
			From: dotfilesFilesDir,
			Extra: commands.ApplyArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 2, Lognl: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Applying %s to %s ...", f1.Name(), strings.Replace(f1.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Applying %s to %s ...", f2.Name(), strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✓"}))
	})

	It("should return an error if something appens while applying directory", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		f1, _ := os.CreateTemp(dotfilesFilesDir, "1-*")
		defer f1.Close()
		f2, _ := os.CreateTemp(dotfilesFilesDir, "2-*")
		defer f2.Close()

		core.RecreateFile(f2.Name(), strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1))
		os.Chmod(strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1), 0o400)

		result, err := cmd.Apply(commands.ApplyArgs{
			From: dotfilesFilesDir,
			Extra: commands.ApplyArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		unwrapErrors, _ := err.(interface{ Unwrap() []error })
		Expect(unwrapErrors.Unwrap()[0]).To(MatchError("error applying directory"))
		unwrapErrors, _ = unwrapErrors.Unwrap()[1].(interface{ Unwrap() []error })
		Expect(
			unwrapErrors.Unwrap()[0],
		).To(MatchError(fmt.Sprintf("error applying %s to %s", f2.Name(), strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1))))
		Expect(
			unwrapErrors.Unwrap()[1],
		).To(MatchError(fmt.Sprintf("open %s: permission denied", strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1))))
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 2, Lognl: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Applying %s to %s ...", f1.Name(), strings.Replace(f1.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Applying %s to %s ...", f2.Name(), strings.Replace(f2.Name(), dotfilesFilesDir, homedir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✕"}))
	})
})
