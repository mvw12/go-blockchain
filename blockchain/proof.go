package blockchain

import (
  "bytes"
  "crypto/sha256"
  "encoding/binary"
  "fmt"
  "log"
  "math"
  "math/big"
)

const Difficulty = 12

type ProofOfWork struct {
  Block *Block
  Target *big.Int
}

// initialise PoW struct
func NewProofOfWork(b *Block) *ProofOfWork {
  target := big.NewInt(1)
  target.Lsh(target, uint(256 - Difficulty))

  pow := &ProofOfWork{b, target}

  return pow
}

// utility function to convert int to []byte
func ToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  err := binary.Write(buff, binary.BigEndian, num)
  if err != nil {
    log.Panic(err)
  }
  return buff.Bytes()
}

// method to initialise data of the block
func (pow *ProofOfWork) InitData(nonce int) []byte {
  data := bytes.Join(
    [][]byte{
      pow.Block.PrevHash,
      pow.Block.HashTransactions(),
      ToHex(int64(nonce)),
      ToHex(int64(Difficulty)),
    },
    []byte{},
  )
  return data
}

// running the PoW algorithm
func (pow *ProofOfWork) Run() (int, []byte) {
  var intHash big.Int
  var hash [32]byte

  nonce := 0

  for nonce < math.MaxInt64 {
    data := pow.InitData(nonce)
    hash = sha256.Sum256(data)

    fmt.Printf("\r%x", hash)
    intHash.SetBytes(hash[:])

    if intHash.Cmp(pow.Target) == -1 {
      break
    } else {
      nonce++
    }
  }
  fmt.Println()

  return nonce, hash[:]
}

// method to validate the hash by comparing it with the target from PoW
func (pow *ProofOfWork) Validate() bool {
  var intHash big.Int

  data := pow.InitData(pow.Block.Nonce)

  hash := sha256.Sum256(data)
  intHash.SetBytes(hash[:])

  return intHash.Cmp(pow.Target) == -1
}
