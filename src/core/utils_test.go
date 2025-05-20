package core_test

import (
	"errors"
	"os"

	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/core"
	"github.com/m4rc3l05/dots/src/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	workingDir string
	logger     *testing.SpyLogger
)

var _ = BeforeEach(func() {
	workingDir, _ = os.MkdirTemp(os.TempDir(), "wd-*")
	logger = testing.MakeSpyLogger()
})

var _ = AfterEach(func() {
	os.RemoveAll(workingDir)
})

var _ = Describe("RecreateFile()", func() {
	It("should recreate a file from one place to another", func() {
		os.MkdirAll(workingDir+"/1/2", os.ModePerm)
		f, _ := os.Create(workingDir + "/1/2/3")
		f.WriteString("foo")

		defer f.Close()

		err := core.RecreateFile(workingDir+"/1/2/3", workingDir+"/4/5/6")

		Expect(err).To(BeNil())
		Expect(string(fsutil.ReadAll(workingDir + "/4/5/6"))).To(Equal("foo"))
	})

	It(
		"should recreate a file from one place to another with the fil existing in destination",
		func() {
			os.MkdirAll(workingDir+"/1/2", os.ModePerm)
			os.MkdirAll(workingDir+"/4/5", os.ModePerm)
			f, _ := os.Create(workingDir + "/1/2/3")
			f.WriteString("foo")
			f2, _ := os.Create(workingDir + "/4/5/6")
			f2.WriteString("bar")

			defer f.Close()
			defer f2.Close()

			err := core.RecreateFile(workingDir+"/1/2/3", workingDir+"/4/5/6")

			Expect(err).To(BeNil())
			Expect(string(fsutil.ReadAll(workingDir + "/4/5/6"))).To(Equal("foo"))
		},
	)
})

var _ = Describe("IsPathReadable()", func() {
	It("should return false if path is not readable", func() {
		f, _ := os.CreateTemp(workingDir, "*")
		defer f.Close()

		os.Chmod(f.Name(), 0o200)

		ok := core.IsPathReadable(f.Name())
		Expect(ok).To(BeFalse())
	})

	It("should return true if path is readable", func() {
		f, _ := os.CreateTemp(workingDir, "*")
		defer f.Close()

		ok := core.IsPathReadable(f.Name())
		Expect(ok).To(BeTrue())
	})
})

var _ = Describe("LogErrors()", func() {
	It("should log correctly a simple error", func() {
		core.LogErrors(logger, errors.New("foo"), 0)

		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Errornl: 1})
		Expect(logger.Calls.Errornl[0].Args).To(Equal([]any{"foo"}))
	})

	It("should log correctly a joined error", func() {
		core.LogErrors(
			logger,
			errors.Join(
				errors.New("foo"),
				errors.New("bar"),
				errors.Join(errors.New("biz"), errors.New("buz")),
			),
			0,
		)

		testing.AssertSpyLoggerCalls(*logger, &testing.SpyLoggerCallNumber{Errornl: 1, Lognl: 3})
		Expect(logger.Calls.Errornl[0].Args).To(Equal([]any{"foo"}))
		Expect(logger.Calls.Lognl[0].Args).To(Equal([]any{"  -> bar"}))
		Expect(logger.Calls.Lognl[1].Args).To(Equal([]any{"  -> biz"}))
		Expect(logger.Calls.Lognl[2].Args).To(Equal([]any{"    -> buz"}))
	})
})
