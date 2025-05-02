package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	system "github.com/gagliardetto/solana-go/programs/system"
	token "github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	ProgramID        = "dbcij3LWUppWqq96dh6gJWwBifmcGfLSB5D4DuSMaqN"
	DevnetRPC        = "https://devnet.helius-rpc.com/?api-key=YOUR_API_KEY"
	NativeMintString = "So11111111111111111111111111111111111111112"
	MetadataProgram  = "metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s"
	PoolAuthority    = "FhVo3mqL8PW5pH5U2CN4XE33DokiyZnUwuGpH2hmHLuM"
	TokenProgram     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
)

// derivePoolPDA derives the dbc pool address
func derivePoolPDA(quoteMint, baseMint, config solana.PublicKey) solana.PublicKey {
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
	pda, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(ProgramID))
	if err != nil {
		log.Fatalf("find pool PDA: %v", err)
	}
	return pda
}

// deriveVaultPDA derives the dbc token vault address
func deriveVaultPDA(pool, mint solana.PublicKey) solana.PublicKey {
	seed := [][]byte{
		[]byte("token_vault"),
		mint.Bytes(),
		pool.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seed, solana.MustPublicKeyFromBase58(ProgramID))
	if err != nil {
		log.Fatalf("find vault PDA: %v", err)
	}
	return pda
}

// deriveEventAuthorityPDA derives the program event authority address
func deriveEventAuthorityPDA() solana.PublicKey {
	seed := [][]byte{
		[]byte("__event_authority"),
	}
	pda, _, err := solana.FindProgramAddress(seed, solana.MustPublicKeyFromBase58(ProgramID))
	if err != nil {
		log.Fatalf("find event authority PDA: %v", err)
	}
	return pda
}

// deriveMintMetadata derives the mint metadata address
func deriveMintMetadata(mint solana.PublicKey) solana.PublicKey {
	seeds := [][]byte{
		[]byte("metadata"),
		solana.MustPublicKeyFromBase58(MetadataProgram).Bytes(),
		mint.Bytes(),
	}
	pda, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(MetadataProgram))
	if err != nil {
		log.Fatalf("find mint metadata PDA: %v", err)
	}
	return pda
}

