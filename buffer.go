/********************************************************/
// 字节对象
// Author 		:Jella
// Version 		:1.0.3(release)
// Dependency	:none
// Example		:
//				buf:=byt.NewBuffer()
//				buf.WriteUTF8String("HelloWorld!")
/********************************************************/

package byt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

var (
	/**
	 * 字节缓冲默认数据容量
	 */
	__capacity__ int = 32
	/**
	 * 最大数据长度（400*1024字节）
	 */
	__maxlength__ int = 400 * 1024
	/**
	 * 编码模式
	 */
	endian string = "big_endian"
)

/**
 * 设置编码模式
 * @param 编码模式（big_endian以及lit_endian）
 */
func SetEndian(s string) {
	if s != "big_endian" && s != "lit_endian" {
		return
	}
	endian = s
}

/**
 * 获取编码模式结构体（内部使用）
 */
func getEndian() binary.ByteOrder {
	if endian == "big_endian" {
		return binary.BigEndian

	} else if endian == "lit_endian" {
		return binary.LittleEndian
	}
	return nil
}

/**
 * 字节缓冲对象
 */
type __buffer__ struct {
	byt    []byte //字节对象
	top    int
	offset int
}

/**
 * 创建一个字节缓冲对象（默认容量为CAPACITY=32）
 * @return 返回创建完毕的字节缓冲对象
 */
func NewBuffer() *__buffer__ {
	return NewBufferWithLen(__capacity__)
}

/**
 * 创建一个字节缓冲对象
 * @param capacity 容量值
 * @return 返回创建完毕的字节缓冲对象
 */
func NewBufferWithLen(capacity int) *__buffer__ {
	if capacity < 1 {
		fmt.Println("[ERR]: 参数 < 1.")
		return nil
	}
	val := make([]byte, capacity)
	return &__buffer__{
		top:    0,
		offset: 0,
		byt:    val,
	}
}

/**
 * 创建一个字节缓冲对象
 * @param b 字节数组
 * @return 返回创建完毕的字节缓冲对象
 */
func NewBufferWithByte(b []byte) *__buffer__ {
	if b == nil {
		return nil
	}
	l := len(b)
	return &__buffer__{
		byt:    b,
		top:    l,
		offset: 0,
	}
}

/**
 * 设置容量
 * @param capa 容量值
 */
func (b *__buffer__) SetCapacity(capa int) {
	l := len(b.byt)
	if capa < l {
		fmt.Println("[ERR]: 参数长度不能小于当前字节容量")
		return
	}
	for ; l < capa; l = (l << 1) + 1 {
	}

	newb := make([]byte, l)
	//拷贝
	copy(newb[0:], b.byt[0:b.top])
	b.byt = newb
}

/**
 * 获取字节缓冲对象的top值
 * @return top值
 */
func (b *__buffer__) GetTop() int {
	return b.top
}

/**
 * 设置字节缓冲对象的top值
 * @param top值
 */
func (b *__buffer__) SetTop(t int) {
	if t < b.offset {
		fmt.Println("[ERR]: 参数不能小于当前的字节缓冲偏移量")
		return
	}
	if t > len(b.byt) {
		b.SetCapacity(t)
	}
	b.top = t
}

/**
 * 获取当前偏移位置
 * @return 偏移位置
 */
func (b *__buffer__) GetOffset() int {
	return b.offset
}

/**
 * 设置偏移位置
 * @param offs 偏移位置
 */
func (b *__buffer__) SetOffet(offs int) {
	if offs < 0 || offs > b.top {
		fmt.Println("[ERR]: 设置偏移位置不合法.")
		return
	}
	b.offset = offs
}

/**
 * 剩余可读取的内容长度
 */
func (b *__buffer__) Remaining() int {
	return b.top - b.offset
}

/**
* 剩余可读取的内容长度是否大于0
* @return true：大于0；false：小于0
 */
