package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nu-kotov/GophKeeper/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestRequestLogger_BasicFlow(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger.Log = zap.New(core)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	logged := RequestLogger(handler)
	logged.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "OK", string(body))

	logOutput := buf.String()
	assert.Contains(t, logOutput, `"method":"GET"`)
	assert.Contains(t, logOutput, `"uri":"/test"`)
	assert.Contains(t, logOutput, `"status":200`)
	assert.Contains(t, logOutput, `"size":2`)
}

func TestRequestLogger_HandlerError(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger.Log = zap.New(core)

	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal", http.StatusInternalServerError)
	}

	req := httptest.NewRequest(http.MethodPost, "/fail", nil)
	rec := httptest.NewRecorder()

	RequestLogger(handler).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, `"method":"POST"`)
	assert.Contains(t, logOutput, `"uri":"/fail"`)
	assert.Contains(t, logOutput, `"status":500`)
	assert.Contains(t, logOutput, `"size"`)
}

func TestRequestLogger_LogsDuration(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger.Log = zap.New(core)

	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	rec := httptest.NewRecorder()

	RequestLogger(handler).ServeHTTP(rec, req)

	logOutput := buf.String()
	assert.Contains(t, logOutput, `"duration"`)
	assert.Contains(t, logOutput, `"uri":"/slow"`)
}
