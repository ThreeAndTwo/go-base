package logging

type AlertSyncer struct {
	channel AlertChannelI
}

func (w *AlertSyncer) setChannel(service string, channel AlertChannelI) {
	w.channel = channel
	w.channel.SetService(service)
}

func (w *AlertSyncer) Write(p []byte) (n int, err error) {
	go w.channel.Send(string(p))
	return len(p), nil
}
