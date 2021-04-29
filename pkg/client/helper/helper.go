package helper

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bgentry/speakeasy"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/ci123chain/ci123chain/pkg/util"
)

//const (
//	MinPassLength = 4
//
//	FlagBlocked = "blocked"
//	FlagHeight = "height"
//	FlagHomeDir = "clihome"
//	FlagVerbose = "verbose"
//	FlagNode = "node"
//	FlagAddress = "address"
//	FlagPassword = "password"
//	//FlagWithCrypto 	   = "cryptosuit"
//
//	FlagFile = "file"
//	FlagGas = "gas"
//	FlagPrivateKey = "privateKey"
//	//FlagMsg = "msg"
//	FlagArgs = "args"
//	FlagName = "name"
//	FlagCodeHash = "codeHash"
//	FlagVersion = "version"
//	FlagAuthor = "author"
//	FlagEmail = "email"
//	FlagDescribe = "describe"
//	FlagHash = "codeHash"
//	FlagFunds = "funds"
//	FlagContractAddress = "contractAddress"
//)

func GetPassphrase(addr sdk.AccAddress) (string, error) {
	pass := viper.GetString(util.FlagPassword)
	if pass == "" {
		return getPassphraseFromStdin(addr)
	}
	return pass, nil
}

// Get passphrase from std input
func getPassphraseFromStdin(addr sdk.AccAddress) (string, error) {
	buf := BufferStdin()
	prompt := fmt.Sprintf("Enter password for address: '%s'", addr.Hex())
	return GetPassword(prompt, buf)
}

// Allows for reading prompts for stdin
func BufferStdin() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func GetPasswordFromStd() (string, error) {
	var err error
	buf := BufferStdin()
	pass, err := GetCheckPassword("Enter a passphrase for your account:", "Repeat the passphrase:", buf)
	if err != nil {
		return "", sdkerrors.Wrap(sdkerrors.ErrWrongPassword, err.Error())
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
		return "", sdkerrors.Wrap(sdkerrors.ErrWrongPassword, err.Error())
	}
	pass2, err := GetPassword(prompt2, buf)
	if err != nil {
		return "", sdkerrors.Wrap(sdkerrors.ErrWrongPassword, err.Error())
	}
	if pass != pass2 {
		return "", sdkerrors.Wrap(sdkerrors.ErrWrongPassword, fmt.Sprintf("password does not math"))
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
	if len(pass) < util.MinPassLength {
		return "", fmt.Errorf("Password must be at least %d characters", util.MinPassLength)
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


func ParseAddrs(addrStr string) ([]sdk.AccAddress, error) {
	var addrs []sdk.AccAddress
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

func StrToAddress(addrStr string) (sdk.AccAddress, error) {
	if !common.IsHexAddress(strings.TrimSpace(addrStr)) {
		return sdk.AccAddress{}, errors.New("invalid address provided, please use hex format")
	}
	return sdk.HexToAddress(addrStr), nil
}