// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/bindings/ERC20Handler"
	ERC20 "github.com/stafiprotocol/chainbridge/bindings/ERC20PresetMinterPauser"
)

func DeployMintApproveErc20(client *Client, erc20Handler common.Address, amount *big.Int) (common.Address, error) {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	// Deploy
	erc20Addr, tx, erc20Instance, err := ERC20.DeployERC20PresetMinterPauser(client.Opts, client.Client, "", "")
	if err != nil {
		return ZeroAddress, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	// Mint
	err = client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	_, err = erc20Instance.Mint(client.Opts, client.Opts.From, amount)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	// Approve
	err = client.LockNonceAndUpdate()
	if err != nil {
		return ZeroAddress, err
	}

	tx, err = erc20Instance.Approve(client.Opts, erc20Handler, amount)
	if err != nil {
		return ZeroAddress, err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return ZeroAddress, err
	}

	client.UnlockNonce()

	return erc20Addr, nil
}

func Erc20Approve(client *Client, erc20Contract, recipient common.Address, amount *big.Int) error {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return err
	}

	tx, err := instance.Approve(client.Opts, recipient, amount)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockNonce()

	return nil
}

func Erc20GetBalance(client *Client, erc20Contract, account common.Address) (*big.Int, error) { //nolint:unused,deadcode
	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return nil, err
	}

	bal, err := instance.BalanceOf(client.CallOpts, account)
	if err != nil {
		return nil, err

	}
	return bal, nil

}

func FundErc20Handler(client *Client, handlerAddress, erc20Address common.Address, amount *big.Int) error {
	err := Erc20Approve(client, erc20Address, handlerAddress, amount)
	if err != nil {
		return err
	}

	instance, err := ERC20Handler.NewERC20Handler(handlerAddress, client.Client)
	if err != nil {
		return err
	}

	client.Opts.Nonce = client.Opts.Nonce.Add(client.Opts.Nonce, big.NewInt(1))
	tx, err := instance.FundERC20(client.Opts, erc20Address, client.Opts.From, amount)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	return nil
}
