package metahash

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
)

type Varint struct {
	buffer bytes.Buffer
}

//NewVarint creates a new variant
func NewVarint() *Varint {
	return new(Varint)
}

//GetBytes returns byte data from variant
func (vt *Varint) GetBytes() []byte {
	return vt.buffer.Bytes()
}

/*
 first byte      value
 0-249 	         the same number
 250 (0xfa) 	 as uint16
 251 (0xfb) 	 as uint32
 252 (0xfc) 	 as uint64
 253 (0xfd) 	 as uint128
 254 (0xfe) 	 as uint256
 255 (0xff) 	 as uint512
*/

//generatteVarint generates a variant data
func generatteVarint(number *big.Int, buffer *bytes.Buffer) (*bytes.Buffer, error) {
	var buff *bytes.Buffer
	if buffer == nil {
		buff = &bytes.Buffer{}
	} else {
		buff = buffer
	}

	bitLen := number.BitLen()
	var b, p int

	switch {
	case number.Sign() < 0:
		return buff, errors.New("error processing negative number")
	case number.Sign() == 0:
		buff.WriteByte(0)
		return buff, nil
	case number.Cmp(big.NewInt(249)) <= 0:
		buff.WriteByte(number.Bytes()[0])
		return buff, nil
	case bitLen > 512:
		return buff, errors.New("the number is too big")
	default:
		for b, p = 16, 250; b < 512 && b < bitLen; b, p = b*2, p+1 {
		}
		//write tag
		buff.WriteByte(byte(p))
		//write number, BE as LE
		numberBytes := number.Bytes()
		for i := (bitLen - 1) / 8; i >= 0; i-- {
			buff.WriteByte(numberBytes[i])
		}
		//write leading zeroes (eg 300 bits number & 512 necessary by format)
		zeroBytes := (b - bitLen) / 8
		for i := 0; i < zeroBytes; i++ {
			buff.WriteByte(0)
		}
		return buff, nil
	}
}

//Append adds data to the variant
func (vt *Varint) Append(number *big.Int) error {
	_, err := generatteVarint(number, &vt.buffer)
	return err
}

//AppendString appends string to the variant
func (vt *Varint) AppendString(str string) error {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	vt.AppendBytes(bytes)
	return nil
}

//AppendBytes appends bytes to the variant
func (vt *Varint) AppendBytes(bytes []byte) {
	if len(bytes) == 0 {
		vt.buffer.WriteByte(0)
	} else {
		vt.buffer.Write(bytes)
	}
}
