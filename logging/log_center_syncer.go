package logging

type LogCenterSyncer struct {
}

func (w *LogCenterSyncer) Write(p []byte) (n int, err error) {
	// TODO: 将日志写到阿里云日志中心
	return len(p), nil
}
