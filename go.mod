module github.com/gabe565/relax-sounds

go 1.19

require (
	github.com/aofei/mimesniffer v1.2.1
	github.com/faiface/beep v1.1.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/spf13/pflag v1.0.5
	github.com/viert/go-lame v0.0.0-20201108052322-bb552596b11d
	golang.org/x/sync v0.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/hajimehoshi/go-mp3 v0.3.3 // indirect
	github.com/icza/bitio v1.1.0 // indirect
	github.com/jfreymuth/oggvorbis v1.0.3 // indirect
	github.com/jfreymuth/vorbis v1.0.2 // indirect
	github.com/mewkiz/flac v1.0.7 // indirect
	github.com/mewkiz/pkg v0.0.0-20211102230744-16a6ce8f1b77 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace github.com/viert/go-lame => github.com/au1/go-lame v0.0.0-20210615194416-b6cef834b0e1

replace github.com/faiface/beep => github.com/gabe565/beep v1.2.0
