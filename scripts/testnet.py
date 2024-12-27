import json

genesis_path = "../export-test.json"
output_genesis = "./genesis_new.json"

def modify_genesis(export_json_path, output_json_path):
    with open(export_json_path, 'r') as f:
        export_data = json.load(f)
    
    # set a new chain id
    export_data['chain_id'] = "uni-420"

    # remove smart contract codes, not sure why but initgenesis fails
    export_data['app_state']['wasm']['codes'] = []
    export_data['app_state']['wasm']['contracts'] = []
    export_data['app_state']['wasm']['sequences'] = []

    # Set all validators to jailed except for kintsugi
    kintsugi_addr = "junovaloper1juczud9nep06t0khghvm643hf9usw45r23gsmr"
    kintsugi_byte_addr = "31E927F677282369B7E57D39FF9C47E3845BFDEA"
    kintsugi_conskey = "junovalcons1x85j0anh9q3kndl905ull8z8uwz9hl02y0n499"

    # replace mainnet pubkey with a new local pubkey
    replace_pubkey = "fscxRe/wWtcPp07H4WfiH89kQYEYBYTjRJi2TbhX+Lk="
    replace_byte_addr = "0FAF0893503B86A775C12F434BBC17EAD232C03F"
    replace_conskey = "junovalcons1p7hs3y6s8wr2wawp9ap5h0qhatfr9spl8qyt2s"

    kintval = {}
    new_total_vp = ""
    for val in export_data['validators']:
        if val['address'] == kintsugi_byte_addr:
            val['pub_key']['value'] = replace_pubkey
            val['address'] = replace_byte_addr
            kintval = val
            new_total_vp = val['power']
            break
    
    # remove all validators except kintsugi
    export_data['validators'] = [kintval]

    # replace pubkey and set all validators to jailed except mine
    for val in export_data['app_state']['staking']['validators']:
        if val['operator_address'] == kintsugi_addr:
            val['consensus_pubkey']['key'] = replace_pubkey
        else:
            val['jailed'] = True         

    # set total vp to the new total vp
    export_data['app_state']['staking']['last_total_power'] = new_total_vp

    # remove all validators last powers except mine
    kintlastpower = {}
    for valvp in export_data['app_state']['staking']['last_validator_powers']:
        if valvp['address'] == kintsugi_addr:   
            kintlastpower = valvp
            break
    
    export_data['app_state']['staking']['last_validator_powers'] = [kintlastpower]

    # reset gov params
    export_data['app_state']['gov']['params']['voting_period'] = "300s"
    export_data['app_state']['gov']['params']['quorum'] = "0.001000000000000000"


    # Add a new account with a large amount
    new_account_addr = "bitsong1e7f6j0006x5csqra593uvm7q850ek72jke37z9"
    new_account_pubkey = "A+KQhHTx5smmpHkZCkv3YLv3u1pfRiOPdZpWwt2+DgXq"
    new_account_balance = 10000000000000000  # 10 billion tokens

    # add to app_state.auth.accounts[]
    new_account = {
        "type": "cosmos-sdk/Account",
        "value": {
            "address": new_account_addr,
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": new_account_pubkey
            },
            "account_number": len(export_data['app_state']['auth']['accounts']) + 1,
            "sequence": 0
        }
    }
    export_data['app_state']['auth']['accounts'].append(new_account)

    # Add the new account to the bank module's balances list
    new_account_balance_obj = {
        "address": new_account_addr,
        "coins": [
            {
                "denom": "ubtsg",
                "amount": str(new_account_balance)
            }
        ]
    }
    export_data['app_state']['bank']['balances'].append(new_account_balance_obj)

    # Update the total supply
    ubtsg_supply = next((coin for coin in export_data['app_state']['bank']['supply'] if coin['denom'] == 'ubtsg'), None)    
    if ubtsg_supply:
        ubtsg_supply['amount'] = str(int(ubtsg_supply['amount']) + new_account_balance)
    else:
        export_data['app_state']['bank']['supply'].append({
            "denom": "ubtsg",
            "amount": str(new_account_balance)
        })

    # Convert the JSON data to a string
    export_data_str = json.dumps(export_data)
    
    # Replace all occurrences of kintsugi_conskey with replace_conskey
    export_data_str = export_data_str.replace(kintsugi_conskey, replace_conskey)
    
    # Convert the string back to JSON
    modified_data = json.loads(export_data_str)
    
    # Save the modified clean genesis file
    with open(output_json_path, 'w') as f:
        json.dump(modified_data, f, indent=2)

        ## add new account 

        ## add new balance for account E

        ## update total supply 
        
    print(f"Modified genesis file saved to {output_json_path}")

modify_genesis(genesis_path, output_genesis);