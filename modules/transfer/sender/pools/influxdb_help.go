package pools

import (
	"errors"
	"strings"

	cpool "github.com/Cepave/open-falcon-backend/modules/transfer/sender/conn_pool"
)

var (
	errInvalidDSNUnescaped = errors.New("Invalid DSN: Did you forget to escape a param value?")
	errInvalidDSNAddr      = errors.New("Invalid DSN: Network Address not terminated (missing closing brace)")
	errInvalidDSNNoSlash   = errors.New("Invalid DSN: Missing the slash separating the database name")
)

func InfuxdbParseDSN(dsn string) (cfg *cpool.InfluxdbConnection, err error) {
	// New config
	cfg = &cpool.InfluxdbConnection{}

	// [username[:password]@][protocol[(address)]]/dbname
	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								cfg.Password = dsn[k+1 : j]
								break
							}
						}
						cfg.Username = dsn[:k]

						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return nil, errInvalidDSNUnescaped
							}
							return nil, errInvalidDSNAddr
						}
						cfg.Address = dsn[k+1 : i-1]
						break
					}
				}
				cfg.Protocol = dsn[j+1 : k]
			}

			// /dbname
			cfg.DBName = dsn[i+1 : len(dsn)]

			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}

	// Set default network if empty
	if cfg.Protocol == "" {
		cfg.Protocol = "tcp"
	}

	// Set default address if empty
	if cfg.Address == "" {
		return nil, errors.New("Empty DSN address")
	}

	return
}
