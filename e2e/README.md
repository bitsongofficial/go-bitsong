# Interchaintest

## Current ICTests
| CLI | Description   |  |   |   |
|---|---|---|---|---|
| `basic`  | Basic network setup and cosmwasm upload & execute  |   |   |   |
| `pfm`  | Test the functionality of packet-forward-middleware   |   |   |   |
| `polytone`  | Test IBC-Callbacks via Polytone deployment  |   |   |   |
| `upgrade`  | Simulate an planned upgrade workflow   |   |   |   |

## Requirements 
##  How to add more tests
- default configs & enviroment setup logic located in `setup.go`

### 1. Define new `x_test.go` file
Ideally the name of the file is the scope of the actions being tested.

### 2. Write Custom Integration Tests
Existing documentation for using interchaintest to write tests can be found [here](https://interchaintest-docs.vercel.app/RunTests/write-custom-tests).

### 3. Add to Make File & CI File
Make sure to add the cli script to the make file helper documentation in `./scripts/makefile/e2e.mk`, along with the make command for the new test. prefix the command with `e2e-`.
```mk
e2e-grpc: rm-testcache
	cd e2e && go test -race -v -run TestBasicGrpcQuery .
```

 To run the integration test in the CI environment, add the make file cli commands to `/.github/interchain-e2e.yml`:
 ```mk
        # names of `make` commands to run tests
        test:
          - "e2e-basic"
          - "e2e-pfm"
          - "e2e-polytone"
          # - "e2e-upgrade"
 ```