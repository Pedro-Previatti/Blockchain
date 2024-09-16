package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// This struct represents a Block in the Blockchain
type Block struct {
	Index        int    `json:"index"`        // Position the Block is on the Blockchain
	Timestamp    string `json:"Timestamp"`    // Time when the block was created
	Data         struct {
		Cpf  string `json:"cpf"`	  // The identifier of the voter
		Vote int    `json:"vote"`	  // The identifier of the vote
	} `json:"data"`         		  // Data stored on the block
	PreviousHash string `json:"previousHash"` // Hash of the previous block
	Hash         string `json:"hash"`         // Hash of the current block
}

// This struct represents a chain of blocks
type Blockchain struct {
	Chain []Block // The slice of blocks
}

// Function to generate a SHA256 hash of the block
func (b *Block) generateHash() string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d%s%s%s%d", b.Index, b.PreviousHash, b.Timestamp, b.Data.Cpf, b.Data.Vote)))
	return hex.EncodeToString(hash.Sum(nil))
}

// NewBlock creates a Block
func NewBlock(index int, cpf string, vote int, previousHash string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().String(),
		PreviousHash: previousHash,
	}
	block.Data.Cpf = cpf
	block.Data.Vote = vote
	block.Hash = block.generateHash()
	return block
}

// createGenesis creates the first block in the blockchain
func (bc *Blockchain) createGenesis() Block {
	return *NewBlock(0, "000.000.000-00", 0, "0")
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{Chain: []Block{}}
	bc.Chain = append(bc.Chain, bc.createGenesis())
	return bc
}

// function to retrieve latest block in the chain
func (bc *Blockchain) getLatest() Block {
	return bc.Chain[len(bc.Chain)-1]
}

// function to add a new block to the chain
func (bc *Blockchain) add(cpf string, vote int) {
	latestBlock := bc.getLatest()
	newBlock := NewBlock(latestBlock.Index+1, cpf, vote, latestBlock.Hash)
	bc.Chain = append(bc.Chain, *newBlock)
}

// function loops through the chain to check if is valid or not
func (bc *Blockchain) isValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		current := bc.Chain[i]
		previous := bc.Chain[i-1]

		// generate the current block's hash and compare
		if current.Hash != current.generateHash() {
			return false
		}

		// check if previous hash matches
		if current.PreviousHash != previous.Hash {
			return false
		}
	}
	return true
}

// turn to json
func (bc *Blockchain) toJSON() (string, error) {
	bytes, err := json.MarshalIndent(bc.Chain, "", "	")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	blockchain := NewBlockchain()

	blockchain.add("000.000.000-00", 22)
	blockchain.add("000.000.000-01", 13)

	jsonBlockchain, err := blockchain.toJSON()
	if err != nil {
		fmt.Println("Error converting blockchain to JSON:", err)
		return
	}
	fmt.Println(jsonBlockchain)

	fmt.Printf("Valid Blockchain: %v\n", blockchain.isValid())
}
