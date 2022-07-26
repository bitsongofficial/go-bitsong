<!--
order: 4
-->

# End-Block

Each abci end block call, the operations to update the pending _merkledrops_ are specified to execute. 
More specifically, since each _merkledrop_ is characterized by an `EndHeight` (i.e., the block height at which the airdrop expires), the module can verify at each block if there is any expired _merkledrop_ at that particular time. To perform the operations, the module is able to retrive the the _merkledrops_ ids by the `EndHeight` and to process those drops. In particular, for each retrived _merkledrop_ they are executed the `withdraw` of the unclaimed tokens and then, the _merkledrop_ is cleaned by the state.

## Withdraw
If at the the `EndHeight` block the _merkledrop_ is still in the store, it means that not all the tokens were claimed. For this reason, the module automatically performs a **withdraw** of the unclaimed tokens to the owner wallet. In particular, the module verifies if the `total amount` is lower than the `claimed amount`, calculates the balance as the unclaimed tokens (i.e., the `total amount` - the `claimed` one). This amount is sent to the owner wallet and a corresponding event of type `EventWithdraw` is emitted.

## Delete completed merkledrop
Once the _merkledrop_ ended and its unclaimed tokens have been withdrawn, it is possible to clean indexes and store from the drop. More specifically, it is removed by the list of _merkledrop_ per `owner`, all the indexes linked to its `id` are deleted together with the `merkledrop` object in the store.