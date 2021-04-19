package commands

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// startCmd represents the start command
// NOTE: This is basically psuedocode
func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [path-name]",
		Aliases: []string{"st"},
		Short:   "Start the listening relayer on a given path",
		Args:    cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start demo-path --max-msgs 3
$ %s start demo-path2 --max-tx-size 10`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}

			if err = ensureKeysExist(c); err != nil {
				return err
			}

			path := config.Paths.MustGet(args[0])
			strategy, err := GetStrategyWithOptions(cmd, path.MustGetStrategy())
			if err != nil {
				return err
			}

			//if relayer.SendToController != nil {
			//	action := relayer.PathAction{
			//		Path: path,
			//		Type: "RELAYER_PATH_START",
			//	}
			//	cont, err := relayer.ControllerUpcall(&action)
			//	if !cont {
			//		return err
			//	}
			//}

			done, err := collactor.RunStrategy(c[src], c[dst], strategy)
			if err != nil {
				return err
			}

			thresholdTime := viper.GetDuration(flagThresholdTime)

			eg := new(errgroup.Group)
			eg.Go(func() error {
				for {
					var timeToExpiry time.Duration
					if err := retry.Do(func() error {
						timeToExpiry, err = UpdateClientsFromChains(c[src], c[dst], thresholdTime)
						if err != nil {
							return err
						}
						return nil
					}, retry.Attempts(5), retry.Delay(time.Millisecond*500), retry.LastErrorOnly(true)); err != nil {
						return err
					}
					time.Sleep(timeToExpiry - thresholdTime)
				}
			})
			if err = eg.Wait(); err != nil {
				return err
			}

			trapSignal(done)
			return nil
		},
	}
	return strategyFlag(updateTimeFlags(cmd))
}


// UpdateClientsFromChains takes src, dst chains, threshold time and update clients based on expiry time
func UpdateClientsFromChains(src, dst *collactor.Chain, thresholdTime time.Duration) (time.Duration, error) {
	var (
		srcTimeExpiry, dstTimeExpiry time.Duration
		err                          error
	)

	eg := new(errgroup.Group)
	eg.Go(func() error {
		srcTimeExpiry, err = collactor.AutoUpdateClient(src, dst, thresholdTime)
		return err
	})
	eg.Go(func() error {
		dstTimeExpiry, err = collactor.AutoUpdateClient(dst, src, thresholdTime)
		return err
	})
	if err := eg.Wait(); err != nil {
		return 0, err
	}

	if srcTimeExpiry <= 0 {
		return 0, fmt.Errorf("client (%s) of chain: %s is expired",
			src.PathEnd.ClientID, src.ChainID)
	}

	if dstTimeExpiry <= 0 {
		return 0, fmt.Errorf("client (%s) of chain: %s is expired",
			dst.PathEnd.ClientID, dst.ChainID)
	}

	minTimeExpiry := math.Min(float64(srcTimeExpiry), float64(dstTimeExpiry))

	return time.Duration(int64(minTimeExpiry)), nil
}

// GetStrategyWithOptions sets strategy specific fields.
func GetStrategyWithOptions(cmd *cobra.Command, strategy collactor.Strategy) (collactor.Strategy, error) {
	switch strategyType := strategy.(type) {
	case *collactor.NaiveStrategy:
		maxTxSize, err := cmd.Flags().GetString(flagMaxTxSize)
		if err != nil {
			return strategyType, err
		}

		txSize, err := strconv.ParseUint(maxTxSize, 10, 64)
		if err != nil {
			return strategyType, err
		}

		// set max size of messages in a relay transaction
		strategyType.MaxTxSize = txSize * MB // in MB

		maxMsgLength, err := cmd.Flags().GetString(flagMaxMsgLength)
		if err != nil {
			return strategyType, err
		}

		msgLen, err := strconv.ParseUint(maxMsgLength, 10, 64)
		if err != nil {
			return strategyType, err
		}

		// set max length messages in relay transaction
		strategyType.MaxMsgLength = msgLen

		return strategyType, nil
	default:
		return strategy, nil
	}
}

// trap signal waits for a SIGINT or SIGTERM and then sends down the done channel
func trapSignal(done func()) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// wait for a signal
	sig := <-sigCh
	fmt.Println("Signal Received", sig.String())
	close(sigCh)

	// call the cleanup func
	done()
}
