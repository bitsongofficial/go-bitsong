RH1=$(bitsongcli keys show faucet -a --keyring-backend=test)
RH2=$(bitsongcli keys show validator -a --keyring-backend=test)
IPFS=https://ipfs.infura.io:5001
#IPFS=http://localhost:5001

# create a new addcontent tx
bitsongcli tx content add $1 $2 $3 $4  \
--stream-price 1ubtsg \
--download-price 10ubtsg \
--right-holder "80:$RH1" \
--right-holder "20:$RH2" \
--ipfs-addr $IPFS \
--generate-only > add_content.json

# right holder 1 sign the tx
bitsongcli tx sign add_content.json --from faucet --keyring-backend=test > add_content_sig1.json

# right holder 2 sign the tx
bitsongcli tx sign add_content_sig1.json --from validator --keyring-backend=test > add_content_sig2.json

# broadcast
bitsongcli tx broadcast add_content_sig2.json --from faucet --keyring-backend=test -b block

# query content
# bitsongcli query content resolve $1
