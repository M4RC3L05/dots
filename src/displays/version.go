package displays

func (d Displays) Version(version string) {
	d.Logger.Lognl("%s", version)
}
