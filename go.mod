module github.com/FractalKnight/chrysalis/chrysalis

go 1.15

require (
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/djherbis/atime v1.0.0
	github.com/FractalKnight/chrysalis v1.0.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/spf13/viper v1.16.0
	github.com/tmc/scp v0.0.0-20170824174625-f7b48647feef
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/xdefrag/viper-etcd v1.1.0
	github.com/xorrior/keyctl v1.0.1-0.20210425144957-8746c535bf58
	go.etcd.io/etcd v3.3.27+incompatible
	golang.org/x/crypto v0.9.0
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.11.0
	golang.org/x/text v0.12.0 // indirect
)

//replace (
//  github.com/FractalKnight/chrysalis/src/bash_executor => ./src/bash_executor
//	github.com/FractalKnight/chrysalis/src/cmd_executor => ./src/cmd_executor
//	github.com/FractalKnight/chrysalis/src/download => ./src/download
//	github.com/FractalKnight/chrysalis/src/link_tcp => ./src/link_tcp
//	github.com/FractalKnight/chrysalis/src/pkg/profiles => ./src/pkg/profiles
//	github.com/FractalKnight/chrysalis/src/pkg/utils/structs => ./src/pkg/utils/structs
//	github.com/FractalKnight/chrysalis/src/powershell_executor => ./src/powershell_executor
//	github.com/FractalKnight/chrysalis/src/sh_executor => ./src/sh_executor
//	github.com/FractalKnight/chrysalis/src/sleep => ./src/sleep
//	github.com/FractalKnight/chrysalis/src/socks => ./src/socks
//	github.com/FractalKnight/chrysalis/src/unlink_tcp => ./src/unlink_tcp
//	github.com/FractalKnight/chrysalis/src/upload => ./src/upload
//	github.com/FractalKnight/chrysalis/src/zsh_executor => ./src/zsh_executor
//)
