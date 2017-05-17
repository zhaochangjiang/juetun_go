package general

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

type PasswordEncyption struct {
}

func (this *PasswordEncyption) Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
func (this *PasswordEncyption) Encyption(password string) string {
	h := md5.New()
	h.Write([]byte(password)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	rString := hex.EncodeToString(cipherStr)
	return rString
}
