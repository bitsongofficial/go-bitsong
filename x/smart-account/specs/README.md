# Smart Account Workflow

This modules standout feature is allowing **an account-by-account customization on how transaction are authenticated.**

Bitsong nodes now make use of their compute resources prior to a block being finalized, in order to handle the minimum processes that are needed to verify whether or not an account is valid to be included in the mempool.

**This logic is now programmable via smart contracts & additional configurations,vastly expanding how transactions can be implemented.**


### Minimum Requirements
- One secp256k1 keypair to register to the smart-account, funded with gas to pay for transsaction fees
- An authentication method of choice
- Optional - an `TxExtension` appended to each message, in order to specific which authenticator to use (defualts to regular key ecdsa authorization)

## Workflow Module Uses To Authorize Transactions

###  1. Inital Gas Limit
To increase cost for abusing this authentication module, a maximum gas limit for actions performed by the module is enforced. If any transaction to be authenticated

### 2. Identify Fee Payer:
**The first signer of the transaction is always the feepayer.** This means that for multiple messages, all fees these messages incur will be paid by the single account. Front end curators should keep this in mind when implementing login workflows.

### 3. Authenticate Each Message
Multiple messages to be authorized at once. For each message to be authorized
- its associated account is identified
- any authenticator registered for this account is fetched
- its message is then either validated or rejected.

### 4. Gas Limit Resets
 Once the authorizations are complete, the module resets the gas limits in preparation of the last  steps, which is what w

### 5. **Track**:
After all messages are authenticated, the `Track` function notifies each executed authenticator . This allows authenticators to store any transaction-specific data they might need for the future.

### 6. Execute Message
The transaction is executed

### 7. Confirm Execution
After all messages are executed, the module calls the `ConfirmExecution` function is called for each of the authenticators that authenticated the tx.
This allows authenticators to enforce rules that depend on the outcome of the message execution, like spending and transaction limits.


## Default Authentication Options
___
Below we review the available ways transactions sent to nodes can be authenticated as to whether or not they are inculded when a block gets finalized

## AllOf
Allof means that you can stack authentication requirements together, requiring  100% of the methods to be valid. An example would be a multi-sig, or even a two step-verification making use of one of the other authentication options.

## AnyOf
**AnyOf will recognize transactions as valid if `1 of n` authenticators that an account registered to be used is successful.**

## Signature Verification
The signature verification authenticator is the default authenticator for all accounts. It verifies that the signer of a message is the same as the account associated with the message.


## CosmWasm Authenticator

**This authenticator option allows us to have custom smart contracts made to handle how we accounts can  have actions authorized for them.**
When an account registers a contract address to be used as an autheentication method, the specific paramaters sent by the account registering is **not** stored in the contract state, but rather the module storage,which keeps things light & keeps the compute resources light when making use of the contract. 
 
### How it work's
Bitsong will make use of the contract sudo entry point, only which can be done called by the chain itself. This means when an account making use of cosmwasm authentications submits a tx to be authenticated, the CosmwasmVM is deterministically processed & either is validated or rejected processed prior to deterministically processing the actual message to perform.


### Modules Go Message Structure
To register a cosmwasm authenticator to an account, use the following format
`MsgAddAuthenticator` arguments:

```text
sender: <bech32_address>
type: "CosmwasmAuthenticatorV1"
data: json_bytes({
    contract: "<contract_address>",
    params: [<byte_array>]
})
```

**params** field is a json encoded value for any parameters to save regarding this authenticator. This is in contrast to saving these paramters into the contract state, which is more expensive when retrieving the state in a latter date.
 
**Contract storage should be used only when the authenticator needs to track any dynamic information required.**

### Minimum Requirements For Creating Custom Authentication Contracts
Any smart contract must have theese entrypoints availal  in the Sudo entry point, or else 100% of the authorization request made to a contract missing one of these will fail.

```rs
#[cw_serde]
pub enum AuthenticatorSudoMsg {
    // These two are used to handle the addition and removal of the authenticator
    OnAuthenticatorAdded(OnAuthenticatorAddedRequest),
    OnAuthenticatorRemoved(OnAuthenticatorRemovedRequest),

    // These three are run during authenticating a transaction, specifically during steps 3,5,& 7 in the authentication process
    // link: https://github.com/permissionlessweb/go-bitsong/blob/d7962e28589e2977280cdffbd2d2ea7e62b181e0/x/smart-account/README.md#transaction-authentication-overvie
    Authenticate(AuthenticationRequest),
    Track(TrackRequest),
    ConfirmExecution(ConfirmExecutionRequest),
}
```

#### Rust library 

A [simple rust library](https://github.com/permissionlessweb/bs-nfts/tree/cosmwasm-std-v2/packages/btsg-auth) can be added to the dependencies of your contract to access type definitions used by this module.

## MessageFilter

**The message filter authentications means that you can register to authenticate by default any specific message with a given pattern.**

**This is a very powerful filter, as it can bypass default authentication for your account to perform actions, so use with care!**

Recognizing these accounts more as a **permissionless utility account** may help visualize how this authenticator can be used. 

For example, a faucet-like account can be created by allowing allowing any spend messages with the specific values set:
```json
{
  "@type": "/cosmos.bank.v1beta1.MsgSend",
  "amount": [
    {
      "denom": "ubtsg",
      "amount": "69"
    }
  ]
}
```

Or a way to mint new tokens during a streaming session:
```json
{
   "@type":"/bitsong.fantoken.v1beta1.MsgMint",
   "sender":"bitsong1...", 
   // ...
}
```


## Risks & Limitations Present

### Registration Of Accounts

### Fee's and Gas Consumptions

### Composite Authenticators

### Composite Ids

### Composite Signatures

