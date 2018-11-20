// Tamanho da chave -> 256 bits
// Tamanho do bloco -> 128 bits

package common

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func Encrypt(key []byte, message []byte) (ciphertext []byte) {
	plaintext := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	ciphertext = aesgcm.Seal(nil, iv, plaintext, nil)
	ciphertext = append(iv, ciphertext...)
	//fmt.Printf("%x\n", ciphertext)

	return
}

func Decrypt(key []byte, ciphertext []byte) (message []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	iv := ciphertext[:12]
	msg := ciphertext[12:]

	plaintext, err := aesgcm.Open(nil, iv, msg, nil)
	if err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0)
	}

	message = plaintext

	return
}

func main() {
	key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Digite a mensagem: ")
	fmt.Print("->")
	msg, _ := reader.ReadString('\n')

	encryptedMsg := Encrypt([]byte(key), []byte(msg))
	decryptedMsg := Decrypt(key, encryptedMsg)

	fmt.Printf("CIPHER KEY: %d\n", len(key))
	fmt.Printf("ENCRYPTED: %d\n", len(encryptedMsg))
	fmt.Printf("DECRYPTED: %d\n", len(decryptedMsg))
}
