package safer

import "github.com/rikaaa0928/awsl/aconn"

func IOSaferFactory(m *uint32, read bool) aconn.IOMid {
	return func(io aconn.IOer) aconn.IOer {
		if m == nil {
			return io
		}
		magic := byte(*m)
		//magic += 128
		//if magic == 0 {
		//	magic = 128
		//}
		magic = Magic(magic)
		return func(bytes []byte) (int, error) {
			if !read {
				Handle(bytes, magic, false)
			}
			n, err := io(bytes)
			if err != nil {
				return n, err
			}
			if read {
				Handle(bytes, magic, true)
			}
			return n, err
		}
	}
}

func Magic(magic byte) byte {
	magic += 128
	if magic == 0 {
		magic = 128
	}
	return magic
}

func Handle(bytes []byte, magic byte, decode bool) {
	for i, v := range bytes {
		if decode {
			bytes[i] = v - magic
		} else {
			bytes[i] = v + magic
		}
	}
}

func HandleStr(str string, magic byte, decode bool) string {
	bs := []byte(str)
	Handle(bs, magic, decode)
	return string(bs)
}
