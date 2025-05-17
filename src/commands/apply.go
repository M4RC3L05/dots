package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/core"
)

type ApplyArgsExtra struct {
	Homedir          string
	DotfilesFilesDir string
}

type ApplyArgs struct {
	From  string
	Extra ApplyArgsExtra
}

func (c Commands) Apply(args ApplyArgs) (bool, error) {
	fromFormatted, err := filepath.Abs(args.From)
	if err != nil {
		return false, err
	}

	args.From = fromFormatted

	if !fsutil.PathExist(args.From) || !core.IsPathReadable(args.From) {
		return false, fmt.Errorf(
			"path %s does not exists or is not readable",
			color.BlueString(args.From),
		)
	}

	if !strings.HasPrefix(args.From, args.Extra.DotfilesFilesDir) {
		return false, fmt.Errorf(
			"path %s is not a subpath of %s",
			color.BlueString(args.From), color.BlueString(args.Extra.DotfilesFilesDir),
		)
	}

	if fsutil.IsFile(args.From) {
		to, err := filepath.Abs(
			strings.Replace(args.From, args.Extra.DotfilesFilesDir, args.Extra.Homedir, 1),
		)
		if err != nil {
			return false, err
		}

		c.Logger.Log(
			"Applying %s to %s ...",
			color.BlueString(args.From),
			color.BlueString(to),
		)

		if err := core.RecreateFile(args.From, to); err != nil {
			c.Logger.Lognl(color.GreenString(" ✕"))

			return false, errors.Join(fmt.Errorf(
				"error applying %s to %s",
				color.BlueString(args.From),
				color.BlueString(to),
			), err)
		}

		c.Logger.Lognl(color.GreenString(" ✓"))

		return true, nil
	}

	var errorsArr []error

	err = filepath.WalkDir(args.From, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() {
			return nil
		}

		to := strings.Replace(path, args.Extra.DotfilesFilesDir, args.Extra.Homedir, 1)

		c.Logger.Log("Applying %s to %s ...", color.BlueString(path), color.BlueString(to))

		if err := core.RecreateFile(path, to); err != nil {
			c.Logger.Lognl(color.RedString(" ✕"))

			errorsArr = append(
				errorsArr,
				errors.Join(fmt.Errorf(
					"error applying %s to %s",
					color.BlueString(path),
					color.BlueString(to),
				), err),
			)
		} else {
			c.Logger.Lognl(color.GreenString(" ✓"))
		}

		return nil
	})
	if err != nil {
		errorsArr = append(errorsArr, err)
	}

	if len(errorsArr) > 0 {
		errorsArr = append([]error{errors.New("error applying directory")}, errorsArr...)
	}

	return len(errorsArr) <= 0, errors.Join(errorsArr...)
}
