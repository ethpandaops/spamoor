package scenario

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ConfigValidator provides validation for scenario configurations
type ConfigValidator struct {
	scenarioName string
	validFields  map[string]FieldInfo
	logger       logrus.FieldLogger
}

// FieldInfo contains metadata about a valid configuration field
type FieldInfo struct {
	Type        reflect.Type
	Required    bool
	DefaultValue interface{}
	Description string
}

// ValidationResult contains the results of configuration validation
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// NewConfigValidator creates a new configuration validator for a specific scenario
func NewConfigValidator(scenarioName string, validFields map[string]FieldInfo, logger logrus.FieldLogger) *ConfigValidator {
	return &ConfigValidator{
		scenarioName: scenarioName,
		validFields:  validFields,
		logger:       logger,
	}
}

// ValidateConfig validates a YAML configuration against the scenario's valid fields
func (cv *ConfigValidator) ValidateConfig(configYAML string) *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	if configYAML == "" {
		return result // Empty config is valid
	}

	// Parse the YAML configuration
	var config map[string]interface{}
	if err := yaml.Unmarshal([]byte(configYAML), &config); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML configuration: %v", err))
		return result
	}

	// Collect all invalid fields
	var invalidFields []string
	var suggestions []string

	// Check each field in the provided configuration
	for fieldName, fieldValue := range config {
		if fieldInfo, exists := cv.validFields[fieldName]; exists {
			// Validate field type and value
			if validationErr := cv.validateFieldValue(fieldName, fieldValue, fieldInfo); validationErr != nil {
				result.Valid = false
				result.Errors = append(result.Errors, validationErr.Error())
			}
		} else {
			// Field doesn't exist in valid fields - this is the main issue we're solving
			invalidFields = append(invalidFields, fmt.Sprintf("'%s' (value: %v)", fieldName, fieldValue))
			
			// Check for common typos (underscore vs dash)
			if suggestion := cv.suggestCorrectField(fieldName); suggestion != "" {
				suggestions = append(suggestions, fmt.Sprintf("'%s' -> '%s'", fieldName, suggestion))
			}
		}
	}

	// Log a single comprehensive warning message if there are invalid fields
	if len(invalidFields) > 0 {
		cv.logInvalidFields(invalidFields, suggestions)
	}

	// Check all required fields are present
	for fieldName, fieldInfo := range cv.validFields {
		if fieldInfo.Required {
			if _, exists := config[fieldName]; !exists {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("Required field '%s' is missing", fieldName))
			}
		}
	}

	return result
}

// validateFieldValue validates a specific field value against its expected type
func (cv *ConfigValidator) validateFieldValue(fieldName string, value interface{}, fieldInfo FieldInfo) error {
	// Type validation would go here - for now just basic checks
	switch fieldInfo.Type.Kind() {
	case reflect.Uint64:
		if _, ok := value.(int); !ok {
			if _, ok := value.(uint64); !ok {
				return fmt.Errorf("Field '%s' must be a positive integer", fieldName)
			}
		}
	case reflect.String:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("Field '%s' must be a string", fieldName)
		}
	case reflect.Bool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("Field '%s' must be a boolean", fieldName)
		}
	}
	
	return nil
}

// logInvalidFields logs a single comprehensive warning about all invalid configuration fields
func (cv *ConfigValidator) logInvalidFields(invalidFields []string, suggestions []string) {
	message := fmt.Sprintf("Invalid configuration fields detected for scenario '%s': %s", 
		cv.scenarioName, strings.Join(invalidFields, ", "))
	
	message += ". These fields will be ignored."
	
	if len(suggestions) > 0 {
		message += fmt.Sprintf(" Possible corrections: %s", strings.Join(suggestions, ", "))
	}
	
	cv.logger.Warnf("%s", message)
}

// suggestCorrectField suggests a correct field name based on common typos
func (cv *ConfigValidator) suggestCorrectField(invalidField string) string {
	// Common transformations to check
	transformations := []func(string) string{
		func(s string) string { return strings.ReplaceAll(s, "-", "_") },  // dash to underscore
		func(s string) string { return strings.ReplaceAll(s, "_", "-") },  // underscore to dash
		strings.ToLower,  // case insensitive
	}
	
	for _, transform := range transformations {
		transformed := transform(invalidField)
		if _, exists := cv.validFields[transformed]; exists {
			return transformed
		}
	}
	
	// If no direct transformation works, find the closest match
	// This is a simple approach - could be enhanced with edit distance algorithms
	invalidLower := strings.ToLower(invalidField)
	for validField := range cv.validFields {
		validLower := strings.ToLower(validField)
		if strings.Contains(validLower, invalidLower) || strings.Contains(invalidLower, validLower) {
			return validField
		}
	}
	
	return ""
}

