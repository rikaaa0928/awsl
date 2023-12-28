package safer

import "github.com/rikaaa0928/awsl/aconn"

func IOSaferFactory(m *uint32, read bool) aconn.IOMid {
	return func(io aconn.IOer) aconn.IOer {
		if m == nil {
			return io
		}
		magic := byte(*m)
		magic += 128
		if magic == 0 {
			magic = 128
		}
		return func(bytes []byte) (int, error) {
			if !read {
				for i, v := range bytes {
					if read {
						bytes[i] = v - magic
					} else {
						bytes[i] = v + magic
					}
				}
			}
			n, err := io(bytes)
			if err != nil {
				return n, err
			}
			if read {
				for i, v := range bytes {
					if read {
						bytes[i] = v - magic
					} else {
						bytes[i] = v + magic
					}
				}
			}
			return n, err
		}
	}
}
