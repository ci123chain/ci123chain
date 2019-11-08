package helper

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-isatty"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"os"
	"strings"
)

const (
	MinPassLength = 4

	FlagHeight = "height"
	FlagHomeDir = "home"
	FlagVerbose = "verbose"
	FlagNode = "node"
	FlagAddress = "address"
	FlagPassword = "password"
	//FlagWithCrypto 	   = "cryptosuit"
)

// Allows for reading prompts for stdin
func BufferStdin() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func GetPasswordFromStd() (string, error) {
	var err error
	buf := BufferStdin()
	pass, err := GetCheckPassword("Enter a passphrase for your types:", "Repeat the passphrase:", buf)
	if err != nil {
		return "", err
	}
	return pass, nil
}


// Prompts for a password twice to verify they match
func GetCheckPassword(prompt, prompt2 string, buf *bufio.Reader) (string, error) {
	if !inputIsTty() {
		return GetPassword(prompt, buf)
	}

	pass, err := GetPassword(prompt, buf)
	if err != nil {
		return "", err
	}
	pass2, err := GetPassword(prompt2, buf)
	if err != nil {
		return "", err
	}
	if pass != pass2 {
		return "", errors.New("Passphrases did not match")
	}
	return pass, nil
}


// Prompts for a password one-time
// Enforces minimum password length
func GetPassword(prompt string, buf *bufio.Reader) (pass string, err error) {
	if inputIsTty() {
		pass, err = speakeasy.Ask(prompt)
	} else {
		pass, err = readLineFromBuf(buf)
	}
	if err != nil {
		return "", err
	}
	if len(pass) < MinPassLength {
		return "", fmt.Errorf("Password must be at least %d characters", MinPassLength)
	}
	return pass, nil
}

// Returns true iff we have an interactive prompt
func inputIsTty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

// reads one line from stdin
func readLineFromBuf(buf *bufio.Reader) (string, error) {
	pass, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(pass), nil
}


func ParseAddrs(addrStr string) ([]types.AccAddress, error) {
	var addrs []types.AccAddress
	as := strings.Split(addrStr, ",")
	for _, a := range as {
		a = strings.TrimSpace(a)
		if a == "" {
			break
		}
		addr, err := StrToAddress(a)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func StrToAddress(addrStr string) (types.AccAddress, error) {
	if !common.IsHexAddress(strings.TrimSpace(addrStr)) {
		return types.AccAddress{}, errors.New("invalid address provided, please use hex format")
	}
	return types.HexToAddress(addrStr), nil
}