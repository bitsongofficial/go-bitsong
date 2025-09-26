# Upgrade From Latest Main-net state

Ensure live network data does not corrupt upgrade integrity.

```sh
# tests from live network data
sh a.sh
# test using cosmovisor + any deterministic preupgrade scripts set in upgradeInfo
sh b.sh
# test using normal cosmovisor + COSMOVISOR_CUSTOM_PREUPGRADE flag
sh c.sh
```
