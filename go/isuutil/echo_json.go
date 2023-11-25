package isuutil

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"net/http"
)

// EchoJSONSerializer はgoccy/go-jsonを使ったEchoのJSONシリアライザです。
// https://twitter.com/fujiwara/status/1440211187581341699
type EchoJSONSerializer struct{}

var _ echo.JSONSerializer = (*EchoJSONSerializer)(nil)

func (j *EchoJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	return enc.Encode(i)
}
func (j *EchoJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	}
	var se *json.SyntaxError
	if errors.As(err, &se) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
