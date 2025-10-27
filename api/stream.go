package api

import (
	"context"
	"sync"

	"github.com/RasmusLindroth/go-mastodon"
)

type MastodonType uint

const (
	StatusType MastodonType = iota
	StatusHistoryType
	UserType
	ProfileType
	NotificationType
	ListsType
	TagType
)

type StreamType uint

const (
	HomeStream StreamType = iota
	LocalStream
	FederatedStream
	DirectStream
	TagStream
	ListStream
)

type StreamID struct {
	Type   StreamType
	Data string
}

func MakeStreamID(st StreamType, data string) StreamID {
	return StreamID{st, data}
}

type Receiver func(mastodon.Event) // always use *Receiver, because == comparison need it

type Stream struct {
	id        StreamID
	receivers []*Receiver
	incoming  chan mastodon.Event
	cancel    context.CancelFunc
	mux       sync.Mutex
}

func (s *Stream) ID() StreamID {
	return s.id
}

func (s *Stream) AddReceiver(r *Receiver) {
	s.receivers = append(s.receivers, r)
}

func (s *Stream) RemoveReceiver(r *Receiver) {
	index := -1
	for i, rec := range s.receivers {
		if rec == r {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	s.receivers = append(s.receivers[:index], s.receivers[index+1:]...)
}

func (s *Stream) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-s.incoming:
			switch e.(type) {
			case *mastodon.UpdateEvent, *mastodon.ConversationEvent, *mastodon.NotificationEvent, *mastodon.DeleteEvent, *mastodon.ErrorEvent:
				for _, r := range s.receivers {
					(*r)(e)
				}
			}
		}
	}
}

func newStream(ctx context.Context, id StreamID, input chan mastodon.Event) *Stream {
	ctx, cancel := context.WithCancel(ctx)
	stream := &Stream{
		id:       id,
		incoming: input,
		cancel:   cancel,
	}
	go stream.listen(ctx)
	return stream
}

func (ac *AccountClient) CreateOrGetStream(ctx context.Context, st StreamType, data string) (stream *Stream, err error) {
	id := MakeStreamID(st, data)

	// get stream
	for _, s := range ac.Streams {
		if s.ID() == id {
			return s, nil
		}
	}

	// create stream
	ch, err := ac.StreamIDtoWebSocketStream(ctx, id)
	if err != nil {
		return nil, err
	}
	stream = newStream(ctx, id, ch)
	ac.Streams[stream.ID()] = stream
	return stream, nil
}

// create mastodon.Event stream channel
func (ac *AccountClient) StreamIDtoWebSocketStream(ctx context.Context, id StreamID) (chan mastodon.Event, error) {
	switch id.Type {
	case HomeStream:
		return ac.WSClient.StreamingWSUser(ctx)
	case LocalStream:
		return ac.WSClient.StreamingWSPublic(ctx, true)
	case FederatedStream:
		return ac.WSClient.StreamingWSPublic(ctx, false)
	case DirectStream:
		return ac.WSClient.StreamingWSDirect(ctx)
	case TagStream:
		return ac.WSClient.StreamingWSHashtag(ctx, id.Data, false)
	case ListStream:
		return ac.WSClient.StreamingWSList(ctx, mastodon.ID(id.Data))
	default:
		panic("invalid StreamType")
	}
}

func (ac *AccountClient) RemoveReceiver(rec *Receiver, st StreamType, data string) {
	id := MakeStreamID(st, data)
	stream, ok := ac.Streams[id]
	if !ok {
		return
	}
	stream.mux.Lock()
	stream.RemoveReceiver(rec)
	if len(stream.receivers) == 0 {
		stream.cancel()
		delete(ac.Streams, id)
	}
	stream.mux.Unlock()
}
