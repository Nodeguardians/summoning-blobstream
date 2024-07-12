# Summoning Blobstream

To follow this workshop there are a few requirements :

- A github account
- Go v1.22.X
- Node >= 18
- Some Go and Solidity knowledge

## Creating an NG Account

Go to the [Node Guardians website](https://nodeguardians.io) and simply login with github.

## Starting the Celestia Blobstream quest

Go to the [Celestia Blobstream](https://nodeguardians.io/adventure/celestia-blobstream/start) quest page and click start !

> [!IMPORTANT]  
> Make sure to be logged in when you start the quest.

You can now start reading Part 1. When done, you can move on to Part 2, after familiarizing yourself with the code you can click on `DEPLOY OBSERVATORY` at the bottom of the screen.

You will probably need to link your wallet to do so, NG will automatically send you some sepolia ETH so that you can get started !

## Posting the blob

While this step is necessary for the quest, we abstracted it for the workshop : the comet blob was already posted to Celestia Mocha testnet in [this](https://testnet.celestia.explorers.guru/transaction/5EC9E71F71C40A16B6F1E71CE4E8ADA205EE07072052DF25E0A44A1ED21DF17A?height=2241721) transaction.

> [!TIP]
> For your convenience, you will find the data in the `comet.json` file present in this very same gist.

## Proving inclusion of the blob

To prove the inclusion of the blob in a celestia block we will need to generate a proof. We provide you with the necessary code in the `go` folder, this code is heavily based on this [celestia tutorial](https://docs.celestia.org/developers/blobstream-proof-queries#full-example-of-proving-that-a-celestia-block-was-committed-to-by-blobstream-x-contract).

Everything should be pre-configured, you simply need to run :

```sh
cd go
go get
cd ..
```

Followed by :

```sh
go run code/get_proof.go code/utils.go
```

And you should see the output listing the following :

- data
- shareProofs
- Namespace
- RowRoots
- RowProofs
- AttestationProofs

## Sending the proof to the blobstream contract

First let's move into the `hardhat` folder :

```sh
cd hardhat
```

And install some dependencies :

```sh
npm install
```

Next we need to compile some solidity files with :

```sh
npx hardhat compile
```

We can now run the `prove-comet` hardhat task that will communicate with your NG instance and solve it using the proof you generated earlier.

To do so, run the following command :

```sh
npx hardhat prove-comet --observatory YOUR_INSTANCE_ADDRESS --path ../../data/proof.json --network sepolia
```

> [!IMPORTANT]  
> Replace `YOUR_INSTANCE_ADDRESS` by the address displayed on Node Guardians

And hopefully you should see the following output :

```sh
Comet is proven!
```

> [!TIP]
> If something goes wrong at this step just let us know and we'll help you debug !
