package webcache

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var Log *log.Entry

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		PadLevelText:  true,
		FullTimestamp: true,
	})

	log.SetReportCaller(true)
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	Log = log.WithField("app", "webcache")
}
