package wallet

import (
  "log"

  "github.com/mr-tron/base58"
)

// encode and decode function for the base58 algorithm
func base58Encode(input []byte) []byte {
  encode := base58.Encode(input)

  return []byte(encode)
}

func base58Decode(input []byte) []byte {
  decode, err := base58.Decode(string(input[:])) // base58.Decode() only accepts a string value
  if err != nil {
    log.Panic(err)
  }
  return decode
}
