package drift

import (
	"os"
	"strconv"
)

// RetentionConfigFromEnv overrides retention config fields from environment variables.
//
// Supported variables:
//
//	DRIFTCTL_RETENTION_HISTORY_DAYS
//	DRIFTCTL_RETENTION_AUDIT_DAYS
//	DRIFTCTL_RETENTION_SNAPSHOT_DAYS
//	DRIFTCTL_RETENTION_TREND_DAYS
func RetentionConfigFromEnv(cfg RetentionConfig) RetentionConfig {
	if v := os.Getenv("DRIFTCTL_RETENTION_HISTORY_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.HistoryDays = n
		}
	}
	if v := os.Getenv("DRIFTCTL_RETENTION_AUDIT_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.AuditDays = n
		}
	}
	if v := os.Getenv("DRIFTCTL_RETENTION_SNAPSHOT_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.SnapshotDays = n
		}
	}
	if v := os.Getenv("DRIFTCTL_RETENTION_TREND_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.TrendDays = n
		}
	}
	return cfg
}
