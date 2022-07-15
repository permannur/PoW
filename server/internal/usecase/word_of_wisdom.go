package usecase

type WordOfWisdom struct {
}

func NewWordOfWisdom() *WordOfWisdom {
	return &WordOfWisdom{}
}

func (w *WordOfWisdom) GetWordOfWisdom() string {
	return "word of wisdom"
}
