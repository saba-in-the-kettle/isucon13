package isuutil

import (
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	socketFile = "/home/isucon/app.sock"
)

// CreateUnixDomainSocketListener はHTTPサーバー用のUnixドメインソケットを作成する。
// **注意** この関数はNginxが動いているs1サーバーだけで動かす必要がある。他のサーバーは普通にTCPでListenするような分岐がmain関数に必要。
// なお、Nginx側の設定は以下のようにする。
//
//	upstream app {
//	   server unix:/home/isucon/app.sock;
//	}
func CreateUnixDomainSocketListener() (net.Listener, error) {
	err := os.Remove(socketFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to remove socket file: %w", err)
	}

	l, err := net.Listen("unix", socketFile)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	// go runユーザとnginxのユーザ（グループ）を同じにすれば777じゃなくてok
	err = os.Chmod(socketFile, 0777)
	if err != nil {
		return nil, fmt.Errorf("failed to chmod: %w", err)
	}

	return l, nil
}
