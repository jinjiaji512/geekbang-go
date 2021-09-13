package protocol

import (
	"errors"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/Terry-Mao/goim/pkg/encoding/binary"
)

//Proto Proto
type Proto struct {
	Ver  int32  `protobuf:"varint,1,opt,name=ver,proto3" json:"ver,omitempty"`
	Op   int32  `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	Seq  int32  `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Body []byte `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
}

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}
