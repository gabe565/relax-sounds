module github.com/gabe565/relax-sounds

go 1.16

require (
	github.com/faiface/beep v1.0.2
	github.com/go-chi/chi/v5 v5.0.3
	github.com/juju/ratelimit v1.0.1
	github.com/viert/go-lame v0.0.0-20201108052322-bb552596b11d
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/viert/go-lame => github.com/au1/go-lame v0.0.0-20210615194416-b6cef834b0e1
