module github.com/FractalKnight/chrysalis

go 1.15

require (
	github.com/djherbis/atime v1.0.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/spf13/viper v1.15.0
	github.com/tmc/scp v0.0.0-20170824174625-f7b48647feef
	github.com/xdefrag/viper-etcd v1.1.0
	github.com/xorrior/keyctl v1.0.1-0.20210425144957-8746c535bf58
	go.etcd.io/etcd v0.0.0-20190225040740-6543273666cb
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.3.0
)

replace (
	"github.com/FractalKnight/chrysalis/bash_executor" => "./chrysalis/bash_executor"
	"github.com/FractalKnight/chrysalis/cmd_executor" => "./cmd_executor"
	"github.com/FractalKnight/chrysalis/download" => "./download"
	"github.com/FractalKnight/chrysalis/pkg/profiles" => "./pkg/profiles"
	"github.com/FractalKnight/chrysalis/pkg/utils/structs" => "./pkg/utils/structs"
	"github.com/FractalKnight/chrysalis/powershell_executor" => "./powershell_executor"
	"github.com/FractalKnight/chrysalis/sh_executor" => "./sh_executor"
	"github.com/FractalKnight/chrysalis/socks" => "./socks"
	"github.com/FractalKnight/chrysalis/upload" => "./upload"
	"github.com/FractalKnight/chrysalis/zsh_executor" => "./zsh_executor"
	"github.com/FractalKnight/chrysalis/link_tcp" => "./link_tcp"
    "github.com/FractalKnight/chrysalis/sleep"  => "./sleep"
    "github.com/FractalKnight/chrysalis/unlink_tcp"  => "./unlink_tcp"
)