func (b *__buffer__) HasRemaining() bool {
	return b.Remaining() > 0
}

/**
 * 字节缓冲对象的数据长度
 */
func (b *__buffer__) Length() int {
	return len(b.byt)
}

/**
* 字节数组对象
* @return 一个byte[]类型对象
 */
func (b *__buffer__) GetByte() []byte {
	return b.byt
}

/**
* 获取字节有效数据总长度与当前偏移量的差值长度字节对象
* @return 字节对象
 */
func (b *__buffer__) GetRemainingByte() []byte {
	data := make([]byte, b.Remaining())
	copy(data[0:], b.byt[b.offset:])
	return data
}

/**
 * 检测参数对象是否是byt.Buffer类型
 * @return true:是；false：否
 */
func (b *__buffer__) Check(val interface{}) bool {
	return reflect.TypeOf(val).String() == reflect.TypeOf(b).String()
}

/**
 * 将偏移指针及top值归0
 */
func (b *__buffer__) Zero() {
	b.top = 0
	b.offset = 0
}

/**
 * Hash值
 */
func (b *__buffer__) Hash() int {
	h := 17
	for i := b.top - 1; i >= 0; i-- {
		h = 65537*h + int(b.byt[i])
	}
	return h
}

/**
 * 相等判断（相当于“==”，并不是“===”）
 * @return true：相等；false：不相等
 */
func (b *__buffer__) Equal(val interface{}) bool {
	if !b.Check(val) {
		return false
	}
	_val_ := val.(__buffer__)
	if _val_.top != b.top {
		return false
	}
	if _val_.offset != b.offset {
		return false
	}
	for i := b.top - 1; i >= 0; i-- {
		if _val_.byt[i] != b.byt[i] {
			return false
		}
	}
	return true
}

/**
 * 释放
 */
func (b *__buffer__) Kill() {
	b.Zero()
	if b.byt != nil {
		b.byt = nil
	}
}

////////////////////////////////////////////////////
//						读						  //
////////////////////////////////////////////////////

/**
 * 读取
 * @param bt 目标字节数组
 * @param pos 读取的内容至目标字节数组中的插入位置
 * @param l 从源数据中读取的长度
 */
func (b *__buffer__) Read(bt []byte, pos int, l int) {
	_pos_ := b.offset
	copy(bt[pos:], b.byt[_pos_:_pos_+l])
	b.offset += l
}

/**
 * 读一个boolean布尔值
 */
func (b *__buffer__) ReadBoolean() bool {
	bol := (b.byt[b.offset] != 0)
	b.offset++
	return bol
}

/**
 * 读取一个无符号的byte值（uint8）
 */
func (b *__buffer__) ReadUnsignedByt() byte {
	_byt_ := b.byt[b.offset]
	b.offset++
	return _byt_
}

/**
 * 读取一个byte值（int8）
 */
func (b *__buffer__) ReadByt() int8 {
	var _val_ int8
	binary.Read(readValue(b, 1), getEndian(), &_val_)
	return _val_
}

/**
 * 读取一个Short值（int16）
 */
func (b *__buffer__) ReadShort() int16 {
	var _val_ int16
	binary.Read(readValue(b, 2), getEndian(), &_val_)
	return _val_
}

/**
 * 读取一个int值（int32）
 */
func (b *__buffer__) ReadInt() int32 {
	var _val_ int32
	binary.Read(readValue(b, 4), getEndian(), &_val_)
	return _val_
}

/**
* 读取一个long值（int64）
*
 */
func (b *__buffer__) ReadLong() int64 {
	var _val_ int64
	binary.Read(readValue(b, 8), getEndian(), &_val_)
	return _val_
}

/**
 * 读取一个float值（float64）
 */
func (b *__buffer__) ReadFloat() float64 {
	var _val_ float64
	binary.Read(readValue(b, 8), getEndian(), &_val_)
	return _val_
}

