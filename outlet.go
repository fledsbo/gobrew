package main

type Outlet struct {
	CodeOn      uint64
	CodeOff     uint64
	Protocol    uint
	PulseLength uint
}

func encode(str []byte) (out uint64) {
	for _, b := range str {
		out <<= 2
		switch b {
		case 'F':
			out |= 1
		case '0':
			out |= 0
		case '1':
			out |= 3
		}
	}
	return
}

func dialCode(group int, id int, on bool) (out []byte) {
	out = []byte{'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F'}
	out[group] = '0'
	out[id+4] = '0'
	if !on {
		out[11] = '0'
	}
	return
}

func NewDialOutlet(group int, id int) (o *Outlet) {
	o = new(Outlet)
	o.Protocol = 0
	o.PulseLength = 350
	o.CodeOn = encode(dialCode(group, id, true))
	o.CodeOff = encode(dialCode(group, id, false))
	return
}
