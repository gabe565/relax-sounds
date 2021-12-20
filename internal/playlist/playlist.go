package playlist

type Entry struct {
	Key    string
	Volume float64
}

func (entry Entry) ToShorthand() ShorthandEntry {
	return ShorthandEntry{entry.Key, entry.Volume}
}

type Playlist []Entry

func (playlist Playlist) ToShorthand() Shorthand {
	shorthand := make(Shorthand, 0, len(playlist))
	for _, entry := range playlist {
		shorthand = append(shorthand, entry.ToShorthand())
	}
	return shorthand
}
