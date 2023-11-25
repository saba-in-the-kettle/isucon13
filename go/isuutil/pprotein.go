package isuutil

import (
	"fmt"
	"net/http"
)

const pproteinCollectEndpoint = "http://watchtower:9000/api/group/collect"

// KickPproteinCollect は PProteinのCollectを開始するAPIを呼び出します。
// この関数をinitializeで呼び出すことで自動でPProteinのCollectが開始されます。
func KickPproteinCollect() error {
	_, err := http.DefaultClient.Get(pproteinCollectEndpoint)
	if err != nil {
		return fmt.Errorf("failed to kick pprotein collect: %w", err)
	}
	return nil
}
