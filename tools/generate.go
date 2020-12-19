package tools

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
)

type GenerateUtil struct {
}

func NewGenerate() *GenerateUtil {
	return &GenerateUtil{}
}

//创建UUID
func (receiver GenerateUtil) GenerateUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

//生成MD5
func (receiver GenerateUtil) GenerateMd5(str string) string {

	has := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", has)
}
