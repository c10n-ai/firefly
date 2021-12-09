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

package contracts

import (
	"context"
	"errors"
	"testing"

	"github.com/hyperledger/firefly/mocks/blockchainmocks"
	"github.com/hyperledger/firefly/mocks/broadcastmocks"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/mocks/publicstoragemocks"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateInvokeContractRequest(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(nil)
	mbi.On("ValidateInvokeContractRequest", mock.Anything, mock.Anything).Return(nil)

	req := &fftypes.InvokeContractRequest{
		Method: &fftypes.FFIMethod{
			Name: "sum",
			Params: []*fftypes.FFIParam{
				{
					Name:         "x",
					Type:         "integer",
					InternalType: "uint256",
				},
				{
					Name:         "y",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
			Returns: []*fftypes.FFIParam{
				{
					Name:         "z",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
		},
		Params: map[string]interface{}{
			"x": float64(1),
			"y": float64(2),
		},
	}
	err := cm.ValidateInvokeContractRequest(context.Background(), req)
	assert.NoError(t, err)
}

func TestValidateInvokeContractRequestMissingInput(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(nil)
	mbi.On("ValidateInvokeContractRequest", mock.Anything, mock.Anything).Return(nil)

	req := &fftypes.InvokeContractRequest{
		Method: &fftypes.FFIMethod{
			Name: "sum",
			Params: []*fftypes.FFIParam{
				{
					Name:         "x",
					Type:         "integer",
					InternalType: "uint256",
				},
				{
					Name:         "y",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
			Returns: []*fftypes.FFIParam{
				{
					Name:         "z",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
		},
		Params: map[string]interface{}{
			"x": float64(1),
		},
	}
	err := cm.ValidateInvokeContractRequest(context.Background(), req)
	assert.Regexp(t, "Missing required input argument 'y'", err)
}

func TestValidateInvokeContractRequestInputWrongType(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(nil)
	mbi.On("ValidateInvokeContractRequest", mock.Anything, mock.Anything).Return(nil)

	req := &fftypes.InvokeContractRequest{
		Method: &fftypes.FFIMethod{
			Name: "sum",
			Params: []*fftypes.FFIParam{
				{
					Name:         "x",
					Type:         "integer",
					InternalType: "uint256",
				},
				{
					Name:         "y",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
			Returns: []*fftypes.FFIParam{
				{
					Name:         "z",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
		},
		Params: map[string]interface{}{
			"x": float64(1),
			"y": "two",
		},
	}
	err := cm.ValidateInvokeContractRequest(context.Background(), req)
	assert.Regexp(t, "Input.*not expected.*integer", err)
}

func TestValidateInvokeContractRequestInvalidMethod(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(errors.New("pop"))

	req := &fftypes.InvokeContractRequest{
		Method: &fftypes.FFIMethod{
			Name: "sum",
			Params: []*fftypes.FFIParam{
				{
					Name:         "x",
					Type:         "integer",
					InternalType: "uint256",
				},
				{
					Name:         "y",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
			Returns: []*fftypes.FFIParam{
				{
					Name:         "z",
					Type:         "integer",
					InternalType: "uint256",
				},
			},
		},
		Params: map[string]interface{}{
			"x": float64(1),
			"y": float64(2),
		},
	}

	err := cm.ValidateInvokeContractRequest(context.Background(), req)
	assert.Regexp(t, err, "pop")
}

func TestValidateFFI(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(nil)
	mbi.On("ValidateFFIEvent", mock.Anything, mock.Anything).Return(nil)

	ffi := &fftypes.FFI{
		Name:      "math",
		Namespace: "default",
		Methods: []*fftypes.FFIMethod{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "x",
						Type:         "integer",
						InternalType: "uint256",
					},
					{
						Name:         "y",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
				Returns: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
		Events: []*fftypes.FFIEvent{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
	}

	err := cm.ValidateFFI(context.Background(), "default", ffi)
	assert.NoError(t, err)
}

func TestValidateFFIBadMethod(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(errors.New("pop"))

	ffi := &fftypes.FFI{
		Name:      "math",
		Namespace: "default",
		Methods: []*fftypes.FFIMethod{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "x",
						Type:         "integer",
						InternalType: "uint256",
					},
					{
						Name:         "y",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
				Returns: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
		Events: []*fftypes.FFIEvent{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
	}

	err := cm.ValidateFFI(context.Background(), "default", ffi)
	assert.Regexp(t, err, "pop")
}

func TestValidateFFIBadEvent(t *testing.T) {
	mdb := &databasemocks.Plugin{}
	mps := &publicstoragemocks.Plugin{}
	mbm := &broadcastmocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	cm := NewContractManager(mdb, mps, mbm, nil, mbi)

	mbi.On("ValidateFFIMethod", mock.Anything, mock.Anything).Return(nil)
	mbi.On("ValidateFFIEvent", mock.Anything, mock.Anything).Return(errors.New("pop"))

	ffi := &fftypes.FFI{
		Name:      "math",
		Namespace: "default",
		Methods: []*fftypes.FFIMethod{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "x",
						Type:         "integer",
						InternalType: "uint256",
					},
					{
						Name:         "y",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
				Returns: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
		Events: []*fftypes.FFIEvent{
			{
				Name: "sum",
				Params: []*fftypes.FFIParam{
					{
						Name:         "z",
						Type:         "integer",
						InternalType: "uint256",
					},
				},
			},
		},
	}

	err := cm.ValidateFFI(context.Background(), "default", ffi)
	assert.Regexp(t, err, "pop")
}