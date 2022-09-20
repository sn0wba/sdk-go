package chain

import (
	"bytes"
	"crypto/rand"
	"io"
	"log"
	"os"
	"path/filepath"

	cosmcrypto "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

const defaultKeyringKeyName = "validator"

var emptyCosmosAddress = cosmtypes.AccAddress{}

func InitCosmosKeyring(
	cosmosKeyringDir string,
	cosmosKeyringAppName string,
	cosmosKeyringBackend string,
	cosmosKeyFrom string,
	cosmosKeyPassphrase string,
	cosmosMnemonic string,
	cosmosUseLedger bool,
) (cosmtypes.AccAddress, keyring.Keyring, error) {
	switch {
	case len(cosmosMnemonic) > 0:
		var keyName string

		if len(keyName) == 0 {
			keyName = defaultKeyringKeyName
		}

		return KeyringForMnemonic(keyName, cosmosMnemonic)

	case len(cosmosKeyFrom) > 0:
		var fromIsAddress bool
		addressFrom, err := cosmtypes.AccAddressFromBech32(cosmosKeyFrom)
		if err == nil {
			fromIsAddress = true
		}

		var passReader io.Reader = os.Stdin
		if len(cosmosKeyPassphrase) > 0 {
			passReader = newPassReader(cosmosKeyPassphrase)
		}

		var absoluteKeyringDir string
		if filepath.IsAbs(cosmosKeyringDir) {
			absoluteKeyringDir = cosmosKeyringDir
		} else {
			absoluteKeyringDir, _ = filepath.Abs(cosmosKeyringDir)
		}

		kb, err := keyring.New(
			cosmosKeyringAppName,
			cosmosKeyringBackend,
			absoluteKeyringDir,
			passReader,
		)
		if err != nil {
			err = errors.Wrap(err, "failed to init keyring")
			return emptyCosmosAddress, nil, err
		}

		var keyInfo keyring.Info
		if fromIsAddress {
			if keyInfo, err = kb.KeyByAddress(addressFrom); err != nil {
				err = errors.Wrapf(err, "couldn't find an entry for the key %s in keybase", addressFrom.String())
				return emptyCosmosAddress, nil, err
			}
		} else {
			if keyInfo, err = kb.Key(cosmosKeyFrom); err != nil {
				err = errors.Wrapf(err, "could not find an entry for the key '%s' in keybase", cosmosKeyFrom)
				return emptyCosmosAddress, nil, err
			}
		}

		switch keyType := keyInfo.GetType(); keyType {
		case keyring.TypeLocal:
			// kb has a key and it's totally usable
			return keyInfo.GetAddress(), kb, nil
		case keyring.TypeLedger:
			// the kb stores references to ledger keys, so we must explicitly
			// check that. kb doesn't know how to scan HD keys - they must be added manually before
			if cosmosUseLedger {
				return keyInfo.GetAddress(), kb, nil
			}
			err := errors.Errorf("'%s' key is a ledger reference, enable ledger option", keyInfo.GetName())
			return emptyCosmosAddress, nil, err
		case keyring.TypeOffline:
			err := errors.Errorf("'%s' key is an offline key, not supported yet", keyInfo.GetName())
			return emptyCosmosAddress, nil, err
		case keyring.TypeMulti:
			err := errors.Errorf("'%s' key is an multisig key, not supported yet", keyInfo.GetName())
			return emptyCosmosAddress, nil, err
		default:
			err := errors.Errorf("'%s' key  has unsupported type: %s", keyInfo.GetName(), keyType)
			return emptyCosmosAddress, nil, err
		}

	default:
		err := errors.New("insufficient cosmos key details provided")
		return emptyCosmosAddress, nil, err
	}
}

func newPassReader(pass string) io.Reader {
	return &passReader{
		pass: pass,
		buf:  new(bytes.Buffer),
	}
}

type passReader struct {
	pass string
	buf  *bytes.Buffer
}

var _ io.Reader = &passReader{}

func (r *passReader) Read(p []byte) (n int, err error) {
	n, err = r.buf.Read(p)
	if err == io.EOF || n == 0 {
		r.buf.WriteString(r.pass + "\n")

		n, err = r.buf.Read(p)
	}

	return
}

// KeyringForPrivKey creates a temporary in-mem keyring for a PrivKey.
// Allows to init Context when the key has been provided in plaintext and parsed.
func KeyringForPrivKey(name string, privKey cryptotypes.PrivKey) (keyring.Keyring, error) {
	kb := keyring.NewInMemory()
	tmpPhrase := randPhrase(64)
	armored := cosmcrypto.EncryptArmorPrivKey(privKey, tmpPhrase, privKey.Type())
	err := kb.ImportPrivKey(name, armored, tmpPhrase)
	if err != nil {
		err = errors.Wrap(err, "failed to import privkey")
		return nil, err
	}

	return kb, nil
}

func KeyringForMnemonic(name, mnemonic string) (cosmtypes.AccAddress, keyring.Keyring, error) {
	kb := keyring.NewInMemory()
	info, err := kb.NewAccount(name, mnemonic, "", "m/44'/118'/0'/0/0", hd.Secp256k1)
	if err != nil {
		err = errors.Wrap(err, "failed to import privkey")
		return emptyCosmosAddress, nil, err
	}

	return info.GetAddress(), kb, nil
}

func randPhrase(size int) string {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	orPanic(err)

	return string(buf)
}

func orPanic(err error) {
	if err != nil {
		log.Panicln()
	}
}
