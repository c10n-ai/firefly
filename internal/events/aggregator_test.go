// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"context"
	"fmt"
	"testing"

	"github.com/kaleido-io/firefly/mocks/databasemocks"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShutdownOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mdi := &databasemocks.Plugin{}
	mdi.On("GetOffset", mock.Anything, fftypes.OffsetTypeAggregator, fftypes.SystemNamespace, aggregatorOffsetName).Return(&fftypes.Offset{
		Type:      fftypes.OffsetTypeAggregator,
		Namespace: fftypes.SystemNamespace,
		Name:      aggregatorOffsetName,
		Current:   12345,
	}, nil)
	mdi.On("GetEvents", mock.Anything, mock.Anything, mock.Anything).Return([]*fftypes.Event{}, nil)
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	err := ag.start()
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), ag.eventPoller.pollingOffset)
	ag.eventPoller.eventNotifier.newEvents <- 12345
	cancel()
	<-ag.eventPoller.closed
}

func TestProcessEventsNoopIncrement(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mdi := &databasemocks.Plugin{}
	var runAsGroupFn func(context.Context) error
	mdi.On("UpdateOffset", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	mdi.On("RunAsGroup", mock.Anything, mock.MatchedBy(
		func(fn func(context.Context) error) bool {
			runAsGroupFn = fn
			return true
		},
	)).Return(nil, nil)
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))

	ev1 := fftypes.NewEvent(fftypes.EventTypeMessageConfirmed, "ns1", fftypes.NewUUID())
	ev1.Sequence = 111
	ev2 := fftypes.NewEvent(fftypes.EventTypeMessageConfirmed, "ns1", fftypes.NewUUID())
	ev2.Sequence = 112
	ev3 := fftypes.NewEvent(fftypes.EventTypeMessageConfirmed, "ns1", fftypes.NewUUID())
	ev3.Sequence = 113
	_, err := ag.processEventRetryAndGroup([]*fftypes.Event{
		ev1, ev2, ev3,
	})
	runAsGroupFn(context.Background())
	assert.NoError(t, err)
	mdi.AssertExpectations(t)
}

func TestProcessEventCheckSequencedReadFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("GetMessageByID", mock.Anything, mock.Anything).Return(msg, fmt.Errorf("pop"))
	ev1 := fftypes.NewEvent(fftypes.EventTypeMessageSequencedBroadcast, "ns1", fftypes.NewUUID())
	ev1.Sequence = 111
	_, err := ag.processEvents(context.Background(), []*fftypes.Event{ev1})
	assert.EqualError(t, err, "pop")
	assert.Equal(t, int64(0), ag.eventPoller.pollingOffset)
	mdi.AssertExpectations(t)
}

func TestProcessEventIgnoredTypeConfirmed(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	ev1 := fftypes.NewEvent(fftypes.EventTypeMessageConfirmed, "ns1", fftypes.NewUUID())
	ev1.Sequence = 111
	repoll, err := ag.processEvent(context.Background(), ev1)
	assert.False(t, repoll)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), ag.eventPoller.pollingOffset)
	mdi.AssertExpectations(t)
}

func TestProcessEventCheckCompleteDataNotAvailable(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("GetMessageByID", mock.Anything, mock.Anything).Return(msg, nil)
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(false, nil)
	ev1 := fftypes.NewEvent(fftypes.EventTypeMessageSequencedBroadcast, "ns1", fftypes.NewUUID())
	ev1.Sequence = 111
	repoll, err := ag.processEvent(context.Background(), ev1)
	assert.False(t, repoll)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), ag.eventPoller.pollingOffset)
	mdi.AssertExpectations(t)
}

func TestProcessEventDataArrivedNoMsgs(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	mdi.On("GetMessagesForData", mock.Anything, mock.Anything, mock.Anything).Return([]*fftypes.Message{}, nil)
	ev1 := fftypes.NewEvent(fftypes.EventTypeDataArrivedBroadcast, "ns1", fftypes.NewUUID())
	ev1.Sequence = 111
	repoll, err := ag.processEvent(context.Background(), ev1)
	assert.False(t, repoll)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), ag.eventPoller.pollingOffset)
	mdi.AssertExpectations(t)
}

func TestProcessDataArrivedError(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	mdi.On("GetMessagesForData", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("pop"))
	repoll, err := ag.processDataArrived(context.Background(), "ns1", fftypes.NewUUID())
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestProcessDataCompleteMessageBlocked(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("GetMessagesForData", mock.Anything, mock.Anything, mock.Anything).Return([]*fftypes.Message{msg}, nil)
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{{ID: fftypes.NewUUID(), Sequence: 111}}, nil)
	repoll, err := ag.processDataArrived(context.Background(), "ns1", fftypes.NewUUID())
	assert.False(t, repoll)
	assert.NoError(t, err)
	mdi.AssertExpectations(t)
}

func TestProcessDataCompleteQueryBlockedFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("GetMessagesForData", mock.Anything, mock.Anything, mock.Anything).Return([]*fftypes.Message{msg}, nil)
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("pop"))
	repoll, err := ag.processDataArrived(context.Background(), "ns1", fftypes.NewUUID())
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteDataAvailFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(false, fmt.Errorf("pop"))
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteUpdateFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{}, nil)
	mdi.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteInsertEventFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{}, nil)
	mdi.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(fmt.Errorf("pop"))
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteGetUnblockedFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{}, nil).Once()
	mdi.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("pop"))
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteInsertUnblockEventFail(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{}, nil).Once()
	mdi.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(nil).Once()
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{
		{ID: fftypes.NewUUID(), Sequence: 111},
	}, nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(fmt.Errorf("pop"))
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.False(t, repoll)
	assert.EqualError(t, err, "pop")
	mdi.AssertExpectations(t)
}

func TestCheckMessageCompleteInsertUnblockOK(t *testing.T) {
	mdi := &databasemocks.Plugin{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ag := newAggregator(ctx, mdi, newEventNotifier(ctx))
	msg := &fftypes.Message{
		Header: fftypes.MessageHeader{
			ID: fftypes.NewUUID(),
		},
	}
	mdi.On("CheckDataAvailable", mock.Anything, mock.Anything).Return(true, nil)
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{}, nil).Once()
	mdi.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(nil).Once()
	mdi.On("GetMessageRefs", mock.Anything, mock.Anything).Return([]*fftypes.MessageRef{
		{ID: fftypes.NewUUID(), Sequence: 111},
	}, nil)
	mdi.On("UpsertEvent", mock.Anything, mock.Anything, false).Return(nil)
	repoll, err := ag.checkMessageComplete(context.Background(), msg)
	assert.True(t, repoll)
	assert.NoError(t, err)
	mdi.AssertExpectations(t)
}