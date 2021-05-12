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

package main

import (
	"strings"

	"poly-bridge/basedef"

	"github.com/urfave/cli"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: 1,
	}

	LogDirFlag = cli.StringFlag{
		Name:  "logdir",
		Usage: "log directory",
		Value: "./logs",
	}

	ConfigPathFlag = cli.StringFlag{
		Name:  "cliconfig",
		Usage: "Server config file `<path>`",
		Value: "./config.json",
	}

	ChainIDFlag = cli.Uint64Flag{
		Name:  "chain",
		Usage: "select chainID",
		Value: basedef.ETHEREUM_CROSSCHAIN_ID,
	}

	NFTNameFlag = cli.StringFlag{
		Name:  "name",
		Usage: "set nft name for deploy nft contract, etc.",
		Value: "",
	}

	NFTSymbolFlag = cli.StringFlag{
		Name:  "symbol",
		Usage: "set nft symbol for deploy nft contract, etc.",
		Value: "",
	}

	DstChainFlag = cli.Uint64Flag{
		Name:  "dstChain",
		Usage: "set dest chain for cross chain",
		Value: 0,
	}

	AssetFlag = cli.StringFlag{
		Name:  "asset",
		Usage: "set asset for cross chain or mint nft",
	}

	DstAssetFlag = cli.StringFlag{
		Name:  "dstAsset",
		Usage: "set dest asset for cross chain",
	}

	OwnerAccountFlag = cli.StringFlag{
		Name:  "owner",
		Usage: "set `owner` account",
	}
	SrcAccountFlag = cli.StringFlag{
		Name:  "from",
		Usage: "set `from` account, or approve `sender` account",
	}
	DstAccountFlag = cli.StringFlag{
		Name:  "to",
		Usage: "set `to` account, or approve `spender` account",
	}

	//FeeTokenFlag = cli.BoolTFlag{
	//	Name:  "feeToken",
	//	Usage: "choose erc20 token to be fee token",
	//}

	//NativeTokenFlag = cli.BoolFlag{
	//	Name:  "nativeToken",
	//	Usage: "choose native token as wrapper fee token",
	//}

	//ERC20TokenFlag = cli.BoolFlag{
	//	Name:  "erc20Token",
	//	Usage: "choose erc20 token to be fee token",
	//}

	AmountFlag = cli.StringFlag{
		Name:  "amount",
		Usage: "transfer amount or fee amount, can also used as approve amount",
		Value: "",
	}

	TokenIdFlag = cli.Uint64Flag{
		Name:  "tokenId",
		Usage: "set token id while mint nft",
	}

	LockIdFlag = cli.Uint64Flag{
		Name:  "lockId",
		Usage: "wrap lock nft item id",
	}

	StartFlag = cli.Uint64Flag{
		Name:  "start",
		Usage: "batch get user tokens info with index start",
	}

	LengthFlag = cli.Uint64Flag{
		Name:  "length",
		Usage: "batch get user tokens info with length",
	}

	MethodCodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "decode method code to params, and code format MUST be hex string",
	}

	AdminIndexFlag = cli.IntFlag{
		Name:  "admin",
		Usage: "admin index in keystore, default value is 0",
		Value: 0,
	}

	AddGasFlag = cli.Uint64Flag{
		Name: "addGas",
		Usage: "set gas price if the estimated gas price is not enough, the value should be nGwei, e.g: 4 denotes add 4000000000wei",
		Value: 0,
	}

	EpochFlag = cli.Uint64Flag{
		Name: "epoch",
		Usage: "set okex epoch",
		Value: 0,
	}

	HexFlag = cli.StringFlag{
		Name: "hexfile",
		Usage: "set ok hex file path",
	}
)

var (
	CmdSample = cli.Command{
		Name:   "sample",
		Usage:  "only used to debug this tool.",
		Action: handleSample,
		Flags: []cli.Flag{
			LogLevelFlag,
			ConfigPathFlag,
			ChainIDFlag,
			NFTNameFlag,
			NFTSymbolFlag,
			DstChainFlag,
			AssetFlag,
			DstAssetFlag,
			SrcAccountFlag,
			DstAccountFlag,
			//FeeTokenFlag,
			//ERC20TokenFlag,
			AmountFlag,
			TokenIdFlag,
		},
	}

	CmdDeployECCDContract = cli.Command{
		Name:   "deployECCD",
		Usage:  "admin account deploy ethereum cross chain data contract.",
		Action: handleCmdDeployECCDContract,
	}

	CmdDeployECCMContract = cli.Command{
		Name:   "deployECCM",
		Usage:  "admin account deploy ethereum cross chain manage contract.",
		Action: handleCmdDeployECCMContract,
	}

	CmdDeployCCMPContract = cli.Command{
		Name:   "deployCCMP",
		Usage:  "admin account deploy ethereum cross chain manager proxy contract.",
		Action: handleCmdDeployCCMPContract,
	}

	CmdBindERC20Asset = cli.Command{
		Name:   "bindToken",
		Usage:  "admin account bind erc20 asset to side chain.",
		Action: handleCmdBindERC20Asset,
		Flags: []cli.Flag{
			AssetFlag,
			DstChainFlag,
			DstAssetFlag,
		},
	}

	CmdTransferECCDOwnership = cli.Command{
		Name:   "transferECCDOwnership",
		Usage:  "admin account transfer ethereum cross chain data contract ownership eccm contract.",
		Action: handleCmdTransferECCDOwnership,
	}

	CmdTransferECCMOwnership = cli.Command{
		Name:   "transferECCMOwnership",
		Usage:  "admin account transfer ethereum cross chain manager contract ownership to ccmp contract.",
		Action: handleCmdTransferECCMOwnership,
	}

	CmdRegisterSideChain = cli.Command{
		Name:   "registerSideChain",
		Usage:  "register side chain in poly.",
		Action: handleCmdRegisterSideChain,
	}

	CmdApproveSideChain = cli.Command{
		Name:   "approveSideChain",
		Usage:  "register side chain in poly.",
		Action: handleCmdApproveSideChain,
	}

	CmdSyncSideChainGenesis2Poly = cli.Command{
		Name:   "syncSideGenesis",
		Usage:  "sync side chain genesis header to poly chain.",
		Action: handleCmdSyncSideChainGenesis2Poly,
		Flags: []cli.Flag{
			HexFlag,
		},
	}

	CmdSyncPolyGenesis2SideChain = cli.Command{
		Name:   "syncPolyGenesis",
		Usage:  "sync poly genesis header to side chain.",
		Action: handleCmdSyncPolyGenesis2SideChain,
	}

	CmdNativeTransfer = cli.Command{
		Name:   "transferNative",
		Usage:  "transfer native token.",
		Action: handleCmdNativeTransfer,
		Flags: []cli.Flag{
			SrcAccountFlag,
			DstAccountFlag,
			AmountFlag,
		},
	}

	CmdNativeBalance = cli.Command{
		Name:   "nativeBalance",
		Usage:  "get native balance.",
		Action: handleGetNativeBalance,
		Flags: []cli.Flag{
			SrcAccountFlag,
		},
	}

	CmdEnv = cli.Command{
		Name:   "env",
		Usage:  "ensure your environment is correct",
		Action: handleCmdEnv,
		Flags: []cli.Flag{
			OwnerAccountFlag,
		},
	}
)

//getFlagName deal with short flag, and return the flag name whether flag name have short name
func getFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
