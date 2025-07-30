package httpserver

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func TestHandlerErr(t *testing.T) {
	tcs := map[string]struct {
		givenErr  error
		expStatus int
		expErr    *Error
	}{
		"err_400": {
			givenErr: &Error{
				Status: http.StatusBadRequest,
				Code:   "bad_request",
				Desc:   "Bad request",
			},
			expStatus: http.StatusBadRequest,
			expErr: &Error{
				Status: http.StatusBadRequest,
				Code:   "bad_request",
				Desc:   "Bad request",
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Given:
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// When:
			HandlerErr(func(w http.ResponseWriter, r *http.Request) error {
				return tc.givenErr
			}).ServeHTTP(w, r)

			// Then:
			require.Equal(t, tc.expStatus, w.Code)
			var actErr Error
			err := ParseJSON(w.Result().Body, &actErr)
			if tc.expErr == nil {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expErr.Code, actErr.Code)
				require.Equal(t, tc.expErr.Desc, actErr.Desc)
			}
		})
	}
}
