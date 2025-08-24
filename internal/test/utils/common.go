package testUtils

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type TestActors struct {
	Admin       solana.PrivateKey
	Operator    solana.PrivateKey
	Partner     solana.PrivateKey
	User        solana.PrivateKey
	PoolCreator solana.PrivateKey
	Config      solana.PublicKey
	Pool        solana.PublicKey
	BaseMint    solana.PrivateKey
}

func newTestActors() *TestActors {
	return &TestActors{
		Admin:       solana.NewWallet().PrivateKey,
		Operator:    solana.NewWallet().PrivateKey,
		Partner:     solana.NewWallet().PrivateKey,
		User:        solana.NewWallet().PrivateKey,
		PoolCreator: solana.NewWallet().PrivateKey,
		BaseMint:    solana.NewWallet().PrivateKey,
		Config:      solana.MustPublicKeyFromBase58("CpFPkZUNjUa5bcoASGWHCCYpx6j8HKmcziWr69BY9brb"),
		// Pool:        solana.MustPublicKeyFromBase58("JBRFeG2d1xd65gSUDWsfRT9eRu126HfXcPYj6rBYvRiM"), // SOL <> FGN27LjGn4KC8rtLvHvYNCwHoDX2Lej9ZWDp3EBinMPw
	}
}

func SetupTestContext(
	t *testing.T,
	conn *rpc.Client,
	wsClient *ws.Client,
) (*TestActors, error) {

	actors := newTestActors()

	// fund rootKeyPair
	if _, err := conn.RequestAirdrop(
		context.Background(),
		actors.Admin.PublicKey(),
		9_998*solana.LAMPORTS_PER_SOL,
		rpc.CommitmentFinalized,
	); err != nil {
		return nil, fmt.Errorf("error: RequestAirdrop - \n%s", err.Error())
	}
	t.Log("got airdrop ✅")

	// fund actors
	pubkeys := []solana.PublicKey{
		actors.Operator.PublicKey(),
		actors.Partner.PublicKey(),
		actors.User.PublicKey(),
		actors.PoolCreator.PublicKey(),
	}

	ixns := make([]solana.Instruction, 0, len(pubkeys))

	for _, pubkey := range pubkeys {
		ix := system.NewTransferInstruction(
			1_000*solana.LAMPORTS_PER_SOL,
			actors.Admin.PublicKey(),
			pubkey,
		).Build()
		ixns = append(ixns, ix)
	}

	// lamports, err := conn.GetMinimumBalanceForRentExemption(
	// 	context.Background(),
	// 	token.MINT_SIZE,
	// 	rpc.CommitmentConfirmed,
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("err from GetMinimumBalanceForRentExemption: %s", err.Error())
	// }

	// createIx := system.NewCreateAccountInstruction(
	// 	lamports,
	// 	token.MINT_SIZE,
	// 	solana.TokenProgramID,
	// 	actors.Admin.PublicKey(),
	// 	actors.BaseMint.PublicKey(),
	// ).Build()
	// initIx := token.NewInitializeMint2Instruction(
	// 	6,
	// 	actors.Admin.PublicKey(),
	// 	solana.PublicKey{},
	// 	actors.BaseMint.PublicKey(),
	// ).Build()

	// ixns = append(ixns, createIx, initIx)

	if _, err := SendAndConfirmTxn(
		conn,
		wsClient,
		ixns,
		actors.Admin,
		actors.BaseMint,
	); err != nil {
		return nil, err
	}
	t.Log("actors funded ✅")

	return actors, nil
}

func ExecuteTransaction(conn *rpc.Client,
	wsClient *ws.Client,
	ixns []solana.Instruction,
	payer solana.PrivateKey,
	signers ...solana.PrivateKey,
) (solana.Signature, error) {
	computebudgetIx := computebudget.NewSetComputeUnitPriceInstruction(400_000).Build()
	newIxns := slices.AppendSeq([]solana.Instruction{computebudgetIx}, slices.Values(ixns))
	return SendAndConfirmTxn(conn, wsClient, newIxns, payer, signers...)
}

func SendAndConfirmTxn(
	conn *rpc.Client,
	wsClient *ws.Client,
	ixns []solana.Instruction,
	payer solana.PrivateKey,
	signers ...solana.PrivateKey,
) (solana.Signature, error) {

	signerMap := make(map[solana.PublicKey]*solana.PrivateKey, 1+len(signers))
	signerMap[payer.PublicKey()] = &payer

	for _, signer := range signers {
		s := signer // avoid loop variable capture
		signerMap[s.PublicKey()] = &s
	}

	blockHash, err := conn.GetLatestBlockhash(
		context.Background(),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("error retrieving blockHash: %s", err.Error())
	}

	txn, err := solana.NewTransaction(
		ixns,
		blockHash.Value.Blockhash,
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("error building txn: %s", err.Error())
	}

	if _, err := txn.Sign(func(pubkey solana.PublicKey) *solana.PrivateKey {
		return signerMap[pubkey]
	}); err != nil {
		return solana.Signature{}, fmt.Errorf("unable to sign transaction: %w", err)
	}

	txnSize, _ := txn.MarshalBinary()
	if size := len(txnSize); size > 1232 {
		return solana.Signature{}, fmt.Errorf("transaction size %d exceeds the limit", size)
	}

	txnSig, err := confirm.SendAndConfirmTransaction(
		context.Background(),
		conn,
		wsClient,
		txn,
	)

	if err != nil {
		return solana.Signature{}, fmt.Errorf("error from sent txn: %s", err.Error())
	}

	return txnSig, nil
}

func MinTo(
	amount uint64,
	wallet, mint, mintAuth, payer solana.PublicKey,
) []solana.Instruction {

	createAtaIx := ata.NewCreateInstruction(
		payer,
		wallet,
		mint,
	).Build()

	ataAddr, _, _ := solana.FindAssociatedTokenAddress(
		wallet,
		mint,
	)

	mintToIx := token.NewMintToInstruction(
		amount,
		mint,
		ataAddr,
		mintAuth,
		nil,
	).Build()

	return []solana.Instruction{
		createAtaIx,
		mintToIx,
	}
}
