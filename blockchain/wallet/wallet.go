package wallet

import (
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/rand"
  "crypto/sha256"
  "log"

  "golang.org/x/crypto/ripemd160"
)

const (
  checksumLength = 4
  version = byte(0x00) // hexadecimal representation of 0
)

type Wallet struct {
  PrivateKey ecdsa.PrivateKey
  PublicKey []byte
}

// function to generate a public and private key
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
  curve := elliptic.P256()

  private, err := ecdsa.GenerateKey(curve, rand.Reader)
  if err != nil {
    log.Panic(err)
  }

  pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

  return *private, pub
}

// function to hash the public key using sha256, followed by ripemd160 to generate PublicHash
func PublicKeyHash(publicKey []byte) []byte {
  hashedPublicKey := sha256.Sum256(publicKey)

  hasher := ripemd160.New()
  _, err := hasher.Write(hashedPublicKey[:])
  if err != nil {
    log.Panic(err)
  }
  publicRipeMd := hasher.Sum(nil)

  return publicRipeMd
}

// function to hash the versioned hash twice using sha256 - first 4 bytes are the checksum
func Checksum(ripeMdHash []byte) []byte {
  firstHash := sha256.Sum256(ripeMdHash)
  secondHash := sha256.Sum256(firstHash[:])

  return secondHash[:checksumLength]
}

// function to generate the wallet address
func (w *Wallet) Address() []byte {
  pubHash := PublicKeyHash(w.PublicKey)
  versionedHash := append([]byte{version}, pubHash...) // obtained by adding version at the start of PublicHash
  checksum := Checksum(versionedHash)
  finalHash := append(versionedHash, checksum...)
  address := base58Encode(finalHash)
  return address
}

// function to initialise a wallet
func MakeWallet() *Wallet {
  privateKey, publicKey := NewKeyPair()
  wallet := Wallet{privateKey, publicKey}
  return &wallet
}
