package commands_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/commands"
	"github.com/m4rc3l05/dots/src/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("diff()", func() {
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

	It("should normalize `fromDir`", func() {
		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: "/foo/bar/..",
			ToDir:   "/baz/biz",
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError("path /foo does not exists or is not a directory or is not readable"))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should normalize `toDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   "/baz/biz/..",
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError("path /baz does not exists or is not a directory or is not readable"))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `fromDir` is not a directory", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		fromFile, _ := os.CreateTemp(fromDir, "foo")

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromFile.Name(),
			ToDir:   toDir,
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not a directory or is not readable", fromFile.Name())))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `fromDir` is not readable", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")

		os.Chmod(fromDir, 0o300)

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not a directory or is not readable", fromDir)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `toDir` is not a directory", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		toFile, _ := os.CreateTemp(toDir, "foo")

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toFile.Name(),
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not a directory or is not readable", toFile.Name())))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should return an error if `toDir` is not readable", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")

		os.Chmod(toDir, 0o300)

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeFalse())
		Expect(
			err,
		).To(MatchError(fmt.Sprintf("path %s does not exists or is not a directory or is not readable", toDir)))
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should not diff anything if no files exists on `fromDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, nil)
	})

	It("should not diff anything if no readable file exists on `fromDir`", func() {
		fromDir, _ := os.MkdirTemp(workingDir, "*")
		toDir, _ := os.MkdirTemp(workingDir, "*")
		from, _ := os.CreateTemp(fromDir, "*")

		os.Chmod(from.Name(), 0o300)

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Warnnl: 1})
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

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Warnnl: 1})
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

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Warnnl: 1})
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

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeTrue())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Lognl: 1, Log: 1})
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

		result, err := cmd.Diff(commands.DiffArgs{
			FromDir: fromDir,
			ToDir:   toDir,
		})

		Expect(result).To(BeFalse())
		Expect(err).To(BeNil())
		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Lognl: 9, Log: 2})
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
