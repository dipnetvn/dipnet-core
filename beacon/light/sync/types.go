// Copyright 2023 The dipnet-core Authors
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

package sync

import (
	"github.com/dipnetvn/dipnet-core/beacon/light/request"
	"github.com/dipnetvn/dipnet-core/beacon/types"
	"github.com/dipnetvn/dipnet-core/common"
)

var (
	EvNewHead             = &request.EventType{Name: "newHead"}             // data: types.HeadInfo
	EvNewOptimisticUpdate = &request.EventType{Name: "newOptimisticUpdate"} // data: types.OptimisticUpdate
	EvNewFinalityUpdate   = &request.EventType{Name: "newFinalityUpdate"}   // data: types.FinalityUpdate
)

type (
	ReqUpdates struct {
		FirstPeriod, Count uint64
	}
	RespUpdates struct {
		Updates    []*types.LightClientUpdate
		Committees []*types.SerializedSyncCommittee
	}
	ReqHeader  common.Hash
	RespHeader struct {
		Header               types.Header
		Canonical, Finalized bool
	}
	ReqCheckpointData common.Hash
	ReqBeaconBlock    common.Hash
	ReqFinality       struct{}
)
