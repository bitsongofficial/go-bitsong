<p align="center" background="black"><img src="bitsong-logo.png" width="398"></p>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/BitSongOfficial/go-bitsong/blob/master/LICENSE)

**BitSong** is a new music streaming platform based on [Tendermint](https://github.com/tendermint/tendermint) consensus algorythm and the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits. Please make sure you study these projects as well if you are not already familiar.

**BitSong** is a project dedicated to musicians and their fans, which aims to overcome the bureaucratic and economic obstacles within this industry and reward artists and users for simply using the platform.

On the **BitSong** platform you (artist) will be able to produce songs in which an advertiser can attach advertisements and users can access from any device. Funds via the Bitsong token `$BTSG` will be credited to the artist wallet immediately and they will be able to withdraw or convert as they see fit.

**Artists** need no longer to wait several months before a record label sends various reports, they can check the progress in real time directly within the Wallet.

_NOTE: This is alpha software. Please contact us if you aim to run it in production._

**Note**: Requires [Go 1.12.4+](https://golang.org/dl/)

# Welcome to StackEdit!

**BitSong** è una nuova piattaforma di streaming musicale basata sull'algoritmo di consenso [Tendermint](https://github.com/tendermint/tendermint)  ed la [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) toolkits. Nel caso in cui non siete gia familiari con questi progetti, è consigliato di andarseli a studiare prima.

**BitSong** è un progetto dedicato a musicisti ed i loro fan, che mira a superare gli ostacoli burocratici ed economici all'interno di questo settore e premiare artisti e utenti semplicemente utilizzando la piattaforma.

Sulla piattaforma di **BitSong** tu (artista) sarai in grado di produrre canzoni in cui un inserzionista può allegare pubblicità e gli utenti possono accedere da qualsiasi dispositivo. I fondi tramite il token Bitsong `$BTSG` saranno accreditati immediatamente al portafoglio dell'artista che sarà in grado di prelevarli o convertirli come lo ritiene più opportuno.

**Gli Artisti** non devono più aspettare diversi mesi prima che un'etichetta discografica invii vari rapporti, possono controllare i progressi in tempo reale direttamente all'interno del Wallet.

_NOTA: Questo è un software Alpha. Vi preghiamo di contattarci se avete intenzione di eseguirlo in produzione._

**Nota**: Richiede [Go 1.12.4+](https://golang.org/dl/)

# Installare la Blockchain di BitSong

Esistono molti modi per installare il nodo TestSense Blockchain di BitSong sulla tua macchina.

## Dalla Source
1. **Installa Go** seguendo il [official docs](https://golang.org/doc/install).Ricorda di impostare le variabili di ambiente `$ GOPATH`,` $ GOBIN` e `$ PATH`, ad esempio:
	```bash
	mkdir -p $HOME/go/bin
	echo  "export GOPATH=$HOME/go" >> ~/.bash_profile
	echo  "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile
	echo  "export PATH=\$PATH:\$GOBIN" >> ~/.bash_profile
	echo  "export GO111MODULE=on" >> ~/.bash_profile
	source ~/.bash_profile
	```
2. **Clona il codice sorgente di BitSong sul tuo computer**
	```bash
	mkdir -p $GOPATH/src/github.com/BitSongOfficial
	cd $GOPATH/src/github.com/BitSongOfficial
	git clone https://github.com/BitSongOfficial/go-bitsong.git
	cd go-bitsong
	```
  3. **Compilare**
		```bash
		# Install the app into your $GOBIN
		make install
		# Now you should be able to run the following commands:
		bitsongd help
		bitsongcli help
		```
		L'ultima versione di `go-bitsong version` è installata correttamente.
3. **Avvia BitSong**
	```bash
	bitsongd start
	```

## Installa su Digital Ocean
1. **Clona il repository**
    ```bash
	git clone https://github.com/BitSongOfficial/go-bitsong.git
    chmod +x go-bitsong/scripts/install/install_ubuntu.sh
	```
2. **Avvia lo script**
    ```bash
    go-bitsong/scripts/install/install_ubuntu.sh
    source ~/.profile
	```
3. Ora dovresti essere in grado di eseguire i seguenti comandi:
	```bash
	bitsongd help
	bitsongcli help
	```
    L'ultima versione di `go-bitsongd version` è installata correttamente.

## Esecuzione della Testnet e utilizzo dei comandi

Per inizializzare la configurazione e un file `genesis.json` per l'applicazione e un account per le transazioni, iniziare eseguendo:

> _*NOTA*_: nei comandi sottostanti gli indirizzi vengono estratti utilizzando le utilità del terminale. Puoi anche inserire le stringhe non elaborate, salvate dalla creazione delle key, mostrate sotto. I comandi richiedono [`jq`](https://stedolan.github.io/jq/download/) da installare sulla tua macchina.

> _*NOTA*_: Se hai già eseguito il tutorial, puoi iniziare da zero con un `bitsongd unsafe-reset-all` o eliminando entrambe le cartelle home `rm -rf ~/.bitsong*`

>  _*NOTA*_: Se si dispone dell'app Cosmos per il ledger e si desidera utilizzarlo, quando si crea la chiave con `bitsongcli keys add jack` aggiungere `--ledger` alla fine. Questo  è tutto cio che vi serve. Quando firmi, `jack` sarà riconosciuto come un tasto Ledger e richiederà un dispositivo.

```bash
# Inizializza i file di configurazione e il file di genesi
bitsongd init --chain-id bitsong-test-network-1

# Copia qui l'output `Address` e salvalo per un uso successivo
# [opzionale] aggiungi "--ledger" alla fine per usare un Ledger Nano S
bitsongcli keys add jack

# Copia qui l'output `Address` e salvalo per un uso successivo
bitsongcli keys add alice

# Aggiungi entrambi gli account, con le monete nel file di genesi
bitsongd add-genesis-account $(bitsongcli keys show jack -a) 1000btsg
bitsongd add-genesis-account $(bitsongcli keys show alice -a) 1000btsg

# Configura la tua CLI per eliminare la necessità di flag dell'identificativo della chain
bitsongcli config chain-id bitsong-test-network-1
bitsongcli config output json
bitsongcli config indent true
bitsongcli config trust-node true
```

È ora possibile avviare `bitsongd` chiamando `bitsongd start`. Vedrai che i registri iniziano lo streaming che rappresenta i blocchi prodotti, questo richiederà un paio di secondi.

Apri un altro terminale per eseguire comandi sulla rete appena creata:

```bash
# Innanzitutto controlla gli account per assicurarti che abbiano fondi
bitsongcli query account $(bitsongcli keys show jack -a)
bitsongcli query account $(bitsongcli keys show alice -a)
```

# Transazioni
Ora puoi iniziare la prima transazione

```bash
bitsongcli tx send --from=$(bitsongcli keys show jack -a)  $(bitsongcli keys show alice -a) 10btsg
```

# Richiesta
Richiedi un account
```bash
bitsongcli query account $(bitsongcli keys show jack -a)
```

## Risorse
- [Website Ufficiale](https://bitsong.io)

### Community
- [Telegram Channel (English)](https://t.me/BitSongOfficial)
- [Facebook](https://www.facebook.com/BitSongOfficial)
- [Twitter](https://twitter.com/BitSongOfficial)
- [Medium](https://medium.com/@BitSongOfficial)
- [Reddit](https://www.reddit.com/r/bitsong/)
- [BitcoinTalk ANN](https://bitcointalk.org/index.php?topic=2850943)
- [Linkedin](https://www.linkedin.com/company/bitsong)
- [Instagram](https://www.instagram.com/bitsong_official/)

## Licenza

MIT License

## Versioning

### SemVer

BitSong utilizza [SemVer](http://semver.org/) per determinare quando e come cambia la versione.
Secondo SemVer, qualsiasi cosa nell'API pubblica può cambiare in qualsiasi momento prima della versione 1.0.0

Per fornire una certa stabilità agli utenti BitSong in questi giorni 0.X.X, la versione MINOR viene utilizzata per segnalare le variazioni di interruzione attraverso un sottoinsieme dell'API pubblica totale. Questo sottoinsieme include tutte le interfacce esposte ad altri processi, ma non include le API Go in-process.
