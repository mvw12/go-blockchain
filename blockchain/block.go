package blockchain

import (
  "bytes"
  "encoding/gob"
  "log"
)

type Block struct {
  Hash []byte
  Data []byte
  PrevHash []byte
  Nonce int
}

// function to create a block - calls utility functions that help compute hash and nonce
func CreateBlock(Data string, PrevHash []byte) *Block {
  block := &Block{[]byte{}, []byte(Data), PrevHash, 0}

  pow := NewProofOfWork(block)
  nonce, hash := pow.Run()

  block.Hash = hash[:]
  block.Nonce = nonce

  return block
}

// create the genesis block
func Genesis() *Block {
  return CreateBlock("Genesis", []byte{})
}

// utility function to handle errors
func Handle(err error) {
  if err != nil {
    log.Panic(err)
  }
}

// method to encode the Block struct into bytes
func (b *Block) Serialize() []byte {
  var res bytes.Buffer
  encoder := gob.NewEncoder(&res)

  err := encoder.Encode(b)

  Handle(err)

  return res.Bytes()
}

// function to decode the bytes back into Block struct
func Deserialize(data []byte) *Block {
  var block Block

  decoder := gob.NewDecoder(bytes.NewReader(data))

  err := decoder.Decode(&block)

  Handle(err)

  return &block
}
