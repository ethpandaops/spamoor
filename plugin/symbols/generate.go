package symbols

//go:generate rm -f github_*.go
//go:generate rm -f gopkg_*.go
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethpandaops/spamoor/spamoor
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethpandaops/spamoor/scenario
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethpandaops/spamoor/txbuilder
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethpandaops/spamoor/utils
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/sirupsen/logrus
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/spf13/pflag
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/holiman/uint256
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/common
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/core/types
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/accounts/abi
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/accounts/abi/bind
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/accounts/abi/bind/v2
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/crypto
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/ethereum/go-ethereum/event
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gopkg.in/yaml.v3
