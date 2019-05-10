# Validators Overview

## Introduction

BitSong is based on Tendermint, which relies on a set of validators that are responsible for committing new blocks in the blockchain. These validators participate in the consensus protocol by broadcasting votes which contain cryptographic signatures signed by each validator's private key.

## Hardware

There currently exists no appropriate cloud solution for validator key management. For this reason, validators must set up a physical operation secured with restricted access. A good starting place, for example, would be co-locating in secure data centers.

Validators should expect to equip their datacenter location with redundant power, connectivity, and storage backups. Expect to have several redundant networking boxes for fiber, firewall and switching and then small servers with redundant hard drive and failover. Hardware can be on the low end of datacenter gear to start out with.

We anticipate that network requirements will be low initially. The current testnet requires minimal resources. Then bandwidth, CPU and memory requirements will rise as the network grows. Large hard drives are recommended for storing years of blockchain history.

## Set Up a Website

Set up a dedicated validator's website and signal your intention to become a validator by contacting BitSong Group LTD (hello at bitsong dot io). This is important since delegators will want to have information about the entity they are delegating their `BTSG` to.

## Seek Legal Advice

Seek legal advice if you intend to run a Validator.