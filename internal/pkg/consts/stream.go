package consts

import (
	"time"
)

const (
	MonitorTypeStop        = 0
	MonitorTypeStart       = 1
	MonitorTypeCron        = 2
	MonitorTypeIntelligent = 3

	DurationExpired = 24 * time.Hour

	DefaultInterval = 30
	MaxInterval     = 900

	RequestFailedThreshold = 10

	RequestsFailedWarningCKP = "live:requests:failed:warning:"
	RequestsFailedWarningTtl = 1 * time.Hour

	MaxTopCount = 3
)
