package proposer

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/taikochain/taiko-client/bindings"
	"github.com/taikochain/taiko-client/cmd/flags"
	"github.com/urfave/cli/v2"
)

func TestNewConfigFromCliContext(t *testing.T) {
	l1Endpoint := os.Getenv("L1_NODE_ENDPOINT")
	l2Endpoint := os.Getenv("L2_NODE_ENDPOINT")
	taikoL1 := os.Getenv("TAIKO_L1_ADDRESS")
	taikoL2 := os.Getenv("TAIKO_L2_ADDRESS")
	proposeInterval := "10s"

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: flags.L1NodeEndpoint.Name},
		&cli.StringFlag{Name: flags.L2NodeEndpoint.Name},
		&cli.StringFlag{Name: flags.TaikoL1Address.Name},
		&cli.StringFlag{Name: flags.TaikoL2Address.Name},
		&cli.StringFlag{Name: flags.L1ProposerPrivKey.Name},
		&cli.StringFlag{Name: flags.L2SuggestedFeeRecipient.Name},
		&cli.StringFlag{Name: flags.ProposeInterval.Name},
	}
	app.Action = func(ctx *cli.Context) error {
		c, err := NewConfigFromCliContext(ctx)
		require.Nil(t, err)
		require.Equal(t, l1Endpoint, c.L1Endpoint)
		require.Equal(t, l2Endpoint, c.L2Endpoint)
		require.Equal(t, taikoL1, c.TaikoL1Address.String())
		require.Equal(t, taikoL2, c.TaikoL2Address.String())
		require.Equal(t, bindings.GoldenTouchAddress, crypto.PubkeyToAddress(c.L1ProposerPrivKey.PublicKey))
		require.Equal(t, bindings.GoldenTouchAddress, c.L2SuggestedFeeRecipient)
		require.Equal(t, float64(10), c.ProposeInterval.Seconds())
		require.Nil(t, new(Proposer).InitFromCli(context.Background(), ctx))

		return err
	}

	require.Nil(t, app.Run([]string{
		"TestNewConfigFromCliContext",
		"-" + flags.L1NodeEndpoint.Name, l1Endpoint,
		"-" + flags.L2NodeEndpoint.Name, l2Endpoint,
		"-" + flags.TaikoL1Address.Name, taikoL1,
		"-" + flags.TaikoL2Address.Name, taikoL2,
		"-" + flags.L1ProposerPrivKey.Name, bindings.GoldenTouchPrivKey[2:],
		"-" + flags.L2SuggestedFeeRecipient.Name, bindings.GoldenTouchAddress.Hex(),
		"-" + flags.ProposeInterval.Name, proposeInterval,
	}))
}