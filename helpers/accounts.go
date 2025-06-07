package helpers

import (
	"bytes"
	"log"

	"github.com/dannwee/dbc-go/common"
	"github.com/gagliardetto/solana-go"
)

// Derives the dbc pool address
func DeriveDbcPoolPDA(quoteMint, baseMint, config solana.PublicKey) solana.PublicKey {
	// pda order: the larger public key bytes goes first
	var mintA, mintB solana.PublicKey
	if bytes.Compare(quoteMint.Bytes(), baseMint.Bytes()) > 0 {
		mintA = quoteMint
		mintB = baseMint
	} else {
		mintA = baseMint
		mintB = quoteMint
	}
	seeds := [][]byte{
		[]byte("pool"),
		config.Bytes(),
		mintA.Bytes(),
		mintB.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(common.DbcProgramID))
	if err != nil {
		log.Fatalf("find pool PDA: %v", err)
	}
	return pda
}

// Derives the dbc token vault address
func DeriveTokenVaultPDA(pool, mint solana.PublicKey) solana.PublicKey {
	seed := [][]byte{
		[]byte("token_vault"),
		mint.Bytes(),
		pool.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seed, solana.MustPublicKeyFromBase58(common.DbcProgramID))
	if err != nil {
		log.Fatalf("find vault PDA: %v", err)
	}
	return pda
}

// Derives the event authority PDA
func DeriveEventAuthorityPDA() solana.PublicKey {
	seeds := [][]byte{[]byte("__event_authority")}
	address, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(common.DbcProgramID))
	if err != nil {
		panic(err)
	}
	return address
}

// Derives the pool authority PDA
func DerivePoolAuthorityPDA() solana.PublicKey {
	seeds := [][]byte{[]byte("pool_authority")}
	address, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(common.DbcProgramID))
	if err != nil {
		panic(err)
	}
	return address
}

// Derives the mint metadata address
func DeriveMintMetadataPDA(mint solana.PublicKey) solana.PublicKey {
	seeds := [][]byte{
		[]byte("metadata"),
		solana.MustPublicKeyFromBase58(common.MetadataProgram).Bytes(),
		mint.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(common.MetadataProgram))
	if err != nil {
		log.Fatalf("find mint metadata PDA: %v", err)
	}
	return pda
}
