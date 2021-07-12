package wallet

import (
  "bytes"
  "crypto/elliptic"
  "encoding/gob"
  "fmt"
  "io/ioutil"
  "log"
  "os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
  Wallets map[string]*Wallet
}

// method to save a file
func (ws *Wallets) SaveFile() {
  var content bytes.Buffer

  gob.Register(elliptic.P256())

  encoder := gob.NewEncoder(&content)
  err := encoder.Encode(ws)
  if err != nil {
    log.Panic(err)
  }

  err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
  if err != nil {
    log.Panic(err)
  }
}

// method to load a file
func (ws *Wallets) LoadFile() error {
  if _, err := os.Stat(walletFile); os.IsNotExist(err) {
    return err
  }

  var wallets Wallets

  fileContent, err := ioutil.ReadFile(walletFile)

  gob.Register(elliptic.P256())
  decoder := gob.NewDecoder(bytes.NewReader(fileContent))
  err = decoder.Decode(&wallets)
  if err != nil {
    return err
  }

  ws.Wallets = wallets.Wallets

  return nil
}

// create a wallet and add it to the list of wallets
func CreateWallets() (*Wallets, error) {
  wallets := Wallets{}
  wallets.Wallets = make(map[string]*Wallet)
  err := wallets.LoadFile()

  return &wallets, err
}

func (ws *Wallets) AddWallet() string {
  wallet := MakeWallet()
  address := fmt.Sprintf("%s", wallet.Address())

  ws.Wallets[address] = wallet

  return address
}

// method to look up a wallet by its address
func (ws Wallets) GetWallet(address string) Wallet {
  return *ws.Wallets[address]
}

// method to list all of the wallets' addresses
func (ws *Wallets) GetAllAddresses() []string {
  var addresses []string

  for address := range ws.Wallets {
    addresses = append(addresses, address)
  }
  return addresses
}
