package flag

import "flag"

var (
	SafeRoot string
	LogPath  string
)

func init() {
	flag.StringVar(&SafeRoot, "safe-root", "./safeFiles", "safe-file-root")
	flag.StringVar(&LogPath, "log-path", "log.txt", "file to write logs")
}
