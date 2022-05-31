package global

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	targetIndex := length - unpadding
	if targetIndex < 0 {
		targetIndex = 0
	}
	return origData[: targetIndex]
}

func AesEncrypt(origData, key []byte) ( string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	outData := decimalByte2HexString( crypted )
	
	return outData, nil
}

func AesDecrypt( result string, key []byte) ([]byte, error) {
	crypted, err := hexString2DecimalByte( result )
	if err != nil {
		return nil, err
	}



	//fmt.Println("cryptedcryptedcrypted 1", crypted )

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//fmt.Println("cryptedcryptedcrypted 2", block )


	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))

	if len( crypted ) == 0 || len( origData ) == 0 {
		return nil, errors.New( "解密失败" )
	}

	if len( crypted ) % blockSize != 0 {
		return nil, errors.New( "[crypto]: input not full blocks" )
	}

	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

func decimalByte2HexString( src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src )
	return string( dst )
}

func hexString2DecimalByte( srcString string ) ( []byte, error ) {
	src := []byte( srcString )

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return []byte{}, nil
	}

	return dst[:n], nil
	//return fmt.Sprint("%s\n", dst[:n]), nil
}
