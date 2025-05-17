package commands

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/core"
)

type DiffArgs struct {
	FromDir string
	ToDir   string
}

func (c Commands) Diff(args DiffArgs) (bool, error) {
	hasFilesWithChanges := false

	fromFormatted, err := filepath.Abs(args.FromDir)
	if err != nil {
		return false, err
	}

	args.FromDir = fromFormatted

	toFormatted, err := filepath.Abs(args.ToDir)
	if err != nil {
		return false, err
	}

	args.ToDir = toFormatted

	if !fsutil.IsDir(args.FromDir) || !core.IsPathReadable(args.FromDir) {
		return false, fmt.Errorf(
			"path %s does not exists or is not a directory or is not readable",
			args.FromDir,
		)
	}

	if !fsutil.IsDir(args.ToDir) || !core.IsPathReadable(args.ToDir) {
		return false, fmt.Errorf(
			"path %s does not exists or is not a directory or is not readable",
			args.ToDir,
		)
	}

	err = filepath.WalkDir(args.FromDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() {
			return nil
		}

		from := path
		to := strings.Replace(from, args.FromDir, args.ToDir, 1)

		if !fsutil.FileExist(from) || !core.IsPathReadable(from) {
			c.Logger.Warnnl(
				"File %s does not exists or is not a file or is not readable, skipping...",
				color.BlueString(from),
			)

			return nil
		}

		if !fsutil.FileExist(to) || !core.IsPathReadable(to) {
			c.Logger.Warnnl(
				"File %s does not exists or is not a file or is not readable, skipping...",
				color.BlueString(to),
			)

			return nil
		}

		c.Logger.Log(
			"Diffing %s against %s ...",
			color.BlueString(from),
			color.BlueString(to),
		)

		diffs := udiff.Unified(from, to, string(fsutil.ReadFile(from)), string(fsutil.ReadFile(to)))

		if len(diffs) > 0 {
			hasFilesWithChanges = true

			c.Logger.Lognl(color.RedString(" ✕"))

			for line := range strings.SplitSeq(strings.TrimSpace(diffs), "\n") {
				if strings.HasPrefix(line, "@@") && strings.HasSuffix(line, "@@") {
					c.Logger.Lognl(color.CyanString(line))
				} else if len(line) > 1 && line[0] == '+' && line[1] != '+' {
					c.Logger.Lognl(color.GreenString(line))
				} else if len(line) > 1 && line[0] == '-' && line[1] != '-' {
					c.Logger.Lognl(color.RedString(line))
				} else {
					c.Logger.Lognl(line)
				}
			}
		} else {
			c.Logger.Lognl(color.GreenString(" ✓"))
		}

		return nil
	})

	return err == nil && !hasFilesWithChanges, err
}
