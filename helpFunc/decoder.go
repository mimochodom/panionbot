package helpFunc

import "golang.org/x/text/encoding/charmap"

func DecoderToWin1251(title string) string {
	dec := charmap.Windows1251.NewDecoder()
	// Разные кодировки = разные длины символов.
	newBody := make([]byte, len(title)*2)
	n, _, err := dec.Transform(newBody, []byte(title), false)
	if err != nil {
		panic(err)
	}
	newBody = newBody[:n]
	return string(newBody)
}
