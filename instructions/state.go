package instructions

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dannwee/dbc-go/common"
	"github.com/dannwee/dbc-go/helpers"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
)

// GetPoolConfig fetches and deserializes pool configuration data from the Solana blockchain
func GetPoolConfig(ctx context.Context, configAddress solana.PublicKey, rpcClient *solRpc.Client) (*common.PoolConfig, error) {
	account, err := rpcClient.GetAccountInfo(ctx, configAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool config account: %w", err)
	}

	if account == nil || account.Value == nil {
		return nil, fmt.Errorf("pool config account not found")
	}

	data := account.Value.Data.GetBinary()

	if len(data) < 8 {
		return nil, fmt.Errorf("data too short")
	}

	expectedDiscriminator := []byte{26, 108, 14, 123, 116, 230, 129, 43}
	if !bytes.Equal(data[:8], expectedDiscriminator) {
		return nil, fmt.Errorf("invalid discriminator, not a pool config account")
	}

	return helpers.DeserializePoolConfig(data)
}
