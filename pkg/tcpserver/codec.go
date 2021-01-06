package tcpserver

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	errorset "github.com/panjf2000/gnet/errors"
)

var CRLFByte = byte('\n')

type (
	ICodec interface {
		Encode(c Conn, buf []byte) ([]byte, error)
		Decode(c Conn) ([]byte, error)
	}

	BuiltInFrameCodec struct {
	}

	LineBasedFrameCodec struct {
	}

	DelimiterBasedFrameCodec struct {
		delimiter byte
	}

	FixedLengthFrameCodec struct {
		frameLength int
	}

	LengthFieldBasedFrameCodec struct {
		encoderConfig EncoderConfig
		decoderConfig DecoderConfig
	}
)

func (cc *BuiltInFrameCodec) Encode(c Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (cc *BuiltInFrameCodec) Decode(c Conn) ([]byte, error) {
	buf := c.Read()
	if len(buf) == 0 {
		return nil, nil
	}
	c.ResetBuffer()
	return buf, nil
}

func (cc *LineBasedFrameCodec) Encode(c Conn, buf []byte) ([]byte, error) {
	return append(buf, CRLFByte), nil
}

func (cc *LineBasedFrameCodec) Decode(c Conn) ([]byte, error) {
	buf := c.Read()
	idx := bytes.IndexByte(buf, CRLFByte)
	if idx == -1 {
		return nil, errorset.ErrCRLFNotFound
	}
	c.ShiftN(idx + 1)
	return buf[:idx], nil
}

func NewDelimiterBasedFrameCodec(delimiter byte) *DelimiterBasedFrameCodec {
	return &DelimiterBasedFrameCodec{delimiter}
}

func (cc *DelimiterBasedFrameCodec) Encode(c Conn, buf []byte) ([]byte, error) {
	return append(buf, cc.delimiter), nil
}

func (cc *DelimiterBasedFrameCodec) Decode(c Conn) ([]byte, error) {
	buf := c.Read()
	idx := bytes.IndexByte(buf, cc.delimiter)
	if idx == -1 {
		return nil, errorset.ErrDelimiterNotFound
	}
	c.ShiftN(idx + 1)
	return buf[:idx], nil
}

func NewFixedLengthFrameCodec(frameLength int) *FixedLengthFrameCodec {
	return &FixedLengthFrameCodec{frameLength}
}

func (cc *FixedLengthFrameCodec) Encode(c Conn, buf []byte) ([]byte, error) {
	if len(buf)%cc.frameLength != 0 {
		return nil, errorset.ErrInvalidFixedLength
	}
	return buf, nil
}

func (cc *FixedLengthFrameCodec) Decode(c Conn) ([]byte, error) {
	size, buf := c.ReadN(cc.frameLength)
	if size == 0 {
		return nil, errorset.ErrUnexpectedEOF
	}
	c.ShiftN(size)
	return buf, nil
}

func NewLengthFieldBasedFrameCodec(ec EncoderConfig, dc DecoderConfig) *LengthFieldBasedFrameCodec {
	return &LengthFieldBasedFrameCodec{encoderConfig: ec, decoderConfig: dc}
}

type EncoderConfig struct {
	ByteOrder binary.ByteOrder
	LengthFieldLength int
	LengthAdjustment int
	LengthIncludesLengthFieldLength bool
}

type DecoderConfig struct {
	ByteOrder binary.ByteOrder
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment int
	InitialBytesToStrip int
}

func (cc *LengthFieldBasedFrameCodec) Encode(c Conn, buf []byte) (out []byte, err error) {
	length := len(buf) + cc.encoderConfig.LengthAdjustment
	if cc.encoderConfig.LengthIncludesLengthFieldLength {
		length += cc.encoderConfig.LengthFieldLength
	}

	if length < 0 {
		return nil, errorset.ErrTooLessLength
	}

	switch cc.encoderConfig.LengthFieldLength {
	case 1:
		if length >= 256 {
			return nil, fmt.Errorf("length does not fit into a byte: %d", length)
		}
		out = []byte{byte(length)}
	case 2:
		if length >= 65536 {
			return nil, fmt.Errorf("length does not fit into a short integer: %d", length)
		}
		out = make([]byte, 2)
		cc.encoderConfig.ByteOrder.PutUint16(out, uint16(length))
	case 3:
		if length >= 16777216 {
			return nil, fmt.Errorf("length does not fit into a medium integer: %d", length)
		}
		out = writeUint24(cc.encoderConfig.ByteOrder, length)
	case 4:
		out = make([]byte, 4)
		cc.encoderConfig.ByteOrder.PutUint32(out, uint32(length))
	case 8:
		out = make([]byte, 8)
		cc.encoderConfig.ByteOrder.PutUint64(out, uint64(length))
	default:
		return nil, errorset.ErrUnsupportedLength
	}

	out = append(out, buf...)
	return
}

type innerBuffer []byte

func (in *innerBuffer) readN(n int) (buf []byte, err error) {
	if n == 0 {
		return nil, nil
	}

	if n < 0 {
		return nil, errors.New("negative length is invalid")
	} else if n > len(*in) {
		return nil, errors.New("exceeding buffer length")
	}
	buf = (*in)[:n]
	*in = (*in)[n:]
	return
}

func (cc *LengthFieldBasedFrameCodec) Decode(c Conn) ([]byte, error) {
	var (
		in     innerBuffer
		header []byte
		err    error
	)
	in = c.Read()
	if cc.decoderConfig.LengthFieldOffset > 0 {
		header, err = in.readN(cc.decoderConfig.LengthFieldOffset)
		if err != nil {
			return nil, errorset.ErrUnexpectedEOF
		}
	}

	lenBuf, frameLength, err := cc.getUnadjustedFrameLength(&in)
	if err != nil {
		return nil, err
	}

	msgLength := int(frameLength) + cc.decoderConfig.LengthAdjustment
	msg, err := in.readN(msgLength)
	if err != nil {
		return nil, errorset.ErrUnexpectedEOF
	}

	fullMessage := make([]byte, len(header)+len(lenBuf)+msgLength)
	copy(fullMessage, header)
	copy(fullMessage[len(header):], lenBuf)
	copy(fullMessage[len(header)+len(lenBuf):], msg)
	c.ShiftN(len(fullMessage))
	return fullMessage[cc.decoderConfig.InitialBytesToStrip:], nil
}

func (cc *LengthFieldBasedFrameCodec) getUnadjustedFrameLength(in *innerBuffer) ([]byte, uint64, error) {
	switch cc.decoderConfig.LengthFieldLength {
	case 1:
		b, err := in.readN(1)
		if err != nil {
			return nil, 0, errorset.ErrUnexpectedEOF
		}
		return b, uint64(b[0]), nil
	case 2:
		lenBuf, err := in.readN(2)
		if err != nil {
			return nil, 0, errorset.ErrUnexpectedEOF
		}
		return lenBuf, uint64(cc.decoderConfig.ByteOrder.Uint16(lenBuf)), nil
	case 3:
		lenBuf, err := in.readN(3)
		if err != nil {
			return nil, 0, errorset.ErrUnexpectedEOF
		}
		return lenBuf, readUint24(cc.decoderConfig.ByteOrder, lenBuf), nil
	case 4:
		lenBuf, err := in.readN(4)
		if err != nil {
			return nil, 0, errorset.ErrUnexpectedEOF
		}
		return lenBuf, uint64(cc.decoderConfig.ByteOrder.Uint32(lenBuf)), nil
	case 8:
		lenBuf, err := in.readN(8)
		if err != nil {
			return nil, 0, errorset.ErrUnexpectedEOF
		}
		return lenBuf, cc.decoderConfig.ByteOrder.Uint64(lenBuf), nil
	default:
		return nil, 0, errorset.ErrUnsupportedLength
	}
}

func readUint24(byteOrder binary.ByteOrder, b []byte) uint64 {
	_ = b[2]
	if byteOrder == binary.LittleEndian {
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16
	}
	return uint64(b[2]) | uint64(b[1])<<8 | uint64(b[0])<<16
}

func writeUint24(byteOrder binary.ByteOrder, v int) []byte {
	b := make([]byte, 3)
	if byteOrder == binary.LittleEndian {
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
	} else {
		b[2] = byte(v)
		b[1] = byte(v >> 8)
		b[0] = byte(v >> 16)
	}
	return b
}
