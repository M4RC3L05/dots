package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/gookit/goutil/fsutil"
)

func genEnvOrNil(key string) *string {
	env, exists := os.LookupEnv(key)

	if !exists {
		return nil
	}

	return &env
}

func resolveHomeDirFromEnv() *string {
	if runtime.GOOS == "windows" {
		return genEnvOrNil("USERPROFILE")
	} else {
		return genEnvOrNil("HOME")
	}
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

	statSys, ok := stat.Sys().(*syscall.Stat_t)

	if !ok {
		return true
	}

	if os.Getuid() == int(statSys.Uid) {
		return (stat.Mode() & 0o400) == 0o400
	} else if os.Getgid() == int(statSys.Gid) {
		return (stat.Mode() & 0o040) == 0o040
	}

	return (stat.Mode() & 0o004) == 0o004
}

func ResolveHomeDir() (string, error) {
	homeDirFromEnv := resolveHomeDirFromEnv()

	if homeDirFromEnv == nil {
		return "", errors.New("could not determine users home directory")
	}

	homedir, err := filepath.Abs(*homeDirFromEnv)
	if err != nil {
		return "", err
	}

	if !fsutil.DirExist(homedir) || !IsPathReadable(homedir) {
		return "", fmt.Errorf(
			"homedir \"%s\" does not exists or is not a directory or is not readable",
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
