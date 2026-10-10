package httpapi

import (
	"errors"
	"net/http"

	"example.edu/6cc545/cw1/internal/model"
	"example.edu/6cc545/cw1/internal/state"
)

func Peers(s *state.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			WriteJSON(w, http.StatusOK, s.Snapshot())

		case http.MethodPost:
			var peer model.Peer

			if err := DecodeOne(w, r, &peer); err != nil {
				RejectInput(w, err)
				return
			}

			if err := s.Add(peer); err != nil {
				if errors.Is(err, state.ErrDuplicate) {
					ErrorJSON(w, http.StatusConflict, "duplicate_peer")
					return
				}

				ErrorJSON(w, http.StatusBadRequest, "invalid_peer")
				return
			}

			WriteJSON(w, http.StatusCreated, peer)

		default:
			w.Header().Set("Allow", "GET, POST")
			ErrorJSON(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	}
}

func Files(c *state.Catalog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			ErrorJSON(w, http.StatusMethodNotAllowed, "method_not_allowed")
			return
		}

		WriteJSON(w, http.StatusOK, c.Snapshot())
	}
}
