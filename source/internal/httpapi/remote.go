package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"example.edu/6cc545/cw1/internal/rpcfiles"
	"example.edu/6cc545/cw1/internal/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RemoteFiles(reg *state.Registry, nodeID string, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			ErrorJSON(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}

		peerID := r.URL.Query().Get("peer_id")
		if peerID == "" || strings.TrimSpace(peerID) != peerID {
			ErrorJSON(w, http.StatusBadRequest, "invalid_peer_id")
			return
		}

		peer, err := reg.Get(peerID)
		if err != nil {
			if errors.Is(err, state.ErrMissing) {
				ErrorJSON(w, http.StatusNotFound, "peer_not_found")
				return
			}
			ErrorJSON(w, http.StatusInternalServerError, "registry_error")
			return
		}

		response, err := rpcfiles.Fetch(r.Context(), peer.GRPCAddress, nodeID, timeout)
		if err != nil {
			switch status.Code(err) {
			case codes.InvalidArgument:
				ErrorJSON(w, http.StatusBadRequest, "invalid_remote_request")
			case codes.Unavailable:
				ErrorJSON(w, http.StatusServiceUnavailable, "remote_unavailable")
			case codes.DeadlineExceeded:
				ErrorJSON(w, http.StatusGatewayTimeout, "remote_timeout")
			case codes.Canceled:
				ErrorJSON(w, http.StatusGatewayTimeout, "remote_cancelled")
			default:
				ErrorJSON(w, http.StatusBadGateway, "remote_request_failed")
			}
			return
		}

		files := make([]map[string]any, 0, len(response.GetFiles()))
		for _, file := range response.GetFiles() {
			if file == nil {
				continue
			}
			files = append(files, map[string]any{
				"name":          file.GetName(),
				"size":          file.GetSize(),
				"sha256":        file.GetSha256(),
				"modified_unix": file.GetModifiedUnix(),
			})
		}

		WriteJSON(w, http.StatusOK, map[string]any{
			"node_id": response.GetNodeId(),
			"files":   files,
		})
	}
}
