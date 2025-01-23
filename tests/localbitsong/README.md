# LocalBitSong

LocalBitSong is a complete BitSong testnet containerized with Docker and orchestrated with a simple docker-compose file. LocalBitSong comes preconfigured with opinionated, sensible defaults for a standard testing environment.

LocalBitSong comes in two flavors:

1. No initial state: brand new testnet with no initial state. 
2. TBD - With mainnet state: creates a testnet from a mainnet state export

Both ways, the chain-id for LocalBitSong is set to 'LocalBitSong'.

## Prerequisites

Ensure you have docker and docker-compose installed:

```sh
# Docker
sudo apt-get remove docker docker-engine docker.io
sudo apt-get update
sudo apt install docker.io -y

# Docker compose
sudo apt install docker-compose -y
```

## 1. LocalBitSong - No Initial State

The following commands must be executed from the root folder of the BitSong repository.

1. Make any change to the bitsong code that you want to test

2. Initialize LocalBitSong:

```bash
make localnet-init
```

The command:

- Builds a local docker image with the latest changes
- Cleans the `$HOME/.bitsongd-local` folder

3. Start LocalBitSong:

```bash
make localnet-start
```

> Note
>
> You can also start LocalBitSong in detach mode with:
>
> `make localnet-startd`

4. (optional) Add your validator wallet and 10 other preloaded wallets automatically:

```bash
make localnet-keys
```

- These keys are added to your `--keyring-backend test`
- If the keys are already on your keyring, you will get an `"Error: aborted"`
- Ensure you use the name of the account as listed in the table below, as well as ensure you append the `--keyring-backend test` to your txs
- Example: `bitsongd tx bank send lo-test2 bitsong1cyyzpxplxdzkeea7kwsydadg87357qnahakaks --keyring-backend test --chain-id LocalBitSong`

5. You can stop chain, keeping the state with

```bash
make localnet-stop
```

6. When you are done you can clean up the environment with:

```bash
make localnet-clean
```

## LocalBitSong Accounts and Keys

LocalBitSong is pre-configured with one validator and 9 accounts with BTSG balances.

| Account   | Address                                                                                                | Mnemonic                                                                                                                                                                   |
|-----------|--------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| val    | `bitsong1gws6wz8q5kyyu4gqze48fwlmm4m0mdjz0620gw`<br/>`bitsongvaloper1gws6wz8q5kyyu4gqze48fwlmm4m0mdjzw7kxcn` | `bottom loan skill merry east cradle onion journey palm apology verb edit desert impose absurd oil bubble sweet glove shallow size build burst effort`                    |
| lo-test1  | `bitsong1regz7kj3ylg2dn9rl8vwrhclkgz528mf0tfsck`                                                          | `notice oak worry limit wrap speak medal online prefer cluster roof addict wrist behave treat actual wasp year salad speed social layer crew genius`                       |
| lo-test2  | `bitsong1hvrhhex6wfxh7r7nnc3y39p0qlmff6v9t5rc25`                                                          | `quality vacuum heart guard buzz spike sight swarm shove special gym robust assume sudden deposit grid alcohol choice devote leader tilt noodle tide penalty`              |
| lo-test3  | `bitsong175vgzztymvvcxvqun54nlu9dq6856thgvyl5sa`                                                          | `symbol force gallery make bulk round subway violin worry mixture penalty kingdom boring survey tool fringe patrol sausage hard admit remember broken alien absorb`        |
| lo-test4  | `bitsong1t8nznzj4sd6zzutwdmslgy4dcxyd2jafz7822x`                                                          | `bounce success option birth apple portion aunt rural episode solution hockey pencil lend session cause hedgehog slender journey system canvas decorate razor catch empty` |
| lo-test5  | `bitsong14vdrvstsffj8mq5e4fhm6y2hpfxtedajczsj5d`                                                          | `second render cat sing soup reward cluster island bench diet lumber grocery repeat balcony perfect diesel stumble piano distance caught occur example ozone loyal`        |
| lo-test6  | `bitsong1vwe5hay74v0vhuzdhadteyqfasu5d7tdf83pyy`                                                          | `spatial forest elevator battle also spoon fun skirt flight initial nasty transfer glory palm drama gossip remove fan joke shove label dune debate quick`                  |
| lo-test7  | `bitsong16866dezn6ez2qpmpcrrv9cyud8v8c7ufnzwhhh`                                                          | `noble width taxi input there patrol clown public spell aunt wish punch moment will misery eight excess arena pen turtle minimum grain vague inmate`                       |
| lo-test8  | `bitsong1tlwh75lvu35nw9vcg557mxhspz5s88t6vzscd8`                                                          | `cream sport mango believe inhale text fish rely elegant below earth april wall rug ritual blossom cherry detail length blind digital proof identify ride`                 |
| lo-test9  | `bitsong16z9wj8n5f3zgzwspw0r9sj9v7k7hdasqj95us9`                                                          | `index light average senior silent limit usual local involve delay update rack cause inmate wall render magnet common feature laundry exact casual resource hundred`       |
| lo-test10 | `bitsong1gulaxnca7rped0grw0lz4h4zy0xn3ttvmlad8x`                                                          | `prefer forget visit mistake mixture feel eyebrow autumn shop pair address airport diesel street pass vague innocent poem method awful require hurry unhappy shoulder`     |

To list all keys in the keyring named `test`
```bash
bitsongd keys list --keyring-backend test
```

To import an account into the keyring `test`. NOTE: replace the address with any of the above user accounts. 
```bash
bitsongd keys add bitsong1regz7kj3ylg2dn9rl8vwrhclkgz528mf0tfsck --keyring-backend test --recover
```

## FAQ

Q: How do I enable pprof server in LocalBitSong?

A: everything but the Dockerfile is already configured. Since we use a production Dockerfile in LocalBitSong, we don't want to expose the pprof server there by default. As a result, if you would like to use pprof, make sure to add `EXPOSE 6060` to the Dockerfile and rebuild the LocalBitSong image.
