package tgclient

const (
	// Number of message workers which process the messages in parallel.
	MAX_MESSAGE_WORKERS = 4
	// Maximum number of updates (messages) holded in update channel buffer and consumed by message workers.
	UPDATE_CHANNEL_BUFFER = 10
)
