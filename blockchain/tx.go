package blockchain

type TxOutput struct {
  Value int // amount of coins in a tx
  PubKey string // sender of coins is identified by their Public Key
}

type TxInput struct {
  ID []byte // used to identify a Tx
  Out int // index of the specific output within a Tx
  Sig string // script that adds data to PubKey
}

// utility functions to validate Sig and PubKey
func (in *TxInput) CanUnlock(data string) bool {
  return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
  return out.PubKey == data
}
