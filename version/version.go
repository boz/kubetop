package version

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func Date() string {
	return date
}
