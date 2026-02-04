# Release vX.Y.Z

> **Release Date**: YYYY-MM-DD
> **Upgrade Type**: State-breaking / Non-breaking
> **Upgrade Height**: TBD (if state-breaking)

## Overview

Brief description of this release and its main purpose.

## Upgrade Instructions

### For Validators

```bash
# Using cosmovisor (recommended)
# Binary will be downloaded automatically at upgrade height

# Manual upgrade
cd go-bitsong
git fetch --tags
git checkout vX.Y.Z
make install
```

### Upgrade Height

- **Mainnet (bitsong-2b)**: TBD
- **Testnet**: TBD

## Changes

### Security

- Description ([CSA-XXXX](../../operations/security/advisories/CSA-XXXX.md))

### State-Breaking Changes

- Description ([ADR-XXX](../../adr/XXX-title.md))

### Features

- Description (#PR)

### Bug Fixes

- Description (#PR)

### Dependencies

| Dependency | Previous | New |
|------------|----------|-----|
| cosmos-sdk | vX.Y.Z | vX.Y.Z |
| cometbft | vX.Y.Z | vX.Y.Z |
| ibc-go | vX.Y.Z | vX.Y.Z |
| wasmd | vX.Y.Z | vX.Y.Z |

## Binaries

| Platform | Architecture | Download | Checksum |
|----------|--------------|----------|----------|
| Linux | amd64 | [Download](URL) | `sha256:...` |
| Linux | arm64 | [Download](URL) | `sha256:...` |
| macOS | amd64 | [Download](URL) | `sha256:...` |
| macOS | arm64 | [Download](URL) | `sha256:...` |
| Windows | amd64 | [Download](URL) | `sha256:...` |

## Verification

```bash
# Verify binary checksum
sha256sum bitsongd-vX.Y.Z-linux-amd64.tar.gz

# Verify version after upgrade
bitsongd version --long
```

## Related Documents

- [Upgrade Proposal](link-to-proposal)
- [ADR-XXX: Title](../../adr/XXX-title.md)
- [Security Advisory](../../operations/security/advisories/CSA-XXXX.md) (if applicable)
