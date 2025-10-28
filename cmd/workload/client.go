// Copyright 2025 The dipnet-core Authors
// This file is part of dipnet-core.
//
// dipnet-core is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dipnet-core is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dipnet-core. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"

	"github.com/dipnetvn/dipnet-core/common"
	"github.com/dipnetvn/dipnet-core/common/hexutil"
	"github.com/dipnetvn/dipnet-core/core/types"
	"github.com/dipnetvn/dipnet-core/ethclient"
	dipnetclient "github.com/dipnetvn/dipnet-core/ethclient/gethclient"
	"github.com/dipnetvn/dipnet-core/rpc"
	"github.com/urfave/cli/v2"
)

type client struct {
	Eth    *ethclient.Client
	DipNet *dipnetclient.Client
	RPC    *rpc.Client
}

func makeClient(ctx *cli.Context) *client {
	if ctx.NArg() < 1 {
		exit("missing RPC endpoint URL as command-line argument")
	}
	url := ctx.Args().First()
	cl, err := rpc.Dial(url)
	if err != nil {
		exit(fmt.Errorf("could not create RPC client at %s: %v", url, err))
	}
	return &client{
		RPC:    cl,
		Eth:    ethclient.NewClient(cl),
		DipNet: dipnetclient.New(cl),
	}
}

type simpleBlock struct {
	Number hexutil.Uint64 `json:"number"`
	Hash   common.Hash    `json:"hash"`
}

type simpleTransaction struct {
	Hash             common.Hash    `json:"hash"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`
}

func (c *client) getBlockByHash(ctx context.Context, arg common.Hash, fullTx bool) (*simpleBlock, error) {
	var r *simpleBlock
	err := c.RPC.CallContext(ctx, &r, "eth_getBlockByHash", arg, fullTx)
	return r, err
}

func (c *client) getBlockByNumber(ctx context.Context, arg uint64, fullTx bool) (*simpleBlock, error) {
	var r *simpleBlock
	err := c.RPC.CallContext(ctx, &r, "eth_getBlockByNumber", hexutil.Uint64(arg), fullTx)
	return r, err
}

func (c *client) getTransactionByBlockHashAndIndex(ctx context.Context, block common.Hash, index uint64) (*simpleTransaction, error) {
	var r *simpleTransaction
	err := c.RPC.CallContext(ctx, &r, "eth_getTransactionByBlockHashAndIndex", block, hexutil.Uint64(index))
	return r, err
}

func (c *client) getTransactionByBlockNumberAndIndex(ctx context.Context, block uint64, index uint64) (*simpleTransaction, error) {
	var r *simpleTransaction
	err := c.RPC.CallContext(ctx, &r, "eth_getTransactionByBlockNumberAndIndex", hexutil.Uint64(block), hexutil.Uint64(index))
	return r, err
}

func (c *client) getBlockTransactionCountByHash(ctx context.Context, block common.Hash) (uint64, error) {
	var r hexutil.Uint64
	err := c.RPC.CallContext(ctx, &r, "eth_getBlockTransactionCountByHash", block)
	return uint64(r), err
}

func (c *client) getBlockTransactionCountByNumber(ctx context.Context, block uint64) (uint64, error) {
	var r hexutil.Uint64
	err := c.RPC.CallContext(ctx, &r, "eth_getBlockTransactionCountByNumber", hexutil.Uint64(block))
	return uint64(r), err
}

func (c *client) getBlockReceipts(ctx context.Context, arg any) ([]*types.Receipt, error) {
	var result []*types.Receipt
	err := c.RPC.CallContext(ctx, &result, "eth_getBlockReceipts", arg)
	return result, err
}
