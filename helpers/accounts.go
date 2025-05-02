package helpers

import (
	"bytes"
	"log"

	"github.com/gagliardetto/solana-go"

	"github.com/dannwee/dbc-go/common"
)

// DeriveDbcPoolPDA derives the dbc pool address
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

// DeriveTokenVaultPDA derives the dbc token vault address
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

// DeriveEventAuthorityPDA derives the program event authority address
func DeriveEventAuthorityPDA() solana.PublicKey {
	seed := [][]byte{
		[]byte("__event_authority"),
	}
	pda, _, err := solana.FindProgramAddress(seed, solana.MustPublicKeyFromBase58(common.DbcProgramID))
	if err != nil {
		log.Fatalf("find event authority PDA: %v", err)
	}
	return pda
}

// DeriveMintMetadataPDA derives the mint metadata address
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
