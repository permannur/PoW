package pow

type Option func(*Pow)

func MaxAttemptsCount(count uint32) Option {
	return func(p *Pow) {
		p.maxAttemptsCount = count
	}
}
