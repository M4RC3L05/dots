package commands_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
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
		result, err := cmd.Adopt(commands.AdoptArgs{
			From: workingDir + "/foo/bar/..",
			Extra: commands.AdoptArgsExtra{
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

	It("should return an error if `from` does not exists", func() {
		result, err := cmd.Adopt(commands.AdoptArgs{
			From: workingDir + "/foo/bar/..",
			Extra: commands.AdoptArgsExtra{
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

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from,
			Extra: commands.AdoptArgsExtra{
				Homedir:          workingDir,
				DotfilesFilesDir: workingDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(err).To(MatchError(fmt.Sprintf("path %s does not exists or is not readable", from)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `from` is not a subdirectory of `homedir`", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.MkdirTemp(workingDir, "*")

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from,
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(err).To(MatchError(fmt.Sprintf("path %s is not a subpath of %s", from, homedir)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `from` is a subdirectory of `dotfilesFilesDir`", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(homedir, "*")
		from, _ := os.MkdirTemp(dotfilesFilesDir, "*")

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from,
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s can not be a subpath of %s", from, dotfilesFilesDir)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should adopt a single file", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(homedir, "*")
		defer from.Close()

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from.Name(),
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 1, Lognl: 1})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Adopting %s to %s ...", from.Name(), strings.Replace(from.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
	})

	It("should return an error if something appens while adopting file", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(homedir, "*")
		defer from.Close()
		to := strings.Replace(from.Name(), homedir, dotfilesFilesDir, 1)
		toF, _ := os.Create(to)
		defer toF.Close()

		os.Chmod(to, 0o400)

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from.Name(),
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		unwrapErrors, _ := err.(interface{ Unwrap() []error })
		Expect(
			unwrapErrors.Unwrap()[0],
		).To(MatchError(fmt.Sprintf("error adopting %s to %s", from.Name(), to)))
		Expect(
			unwrapErrors.Unwrap()[1],
		).To(MatchError(fmt.Sprintf("open %s: permission denied", to)))
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 1, Lognl: 1})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Adopting %s to %s ...", from.Name(), strings.Replace(from.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✕"}))
	})

	It("should adopt a directory when `from` is the same as `dotfilesFilesDir`", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		f1, _ := os.CreateTemp(homedir, "1-*")
		defer f1.Close()
		f2, _ := os.CreateTemp(homedir, "2-*")
		defer f2.Close()

		fsutil.CopyFile(f1.Name(), strings.Replace(f1.Name(), homedir, dotfilesFilesDir, 1))
		fsutil.CopyFile(f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1))

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: dotfilesFilesDir,
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Lognl: 2, Log: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f1.Name(), strings.Replace(f1.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✓"}))
	})

	It(
		"should not adopt a directory when `from` is the same as `dotfilesFilesDir` and no files exists on `dotfilesFilesDir`",
		func() {
			homedir, _ := os.MkdirTemp(workingDir, "*")
			dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")

			result, err := cmd.Adopt(commands.AdoptArgs{
				From: dotfilesFilesDir,
				Extra: commands.AdoptArgsExtra{
					Homedir:          homedir,
					DotfilesFilesDir: dotfilesFilesDir,
				},
			})

			Expect(result).To(BeTrue())
			Expect(err).To(BeNil())
			testing.AssertSpyLoggerCalls(*logger, nil)
		},
	)

	It("should adopt a directory", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.MkdirTemp(homedir, "*")
		f1, _ := os.CreateTemp(from, "1-*")
		defer f1.Close()
		f2, _ := os.CreateTemp(from, "2-*")
		defer f2.Close()

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from,
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 2, Lognl: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f1.Name(), strings.Replace(f1.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✓"}))
	})

	It("should return an error if something appens while adopting directory", func() {
		homedir, _ := os.MkdirTemp(workingDir, "*")
		dotfilesFilesDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.MkdirTemp(homedir, "*")
		f1, _ := os.CreateTemp(from, "1-*")
		defer f1.Close()
		f2, _ := os.CreateTemp(from, "2-*")
		defer f2.Close()

		core.RecreateFile(f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1))
		os.Chmod(strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1), 0o400)

		result, err := cmd.Adopt(commands.AdoptArgs{
			From: from,
			Extra: commands.AdoptArgsExtra{
				Homedir:          homedir,
				DotfilesFilesDir: dotfilesFilesDir,
			},
		})

		Expect(result).To(BeFalse())
		unr1err, _ := err.(interface {
			Unwrap() []error
		})
		errorsArr := unr1err.Unwrap()

		Expect(errorsArr[0]).To(MatchError("error adopting directory"))
		unwrapErrors, _ := errorsArr[1].(interface{ Unwrap() []error })
		Expect(
			unwrapErrors.Unwrap()[0],
		).To(MatchError(fmt.Sprintf("error adopting %s to %s", f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1))))
		Expect(
			unwrapErrors.Unwrap()[1],
		).To(MatchError(fmt.Sprintf("open %s: permission denied", strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1))))
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Log: 2, Lognl: 2})
		Expect(
			logger.Calls.Log[0].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f1.Name(), strings.Replace(f1.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(
			logger.Calls.Log[1].Args,
		).To(Equal([]any{"Adopting %s to %s ...", f2.Name(), strings.Replace(f2.Name(), homedir, dotfilesFilesDir, 1)}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{" ✓"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{" ✕"}))
	})
})
