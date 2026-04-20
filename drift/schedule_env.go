package drift

import (
	"os"
	"strconv"
)

// ScheduleConfigFromEnv builds a ScheduleConfig from environment variables.
// Variables:
//
//	DRIFTCTL_SCHEDULE_ENABLED      = "true" | "false"
//	DRIFTCTL_SCHEDULE_INTERVAL     = minutes (int)
//	DRIFTCTL_SCHEDULE_STATE_FILE   = path to state file
//	DRIFTCTL_SCHEDULE_OUTPUT_FORMAT = text | json | csv
func ScheduleConfigFromEnv(base ScheduleConfig) ScheduleConfig {
	if v := os.Getenv("DRIFTCTL_SCHEDULE_ENABLED"); v != "" {
		base.Enabled = v == "true" || v == "1"
	}
	if v := os.Getenv("DRIFTCTL_SCHEDULE_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			base.IntervalMins = n
		}
	}
	if v := os.Getenv("DRIFTCTL_SCHEDULE_STATE_FILE"); v != "" {
		base.StateFile = v
	}
	if v := os.Getenv("DRIFTCTL_SCHEDULE_OUTPUT_FORMAT"); v != "" {
		base.OutputFormat = v
	}
	return base
}
