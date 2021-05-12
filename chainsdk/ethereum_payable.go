/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

// Notice: functions in this file only used for deploy_tool and test cases.

package chainsdk

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	xecdsa "poly-bridge/utils/ecdsa"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	polycm "github.com/polynetwork/poly/common"
)

var NativeFeeToken = common.HexToAddress("0x0000000000000000000000000000000000000000")

func (s *EthereumSdk) TransferNative(
	key *ecdsa.PrivateKey,
	to common.Address,
	amount *big.Int,
) (common.Hash, error) {

	from := xecdsa.Key2address(key)
	nonce, err := s.NonceAt(from)
	if err != nil {
		return EmptyHash, err
	}

	gasPrice, err := s.SuggestGasPrice()
	if err != nil {
		return EmptyHash, err
	}

	gasLimit, err := s.EstimateGas(ethereum.CallMsg{
		From: from, To: &to, Gas: 0, GasPrice: gasPrice,
		Value: amount, Data: []byte{},
	})
	if err != nil {
		return EmptyHash, err
	}

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, key)
	if err != nil {
		return EmptyHash, err
	}
	if err := s.SendRawTransaction(signedTx); err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(signedTx.Hash()); err != nil {
		return EmptyHash, err
	}
	return signedTx.Hash(), nil
}

func (s *EthereumSdk) GetNativeBalance(owner common.Address) (*big.Int, error) {
	return s.rawClient.BalanceAt(context.Background(), owner, nil)
}

type WrapLockMethod struct {
	FromAsset common.Address
	ToChainId uint64
	ToAddress common.Address
	TokenId   *big.Int
	FeeToken  common.Address
	Fee       *big.Int
	Id        *big.Int
}

func assembleSafeTransferCallData(toAddress common.Address, chainID uint64) []byte {
	sink := polycm.NewZeroCopySink(nil)
	sink.WriteVarBytes(toAddress.Bytes())
	sink.WriteUint64(chainID)
	return sink.Bytes()
}

func filterTokenInfo(enc []byte) map[string]string {
	source := polycm.NewZeroCopySource(enc)
	var (
		num     polycm.Uint256
		url     string
		tokenId *big.Int
		eof     bool
		res     = make(map[string]string)
	)
	for {
		if num, eof = source.NextHash(); !eof {
			bz := polycm.ToArrayReverse(num[:])
			tokenId = new(big.Int).SetBytes(bz)
		} else {
			break
		}
		if url, eof = source.NextString(); !eof {
			res[tokenId.String()] = url
		} else {
			break
		}
	}

	return res
}
