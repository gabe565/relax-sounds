module github.com/gabe565/relax-sounds

go 1.17

require (
	github.com/faiface/beep v1.1.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/spf13/pflag v1.0.5
	github.com/viert/go-lame v0.0.0-20201108052322-bb552596b11d
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require (
	github.com/jfreymuth/oggvorbis v1.0.3 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace github.com/viert/go-lame => github.com/au1/go-lame v0.0.0-20210615194416-b6cef834b0e1

replace github.com/faiface/beep => github.com/gabe565/beep v1.1.1-0.20211227183743-02af4fbda67b
