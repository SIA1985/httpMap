package storage

import "blockchain/internal/block"

type Storage struct {
	data map[[32]byte]block.Block
}