func main() {
	ctx := context.Background()
	client := rpc.New(DevnetRPC)

	// 1) load payer and pool creator PKs
	payer := solana.MustPrivateKeyFromBase58("YOUR_PAYER_PK")
	poolCreator := solana.MustPrivateKeyFromBase58("YOUR_POOL_CREATOR_PK")
	tokenBaseProgram := solana.MustPublicKeyFromBase58(TokenProgram)

	// 2) config key (generate on launch.meteora.ag)
	config := solana.MustPublicKeyFromBase58("YOUR_CONFIG_KEY")

	// 3) generate baseMint (can be vanity)
	baseMintWallet := solana.NewWallet()
	baseMint := baseMintWallet.PublicKey()
	fmt.Println("Base mint:", baseMint)

	// 4) quote mint = wrapped SOL
	quoteMint := solana.MustPublicKeyFromBase58(NativeMintString)

	userInputTokenAccount, _, _ := solana.FindAssociatedTokenAddress(
		payer.PublicKey(),
		quoteMint,
	)
	userOutputTokenAccount, _, _ := solana.FindAssociatedTokenAddress(
		payer.PublicKey(),
		baseMint,
	)

	// 5) derive PDAs
	pool := derivePoolPDA(quoteMint, baseMint, config)
	baseVault := deriveVaultPDA(pool, baseMint)
	quoteVault := deriveVaultPDA(pool, quoteMint)
	eventAuthority := deriveEventAuthorityPDA()
	mintMetadata := deriveMintMetadata(baseMint)

	// ==== Build initialize_virtual_pool_with_spl_token instruction ====
	disc := []byte{140, 85, 215, 176, 102, 54, 104, 79}
	name := "test"
	symbol := "TEST"
	uri := "https://test.fun"

	packString := func(s string) []byte {
		b := make([]byte, 4+len(s))
		binary.LittleEndian.PutUint32(b[:4], uint32(len(s)))
		copy(b[4:], []byte(s))
		return b
	}
	data := append(append(append(disc, packString(name)...), packString(symbol)...), packString(uri)...)

	// define static addresses
	poolAuthority := solana.MustPublicKeyFromBase58(PoolAuthority)
	tokenQuoteProgram := solana.MustPublicKeyFromBase58(TokenProgram)
	tokenProgram := solana.MustPublicKeyFromBase58(TokenProgram)

	acctMeta := solana.AccountMetaSlice{
		// 1. config
		{PublicKey: config, IsSigner: false, IsWritable: false},
		// 2. pool_authority
		{PublicKey: poolAuthority, IsSigner: false, IsWritable: false},
		// 3. creator (signer)
		{PublicKey: poolCreator.PublicKey(), IsSigner: true, IsWritable: false},
		// 4. base_mint (signer, writable)
		{PublicKey: baseMint, IsSigner: true, IsWritable: true},
		// 5. quote_mint
		{PublicKey: quoteMint, IsSigner: false, IsWritable: false},
		// 6. pool (writable)
		{PublicKey: pool, IsSigner: false, IsWritable: true},
		// 7. base_vault (writable)
		{PublicKey: baseVault, IsSigner: false, IsWritable: true},
		// 8. quote_vault (writable)
		{PublicKey: quoteVault, IsSigner: false, IsWritable: true},
		// 9. mint_metadata (writable)
		{PublicKey: mintMetadata, IsSigner: false, IsWritable: true},
		// 10. metadata_program
		{PublicKey: solana.MustPublicKeyFromBase58(MetadataProgram), IsSigner: false, IsWritable: false},
		// 11. payer (signer, writable)
		{PublicKey: payer.PublicKey(), IsSigner: true, IsWritable: true},
		// 12. token_quote_program
		{PublicKey: tokenQuoteProgram, IsSigner: false, IsWritable: false},
		// 13. token_program
		{PublicKey: tokenProgram, IsSigner: false, IsWritable: false},
		// 14. system_program
		{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		// 15. event_authority (PDA, same as pool_authority)
		{PublicKey: eventAuthority, IsSigner: false, IsWritable: false},
		// 16. program (ProgramID)
		{PublicKey: solana.MustPublicKeyFromBase58(ProgramID), IsSigner: false, IsWritable: false},
	}
	ixInit := solana.NewInstruction(
		solana.MustPublicKeyFromBase58(ProgramID),
		acctMeta,
		data,
	)

	// wrap and swap quote mint (0.01 SOL)
	amountIn := uint64(1e7)
	rentExemptAmount := uint64(2039280) // minimum rent-exempt balance for WSOL account
	totalAmount := amountIn + rentExemptAmount

	// create WSOL associated token account (ATA)
	createWSOLIx := associatedtokenaccount.NewCreateInstruction(
		payer.PublicKey(),
		payer.PublicKey(),
		quoteMint,
	).Build()

	// wrap SOL by transferring lamports into the WSOL ATA
	wrapSOLIx := system.NewTransferInstruction(
		totalAmount,
		payer.PublicKey(),
		userInputTokenAccount,
	).Build()

	// sync the WSOL account to update its balance
	syncNativeIx := token.NewSyncNativeInstruction(
		userInputTokenAccount,
	).Build()

	// create base-mint ATA for swap output
	createBaseAtaIx := associatedtokenaccount.NewCreateInstruction(
		payer.PublicKey(),
		payer.PublicKey(),
		baseMint,
	).Build()

	// ==== Build swap instruction (buy) ====
	swapDisc := []byte{248, 198, 158, 145, 225, 117, 135, 200}
	minOut := uint64(1)
	buf := make([]byte, 8+8+8)
	copy(buf, swapDisc)
	binary.LittleEndian.PutUint64(buf[8:], amountIn)
	binary.LittleEndian.PutUint64(buf[16:], minOut)

	acctMetaSwap := solana.AccountMetaSlice{
		// 1. pool_authority
		{PublicKey: poolAuthority, IsSigner: false, IsWritable: false},
		// 2. config
		{PublicKey: config, IsSigner: false, IsWritable: false},
		// 3. pool
		{PublicKey: pool, IsSigner: false, IsWritable: true},
		// 4. input_token_account (user's token account for input token)
		{PublicKey: userInputTokenAccount, IsSigner: false, IsWritable: true},
		// 5. output_token_account (user's token account for output token)
		{PublicKey: userOutputTokenAccount, IsSigner: false, IsWritable: true},
		// 6. base_vault
		{PublicKey: baseVault, IsSigner: false, IsWritable: true},
		// 7. quote_vault
		{PublicKey: quoteVault, IsSigner: false, IsWritable: true},
		// 8. base_mint
		{PublicKey: baseMint, IsSigner: false, IsWritable: false},
		// 9. quote_mint
		{PublicKey: quoteMint, IsSigner: false, IsWritable: false},
		// 10. payer
		{PublicKey: payer.PublicKey(), IsSigner: true, IsWritable: true},
		// 11. token_base_program
		{PublicKey: tokenBaseProgram, IsSigner: false, IsWritable: false},
		// 12. token_quote_program
		{PublicKey: tokenQuoteProgram, IsSigner: false, IsWritable: false},
		// 13. referral_token_account (optional; use a valid token account)
		{PublicKey: userInputTokenAccount, IsSigner: false, IsWritable: true},
		// 14. event_authority
		{PublicKey: eventAuthority, IsSigner: false, IsWritable: false},
		// 15. program
		{PublicKey: solana.MustPublicKeyFromBase58(ProgramID), IsSigner: false, IsWritable: false},
	}
	ixSwap := solana.NewInstruction(
		solana.MustPublicKeyFromBase58(ProgramID),
		acctMetaSwap,
		buf,
	)

	// close the WSOL account after swap to recover rent
	closeWSOLIx := token.NewCloseAccountInstruction(
		userInputTokenAccount,
		payer.PublicKey(),
		payer.PublicKey(),
		[]solana.PublicKey{},
	).Build()

	// 6) assemble transaction
	bh, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatalf("GetLatestBlockhash: %v", err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			ixInit,          // your pool init
			createWSOLIx,    // create WSOL ATA
			wrapSOLIx,       // wrap SOL
			syncNativeIx,    // sync WSOL balance
			createBaseAtaIx, // create output base mint ATA
			ixSwap,          // your swap
			closeWSOLIx,     // close WSOL ATA after swap
		},
		bh.Value.Blockhash,
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		log.Fatalf("NewTransaction: %v", err)
	}
	// 7) sign with payer, poolCreator, baseMint
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		switch {
		case key.Equals(payer.PublicKey()):
			return &payer
		case key.Equals(poolCreator.PublicKey()):
			return &poolCreator
		case key.Equals(baseMint):
			return &baseMintWallet.PrivateKey
		default:
			return nil
		}
	})
	if err != nil {
		log.Fatalf("Sign: %v", err)
	}

	// 8) send & confirm
	sig, err := client.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("SendTransaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", sig)

	// wait for confirmation by polling
	for i := 0; i < 30; i++ { // try for 30 secs
		time.Sleep(time.Second)
		resp, err := client.GetTransaction(ctx, sig, &rpc.GetTransactionOpts{
			Commitment: rpc.CommitmentFinalized,
		})
		if err != nil {
			continue
		}
		if resp != nil {
			if resp.Meta != nil && resp.Meta.Err != nil {
				log.Fatalf("Transaction failed: %v", resp.Meta.Err)
			}
			fmt.Printf("Transaction confirmed: %s\n", `https://solscan.io/tx/`+sig.String())
			return
		}
	}
	log.Fatalf("Transaction confirmation timeout")
}
