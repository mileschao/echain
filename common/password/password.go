package password

import (
	"bytes"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

// WaitForPwdInput wait for user to input password
func WaitForPwdInput() ([]byte, error) {
	fmt.Printf("Please input your Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

// PwdInputAndConfirm wait for user to input password and its confirmation
func PwdInputAndConfirm() ([]byte, error) {
	fmt.Printf("Please input your Password:")
	first, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(first) {
		fmt.Println("You have to input password.")
		os.Exit(1)
	}

	fmt.Printf("Re-enter Password:")
	second, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if 0 == len(second) {
		fmt.Println("You have to input password.")
		os.Exit(1)
	}

	if !bytes.Equal(first, second) {
		fmt.Println("Unmatched Password")
		os.Exit(1)
	}
	return first, nil
}
