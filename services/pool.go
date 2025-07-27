package services

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/generated/dbc"
	"dbcGoSDK/types"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type PoolService struct {
	state *StateService
}

func NewPoolService(
	conn *rpc.Client,
	commitment rpc.CommitmentType,
) *PoolService {
	return &PoolService{
		state: NewStateService(conn, commitment),
	}
}

// initializeSplPool initialize a pool with SPL token.
func (p *PoolService) initializeSplPool(
	param types.InitializePoolBaseParam,
) (*dbc.Instruction, error) {
	initializeVirtualPoolWithSplTokenPtr := dbc.NewInitializeVirtualPoolWithSplTokenInstruction(
		dbc.InitializePoolParameters{
			Name:   param.Name,
			Symbol: param.Symbol,
			Uri:    param.URI,
		},
		param.Config,
		p.state.GetPoolAuthority(),
		param.PoolCreator,
		param.BaseMint,
		param.QuoteMint,
		param.Pool,
		param.BaseVault,
		param.QuoteVault,
		param.MintMetadata,
		constants.MetaplexProgramId,
		param.Payer,
		solana.TokenProgramID,
		solana.TokenProgramID,
		solana.SystemProgramID,
		solana.PublicKey{},
		constants.DBCProgramId,
	)
	eventAuthPDA, _, err := initializeVirtualPoolWithSplTokenPtr.FindEventAuthorityAddress()
	if err != nil {
		return nil, fmt.Errorf("err deriving eventAuthPDA: %w", err)
	}

	return initializeVirtualPoolWithSplTokenPtr.
		SetEventAuthorityAccount(eventAuthPDA).ValidateAndBuild()
}

// initializeToken2022Pool initialize a pool with token22.
func (p *PoolService) initializeToken2022Pool(
	param types.InitializePoolBaseParam,
) (*dbc.Instruction, error) {
	initializeVirtualPoolWithToken2022Ptr := dbc.NewInitializeVirtualPoolWithToken2022Instruction(
		dbc.InitializePoolParameters{
			Name:   param.Name,
			Symbol: param.Symbol,
			Uri:    param.URI,
		},
		param.Config,
		p.state.GetPoolAuthority(),
		param.PoolCreator,
		param.BaseMint,
		param.QuoteMint,
		param.Pool,
		param.BaseVault,
		param.QuoteVault,
		param.Payer,
		solana.TokenProgramID,
		solana.Token2022ProgramID,
		solana.SystemProgramID,
		solana.PublicKey{},
		constants.DBCProgramId,
	)

	eventAuthPDA, _, err := initializeVirtualPoolWithToken2022Ptr.FindEventAuthorityAddress()
	if err != nil {
		return nil, fmt.Errorf("err deriving eventAuthPDA: %w", err)
	}

	return initializeVirtualPoolWithToken2022Ptr.
		SetEventAuthorityAccount(eventAuthPDA).ValidateAndBuild()
}

// func (p *PoolService) initializeToken2022Pool()
