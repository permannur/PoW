package pow

type Option func(*Pow)

func LeadingZeroCount(count byte) Option {
	return func(p *Pow) {
		p.leadingZeroCount = count
	}
}

func Generator(generator func() []byte) Option {
	return func(p *Pow) {
		p.generator = generator
	}
}
