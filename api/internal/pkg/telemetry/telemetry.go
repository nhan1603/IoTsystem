package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"
)

func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func NewRunID() string {
	return time.Now().UTC().Format("20060102T150405.000Z07:00")
}

func FingerprintEnv(keys ...string) string {
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+os.Getenv(k))
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:8]) // short
}

func LogRunStart(lg *slog.Logger, runID string, fields map[string]any) {
	attrs := []any{"run_id", runID}
	for k, v := range fields {
		attrs = append(attrs, k, v)
	}
	lg.Info("run_start", attrs...)
}

func LogRebalance(lg *slog.Logger, runID string, member string, claims map[string][]int32) {
	lg.Info("rebalance", "run_id", runID, "member_id", member, "claims", claims)
}

func LogBatchProcessed(lg *slog.Logger, runID string, partition int32, first, last int64, n int, latencyMs int64, sampled bool) {
	if sampled { // call with sampled=true when you want to emit
		lg.Info("batch_processed", "run_id", runID, "partition", partition,
			"offset_first", first, "offset_last", last, "batch_size", n, "latency_ms", latencyMs)
	}
}

func LogOffsetCommit(lg *slog.Logger, runID string, partition int32, offset int64, durMs int64, groupLag int64) {
	lg.Info("offset_commit", "run_id", runID, "partition", partition, "offset", offset, "duration_ms", durMs, "group_lag", groupLag)
}

func LogAnomaly(lg *slog.Logger, runID, kind string, value any, threshold any, ctx map[string]any) {
	attrs := []any{"run_id", runID, "type", kind, "value", value, "threshold", threshold}
	for k, v := range ctx {
		attrs = append(attrs, k, v)
	}
	lg.Warn("anomaly", attrs...)
}

func LogShutdown(lg *slog.Logger, runID string, lastOffsets map[string]int64, ok bool) {
	lg.Info("shutdown", "run_id", runID, "last_offsets", lastOffsets, "metrics_flushed", ok)
}

func LogProducerSummary(lg *slog.Logger, runID string, sent, errors, retries int, avgThroughput float64) {
	lg.Info("producer_summary", "run_id", runID, "messages_sent", sent, "send_errors", errors, "retries", retries, "throughput_avg_rps", avgThroughput)
}
