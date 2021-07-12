# go-blockchain

This is a basic blockchain project created using Golang. It uses a simple PoW mechanism. Currently, the transactions support the receiving and sending of coins. There is also a basic wallet functionality, which hasn't yet been integrated with the transactions. 

Open the repository in a command line interface and run the following commands:

### go run main.go getbalance -address ADDRESS
  to get the balance for a particular address
  
### go run main.go createblockchain -address ADDRESS
  to create a blockchain and reward the creator with a mining fee

### go run main.go printchain
  to print all the blocks in the chain

### go run main.go send -from FROM -to TO -amount AMOUNT
  to send a specified amount of coins from one address to another
  
### go run main.go createwallet
  to create a new wallet
  
### go run main.go listaddresses
  to list all the addresses in the wallet file
  
#### All blockchain data will be stored in the tmp/blocks directory, while wallet data can be found in tmp/wallets.data
