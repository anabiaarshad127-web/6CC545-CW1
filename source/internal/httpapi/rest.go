package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"example.edu/6cc545/cw1/internal/state"
	"io"
	"mime"
	"net/http"
)

const MaxBody = 4096

// DecodeOne enforces a single bounded JSON object. No state is mutated here.

func DecodeOne(w http.ResponseWriter, r *http.Request, dst any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return &InputError{http.StatusUnsupportedMediaType, "unsupported_media_type"}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxBody)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return &InputError{http.StatusRequestEntityTooLarge, "body_too_large"}
		}
		return &InputError{http.StatusBadRequest, "invalid_json"}
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return &InputError{http.StatusBadRequest, "invalid_json"}
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return &InputError{http.StatusBadRequest, "invalid_json"}
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return &InputError{http.StatusBadRequest, "invalid_json"}
	}

	return nil
}

type InputError struct {
	Status int
	Label  string
}

func (e *InputError) Error() string { return e.Label }
func RejectInput(w http.ResponseWriter, e error) {
	var in *InputError
	if errors.As(e, &in) {
		ErrorJSON(w, in.Status, in.Label)
		return
	}
	ErrorJSON(w, 400, "invalid_input")
}

func StatusResource(s *state.Status) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			WriteJSON(w, http.StatusOK, map[string]string{
				"status": s.Get(),
			})

		case http.MethodPost:
			var input struct {
				Status string `json:"status"`
			}

			if err := DecodeOne(w, r, &input); err != nil {
				RejectInput(w, err)
				return
			}

			if err := s.Set(input.Status); err != nil {
				ErrorJSON(w, http.StatusBadRequest, "invalid_status")
				return
			}

			WriteJSON(w, http.StatusOK, map[string]string{
				"status": s.Get(),
			})

		default:
			w.Header().Set("Allow", "GET, POST")
			ErrorJSON(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	}
}

// Preserve imports while student bodies are replaced. These declarations are not solutions.
var _ = bytes.NewReader
var _ = json.NewDecoder
var _ = io.EOF
var _ = mime.ParseMediaType
