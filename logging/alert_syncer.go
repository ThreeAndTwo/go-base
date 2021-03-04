package logging

type AlertSyncer struct {
	channel AlertChannelI
}

func (w *AlertSyncer) setChannel(service string, channel AlertChannelI) {
	w.channel = channel
	w.channel.SetServiceName(service)
}

func (w *AlertSyncer) Write(p []byte) (n int, err error) {
	w.channel.SetMsg(string(p))
	go w.channel.Send()
	return len(p), nil
}