// GetScenarioValidFields returns the valid fields for built-in scenarios
func GetScenarioValidFields(scenarioName string) map[string]FieldInfo {
	// Get scenario-specific fields
	var scenarioFields map[string]FieldInfo
	switch scenarioName {
	case "eoatx":
		scenarioFields = getEOATxValidFields()
	case "erctx":
		scenarioFields = getERCTxValidFields()
	case "blobs":
		scenarioFields = getBlobsValidFields()
	case "blob-combined":
		scenarioFields = getBlobCombinedValidFields()
	case "blob-conflicting":
		scenarioFields = getBlobConflictingValidFields()
	case "blob-replacements":
		scenarioFields = getBlobReplacementsValidFields()
	case "deploytx":
		scenarioFields = getDeployTxValidFields()
	case "calltx":
		scenarioFields = getCallTxValidFields()
	case "gasburnertx":
		scenarioFields = getGasBurnerTxValidFields()
	case "geastx":
		scenarioFields = getGEASTxValidFields()
	case "setcodetx":
		scenarioFields = getSetCodeTxValidFields()
	case "storagespam":
		scenarioFields = getStorageSpamValidFields()
	case "factorydeploytx":
		scenarioFields = getFactoryDeployTxValidFields()
	case "deploy-destruct":
		scenarioFields = getDeployDestructValidFields()
	case "uniswap-swaps":
		scenarioFields = getUniswapSwapsValidFields()
	case "xentoken":
		scenarioFields = getXenTokenValidFields()
	case "wallets":
		scenarioFields = getWalletsValidFields()
	default:
		scenarioFields = make(map[string]FieldInfo)
	}
	
	// Add common wallet configuration fields that are valid for all scenarios
	walletFields := getCommonWalletFields()
	for key, value := range walletFields {
		scenarioFields[key] = value
	}
	
	return scenarioFields
}

