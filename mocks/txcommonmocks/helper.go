// Code generated by mockery v2.30.1. DO NOT EDIT.

package txcommonmocks

import (
	context "context"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"
	core "github.com/hyperledger/firefly/pkg/core"

	mock "github.com/stretchr/testify/mock"
)

// Helper is an autogenerated mock type for the Helper type
type Helper struct {
	mock.Mock
}

// AddBlockchainTX provides a mock function with given fields: ctx, tx, blockchainTXID
func (_m *Helper) AddBlockchainTX(ctx context.Context, tx *core.Transaction, blockchainTXID string) error {
	ret := _m.Called(ctx, tx, blockchainTXID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.Transaction, string) error); ok {
		r0 = rf(ctx, tx, blockchainTXID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOperationInTransaction provides a mock function with given fields: ctx, tx, opType
func (_m *Helper) FindOperationInTransaction(ctx context.Context, tx *fftypes.UUID, opType fftypes.FFEnum) (*core.Operation, error) {
	ret := _m.Called(ctx, tx, opType)

	var r0 *core.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, fftypes.FFEnum) (*core.Operation, error)); ok {
		return rf(ctx, tx, opType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, fftypes.FFEnum) *core.Operation); ok {
		r0 = rf(ctx, tx, opType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, fftypes.FFEnum) error); ok {
		r1 = rf(ctx, tx, opType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockchainEventByIDCached provides a mock function with given fields: ctx, id
func (_m *Helper) GetBlockchainEventByIDCached(ctx context.Context, id *fftypes.UUID) (*core.BlockchainEvent, error) {
	ret := _m.Called(ctx, id)

	var r0 *core.BlockchainEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) (*core.BlockchainEvent, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *core.BlockchainEvent); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.BlockchainEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByIDCached provides a mock function with given fields: ctx, id
func (_m *Helper) GetTransactionByIDCached(ctx context.Context, id *fftypes.UUID) (*core.Transaction, error) {
	ret := _m.Called(ctx, id)

	var r0 *core.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) (*core.Transaction, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *core.Transaction); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertNewBlockchainEvents provides a mock function with given fields: ctx, events
func (_m *Helper) InsertNewBlockchainEvents(ctx context.Context, events []*core.BlockchainEvent) ([]*core.BlockchainEvent, error) {
	ret := _m.Called(ctx, events)

	var r0 []*core.BlockchainEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []*core.BlockchainEvent) ([]*core.BlockchainEvent, error)); ok {
		return rf(ctx, events)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []*core.BlockchainEvent) []*core.BlockchainEvent); ok {
		r0 = rf(ctx, events)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.BlockchainEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []*core.BlockchainEvent) error); ok {
		r1 = rf(ctx, events)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistTransaction provides a mock function with given fields: ctx, id, txType, blockchainTXID
func (_m *Helper) PersistTransaction(ctx context.Context, id *fftypes.UUID, txType fftypes.FFEnum, blockchainTXID string) (bool, error) {
	ret := _m.Called(ctx, id, txType, blockchainTXID)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, fftypes.FFEnum, string) (bool, error)); ok {
		return rf(ctx, id, txType, blockchainTXID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, fftypes.FFEnum, string) bool); ok {
		r0 = rf(ctx, id, txType, blockchainTXID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, fftypes.FFEnum, string) error); ok {
		r1 = rf(ctx, id, txType, blockchainTXID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitNewTransaction provides a mock function with given fields: ctx, txType, idempotencyKey
func (_m *Helper) SubmitNewTransaction(ctx context.Context, txType fftypes.FFEnum, idempotencyKey core.IdempotencyKey) (*fftypes.UUID, error) {
	ret := _m.Called(ctx, txType, idempotencyKey)

	var r0 *fftypes.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, fftypes.FFEnum, core.IdempotencyKey) (*fftypes.UUID, error)); ok {
		return rf(ctx, txType, idempotencyKey)
	}
	if rf, ok := ret.Get(0).(func(context.Context, fftypes.FFEnum, core.IdempotencyKey) *fftypes.UUID); ok {
		r0 = rf(ctx, txType, idempotencyKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, fftypes.FFEnum, core.IdempotencyKey) error); ok {
		r1 = rf(ctx, txType, idempotencyKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHelper creates a new instance of Helper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHelper(t interface {
	mock.TestingT
	Cleanup(func())
}) *Helper {
	mock := &Helper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
