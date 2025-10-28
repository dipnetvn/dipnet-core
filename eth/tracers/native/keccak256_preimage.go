package native

// Copyright 2021 The dipnet-core Authors
// This file is part of the dipnet-core library.
//
// The dipnet-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The dipnet-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dipnet-core library. If not, see <http://www.gnu.org/licenses/>.

import (
	"encoding/json"

	"github.com/dipnetvn/dipnet-core/common"
	"github.com/dipnetvn/dipnet-core/common/hexutil"
	"github.com/dipnetvn/dipnet-core/core/tracing"
	"github.com/dipnetvn/dipnet-core/core/vm"
	"github.com/dipnetvn/dipnet-core/crypto"
	"github.com/dipnetvn/dipnet-core/eth/tracers"
	"github.com/dipnetvn/dipnet-core/eth/tracers/internal"
	"github.com/dipnetvn/dipnet-core/log"
	"github.com/dipnetvn/dipnet-core/params"
)

func init() {
	tracers.DefaultDirectory.Register("keccak256PreimageTracer", newKeccak256PreimageTracer, false)
}

// keccak256PreimageTracer is a native tracer that collects preimages of all KECCAK256 operations.
// This tracer is particularly useful for analyzing smart contract execution patterns,
// especially when debugging storage access in Solidity mappings and dynamic arrays.
type keccak256PreimageTracer struct {
	computedHashes map[common.Hash]hexutil.Bytes
}

// newKeccak256PreimageTracer returns a new keccak256PreimageTracer instance.
func newKeccak256PreimageTracer(ctx *tracers.Context, cfg json.RawMessage, chainConfig *params.ChainConfig) (*tracers.Tracer, error) {
	t := &keccak256PreimageTracer{
		computedHashes: make(map[common.Hash]hexutil.Bytes),
	}
	return &tracers.Tracer{
		Hooks: &tracing.Hooks{
			OnOpcode: t.OnOpcode,
		},
		GetResult: t.GetResult,
	}, nil
}

func (t *keccak256PreimageTracer) OnOpcode(pc uint64, op byte, gas, cost uint64, scope tracing.OpContext, rData []byte, depth int, err error) {
	if op == byte(vm.KECCAK256) {
		sd := scope.StackData()
		// it turns out that sometimes the stack is empty, evm will fail in this case, but we should not panic here
		if len(sd) < 2 {
			return
		}

		dataOffset := internal.StackBack(sd, 0).Uint64()
		dataLength := internal.StackBack(sd, 1).Uint64()
		preimage, err := internal.GetMemoryCopyPadded(scope.MemoryData(), int64(dataOffset), int64(dataLength))
		if err != nil {
			log.Warn("keccak256PreimageTracer: failed to copy keccak preimage from memory", "err", err)
			return
		}

		hash := crypto.Keccak256(preimage)

		t.computedHashes[common.Hash(hash)] = hexutil.Bytes(preimage)
	}
}

// GetResult returns the collected keccak256 preimages as a JSON object mapping hashes to preimages.
func (t *keccak256PreimageTracer) GetResult() (json.RawMessage, error) {
	msg, err := json.Marshal(t.computedHashes)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
