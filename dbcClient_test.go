package dbcgosdk_test

import (
	"context"
	dbcgosdk "dbcGoSDK"
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/stretchr/testify/assert"

	"dbcGoSDK/helpers"
	testUtils "dbcGoSDK/internal/test/utils"
	"dbcGoSDK/types"
)

const (
	surfPoolRPCClient = "http://127.0.0.1:8899"
	surfPoolWSlient   = "ws://127.0.0.1:8900"
)

func TestSwap(t *testing.T) {
	conn := rpc.New(surfPoolRPCClient)
	wsClient, err := ws.Connect(context.Background(), surfPoolWSlient)
	if err != nil {
		t.Fatalf("err creating ws client: %s", err.Error())
	}

	t.Cleanup(func() {
		conn.Close()
		wsClient.Close()
	})

	t.Run("swapV1", func(t *testing.T) {
		actors, err := testUtils.SetupTestContext(
			t,
			conn,
			wsClient,
		)
		if err != nil {
			t.Fatalf("err from SetupTestContext: %s", err.Error())
		}

		dbcClient := dbcgosdk.NewDynamicBondingCurveClient(
			conn,
			rpc.CommitmentConfirmed,
		)
		ctx := context.TODO()

		createPoolIx, err := dbcClient.Pool.CreatePool(
			ctx,
			types.CreatePoolParam{
				PreCreatePoolParam: types.PreCreatePoolParam{
					Name:        "Gorks Favourite Phrase (URMOM)",
					Symbol:      "URMOM",
					BaseMint:    actors.BaseMint.PublicKey(),
					PoolCreator: actors.Admin.PublicKey(),
				},
				Config: actors.Config,
				Payer:  actors.Admin.PublicKey(),
			},
		)
		if err != nil {
			t.Fatalf("CreatePool errored: %s", err.Error())
		}

		_, err = testUtils.ExecuteTransaction(
			conn,
			wsClient,
			[]solana.Instruction{createPoolIx},
			actors.Admin,
			actors.BaseMint,
		)
		if err != nil {
			testUtils.PrettyPrintTxnErrorLog(t, err)
			t.FailNow()
		}

		t.Log("createPool success ✅")

		poolConfigState, err := dbcClient.State.GetPoolConfig(ctx, actors.Config)
		if err != nil {
			t.Fatalf("GetPoolConfig errored: %s", err.Error())
		}

		actors.Pool = helpers.DeriveDbcPoolAddress(poolConfigState.QuoteMint, actors.BaseMint.PublicKey(), actors.Config)

		poolState, err := dbcClient.State.GetPool(ctx, actors.Pool)
		if err != nil {
			t.Fatalf("GetPool errored: %s", err.Error())
		}

		fmt.Printf("%+v\n\n%+v\n\n%+v\n", poolConfigState, actors, poolState)

		// poolConfigState, err := dbcClient.State.GetPoolConfig(
		// 	ctx,
		// 	poolState.Config,
		// )
		// if err != nil {
		// 	t.Fatalf("GetPoolConfig errored: %s", err.Error())
		// }

		// currentPoint, err := helpers.GetCurrentPoint(
		// 	conn,
		// 	types.ActivationType(poolConfigState.ActivationType),
		// )
		// if err != nil {
		// 	t.Fatalf("GetCurrentPoint errored: %s", err.Error())
		// }

		// dbcClient.Pool.SwapQuote(
		// 	types.SwapQuoteParam{
		// 		VirtualPool:  poolState,
		// 		Config:       poolConfigState,
		// 		AmountIn:     big.NewInt(1_000_000_000),
		// 		SlippageBps:  50,
		// 		CurrentPoint: currentPoint,
		// 	},
		// )
		t.Log("all good so far ✅")

		ixn, err := dbcClient.Pool.Swap(
			ctx,
			types.SwapParam{
				AmountIn: 1000000000,
				Owner:    actors.User.PublicKey(),
				Pool:     actors.Pool,
				Payer:    actors.User.PublicKey(),
			},
		)
		if err != nil {
			t.Fatalf("Swap errored: %s", err.Error())
		}

		txnSig, err := testUtils.ExecuteTransaction(
			conn,
			wsClient,
			ixn,
			actors.User,
		)
		if err != nil {
			testUtils.PrettyPrintTxnErrorLog(t, err)
			t.FailNow()
		}

		assert.NotNil(t, txnSig)
	})

	t.Run("Add liquidity with Token 2022", func(t *testing.T) {

	})
}
