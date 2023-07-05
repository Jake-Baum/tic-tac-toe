package lambda

import log "github.com/sirupsen/logrus"

func Initialize() {
	log.SetReportCaller(true)
}