// Define valid fields for each scenario type
func getEOATxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transactions to send"},
		"throughput":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transactions to send per slot"},
		"max_pending":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_limit":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for transactions"},
		"amount":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Transfer amount per transaction in gwei"},
		"data":            {Type: reflect.TypeOf(""), Required: false, Description: "Transaction call data to send"},
		"to":              {Type: reflect.TypeOf(""), Required: false, Description: "Target address to send transactions to"},
		"timeout":         {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"random_amount":   {Type: reflect.TypeOf(true), Required: false, Description: "Use random amounts for transactions"},
		"random_target":   {Type: reflect.TypeOf(true), Required: false, Description: "Use random to addresses for transactions"},
		"self_tx_only":    {Type: reflect.TypeOf(true), Required: false, Description: "Only send transactions to self"},
		"client_group":    {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use for sending transactions"},
		"log_txs":         {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getERCTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transactions to send"},
		"throughput":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transactions to send per slot"},
		"max_pending":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"amount":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Transfer amount per transaction in gwei"},
		"random_amount":   {Type: reflect.TypeOf(true), Required: false, Description: "Use random amounts for transactions"},
		"random_target":   {Type: reflect.TypeOf(true), Required: false, Description: "Use random to addresses for transactions"},
		"timeout":         {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":    {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":         {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getBlobsValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of blob transactions to send"},
		"throughput":                     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of blob transactions per slot"},
		"sidecars":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of blob sidecars per transaction"},
		"max_pending":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":                        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"blob_fee":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max blob fee in gwei"},
		"blob_v1_percent":                {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Percentage of blob transactions with v1 wrapper format"},
		"fulu_activation":                {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Unix timestamp of Fulu activation"},
		"throughput_increment_interval":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Increment throughput every interval in seconds"},
		"timeout":                        {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":                   {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":                        {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getBlobCombinedValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of blob transactions to send"},
		"throughput":                     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of blob transactions per slot"},
		"sidecars":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of blob sidecars per transaction"},
		"max_pending":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"replace":                        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of seconds to wait before replacing a transaction"},
		"max_replacements":               {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of replacement transactions"},
		"rebroadcast":                    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":                        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"blob_fee":                       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max blob fee in gwei"},
		"blob_v1_percent":                {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Percentage with v1 wrapper format"},
		"fulu_activation":                {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Unix timestamp of Fulu activation"},
		"throughput_increment_interval":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Increment throughput every interval in seconds"},
		"timeout":                        {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":                   {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":                        {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getBlobConflictingValidFields() map[string]FieldInfo {
	return getBlobCombinedValidFields() // Same fields as blob-combined
}

func getBlobReplacementsValidFields() map[string]FieldInfo {
	return getBlobCombinedValidFields() // Same fields as blob-combined
}

func getDeployTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of deployment transactions to send"},
		"throughput":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of deployment transactions per slot"},
		"max_pending":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"gas_limit":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for deployment transactions"},
		"base_fee":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"bytecodes":       {Type: reflect.TypeOf(""), Required: false, Description: "Comma-separated list of hex bytecodes to deploy"},
		"bytecodes_file":  {Type: reflect.TypeOf(""), Required: false, Description: "File with bytecodes to deploy"},
		"timeout":         {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":    {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":         {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getCallTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of call transactions to send"},
		"throughput":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of call transactions per slot"},
		"max_pending":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":             {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"deploy_gas_limit":    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for deployment transaction"},
		"gas_limit":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for call transactions"},
		"amount":              {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Transfer amount per transaction in gwei"},
		"random_amount":       {Type: reflect.TypeOf(true), Required: false, Description: "Use random amounts for transactions"},
		"random_target":       {Type: reflect.TypeOf(true), Required: false, Description: "Use random to addresses for transactions"},
		"contract_code":       {Type: reflect.TypeOf(""), Required: false, Description: "Contract code to deploy"},
		"contract_file":       {Type: reflect.TypeOf(""), Required: false, Description: "Contract file to deploy"},
		"contract_address":    {Type: reflect.TypeOf(""), Required: false, Description: "Address of already deployed contract"},
		"contract_args":       {Type: reflect.TypeOf(""), Required: false, Description: "Contract arguments for constructor"},
		"contract_addr_path":  {Type: reflect.TypeOf(""), Required: false, Description: "Path to child contract created during deployment"},
		"call_data":           {Type: reflect.TypeOf(""), Required: false, Description: "Data to pass to function call"},
		"call_abi":            {Type: reflect.TypeOf(""), Required: false, Description: "JSON ABI of the contract for function calls"},
		"call_abi_file":       {Type: reflect.TypeOf(""), Required: false, Description: "JSON ABI file of the contract for function calls"},
		"call_fn_name":        {Type: reflect.TypeOf(""), Required: false, Description: "Function name to call (requires call_abi)"},
		"call_fn_sig":         {Type: reflect.TypeOf(""), Required: false, Description: "Function signature to call (alternative to call_abi)"},
		"call_args":           {Type: reflect.TypeOf(""), Required: false, Description: "JSON array of arguments to pass to function"},
		"timeout":             {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":        {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":             {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getGasBurnerTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of gasburner transactions to send"},
		"throughput":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of gasburner transactions per slot"},
		"max_pending":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_units_to_burn":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of gas units for each tx to cost"},
		"gas_remainder":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Minimum gas units left to do another round"},
		"timeout":            {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"opcodes":            {Type: reflect.TypeOf(""), Required: false, Description: "EAS opcodes to use for burning gas in the gasburner contract"},
		"init_opcodes":       {Type: reflect.TypeOf(""), Required: false, Description: "EAS opcodes to use for the init code of the gasburner contract"},
		"client_group":       {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":            {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getGEASTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of geas transactions to send"},
		"throughput":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of geas transactions per slot"},
		"max_pending":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"amount":             {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Amount to send in geas transactions"},
		"base_fee":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_limit":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max gas limit to use in geas transactions"},
		"deploy_gas_limit":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max gas limit for deployment transaction"},
		"geas_file":          {Type: reflect.TypeOf(""), Required: false, Description: "Path to the geas file to use for execution"},
		"geas_code":          {Type: reflect.TypeOf(""), Required: false, Description: "Geas code to use for execution"},
		"client_group":       {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"timeout":            {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"log_txs":            {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getSetCodeTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transfer transactions to send"},
		"throughput":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transfer transactions per slot"},
		"max_pending":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"min_authorizations":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Minimum number of authorizations per transaction"},
		"max_authorizations":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of authorizations per transaction"},
		"max_delegators":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of random delegators to use"},
		"rebroadcast":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":             {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_limit":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for transactions"},
		"amount":              {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Transfer amount per transaction in gwei"},
		"data":                {Type: reflect.TypeOf(""), Required: false, Description: "Transaction call data to send"},
		"code_addr":           {Type: reflect.TypeOf(""), Required: false, Description: "Code delegation target address"},
		"random_amount":       {Type: reflect.TypeOf(true), Required: false, Description: "Use random amounts for transactions"},
		"random_target":       {Type: reflect.TypeOf(true), Required: false, Description: "Use random to addresses for transactions"},
		"random_code_addr":    {Type: reflect.TypeOf(true), Required: false, Description: "Use random delegation target for transactions"},
		"timeout":             {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":        {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":             {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getStorageSpamValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of gasburner transactions to send"},
		"throughput":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of gasburner transactions per slot"},
		"max_pending":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_units_to_burn":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of gas units for each tx to cost"},
		"timeout":            {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":       {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":            {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getFactoryDeployTxValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of contracts to deploy"},
		"throughput":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of deployment transactions per slot"},
		"max_pending":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":            {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":             {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_limit":           {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for transactions"},
		"factory_address":     {Type: reflect.TypeOf(""), Required: false, Description: "Address of existing CREATE2 factory"},
		"init_code":           {Type: reflect.TypeOf(""), Required: true, Description: "Hex-encoded init code of contract to deploy"},
		"start_salt":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Starting salt value for deployments"},
		"well_known_factory":  {Type: reflect.TypeOf(true), Required: false, Description: "Use well-known factory deployer wallet"},
		"timeout":             {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":        {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":             {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getDeployDestructValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transfer transactions to send"},
		"throughput":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transfer transactions per slot"},
		"max_pending":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"amount":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Transfer amount per transaction in gwei"},
		"gas_limit":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for each deployment test tx"},
		"random_amount":   {Type: reflect.TypeOf(true), Required: false, Description: "Use random amounts for transactions"},
		"timeout":         {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":    {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":         {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getUniswapSwapsValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transfer transactions to send"},
		"throughput":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transfer transactions per slot"},
		"max_pending":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":          {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"pair_count":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of uniswap pairs to deploy"},
		"min_swap_amount":  {Type: reflect.TypeOf(""), Required: false, Description: "Minimum swap amount in wei"},
		"max_swap_amount":  {Type: reflect.TypeOf(""), Required: false, Description: "Maximum swap amount in wei"},
		"buy_ratio":        {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Ratio of buy vs sell swaps 0-100"},
		"slippage":         {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Slippage tolerance in basis points"},
		"sell_threshold":   {Type: reflect.TypeOf(""), Required: false, Description: "DAI balance threshold to force sell in wei"},
		"timeout":          {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":     {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":          {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getXenTokenValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"total_count":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Total number of transfer transactions to send"},
		"throughput":    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of transfer transactions per slot"},
		"max_pending":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of pending transactions"},
		"max_wallets":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"rebroadcast":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Enable reliable rebroadcast"},
		"base_fee":      {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max fee per gas in gwei"},
		"tip_fee":       {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Max tip per gas in gwei"},
		"gas_limit":     {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Gas limit for the sybil attack transaction"},
		"xen_address":   {Type: reflect.TypeOf(""), Required: false, Description: "XEN token address"},
		"claim_term":    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "XEN claim term in days"},
		"timeout":       {Type: reflect.TypeOf(""), Required: false, Description: "Timeout duration"},
		"client_group":  {Type: reflect.TypeOf(""), Required: false, Description: "Client group to use"},
		"log_txs":       {Type: reflect.TypeOf(true), Required: false, Description: "Log all submitted transactions"},
	}
}

func getWalletsValidFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"wallets": {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Maximum number of child wallets to use"},
		"reclaim": {Type: reflect.TypeOf(true), Required: false, Description: "Reclaim funds from wallets"},
	}
}

// getCommonWalletFields returns wallet configuration fields that are valid for all scenarios
func getCommonWalletFields() map[string]FieldInfo {
	return map[string]FieldInfo{
		"seed":            {Type: reflect.TypeOf(""), Required: false, Description: "Seed for deterministic wallet generation"},
		"refill_amount":   {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Amount of ETH to fund/refill each child wallet with"},
		"refill_balance":  {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Min amount of ETH each child wallet should hold before refilling"},
		"refill_interval": {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Interval for child wallet balance check and refilling if needed (in sec)"},
		"wallet_count":    {Type: reflect.TypeOf(uint64(0)), Required: false, Description: "Number of child wallets to create"},
	}
}