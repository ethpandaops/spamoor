package spamoortypes

// ClientSelectionMode defines how clients are selected from the pool.
type ClientSelectionMode uint8

var (
	// SelectClientByIndex selects a client by index (modulo pool size).
	SelectClientByIndex ClientSelectionMode = 0
	// SelectClientRandom selects a random client from the pool.
	SelectClientRandom ClientSelectionMode = 1
	// SelectClientRoundRobin selects clients in round-robin fashion.
	SelectClientRoundRobin ClientSelectionMode = 2
)

// ClientPool defines the interface for managing a pool of Ethereum RPC clients with health
// monitoring and selection strategies. It automatically monitors client health by checking
// block heights and maintains a list of "good" clients that are within 2 blocks of the
// highest observed block height.
type ClientPool interface {
	// Client Management
	PrepareClients() error
	GetClient(mode ClientSelectionMode, input int, group string) Client
	GetAllClients() []Client
	GetAllGoodClients() []Client
}
