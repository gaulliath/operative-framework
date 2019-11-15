package faker

import (
	"strings"
)

type FakeLorem interface {
	Character() string         // => "c"
	Characters(num int) string // => "wqFyJIrXYfVP7cL9M"
	Word() string              // => "veritatis"
	Words(num int) []string    // => "omnis libero neque"
	Sentence(words int) string // => "Necessitatibus cum autem."

	// Sentences returns a slice of "num" sentences, 3 to 11 words each.
	Sentences(num int) []string

	// Paragraph returns a random text of "sentences" sentences length.
	Paragraph(sentences int) string

	// Paragraphs returns a slice of "num" paragraphs, 3 to 11 sentences each.
	Paragraphs(num int) []string

	// String returns a random sentence 3 to 11 words in length.
	String() string
}

type fakeLorem struct{}

func Lorem() FakeLorem {
	return fakeLorem{}
}

// Example:
//	Lorem{}.Character() // c
func (l fakeLorem) Character() string {
	return l.Characters(1)
}

// Example:
//	Lorem{}.Characters(17) // wqFyJIrXYfVP7cL9M
func (l fakeLorem) Characters(num int) string {
	perm := localRand.Perm(len(digitsAndLetters) * (num/len(digitsAndLetters) + 1))
	res := make([]byte, num)
	for i := range res {
		res[i] = digitsAndLetters[perm[i]%len(digitsAndLetters)]
	}
	return string(res)
}

// Example:
//	Lorem{}.Word() // veritatis
func (l fakeLorem) Word() string {
	return Fetch("lorem.words")
}

// Example:
//	Lorem{}.Words(3) // [omnis libero neque]
func (l fakeLorem) Words(num int) []string {
	words := valueAt("lorem.words").([]string)
	perm := localRand.Perm(len(words) * (num/len(words) + 1))
	res := make([]string, num)
	for i := range res {
		res[i] = words[perm[i]%len(words)]
	}
	return res
}

// Sentence returns random capitalized sentence of "words" words length.
// Example:
//	Lorem{}.Sentence(3) // Necessitatibus cum autem.
func (l fakeLorem) Sentence(words int) string {
	s := strings.Join(l.Words(words), " ")
	return strings.ToTitle(s[:1]) + s[1:] + "."
}

// Sentences returns a slice of "num" sentences, 3 to 11 words each.
func (l fakeLorem) Sentences(num int) []string {
	res := make([]string, num)
	for i := range res {
		res[i] = l.Sentence(localRand.Intn(9) + 3) // 3 to 11 words
	}
	return res
}

// Paragraph returns a random text of "sentences" sentences length.
func (l fakeLorem) Paragraph(sentences int) string {
	return strings.Join(l.Sentences(sentences), " ")
}

// Paragraphs returns a slice of "num" paragraphs, 3 to 11 sentences each.
func (l fakeLorem) Paragraphs(num int) []string {
	res := make([]string, num)
	for i := range res {
		res[i] = l.Paragraph(localRand.Intn(9) + 3) // 3 to 11 sentences
	}
	return res
}

// String returns a random sentence 3 to 11 words in length.
func (l fakeLorem) String() string {
	return l.Sentence(localRand.Intn(9) + 3) // 3 to 11 words
}
