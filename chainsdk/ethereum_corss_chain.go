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

// Notice: functions in this file only used for nft_deploy_tool and test cases.

package chainsdk

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	log "github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"poly-bridge/go_abi/eccd_abi"
	"poly-bridge/go_abi/eccm_abi"
	"poly-bridge/go_abi/eccmp_abi"
	erc20lp "poly-bridge/go_abi/lock_proxy_abi"
	xecdsa "poly-bridge/utils/ecdsa"
)

var (
	EmptyAddress          = common.Address{}
	EmptyHash             = common.Hash{}
	DefaultDeployGasLimit uint64 = 5000000
	DefaultGasLimit       uint64 = 300000
	DefaultAddGasPrice    = big.NewInt(0)
)

func (s *EthereumSdk) DeployECCDContract(key *ecdsa.PrivateKey) (common.Address, error) {
	auth, err := s.makeAuth(key, DefaultDeployGasLimit)
	if err != nil {
		return EmptyAddress, fmt.Errorf("make auth failed")
	}
	contractAddr, tx, _, err := eccd_abi.DeployEthCrossChainData(auth, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyAddress, err
	}
	return contractAddr, nil
}

func (s *EthereumSdk) DeployECCMContract(
	key *ecdsa.PrivateKey,
	eccd common.Address,
	chainID uint64,
) (common.Address, error) {

	auth, err := s.makeAuth(key, DefaultDeployGasLimit)
	if err != nil {
		return EmptyAddress, err
	}
	contractAddress, tx, _, err := eccm_abi.DeployEthCrossChainManager(auth, s.backend(), eccd, chainID)
	if err != nil {
		return EmptyAddress, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyAddress, err
	}
	return contractAddress, nil
}

func (s *EthereumSdk) DeployECCMPContract(key *ecdsa.PrivateKey, eccmAddress common.Address) (common.Address, error) {
	auth, err := s.makeAuth(key, DefaultDeployGasLimit)
	if err != nil {
		return EmptyAddress, err
	}
	contractAddress, tx, _, err := eccmp_abi.DeployEthCrossChainManagerProxy(auth, s.backend(), eccmAddress)
	if err != nil {
		return EmptyAddress, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyAddress, err
	}
	return contractAddress, nil
}

