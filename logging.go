package kitmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Errorer interface.
type Errorer interface {
	Err() error
}

// RequestValuer interface.
type RequestValuer interface {
	KeyValues(req interface{}, keyvals ...interface{}) []interface{}
}

// NewLogging create a new endpoint.Middleware instance used to log request.
func NewLogging(logger log.Logger, val RequestValuer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
			defer func(begin time.Time) {
				keyvals := []interface{}{"took", time.Since(begin)}
				if err != nil {
					keyvals = append(keyvals, "err", err)
				}
				if r, ok := resp.(Errorer); ok && r.Err() != nil {
					keyvals = append(keyvals, "err", r.Err())
				}

				logger.Log(val.KeyValues(request, keyvals...)...)
			}(time.Now())

			resp, err = next(ctx, request)
			return
		}
	}
}

// NewDefaultRequestValuer create a default RequestValuer implementation.
func NewDefaultRequestValuer() RequestValuer {
	return &defaultValuer{key: "req"}
}

type defaultValuer struct {
	key string
}

func (f *defaultValuer) KeyValues(req interface{}, keyvals ...interface{}) []interface{} {
	if req == nil {
		return keyvals
	}

	b := bytes.NewBuffer([]byte(`[`))

	t := reflect.TypeOf(req)
	v := reflect.ValueOf(req)

	var values []string

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			tf := t.Field(i)
			name, ok := tf.Tag.Lookup("val")
			if ok && name == "-" {
				continue
			} else if !ok || (ok && name == "") {
				name = tf.Name
			}

			msg := fmt.Sprintf("%s:%v", name, v.Field(i))
			values = append(values, msg)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			v.MapIndex(key)
			msg := fmt.Sprintf("%s:%v", key, v.MapIndex(key))
			values = append(values, msg)
		}
	}

	b.WriteString(strings.Join(values, " "))
	b.WriteString("]")

	return append(keyvals, f.key, b.String())
}
