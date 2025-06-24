package taskrunner

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/spamoor"
)

// TaskExecutionContext contains context information for task execution
type TaskExecutionContext struct {
	BaseFee float64         // Max fee per gas in gwei
	TipFee  float64         // Max tip per gas in gwei
	TxPool  *spamoor.TxPool // For fee calculation
}

// Task represents a single executable task in the TaskRunner scenario
type Task interface {
	// GetType returns the task type identifier
	GetType() string
	// GetName returns the task name (optional, for registry)
	GetName() string
	// Validate checks if the task configuration is valid
	Validate() error
	// BuildTransaction creates a transaction for this task
	BuildTransaction(ctx context.Context, wallet *spamoor.Wallet, registry *ContractRegistry, execCtx *TaskExecutionContext) (*types.Transaction, error)
}

// TaskConfig represents the generic configuration for any task
type TaskConfig struct {
	Type string                 `yaml:"type" json:"type"`
	Name string                 `yaml:"name" json:"name"`
	Data map[string]interface{} `yaml:"data" json:"data"`
}

// BaseTask provides common functionality for all task types
type BaseTask struct {
	Type string
	Name string
}

// GetType returns the task type
func (t *BaseTask) GetType() string {
	return t.Type
}

// GetName returns the task name
func (t *BaseTask) GetName() string {
	return t.Name
}

// TaskFactory is a function that creates a task from configuration
type TaskFactory func(name string, data map[string]interface{}) (Task, error)

// taskRegistry stores registered task types
var taskRegistry = map[string]TaskFactory{
	"deploy": NewDeployTask,
	"call":   NewCallTask,
}

// CreateTask creates a task instance from configuration
func CreateTask(config *TaskConfig) (Task, error) {
	factory, exists := taskRegistry[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown task type: %s", config.Type)
	}

	return factory(config.Name, config.Data)
}

// parseValue parses a configuration value with proper type handling
func parseValue(key string, data map[string]interface{}, target interface{}) error {
	value, exists := data[key]
	if !exists {
		return nil // Optional field
	}

	switch t := target.(type) {
	case *string:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("field %s must be a string", key)
		}
		*t = str

	case *uint64:
		switch v := value.(type) {
		case float64:
			*t = uint64(v)
		case int:
			*t = uint64(v)
		case int64:
			*t = uint64(v)
		case uint64:
			*t = v
		default:
			return fmt.Errorf("field %s must be a number", key)
		}

	case **big.Int:
		switch v := value.(type) {
		case string:
			bigInt, ok := new(big.Int).SetString(v, 0)
			if !ok {
				return fmt.Errorf("field %s: invalid big integer value", key)
			}
			*t = bigInt
		case float64:
			*t = big.NewInt(int64(v))
		case int:
			*t = big.NewInt(int64(v))
		case int64:
			*t = big.NewInt(v)
		default:
			return fmt.Errorf("field %s must be a number or string", key)
		}

	case *[]interface{}:
		arr, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("field %s must be an array", key)
		}
		*t = arr

	default:
		return fmt.Errorf("unsupported target type for field %s", key)
	}

	return nil
}

// getRequiredString gets a required string field from configuration
func getRequiredString(key string, data map[string]interface{}) (string, error) {
	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("required field '%s' not found", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("field '%s' must be a string", key)
	}

	return str, nil
}
