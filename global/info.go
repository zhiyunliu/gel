package global

var (
	GitCommit   = "unknown"
	BuildTime   = "unknown"
	Version     = "unknown"
	DisplayName = ""
	Usage       = "unknown"
)

var (
	running bool
)

func IsRunning() bool {
	return running
}

func StartRunning() {
	running = true
}
