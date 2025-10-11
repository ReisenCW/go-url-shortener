package shortener

import (
	"crypto/sha256"
	"os"
	"fmt"
	"math/big"
	"github.com/itchyny/base58-go"
)

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	// sha256只处理“字节数据”，不直接处理字符串, 因此要先把string转化成byte类型的slice
	algorithm.Write([]byte(input))
	// Sum(nil):不额外追加其他数据，直接输出当前计算的哈希结果
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func GenerateShortLink(initialLink string, userId string) string {
	urlHashBytes := sha256Of(initialLink + userId)
	// 把256位hash值看作大整数并截断为uint64
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()

	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	// 取前8字符作为短url
	return finalString[:8]
}