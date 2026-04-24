package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var (
	forecastPeriods    int
	forecastExportFmt  string
	forecastExportPath string
)

func init() {
	forecastCmd := &cobra.Command{
		Use:   "forecast",
		Short: "Predict future drift scores from trend history",
		RunE:  runForecast,
	}
	forecastCmd.Flags().StringVar(&trendFile, "trend-file", "drift-trend.json", "Path to trend history file")
	forecastCmd.Flags().IntVar(&forecastPeriods, "periods", 3, "Number of future periods to forecast")
	forecastCmd.Flags().StringVar(&forecastExportFmt, "export-format", "", "Export format: json or csv")
	forecastCmd.Flags().StringVar(&forecastExportPath, "export-path", "", "Path for exported forecast file")
	rootCmd.AddCommand(forecastCmd)
}

func runForecast(cmd *cobra.Command, args []string) error {
	history, err := drift.LoadTrend(trendFile)
	if err != nil {
		return fmt.Errorf("loading trend: %w", err)
	}

	result := drift.Forecast(history, forecastPeriods)
	drift.FprintForecast(os.Stdout, result)

	if forecastExportFmt != "" && forecastExportPath != "" {
		if err := drift.ExportForecast(result, forecastExportFmt, forecastExportPath); err != nil {
			return fmt.Errorf("exporting forecast: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Forecast exported to %s\n", forecastExportPath)
	}

	if drift.ForecastHasDrift(result) {
		os.Exit(1)
	}
	return nil
}
