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

package broadcast

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kaleido-io/firefly/mocks/databasemocks"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleSystemBroadcastNSOk(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	ns := &fftypes.Namespace{
		ID:   fftypes.NewUUID(),
		Name: "ns1",
	}
	b, err := json.Marshal(&ns)
	assert.NoError(t, err)
	data := &fftypes.Data{
		Value: fftypes.Byteable(b),
	}

	mbi := bm.database.(*databasemocks.Plugin)
	mbi.On("GetNamespace", mock.Anything, "ns1").Return(nil, nil)
	mbi.On("UpsertNamespace", mock.Anything, mock.Anything, true).Return(nil)
	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.True(t, valid)
	assert.NoError(t, err)

	mbi.AssertExpectations(t)
}

func TestHandleSystemBroadcastNSUpsertFail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	ns := &fftypes.Namespace{
		ID:   fftypes.NewUUID(),
		Name: "ns1",
	}
	b, err := json.Marshal(&ns)
	assert.NoError(t, err)
	data := &fftypes.Data{
		Value: fftypes.Byteable(b),
	}

	mbi := bm.database.(*databasemocks.Plugin)
	mbi.On("GetNamespace", mock.Anything, "ns1").Return(nil, nil)
	mbi.On("UpsertNamespace", mock.Anything, mock.Anything, true).Return(fmt.Errorf("pop"))
	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.False(t, valid)
	assert.EqualError(t, err, "pop")

	mbi.AssertExpectations(t)
}

func TestHandleSystemBroadcastNSMissingData(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{})
	assert.False(t, valid)
	assert.NoError(t, err)
}

func TestHandleSystemBroadcastNSBadID(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	ns := &fftypes.Namespace{}
	b, err := json.Marshal(&ns)
	assert.NoError(t, err)
	data := &fftypes.Data{
		Value: fftypes.Byteable(b),
	}

	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.False(t, valid)
	assert.NoError(t, err)
}

func TestHandleSystemBroadcastNSBadData(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	data := &fftypes.Data{
		Value: fftypes.Byteable(`!{json`),
	}

	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.False(t, valid)
	assert.NoError(t, err)
}

func TestHandleSystemBroadcastDuplicate(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	ns := &fftypes.Namespace{
		ID:   fftypes.NewUUID(),
		Name: "ns1",
	}
	b, err := json.Marshal(&ns)
	assert.NoError(t, err)
	data := &fftypes.Data{
		Value: fftypes.Byteable(b),
	}

	mbi := bm.database.(*databasemocks.Plugin)
	mbi.On("GetNamespace", mock.Anything, "ns1").Return(ns, nil)
	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.False(t, valid)
	assert.NoError(t, err)

	mbi.AssertExpectations(t)
}

func TestHandleSystemBroadcastDupCheckFail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	ns := &fftypes.Namespace{
		ID:   fftypes.NewUUID(),
		Name: "ns1",
	}
	b, err := json.Marshal(&ns)
	assert.NoError(t, err)
	data := &fftypes.Data{
		Value: fftypes.Byteable(b),
	}

	mbi := bm.database.(*databasemocks.Plugin)
	mbi.On("GetNamespace", mock.Anything, "ns1").Return(nil, fmt.Errorf("pop"))
	valid, err := bm.HandleSystemBroadcast(context.Background(), &fftypes.Message{
		Header: fftypes.MessageHeader{
			Topic: fftypes.SystemTopicDefineNamespace,
		},
	}, []*fftypes.Data{data})
	assert.False(t, valid)
	assert.EqualError(t, err, "pop")

	mbi.AssertExpectations(t)
}
