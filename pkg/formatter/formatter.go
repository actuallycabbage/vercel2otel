package formatter

import (
	"strings"
	"time"
	"vercel2otel/pkg/vercel"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func VercelSeverityToOtel(in vercel.LogLevel) log.Severity {
	switch in {
	case vercel.StdErr:
		return log.SeverityError
	case vercel.StdOut:
		return log.SeverityUndefined
	default:
		return log.SeverityUndefined
	}
}

func TransformLogItems(logItems []vercel.LogItem) []sdklog.Record {
	var records []sdklog.Record

	for _, logItem := range logItems {
		record := sdklog.Record{}

		record.SetTimestamp(time.UnixMilli(logItem.TimestampMs))
		record.SetSeverity(VercelSeverityToOtel(logItem.Level))
		record.SetBody(log.StringValue(logItem.Message))
		record.AddAttributes(
			log.String("ID", logItem.ID),
			log.String("Message", logItem.Message),
			log.String("Type", logItem.Type),
			log.String("Source", string(logItem.Source)),
			log.String("ProjectID", logItem.ProjectID),
			log.String("DeploymentID", logItem.DeploymentID),
			log.String("BuildID", logItem.BuildID),
			log.String("Host", logItem.Host),
			log.String("Entrypoint", logItem.Entrypoint),
			log.String("RequestID", logItem.RequestID),
			log.Int("StatusCode", logItem.StatusCode),
			log.String("Destination", logItem.Destination),
			log.String("Path", logItem.Path),
			log.String("ExecutionRegion", logItem.ExecutionRegion),
			log.String("Level", string(logItem.Level)),
			log.Int64("ProxyTimestampMs", logItem.Proxy.TimestampMs),
			log.String("ProxyMethod", logItem.Proxy.Method),
			log.String("ProxyScheme", logItem.Proxy.Scheme),
			log.String("ProxyHost", logItem.Proxy.Host),
			log.String("ProxyPath", logItem.Proxy.Path),
			log.String("ProxyUserAgent", strings.Join(logItem.Proxy.UserAgent, ",")),
			log.String("ProxyReferer", logItem.Proxy.Referer),
			log.Int("ProxyStatusCode", logItem.Proxy.StatusCode),
			log.String("ProxyClientIP", logItem.Proxy.ClientIP),
			log.String("ProxyRegion", logItem.Proxy.Region),
			log.String("ProxyCacheID", logItem.Proxy.CacheID),
			log.String("ProxyVercelCache", logItem.Proxy.VercelCache),
		)

		records = append(records, record)
	}

	return records
}
