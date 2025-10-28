// Copyright 2022 The dipnet-core Authors
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

package tracetest

import (
	"encoding/json"
	"math/big"
	"strings"
	"unicode"

	"github.com/dipnetvn/dipnet-core/common"
	"github.com/dipnetvn/dipnet-core/common/math"
	"github.com/dipnetvn/dipnet-core/consensus/misc/eip4844"
	"github.com/dipnetvn/dipnet-core/core"
	"github.com/dipnetvn/dipnet-core/core/types"
	"github.com/dipnetvn/dipnet-core/core/vm"

	// Force-load native and js packages, to trigger registration
	_ "github.com/dipnetvn/dipnet-core/eth/tracers/js"
	_ "github.com/dipnetvn/dipnet-core/eth/tracers/native"
)

// camel converts a snake cased input string into a camel cased output.
func camel(str string) string {
	pieces := strings.Split(str, "_")
	for i := 1; i < len(pieces); i++ {
		pieces[i] = string(unicode.ToUpper(rune(pieces[i][0]))) + pieces[i][1:]
	}
	return strings.Join(pieces, "")
}

// traceContext defines a context used to construct the block context
type traceContext struct {
	Number     math.HexOrDecimal64   `json:"number"`
	Difficulty *math.HexOrDecimal256 `json:"difficulty"`
	Time       math.HexOrDecimal64   `json:"timestamp"`
	GasLimit   math.HexOrDecimal64   `json:"gasLimit"`
	Miner      common.Address        `json:"miner"`
	BaseFee    *math.HexOrDecimal256 `json:"baseFeePerGas"`
}

func (c *traceContext) toBlockContext(genesis *core.Genesis) vm.BlockContext {
	context := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    c.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(c.Number)),
		Time:        uint64(c.Time),
		Difficulty:  (*big.Int)(c.Difficulty),
		GasLimit:    uint64(c.GasLimit),
	}
	if genesis.Config.IsLondon(context.BlockNumber) {
		context.BaseFee = (*big.Int)(c.BaseFee)
	}

	if genesis.Config.TerminalTotalDifficulty != nil && genesis.Config.TerminalTotalDifficulty.Sign() == 0 {
		context.Random = &genesis.Mixhash
	}

	if genesis.ExcessBlobGas != nil && genesis.BlobGasUsed != nil {
		header := &types.Header{Number: genesis.Config.LondonBlock, Time: *genesis.Config.CancunTime}
		excess := eip4844.CalcExcessBlobGas(genesis.Config, header, genesis.Timestamp)
		header.ExcessBlobGas = &excess
		context.BlobBaseFee = eip4844.CalcBlobFee(genesis.Config, header)
	}
	return context
}

// tracerTestEnv defines a tracer test required fields
type tracerTestEnv struct {
	Genesis      *core.Genesis   `json:"genesis"`
	Context      *traceContext   `json:"context"`
	Input        string          `json:"input"`
	TracerConfig json.RawMessage `json:"tracerConfig"`
}