func (s *EthereumSdk) BindERC20Asset(
	key *ecdsa.PrivateKey,
	lockProxyAddr,
	fromAssetHash,
	toAssetHash common.Address,
	targetSideChainId uint64,
) (common.Hash, error) {

	proxy, err := erc20lp.NewLockProxy(lockProxyAddr, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key, DefaultGasLimit)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := proxy.BindAssetHash(auth, fromAssetHash, targetSideChainId, toAssetHash[:])
	if err != nil {
		return EmptyHash, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) TransferECCDOwnership(key *ecdsa.PrivateKey, eccd, eccm common.Address) (common.Hash, error) {

	eccdContract, err := eccd_abi.NewEthCrossChainData(eccd, s.backend())
	if err != nil {
		return EmptyHash, err
	}
	auth, err := s.makeAuth(key, DefaultGasLimit)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := eccdContract.TransferOwnership(auth, eccm)
	if err != nil {
		return EmptyHash, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) GetECCDOwnership(eccdAddr common.Address) (common.Address, error) {

	eccd, err := eccd_abi.NewEthCrossChainData(eccdAddr, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	return eccd.Owner(nil)
}

func (s *EthereumSdk) TransferECCMOwnership(key *ecdsa.PrivateKey, eccm, ccmp common.Address) (common.Hash, error) {

	eccmContract, err := eccm_abi.NewEthCrossChainManager(eccm, s.backend())
	if err != nil {
		return EmptyHash, err
	}
	auth, err := s.makeAuth(key, DefaultGasLimit)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := eccmContract.TransferOwnership(auth, ccmp)
	if err != nil {
		return EmptyHash, fmt.Errorf("TransferECCMOwnership err: %v", err)
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) GetECCMOwnership(eccmAddr common.Address) (common.Address, error) {
	eccm, err := eccm_abi.NewEthCrossChainManager(eccmAddr, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	return eccm.Owner(nil)
}

func (s *EthereumSdk) TransferCCMPOwnership(
	key *ecdsa.PrivateKey,
	ccmpAddr, newOwner common.Address,
) (common.Hash, error) {

	ccmp, err := eccmp_abi.NewEthCrossChainManagerProxy(ccmpAddr, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key, DefaultGasLimit)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := ccmp.TransferOwnership(auth, newOwner)
	if err != nil {
		return EmptyHash, err
	}
	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) GetCCMPOwnership(ccmpAddr common.Address) (common.Address, error) {
	ccmp, err := eccmp_abi.NewEthCrossChainManagerProxy(ccmpAddr, s.backend())
	if err != nil {
		return EmptyAddress, err
	}
	return ccmp.Owner(nil)
}

func (s *EthereumSdk) InitGenesisBlock(key *ecdsa.PrivateKey, eccmAddr common.Address, rawHdr, publickeys []byte) (common.Hash, error) {
	eccm, err := eccm_abi.NewEthCrossChainManager(eccmAddr, s.backend())
	if err != nil {
		return EmptyHash, err
	}

	auth, err := s.makeAuth(key, DefaultGasLimit)
	if err != nil {
		return EmptyHash, err
	}
	tx, err := eccm.InitGenesisBlock(auth, rawHdr, publickeys)
	if err != nil {
		return EmptyHash, err
	}

	if err := s.waitTxConfirm(tx.Hash()); err != nil {
		return EmptyHash, err
	}
	return tx.Hash(), nil
}

func (s *EthereumSdk) dumpTx(hash common.Hash) error {
	tx, err := s.GetTransactionReceipt(hash)
	if err != nil {
		return fmt.Errorf("faild to get receipt %s", hash.Hex())
	}

	if tx.Status == 0 {
		return fmt.Errorf("receipt failed %s", hash.Hex())
	}

	log.Info("txhash %s, block height %d", hash.Hex(), tx.BlockNumber.Uint64())
	for _, event := range tx.Logs {
		log.Info("eventlog address %s", event.Address.Hex())
		log.Info("eventlog data %s", new(big.Int).SetBytes(event.Data).String())
		for i, topic := range event.Topics {
			log.Info("eventlog topic[%d] %s", i, topic.String())
		}
	}
	return nil
}

func (s *EthereumSdk) makeAuth(key *ecdsa.PrivateKey, gasLimit uint64) (*bind.TransactOpts, error) {
	authAddress := xecdsa.Key2address(key)
	nonce, err := s.NonceAt(authAddress)
	if err != nil {
		return nil, fmt.Errorf("makeAuth, addr %s, err %v", authAddress.Hex(), err)
	}

	gasPrice, err := s.SuggestGasPrice()
	if err != nil {
		return nil, fmt.Errorf("makeAuth, get suggest gas price err: %v", err)
	}

	auth := bind.NewKeyedTransactor(key)
	auth.From = authAddress
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(int64(0)) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	//if DefaultAddGasPrice.Cmp(gasPrice) > 0 {
	//	auth.GasPrice = DefaultAddGasPrice
	//} else {
	//	auth.GasPrice = gasPrice
	//}

	return auth, nil
}

func (s *EthereumSdk) waitTxConfirm(hash common.Hash) error {
	ticker := time.NewTicker(time.Second * 1)
	end := time.Now().Add(30 * time.Second)
	for now := range ticker.C {
		_, pending, err := s.TransactionByHash(hash)
		if err != nil {
			log.Debug("failed to call TransactionByHash: %v", err)
			continue
		}
		if !pending {
			break
		}
		if now.Before(end) {
			continue
		}
		log.Info("check your transaction %s on explorer, make sure it's confirmed.", hash.Hex())
		return nil
	}
	if err := s.dumpTx(hash); err != nil {
		log.Error("dump tx %s err: %v", hash.Hex(), err)
	} else {
		log.Info("tx %s confirmed", hash.Hex())
	}
	return nil
}

func (s *EthereumSdk) backend() bind.ContractBackend {
	return s.rawClient
}
