module github.com/gabe565/relax-sounds

go 1.16

require (
	github.com/faiface/beep v1.0.2
	github.com/go-chi/chi/v5 v5.0.3
	github.com/spf13/pflag v1.0.5
	github.com/viert/go-lame v0.0.0-20201108052322-bb552596b11d
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

replace github.com/viert/go-lame => github.com/au1/go-lame v0.0.0-20210615194416-b6cef834b0e1
