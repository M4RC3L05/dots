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

type AdoptArgsExtra struct {
	Homedir          string
	DotfilesFilesDir string
}

type AdoptArgs struct {
	From  string
	Extra AdoptArgsExtra
}

func (c Commands) Adopt(args AdoptArgs) (bool, error) {
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

	if args.From != args.Extra.DotfilesFilesDir {
		if !strings.HasPrefix(args.From, args.Extra.Homedir) {
			return false, fmt.Errorf(
				"path %s is not a subpath of %s",
				color.BlueString(args.From), color.BlueString(args.Extra.Homedir),
			)
		}

		if strings.HasPrefix(args.From, args.Extra.DotfilesFilesDir) {
			return false, fmt.Errorf(
				"path %s can not be a subpath of %s",
				color.BlueString(args.From), color.BlueString(args.Extra.DotfilesFilesDir),
			)
		}
	}

	if fsutil.IsFile(args.From) {
		to := strings.Replace(args.From, args.Extra.Homedir, args.Extra.DotfilesFilesDir, 1)

		c.Logger.Log(
			"Adopting %s to %s ...",
			color.BlueString(args.From),
			color.BlueString(to),
		)

		if err := core.RecreateFile(args.From, to); err != nil {
			c.Logger.Lognl(color.GreenString(" ✕"))

			return false, errors.Join(fmt.Errorf(
				"error adopting %s to %s",
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

		var origin string
		var destination string

		if strings.HasPrefix(path, args.Extra.DotfilesFilesDir) {
			origin = strings.Replace(path, args.Extra.DotfilesFilesDir, args.Extra.Homedir, 1)
			destination = path
		} else {
			origin = path
			destination = strings.Replace(path, args.Extra.Homedir, args.Extra.DotfilesFilesDir, 1)
		}

		c.Logger.Log(
			"Adopting %s to %s ...",
			color.BlueString(origin),
			color.BlueString(destination),
		)

		if err := core.RecreateFile(origin, destination); err != nil {
			c.Logger.Lognl(color.RedString(" ✕"))

			errorsArr = append(
				errorsArr,
				errors.Join(fmt.Errorf(
					"error adopting %s to %s",
					color.BlueString(origin),
					color.BlueString(destination),
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
		errorsArr = append([]error{errors.New("error adopting directory")}, errorsArr...)
	}

	return len(errorsArr) <= 0, errors.Join(errorsArr...)
}
