package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/goutil/fsutil"
)

func genEnvOrNil(key string) *string {
	env, exists := os.LookupEnv(key)

	if !exists {
		return nil
	}

	return &env
}

func resolveDotfilesFilesDirFromEnv(fallback *string) *string {
	fromEnv := genEnvOrNil("DOTS_DOTFILES_FILES_DIR")

	if fromEnv == nil {
		return fallback
	}

	return fromEnv
}

func RecreateFile(from string, to string) error {
	if err := os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return err
	}

	if err := fsutil.CopyFile(from, to); err != nil {
		return err
	}

	return nil
}

func IsPathReadable(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	if stat.Mode() == os.ModeIrregular {
		return true
	}

	file, err := os.Open(path)
	if err != nil {
		return false
	}

	if err = file.Close(); err != nil {
		return true
	}

	return true
}

func ResolveHomeDir() (string, error) {
	homeDirFromEnv, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	homedir, err := filepath.Abs(homeDirFromEnv)
	if err != nil {
		return "", err
	}

	if !fsutil.DirExist(homedir) || !IsPathReadable(homedir) {
		return "", fmt.Errorf(
			"homedir %s does not exists or is not a directory or is not readable",
			homedir,
		)
	}

	return homedir, nil
}

func ResolveDotfilesFilesDir(dotfilesFilesDirPath *string) (string, error) {
	dir := resolveDotfilesFilesDirFromEnv(dotfilesFilesDirPath)

	if dir == nil {
		return "", errors.New("dotfiles files directory path not provided")
	}

	dotfilesFilesDir, err := filepath.Abs(*dir)
	if err != nil {
		return "", err
	}

	if !fsutil.DirExist(dotfilesFilesDir) || !IsPathReadable(dotfilesFilesDir) {
		return "", fmt.Errorf(
			"dotfiles \"%s\" does not exists or is not a directory os is not readable",
			dotfilesFilesDir,
		)
	}

	return dotfilesFilesDir, nil
}

func LogErrors(logger ILogger, err error, n int) {
	wrapArr, ok := err.(interface{ Unwrap() []error })

	if ok {
		errArr := wrapArr.Unwrap()

		if len(errArr) <= 0 {
			return
		}

		if n == 0 {
			logger.Errornl(errArr[0].Error())
		} else {
			logger.Lognl(strings.Repeat(" ", n) + "-> " + errArr[0].Error())
		}

		for _, e := range errArr[1:] {
			LogErrors(logger, e, n+2)
		}

		return
	}

	if n == 0 {
		logger.Errornl(err.Error())
	} else {
		logger.Lognl(strings.Repeat(" ", n) + "-> " + err.Error())
	}
}
