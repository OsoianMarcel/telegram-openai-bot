package stats

import "sync/atomic"

type Stats struct {
	aiAllMessages     uint32
	aiInvalidMessages uint32
	aiErrors          uint32
	aiRequestChars    uint32
	aiResponseChars   uint32
}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) GetAiAllMessages() uint32 {
	return atomic.LoadUint32(&s.aiAllMessages)
}

func (s *Stats) IncAiAllMessages() {
	atomic.AddUint32(&s.aiAllMessages, 1)
}

func (s *Stats) GetAiInvalidMessages() uint32 {
	return atomic.LoadUint32(&s.aiInvalidMessages)
}

func (s *Stats) IncAiInvalidMessages() {
	atomic.AddUint32(&s.aiInvalidMessages, 1)
}

func (s *Stats) GetAiErrors() uint32 {
	return atomic.LoadUint32(&s.aiErrors)
}

func (s *Stats) IncAiErrors() {
	atomic.AddUint32(&s.aiErrors, 1)
}

func (s *Stats) GetAiRequestChars() uint32 {
	return atomic.LoadUint32(&s.aiRequestChars)
}

func (s *Stats) AddAiRequestChars(n uint32) {
	atomic.AddUint32(&s.aiRequestChars, n)
}

func (s *Stats) GetAiResponseChars() uint32 {
	return atomic.LoadUint32(&s.aiResponseChars)
}

func (s *Stats) AddAiResponseChars(n uint32) {
	atomic.AddUint32(&s.aiResponseChars, n)
}
