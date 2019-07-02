# Ledger Nano Support

Using a hardware wallet to store your keys greatly improves the security of your crypto assets. The Ledger device acts as an enclave of the seed and private keys, and the process of signing transaction takes place within it. No private information ever leaves the Ledger device. The following is a short tutorial on using the Cosmos Ledger app with the BitSong CLI.

At the core of a Ledger device there is a mnemonic seed phrase that is used to generate private keys. This phrase is generated when you initialize you Ledger. The mnemonic is compatible with Cosmos and can be used to seed new accounts.

::: danger
Do not lose or share your 24 words with anyone. To prevent theft or loss of funds, it is best to keep multiple copies of your mnemonic stored in safe, secure places. If someone is able to gain access to your mnemonic, they will fully control the accounts associated with them.
:::

## BitSong CLI + Ledger Nano

The tool used to generate addresses and transactions on the BitSong network is `bitsongcli`. Here is how to get started.

### Before you Begin

- [Install the Cosmos app onto your Ledger](https://github.com/cosmos/ledger-cosmos/blob/master/README.md#installing)
- [Install Golang](https://golang.org/doc/install)
- [Install Go-BitSong](./installation.md)

Verify that bitsongcli is installed correctly with the following command

```bash
bitsongcli version --long

➜ cosmos-sdk: 0.34.3
git commit: 67ab0b1e1d1e5b898c8cbdede35ad5196dba01b2
vendor hash: 0341b356ad7168074391ca7507f40b050e667722
build tags: netgo ledger
go version go1.11.5 darwin/amd64

```

### Add your Ledger key

- Connect and unlock your Ledger device.
- Open the Cosmos app on your Ledger.
- Create an account in bitsongcli from your ledger key.

::: tip
Be sure to change the _keyName_ parameter to be a meaningful name. The `ledger` flag tells `bitsongcli` to use your Ledger to seed the account.
:::

```bash
bitsongcli keys add <keyName> --ledger

➜ NAME: TYPE: ADDRESS:     PUBKEY:
<keyName> ledger bitsong1... bitsongpub1...
```

BitSong uses [HD Wallets](./hd-wallets.md). This means you can setup many accounts using the same Ledger seed. To create another account from your Ledger device, run;

```bash
bitsongcli keys add <secondKeyName> --ledger
```

### Confirm your address

Run this command to display your address on the device. Use the `keyName` you gave your ledger key. The `-d` flag is supported in version `1.5.0` and higher.

```bash
bitsongcli keys show <keyName> -d
```

Confirm that the address displayed on the device matches that displayed when you added the key.

### Connect to a full node

Next, you need to configure bitsongcli with the URL of a BitSong full node and the appropriate `chain_id`. In this example we connect to the public load balanced full node operated by BitSong on the `bitsong-testnet-1` chain. But you can point your `bitsongcli` to any BitSong full node. Be sure that the `chain_id` is set to the same chain as the full node.

```bash
bitsongcli config node https://node.bitsong.network:26657
bitsongcli config chain_id bitsong-testnet-1
```

Test your connection with a query such as:

``` bash
`bitsongcli query staking validators`
```

### Sign a transaction

You are now ready to start signing and sending transactions. Send a transaction with bitsongcli using the `tx send` command.

``` bash
bitsongcli tx send --help # to see all available options.
```

::: tip
Be sure to unlock your device with the PIN and open the Cosmos app before trying to run these commands
:::

Use the `keyName` you set for your Ledger key and gaia will connect with the Cosmos Ledger app to then sign your transaction.

```bash
bitsongcli tx send <keyName> <destinationAddress> <amount><denomination>
```

When prompted with `confirm transaction before signing`, Answer `Y`.

Next you will be prompted to review and approve the transaction on your Ledger device. Be sure to inspect the transaction JSON displayed on the screen. You can scroll through each field and each message. Scroll down to read more about the data fields of a standard transaction object.

Now, you are all set to start [sending transactions on the network](./delegator-guide-cli.md#sending-transactions).

### Receive funds

To receive funds to the Cosmos account on your Ledger device, retrieve the address for your Ledger account (the ones with `TYPE ledger`) with this command:

```bash
bitsongcli keys list

➜ NAME: TYPE: ADDRESS:     PUBKEY:
<keyName> ledger bitsong1... bitsongpub1...
```

### Further documentation

Not sure what `bitsongcli` can do? Simply run the command without arguments to output documentation for the commands in supports.

::: tip
The `bitsongcli` help commands are nested. So `$ bitsongcli` will output docs for the top level commands (status, config, query, and tx). You can access documentation for sub commands with further help commands.

For example, to print the `query` commands:

```bash
bitsongcli query --help
```

Or to print the `tx` (transaction) commands:

```bash
bitsongcli tx --help
```
:::

# The BitSong Standard Transaction

Transactions in BitSong embed the [Standard Transaction type](https://godoc.org/github.com/cosmos/cosmos-sdk/x/auth#StdTx) from the Cosmos SDK. The Ledger device displays a serialized JSON representation of this object for you to review before signing the transaction. Here are the fields and what they mean:

- `chain-id`: The chain to which you are broadcasting the tx, such as the `bitsong-testnet-1` testnet.
- `account_number`: The global id of the sending account assigned when the account receives funds for the first time.
- `sequence`: The nonce for this account, incremented with each transaction.
- `fee`: JSON object describing the transaction fee, its gas amount and coin denomination
- `memo`: optional text field used in various ways to tag transactions.
- `msgs_<index>/<field>`: The array of messages included in the transaction. Double click to drill down into nested fields of the JSON.