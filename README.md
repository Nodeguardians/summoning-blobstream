# Summoning Blobstream

To follow this workshop there are a few requirements :

- A github account
- Go v1.22.X
- Some Go and Solidity knowledge

## Creating an NG Account

Go to the [Node Guardians website](https://nodeguardians.io) and simply login with github.

## Starting the Celestia Blobstream quest

Go to the [Celestia Blobstream](https://nodeguardians.io/adventure/celestia-blobstream/start) quest page and click start !

> [!IMPORTANT]  
> Make sure to be logged in when you start the quest.

## Posting the blob

While this step is necessary for the quest, we abstracted it for the workshop : the comet blob was already posted to Celestia Mocha testnet in [this](https://testnet.celestia.explorers.guru/transaction/5EC9E71F71C40A16B6F1E71CE4E8ADA205EE07072052DF25E0A44A1ED21DF17A?height=2241721) transaction.

> [!TIP]
> For your convenience, you will find the data in the `comet.json` file present in this very same gist.

## Proving inclusion of the blob

To prove the inclusion of the blob in a celestia block we will need to generate a proof. We provide you with the necessary code in the `go` folder, this code is heavily based on this [celestia tutorial](https://docs.celestia.org/developers/blobstream-proof-queries#full-example-of-proving-that-a-celestia-block-was-committed-to-by-blobstream-x-contract).

Everything should be pre-configured, you simply need to run :

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
