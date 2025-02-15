package v0214

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitsongofficial/go-bitsong/app/keepers"
	"github.com/bitsongofficial/go-bitsong/app/upgrades"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/viper"
)

func CreateV0214Handler(mm *module.Manager, configurator module.Configurator, bpm upgrades.BaseAppParamManager, k *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(context)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		// Run migrations first
		logger.Info(fmt.Sprintf("pre migrate version map: %v", vm))
		versionMap, err := mm.RunMigrations(sdkCtx, configurator, vm)
		if err != nil {
			return nil, err
		}

		// using viper to grab home variable
		// Get the home directory from the SDK configuration
		nodeHomeDir := viper.GetString("home")

		if nodeHomeDir == "" {
			// use default
			homedir, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			nodeHomeDir = filepath.Join(homedir, ".bitsongd")
		}

		if err := moveWasmFolder(sdkCtx, nodeHomeDir); err != nil {
			return nil, fmt.Errorf("failed to move WASM folder: %w", err)
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", versionMap))
		return versionMap, err
	}
}

// Function to move the entire WASM directory to the new location
func moveWasmFolder(sdkCtx sdk.Context, homeDir string) error {

	// Define source and destination paths
	destDir := filepath.Join(homeDir, ".bitsongd", "wasm")
	srcDir := filepath.Join(homeDir, ".bitsongd", "data", "wasm")
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy the entire directory
	err := copyDirectory(srcDir, destDir)
	if err != nil {
		return fmt.Errorf("failed to copy WASM directory: %w", err)
	}

	sdkCtx.Logger().Info("Successfully moved WASM directory to the new location")

	// now we can prune the old directory
	os.RemoveAll(srcDir)

	return nil
}

// Helper function to copy a directory recursively
func copyDirectory(src string, path string) error {
	// Open the source directory
	srcDir, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source directory: %w", err)
	}
	defer srcDir.Close()

	// Create the destination directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	// Copy each entry in the source directory
	entries, err := srcDir.ReadDir(0)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir.Name(), entry.Name())
		destPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyDirectory(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy subdirectory: %w", err)
			}
		} else {
			// Copy files
			if err := copyFile(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
		}
	}

	return nil
}

// Helper function to copy a single file
func copyFile(src string, dest string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy the file contents
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
