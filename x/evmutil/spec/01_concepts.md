<!--
order: 1
-->

# Concepts

## EVM Gas Denom

In order to use the EVM and be compatible with existing clients, the gas denom used by the EVM must be in 18 decimals. Since `uaeth` has 6 decimals of precision, it cannot be used as the EVM gas denom directly.

To use the Aether token on the EVM, the evmutil module provides an `EvmBankKeeper` that is responsible for the conversion of `uaeth` and `aaeth`. A user's excess `aaeth` balance is stored in the `x/evmutil` store, while its `uaeth` balance remains in the cosmos-sdk `x/bank` module.

## `EvmBankKeeper` Overview

The `EvmBankKeeper` provides access to an account's total `aaeth` balance and the ability to transfer, mint, and burn `aaeth`. If anything other than the `aaeth` denom is requested, the `EvmBankKeeper` will panic.

This keeper implements the `x/evm` module's `BankKeeper` interface to enable the usage of `aaeth` denom on the EVM.

### `x/evm` Parameter Requirements

Since the EVM denom `aaeth` is required to use the `EvmBankKeeper`, it is necessary to set the `EVMDenom` param of the `x/evm` module to `aaeth`.

### Balance Calculation of `aaeth`

The `aaeth` balance of an account is derived from an account's **spendable** `uaeth` balance times 10^12 (to derive its `aaeth` equivalent), plus the account's excess `aaeth` balance that can be accessed via the module `Keeper`.

### `aaeth` <> `uaeth` Conversion

When an account does not have sufficient `aaeth` to cover a transfer or burn, the `EvmBankKeeper` will try to swap 1 `uaeth` to its equivalent `aaeth` amount. It does this by transferring 1 `uaeth` from the sender to the `x/evmutil` module account, then adding the equivalent `aaeth` amount to the sender's balance in the module state.

In reverse, if an account has enough `aaeth` balance for one or more `uaeth`, the excess `aaeth` balance will be converted to `uaeth`. This is done by removing the excess `aaeth` balance in the module store, then transferring the equivalent `uaeth` coins from the `x/evmutil` module account to the target account.

The swap logic ensures that all `aaeth` is backed by the equivalent `uaeth` balance stored in the module account.

## ERC20 token <> sdk.Coin Conversion

`x/evmutil` enables the conversion between ERC20 tokens and sdk.Coins. This done through the use of the `MsgConvertERC20ToCoin` & `MsgConvertCoinToERC20` messages (see **[Messages](03_messages.md)**).

Only ERC20 contract address that are whitelist via the `EnabledConversionPairs` param (see **[Params](05_params.md)**) can be converted via these messages.

## Module Keeper

The module Keeper provides access to an account's excess `aaeth` balance and the ability to update the balance.
