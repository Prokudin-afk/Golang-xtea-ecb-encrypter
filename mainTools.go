package main

import (
	"encoding/binary"
)

func encrypt(key []byte, block []byte, rounds uint32) []byte {
	var k [4]uint32
	var i uint32
	end := make([]byte, 8)                      //пустой массив, вместимость 8, тип-байт
	v0 := binary.LittleEndian.Uint32(block[:4]) //первые 4 байта из полученного массива
	v1 := binary.LittleEndian.Uint32(block[4:]) //вторые 4 байта из полученного массива

	k[0] = binary.LittleEndian.Uint32(key[:4]) //разбиваем ключ по 4 байта и преобразуем в бинарный вид
	k[1] = binary.LittleEndian.Uint32(key[4:8])
	k[2] = binary.LittleEndian.Uint32(key[8:12])
	k[3] = binary.LittleEndian.Uint32(key[12:])

	delta := binary.LittleEndian.Uint32([]byte{0xb9, 0x79, 0x37, 0x9e}) //4 байта в бинарном виде
	mask := binary.LittleEndian.Uint32([]byte{0xff, 0xff, 0xff, 0xff})

	var sum uint32 = 0 //логическое И

	for i = 0; i < rounds; i++ { //выполнится 32 раза
		v0 = (v0 + (((v1<<4 ^ v1>>5) + v1) ^ (sum + k[sum&3]))) & mask
		sum = (sum + delta) & mask
		v1 = (v1 + (((v0<<4 ^ v0>>5) + v0) ^ (sum + k[sum>>11&3]))) & mask
	}

	binary.LittleEndian.PutUint32(end[:4], v0) //перезаписать v0 и v1
	binary.LittleEndian.PutUint32(end[4:], v1)

	return end
}

func decrypt(key []byte, block []byte, rounds uint32) []byte {
	var k [4]uint32
	var i uint32
	end := make([]byte, 8)                      //пустой массив, вместимость 8, тип-байт
	v0 := binary.LittleEndian.Uint32(block[:4]) //первые 4 байта из полученного массива
	v1 := binary.LittleEndian.Uint32(block[4:]) //вторые 4 байта из полученного массива

	k[0] = binary.LittleEndian.Uint32(key[:4]) //разбиваем ключ по 4 байта и преобразуем в бинарный вид
	k[1] = binary.LittleEndian.Uint32(key[4:8])
	k[2] = binary.LittleEndian.Uint32(key[8:12])
	k[3] = binary.LittleEndian.Uint32(key[12:])

	delta := binary.LittleEndian.Uint32([]byte{0xb9, 0x79, 0x37, 0x9e}) //4 байта в бинарном виде
	mask := binary.LittleEndian.Uint32([]byte{0xff, 0xff, 0xff, 0xff})

	sum := (delta * rounds) & mask //логическое И

	for i = 0; i < rounds; i++ { //выполнится 32 раза
		v1 = (v1 - (((v0<<4 ^ v0>>5) + v0) ^ (sum + k[sum>>11&3]))) & mask
		sum = (sum - delta) & mask
		v0 = (v0 - (((v1<<4 ^ v1>>5) + v1) ^ (sum + k[sum&3]))) & mask
	}

	binary.LittleEndian.PutUint32(end[:4], v0) //перезаписать v0 и v1
	binary.LittleEndian.PutUint32(end[4:], v1)

	return end
}
