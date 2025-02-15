// TAKE CAUTION RUNNING THIS TEST! IT WILL MODIFY YOUR .BITSONGD DIRECTORY. DO NOT RUN IF USING PRODUCTION MACHINE.
package v0214_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	apptesting "github.com/bitsongofficial/go-bitsong/app/testing"
	"github.com/cosmos/cosmos-sdk/client/flags"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/suite"
)

const (
	dummyUpgradeHeight = 5
	upgradeName        = "v0214"
	TestCustomNodeDir  = ".custom-bitsong"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
	preModule appmodule.HasPreBlocker
}

func (s *UpgradeTestSuite) SetupTest() {
	s.Setup()
	s.preModule = upgrade.NewAppModule(s.App.AppKeepers.UpgradeKeeper, addresscodec.NewBech32Codec("bitsong"))
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) TestUpgrade() {
	upgradeSetup := func(homeDir string) {

		srcDir := filepath.Join(homeDir, "data", "wasm")
		if err := os.MkdirAll(srcDir, 0o755); err != nil {
			panic(err)
		}

		// Create cache/modules directory
		cacheDir := filepath.Join(srcDir, "cache", "modules")
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			panic(err)
		}

		// Create state/wasm directory and populate it with test wasm files
		stateDir := filepath.Join(srcDir, "state", "wasm")
		if err := os.MkdirAll(stateDir, 0o755); err != nil {
			panic(err)
		}

		// Create 73 simulated wasm files each of 0.5MB
		wasmContents := make([]byte, 500000) // 0.5 MB
		for i := 0; i < 73; i++ {
			filename := fmt.Sprintf("wasm_%d.wasm", i)
			filePath := filepath.Join(stateDir, filename)
			if err := os.WriteFile(filePath, wasmContents, 0o644); err != nil {
				panic(err)
			}
		}

		// Confirm that all 73 wasm fileswere created
		for i := 0; i < 73; i++ {
			filename := fmt.Sprintf("wasm_%d.wasm", i)
			filePath := filepath.Join(stateDir, filename)
			if _, err := os.Stat(filePath); err != nil {
				panic(fmt.Sprintf("Failed to find wasm file: %s - %v", filename, err))
			}
		}
		fmt.Printf("Successfully created 73 wasm files\n")
	}

	testCases := []struct {
		name         string
		pre_upgrade  func()
		upgrade      func()
		post_upgrade func()
	}{
		{
			"test_using_default_home",
			func() {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					panic(err)
				}
				homeDir = filepath.Join(homeDir, ".bitsongd")
				upgradeSetup(homeDir)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})

			},
			func() {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					panic(err)
				}
				homeDir = filepath.Join(homeDir, ".bitsongd")
				oldWasmDir := filepath.Join(homeDir, "data", "wasm")
				newWasmDir := filepath.Join(homeDir, "wasm")
				stateDir := filepath.Join(newWasmDir, "state", "wasm")

				// confirm old wasm directory was removed
				if _, err := os.Stat(oldWasmDir); err == nil {
					panic("directory still exists after upgrade")
				}

				// confirm new wasm directory exists
				if _, err := os.Stat(newWasmDir); err != nil {
					panic("new directory does not exists")
				}

				// confirm exclusive.lock exists
				exclusiveLock := filepath.Join(newWasmDir, "exclusive.lock")
				if _, err := os.Stat(exclusiveLock); err != nil {
					panic("exclusive.lock for cosmwasmVM does not exists in updated wasm data path")
				}
				// confirm caching directory exists
				caching := filepath.Join(newWasmDir, "cache", "modules")
				if _, err := os.Stat(caching); err != nil {
					panic("caching file does not exists in updated wasm data path")
				}

				// Verify all 73 wasm files were moved correctly
				// Confirm that all 73 wasm files exist
				for i := 0; i < 73; i++ {
					filename := fmt.Sprintf("wasm_%d.wasm", i)
					filePath := filepath.Join(stateDir, filename)
					if _, err := os.Stat(filePath); err != nil {
						panic(fmt.Sprintf("Failed to find wasm file: %s - %v", filename, err))
					}
				}
			},
		},
		{
			"test_using_custom_home",
			func() {
				appOptions := viper.New()
				testHomeDir := os.ExpandEnv("$HOME/") + TestCustomNodeDir
				if err := os.MkdirAll(testHomeDir, 0o755); err != nil {
					panic(err)
				}

				// set viper variables for test home director
				appOptions.SetDefault(flags.FlagHome, testHomeDir)
				homeDir := filepath.Join(testHomeDir, ".bitsongd")
				upgradeSetup(homeDir)
			},
			func() {
				dummyUpgrade(s)
				s.Require().NotPanics(func() {
					_, err := s.preModule.PreBlock(s.Ctx)
					s.Require().NoError(err)
				})
			},
			func() {
				// confirm old wasm directory was removed
				// confirm new wasm directory exists
				// confirm exclusive.lock exists
				// confirm caching directory exists
				// Confirm that all 73 wasm files exist
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()
			tc.pre_upgrade()
			tc.upgrade()
			tc.post_upgrade()
		})
	}
}

func dummyUpgrade(s *UpgradeTestSuite) {
	s.Ctx = s.Ctx.WithBlockHeight(dummyUpgradeHeight - 1)
	plan := upgradetypes.Plan{Name: upgradeName, Height: dummyUpgradeHeight}
	err := s.App.AppKeepers.UpgradeKeeper.ScheduleUpgrade(s.Ctx, plan)
	s.Require().NoError(err)
	_, err = s.App.AppKeepers.UpgradeKeeper.GetUpgradePlan(s.Ctx)
	s.Require().NoError(err)
	s.Ctx = s.Ctx.WithHeaderInfo(header.Info{Height: dummyUpgradeHeight, Time: s.Ctx.BlockTime().Add(time.Second)}).WithBlockHeight(dummyUpgradeHeight)
}
