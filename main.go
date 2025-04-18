package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const (
	chain_len = 100
)

func main() {
	temp_string := "1234"
	fmt.Println(getHash(temp_string))

	passwords := []string{"1234", "2345", "1456", "3456", "7890", "0000", "1111", "9999", "1212", "2025"}
	rainbowTable := getRainbowTable(passwords)
	fmt.Println(rainbowTable)

	hash1 := getHash("0000")
	pass, is_cracked := crack(hash1, rainbowTable)
	print_result(pass, is_cracked)
}

func print_result(pass string, is_cracked bool) {
	if is_cracked {
		fmt.Printf("The password has been cracked: %s", pass)
	} else {
		fmt.Println("Failed cracking the password!")
	}
}


func getHash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
	// byte to string
    return hex.EncodeToString(hasher.Sum(nil))
}

func reduction(hash string, step int) string {
	sum := 0
	for i := 0; i < len(hash); i++ {
		sum += int(hash[i])
	}
	sum += step
	return fmt.Sprintf("%04d", sum%10000)
}

func generateChain(start string) (string, string) {
	password := start
	for i := 0; i < chain_len; i++ {
		hash := getHash(password)
		password = reduction(hash, i)
	}
	return start, password
}

func getRainbowTable(passwords []string) map[string]string {
	table := map[string]string{}
	for _, pass := range passwords {
		_, end := generateChain(pass)
		table[pass] = end 
	}
	return table
}

func crack(targetHash string, table map[string]string) (string, bool) {
	for i := chain_len - 1; i >= 0; i-- {
		hash := targetHash
		for j := i; j < chain_len; j++ {
			password := reduction(hash, j)
			if start, ok := table[password]; ok {
				p := start
				for k := 0; k < chain_len; k++ {
					h := getHash(p)
					if h == targetHash {
						return p, true
					}
					p = reduction(h, k)
				}
			}
			hash = getHash(reduction(hash, j))
		}
	}
	return "", false
}