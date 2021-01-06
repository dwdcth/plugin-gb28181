module github.com/Monibuca/plugin-gb28181

go 1.13

require (
	github.com/Monibuca/engine/v2 v2.2.5
	github.com/Monibuca/plugin-rtp v1.0.0
	github.com/ReneKroon/ttlcache/v2 v2.1.0
	github.com/basgys/goxml2json v1.1.0
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/golang-module/carbon v1.2.4
	github.com/json-iterator/go v1.1.10
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pion/rtp v1.6.0 // indirect
	github.com/shirou/gopsutil v2.20.8+incompatible // indirect
	golang.org/x/net v0.0.0-20201029221708-28c70e62bb1d
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/Monibuca/engine/v2 v2.2.5 => github.com/dwdcth/engine/v2 v2.2.5-fix