/**
 * 读取一个长度值
 */
func (b *__buffer__) ReadLength() int {
	var n uint8 = b.byt[b.offset] & 0xff
	if n >= 0x80 {
		b.offset++
		return int(n - 0x80)

	} else if n >= 0x40 {
		return int(b.ReadShort() - 0x4000)

	} else if n >= 0x20 {
		return int(b.ReadInt() - 0x20000000)
	}
	fmt.Println("[ERR]: ReadLength 错误.")

	return -1
}

/**
 * 读取一个utf8字符串
 */
func (b *__buffer__) ReadUTF8String() string {
	_len := b.ReadLength() - 1
	if _len < 0 {
		return ""
	}

	if _len == 0 {
		return ""
	}

	if _len > __maxlength__ {
		fmt.Println("[ERR]: 读取错误.")
		return ""
	}

	var (
		i   int
		c   int
		cc  int
		ccc int
	)

	_news := ""
	_pos := b.offset
	_end := _pos + _len
	for _pos < _end {
		c = int(b.byt[_pos] & 0xff)
		i = c >> 4
		if i < 8 {
			// 0xxx xxxx
			_pos++
			_news += string(c)

		} else if i == 12 || i == 13 {
			// 110x xxxx 10xx xxxx
			_pos += 2
			if _pos > _end {
				break
			}
			cc = int(b.byt[_pos-1])
			if (cc & 0xC0) != 0x80 {
				break
			}
			_news += string((((c & 0x1f) << 6) | (cc & 0x3f)))

		} else if i == 14 {
			// 1110 xxxx 10xx xxxx 10xx
			// xxxx
			_pos += 3
			if _pos > _end {
				break
			}
			cc = int(b.byt[_pos-2])
			ccc = int(b.byt[_pos-1])
			if ((cc & 0xC0) != 0x80) || ((ccc & 0xC0) != 0x80) {
				break
			}
			_news += string((((c & 0x0f) << 12) | ((cc & 0x3f) << 6) | (ccc & 0x3f)))

		} else {
			// 10xx xxxx 1111 xxxx
			break
		}
	}
	b.offset += _len
	return _news
}

/**
 * 读取一个字节数组
 */
func (b *__buffer__) ReadData() []byte {
	_len := b.ReadLength() - 1
	if _len < 0 {
		fmt.Println("[ERR]: ReadData发生错误. len < 0.")
		return nil
	}
	if _len > __maxlength__ {
		fmt.Println("[ERR]: ReadData发生错误. len = " + strconv.Itoa(_len))
		return nil
	}

	_b := make([]byte, _len)
	b.Read(_b, 0, _len)

	return _b
}

////////////////////////////////////////////////////
//						写						  //
////////////////////////////////////////////////////

/**
 * 写入
 * @param data 要被写入的[]byte字节数组
 * @param pos 源位置
 * @param l 源长度
 */
func (b *__buffer__) Write(data []byte, pos int, l int) {
	_l := b.top + l
	if len(b.byt) < _l {
		b.SetCapacity(_l)
	}
	//拷贝
	copy(b.byt[b.top:], data[pos:l])
	b.top += l
}

/**
 * 写一个Boolean值
 * @param val 布尔值
 */
func (b *__buffer__) WriteBoolean(val bool) {
	if len(b.byt) < b.top+1 {
		b.SetCapacity(b.top + __capacity__)
	}

	_val_ := 0
	if val {
		_val_ = 1
	}

	b.byt[b.top] = byte(_val_)
	b.top++
}

/**
 * 写一个无符号的Byte值（uint8）
 * @param val byte值
 */
func (b *__buffer__) WriteUnsignedByt(val byte) {
	if len(b.byt) < b.top+1 {
		b.SetCapacity(b.top + __capacity__)
	}
	b.byt[b.top] = val
	b.top++
}

/**
 * 写一个byte值（int8）
 * @param val 值
 */
