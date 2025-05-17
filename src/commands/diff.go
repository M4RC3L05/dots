package commands

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	"github.com/fatih/color"
	"github.com/gookit/goutil/fsutil"
	"github.com/m4rc3l05/dots/src/core"
)

type DiffCmdArgsExtra struct {
	Logger core.ILogger
}

type DiffCmdArgs struct {
	FromDir string
	ToDir   string
	Extra   DiffCmdArgsExtra
}

func DiffCmd(args DiffCmdArgs) (bool, error) {
	hasFilesWithChanges := false

	err := filepath.WalkDir(args.FromDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.Type().IsRegular() {
			return nil
		}

		from := path
		to := strings.Replace(from, args.FromDir, args.ToDir, 1)

		if !fsutil.FileExist(from) || !core.IsPathReadable(from) {
			args.Extra.Logger.Warnnl(
				"File %s does not exists or is not a file or is not readable, skipping...",
				color.BlueString(from),
			)

			return nil
		}

		if !fsutil.FileExist(to) || !core.IsPathReadable(to) {
			args.Extra.Logger.Warnnl(
				"File %s does not exists or is not a file or is not readable, skipping...",
				color.BlueString(to),
			)

			return nil
		}

		args.Extra.Logger.Log(
			"Diffing %s against %s ...",
			color.BlueString(from),
			color.BlueString(to),
		)

		diffs := udiff.Unified(from, to, string(fsutil.ReadFile(from)), string(fsutil.ReadFile(to)))

		if len(diffs) > 0 {
			hasFilesWithChanges = true

			args.Extra.Logger.Lognl(color.RedString(" ✕"))

			for line := range strings.SplitSeq(strings.TrimSpace(diffs), "\n") {
				if strings.HasPrefix(line, "@@") && strings.HasSuffix(line, "@@") {
					args.Extra.Logger.Lognl(color.CyanString(line))
				} else if len(line) > 1 && line[0] == '+' && line[1] != '+' {
					args.Extra.Logger.Lognl(color.GreenString(line))
				} else if len(line) > 1 && line[0] == '-' && line[1] != '-' {
					args.Extra.Logger.Lognl(color.RedString(line))
				} else {
					args.Extra.Logger.Lognl(line)
				}
			}
		} else {
			args.Extra.Logger.Lognl(color.GreenString(" ✓"))
		}

		return nil
	})

	return !hasFilesWithChanges, err
}
