package blockchain

import (
  "bytes"
  "crypto/sha256"
  "encoding/gob"
  "encoding/hex"
  "fmt"
  "log"
)

const reward = 100

type Transaction struct {
  ID []byte
  Inputs []TxInput
  Outputs []TxOutput
}

// to initiate the first (coinbase) transaction
func CoinbaseTx(toAddress, data string) *Transaction {
  if data == "" {
    data = fmt.Sprintf("Coins to %s", toAddress)
  }

  txIn := TxInput{[]byte{}, -1, data}

  txOut := TxOutput{reward, toAddress}

  tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

  return &tx
}

// utility function to encode and hash Transaction ID
func (tx *Transaction) SetID() {
  var encoded bytes.Buffer
  var hash [32]byte

  encoder := gob.NewEncoder(&encoded)
  err := encoder.Encode(tx)
  Handle(err)

  hash = sha256.Sum256(encoded.Bytes())
  tx.ID = hash[:]
}

//function to check if a tx was the coinbase tx
func (tx *Transaction) IsCoinbase() bool {
  return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

/* function to find spendable outputs, check if balance is sufficient before
spending, make inputs that point to the outputs that are being spent, make new
outputs for the leftover money, initialise a fresh transaction with all the new
inputs and outputs, set a new ID and return it */
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
  var inputs []TxInput
  var outputs []TxOutput

  acc, validOutputs := chain.FindSpendableOutputs(from, amount) // finding the spendable outputs

  // checking if balance is sufficient
  if acc < amount {
    log.Panic("Error: Not enough funds!")
  }

  // making inputs which point to the outputs that are being spent
  for txid, outs := range validOutputs {
    txID, err := hex.DecodeString(txid)
    Handle(err)

    for _, out := range outs {
      input := TxInput{txID, out, from}
      inputs = append(inputs, input)
    }
  }

  outputs = append(outputs, TxOutput{amount, to})

  // new outputs for the leftover money 'acc - amount'
  if acc > amount {
    outputs = append(outputs, TxOutput{acc - amount, from})
  }

  //init a new transaction with all the new inputs and outputs made so far
  tx := Transaction{nil, inputs, outputs}

  // set and return a new ID for the transaction
  tx.SetID()

  return &tx
}