func (b *__buffer__) WriteByt(val int8) {
	writeValue(b, 1, val)
}

/**
 * 写一个short值（int16）
 * @param val short值
 */
func (b *__buffer__) WriteShort(val int16) {
	writeValue(b, 2, val)
}

/**
 * 写一个int值（int32）
 * @param int32类型的值
 */
func (b *__buffer__) WriteInt(val int32) {
	writeValue(b, 4, val)
}

/**
* 写一个long值（int64）
* @param val long值
 */
func (b *__buffer__) WriteLong(val int64) {
	writeValue(b, 8, val)
}

/**
 * 写一个float值（float64）
 * @param val 浮点数
 */
func (b *__buffer__) WriteFloat(val float64) {
	writeValue(b, 8, val)
}

/**
 * 写一个长度值
 * @param val 长度值
 */
func (b *__buffer__) WriteLength(val int) {
	if val >= 0x20000000 || val < 0 {
		fmt.Println("[ERR]: WriteLength 长度错误.")
		return
	}
	if val >= 0x4000 { //0100.0000.0000.0000 16位int值
		b.WriteInt(int32(val + 0x20000000)) //0010.0000.0000.0000.0000.0000.0000.0000 32位int值

	} else if val >= 0x80 { //1000.0000 //8位int值
		b.WriteShort(int16(val + 0x4000))

	} else {
		b.WriteByt(int8(val + 0x80))
	}
}

/**
 * 写一个utf8字符串
 * @param s 字符串值
 */
func (b *__buffer__) WriteUTF8String(s string) {
	// fmt.Println(s, len(s))

	//写入字符串长度
	_len := len(s)
	b.WriteLength(_len + 1)

	//根据长度及当前偏移位置进行扩容判断
	pos := b.top
	if len(b.byt) < pos+_len {
		b.SetCapacity(pos + _len)
	}

	//写入字符
	for _, c := range s {
		if (c >= 0x0001) && (c <= 0x007f) {
			b.byt[pos] = byte(c)
			pos++

		} else if c > 0x07ff {
			b.byt[pos] = byte(0xe0 | ((c >> 12) & 0x0f))
			pos++
			b.byt[pos] = byte(0x80 | ((c >> 6) & 0x3f))
			pos++
			b.byt[pos] = byte(0x80 | (c & 0x3f))
			pos++

		} else {
			b.byt[pos] = byte(0xc0 | ((c >> 6) & 0x1f))
			pos++
			b.byt[pos] = byte(0x80 | (c & 0x3f))
			pos++
		}
	}
	b.top += _len
}

/**
 * 写一个字节数组
 * @param bt 字节数组
 */
func (b *__buffer__) WriteData(bt []byte) {
	_len := len(bt)
	b.WriteLength(_len + 1)
	b.Write(bt, 0, _len)
}

////////////////////////////////////////////////////////////////////////
//内部函数

func getStringByteLength(s string) int {
	var _len_ int = 0
	for _, c := range s {
		if (c >= 0x0001) && (c <= 0x007f) {
			_len_++
		} else if c > 0x07ff {
			_len_ += 3
		} else {
			_len_ += 2
		}
	}
	return _len_
}

func writeValue(b *__buffer__, n int, val interface{}) {
	_pos_ := b.top
	if len(b.byt) < _pos_+n {
		b.SetCapacity(_pos_ + __capacity__)
	}
	_b_ := bytes.NewBuffer([]byte{})
	binary.Write(_b_, getEndian(), val)
	_bb_ := _b_.Bytes()
	copy(b.byt[_pos_:], _bb_[0:])
	b.top += n
}

func readValue(b *__buffer__, blen int) *bytes.Buffer {
	_pos_ := b.offset

	bt := make([]byte, blen)
	copy(bt[0:], b.byt[_pos_:_pos_+blen])

	b.offset += blen

	return bytes.NewBuffer(bt)
}
