package blockchain

import (
  "encoding/hex"
  "fmt"
  "github.com/dgraph-io/badger" // database
  "os"
  "runtime"
)

const (
   dbPath = "./tmp/blocks"
   dbFile = "./tmp/blocks/MANIFEST" // used to verify that blockchain exists
   genesisData = "First Transaction from Genesis" // data for genesis block
 )

type BlockChain struct {
  LastHash []byte
  Database *badger.DB
}

// function to create blockchain when it doesn't exist
func InitBlockChain(address string) *BlockChain {
  var lastHash []byte

  if DBExists() {
    fmt.Println("blockchain already exists")
    runtime.Goexit()
  }

  opts := badger.DefaultOptions(dbPath)
  db, err := badger.Open(opts)
  Handle(err)

  err = db.Update(func(txn *badger.Txn) error { // "lh" = last hash key
      fmt.Println("No existing blockchain found")
      cbtx := CoinbaseTx(address, genesisData)
      genesis := Genesis(cbtx) // create the genesis block since blockchain does not exist
      fmt.Println("Genesis Created")
      err = txn.Set(genesis.Hash, genesis.Serialize())
      Handle(err)
      err = txn.Set([]byte("lh"), genesis.Hash) // initialise the "lh" key for next block

      lastHash = genesis.Hash

      return err
    })

    Handle(err)

    blockchain := BlockChain{lastHash, db}
    return &blockchain
}

// blockchain already exists
func ContinueBlockChain(address string) *BlockChain {
  if !DBExists() {
    fmt.Println("No blockchain found, please create one first")
    runtime.Goexit()
  }

  var lastHash []byte

  opts := badger.DefaultOptions(dbPath)
  db, err := badger.Open(opts)
  Handle(err)

  err = db.Update(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("lh"))
    Handle(err)
    err = item.Value(func(val []byte) error {
      lastHash = val
      return nil
    })
    Handle(err)
    return err
  })
  Handle(err)

  chain := BlockChain{lastHash, db}
  return &chain
}

// utility function to check if the database has been initialised
func DBExists() bool {
  if _, err := os.Stat(dbFile); os.IsNotExist(err) {
    return false
  }
  return true
}

// add a block to the chain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
  var lastHash []byte

  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("lh"))
    Handle(err)
    err = item.Value(func(val []byte) error {
      lastHash = val
      return nil
    })
    Handle(err)
    return err
  })
  Handle(err)

  newBlock := CreateBlock(transactions, lastHash)

  err = chain.Database.Update(func(transaction *badger.Txn) error {
    err := transaction.Set(newBlock.Hash, newBlock.Serialize())
    Handle(err)
    err = transaction.Set([]byte("lh"), newBlock.Hash)

    chain.LastHash = newBlock.Hash
    return err
  })
  Handle(err)
}

type BlockChainIterator struct {
  CurrentHash []byte
  Database *badger.DB
}

// utility function to convert BlockChain to BlockChainIterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
  iterator := BlockChainIterator{chain.LastHash, chain.Database}

  return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {
  var block *Block

  err := iterator.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(iterator.CurrentHash)
    Handle(err)

    err = item.Value(func(val []byte) error {
      block = Deserialize(val)
      return nil
    })
    Handle(err)
    return err
  })
  Handle(err)

  iterator.CurrentHash = block.PrevHash

  return block
}

// function to find unspent transactions - helps in finding balance
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
  var unspentTxs []Transaction

  spentTXNs := make(map[string][]int)

  iter := chain.Iterator()

  for {
    block := iter.Next()

    for _, tx := range block.Transactions {
      txID := hex.EncodeToString(tx.ID)

    Outputs:
      for outIdx, out := range tx.Outputs {
        if spentTXNs[txID] != nil {
          for _, spentOut := range spentTXNs[txID] {
            if spentOut == outIdx {
              continue Outputs
            }
          }
        }
        if out.CanBeUnlocked(address) {
          unspentTxs = append(unspentTxs, *tx)
        }
      }
      if tx.IsCoinbase() == false {
        for _, in := range tx.Inputs {
          if in.CanUnlock(address) {
            inTxID := hex.EncodeToString(in.ID)
            spentTXNs[inTxID] = append(spentTXNs[inTxID], in.Out)
          }
        }
      }
      if len(block.PrevHash) == 0 {
        break
      }
    }
    return unspentTxs
  }
}

// function to find the unspent transaction outputs
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
  var UTXOs []TxOutput
  unspentTransactions := chain.FindUnspentTransactions(address)
  for _, tx := range unspentTransactions {
    for _, out := range tx.Outputs {
      if out.CanBeUnlocked(address) {
        UTXOs = append(UTXOs, out)
      }
    }
  }
  return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
  unspentOuts := make(map[string][]int)
  unspentTxs := chain.FindUnspentTransactions(address)
  accumulated := 0

Work:
  for _, tx := range unspentTxs {
    txID := hex.EncodeToString(tx.ID)
    for outIdx, out := range tx.Outputs {
      if out.CanBeUnlocked(address) && accumulated < amount {
        accumulated += out.Value
        unspentOuts[txID] = append(unspentOuts[txID], outIdx)

        if accumulated >= amount {
          break Work
        }
      }
    }
  }
  return accumulated, unspentOuts
}
