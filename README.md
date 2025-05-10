# Meteora DBC Examples in Go

## Overview

This repository contains examples of how to use the Meteora Dynamic Bonding Curve program in Go. Powered by [solana-go](https://github.com/gagliardetto/solana-go).

## Prerequisites

- [Go](https://go.dev/doc/install)

## Usage

1. Install dependencies

```bash
go mod tidy
```

2. Run the examples

```bash
go run examples/<file-name>.go
```

## Examples

- [Create a pool and swap SOL](./examples/create_pool_and_swap_sol.go)
- [Create a pool and swap USDC](./examples/create_pool_and_swap_usdc.go)
- [Fetch pool configuration](./examples/get_pool_config.go)
