package main

import (
  "flag"
  "fmt"
  "os"
  "runtime"
  "strconv"

  "github.com/mvw12/go-blockchain/blockchain"
)

type CommandLine struct {
  blockchain *blockchain.BlockChain
}

// displays the options available to the user
func (cli *CommandLine) printUsage() {
  fmt.Println("Usage: ")
  fmt.Println(" add -block <BLOCK_DATA> - add a block to the chain")
  fmt.Println(" print - prints the blocks in the chain")
}

// validates the input given to cli
func (cli *CommandLine) validateArgs() {
  if len(os.Args) < 2 {
    cli.printUsage()

    runtime.Goexit()
  }
}

// add blocks to the chain via cli
func (cli *CommandLine) addBlock(data string) {
  cli.blockchain.AddBlock(data)
  fmt.Println("Added Block!")
}

// displays the entire contents of the blockchain
func (cli *CommandLine) printChain() {
  iterator := cli.blockchain.Iterator()

  for {
    block := iterator.Next()
    fmt.Printf("Previous hash: %x\n", block.PrevHash)
    fmt.Printf("Data: %s\n", block.Data)
    fmt.Printf("Hash: %x\n", block.Hash)
    pow := blockchain.NewProofOfWork(block)
    fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
    fmt.Println()
    if len(block.PrevHash) == 0 { // happens when we reach genesis block
      break
    }
  }
}

// to start up the command line
func (cli *CommandLine) run() {
  cli.validateArgs()

  addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
  printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
  addBlockData := addBlockCmd.String("block", "", "Block data")

  switch os.Args[1] {

  case "add":
    err := addBlockCmd.Parse(os.Args[2:])
    blockchain.Handle(err)

  case "print":
    err := printChainCmd.Parse(os.Args[2:])
    blockchain.Handle(err)

  default:
    cli.printUsage()
    runtime.Goexit()
  }

  if addBlockCmd.Parsed() { // returns true if object has been called
    if *addBlockData == "" {
      addBlockCmd.Usage()
      runtime.Goexit()
    }
    cli.addBlock(*addBlockData)
  }
  if printChainCmd.Parsed() {
    cli.printChain()
  }
}

func main() {
  defer os.Exit(0)

  chain := blockchain.InitBlockChain()
  defer chain.Database.Close()

  cli := CommandLine{chain}

  cli.run()
}
