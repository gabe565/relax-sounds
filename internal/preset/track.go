package preset

type Track struct {
	Id     string
	Volume float64
}

func (track Track) ToShorthand() ShorthandTrack {
	return ShorthandTrack{track.Id, track.Volume}
}
