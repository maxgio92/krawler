package amazonlinux

const (
	MirrorsDistroVersionRegex = `^(0|[v1-9]\d*)(\.(0|[v1-9]\d*)?)?(\.(0|[v1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
)

var (
	debugScrape = true
)
