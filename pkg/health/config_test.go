package health

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithPeriodicCheckConfig(t *testing.T) {
	// Arrange
	expectedName := "test"
	cfg := checkerConfig{checks: map[string]*Check{}}
	interval := 5 * time.Second
	initialDelay := 1 * time.Minute
	check := Check{Name: expectedName, updateInterval: interval, initialDelay: initialDelay}

	// Act
	WithPeriodicCheck(interval, initialDelay, check)(&cfg)

	// Assert
	assert.Len(t, cfg.checks, 1)
	assert.True(t, reflect.DeepEqual(check, *cfg.checks[expectedName]))
}

func TestWithCheckConfig(t *testing.T) {
	// Arrange
	expectedName := "test"
	cfg := checkerConfig{checks: map[string]*Check{}}
	check := Check{Name: "test"}

	// Act
	WithCheck(check)(&cfg)

	// Assert
	require.Len(t, cfg.checks, 1)
	assert.True(t, reflect.DeepEqual(&check, cfg.checks[expectedName]))
}

func TestWithChecksConfig(t *testing.T) {
	// Arrange
	expectedNames := []string{"test1", "test2"}
	cfg := checkerConfig{checks: map[string]*Check{}}
	checks := []Check{
		{Name: "test1"},
		{Name: "test2"},
	}

	// Act
	WithChecks(checks...)(&cfg)

	// Assert
	require.Len(t, cfg.checks, 2)
	for i, name := range expectedNames {
		assert.True(t, reflect.DeepEqual(&checks[i], cfg.checks[name]))
	}
}

func TestWithCacheDurationConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}
	duration := 5 * time.Hour

	// Act
	WithCacheDuration(duration)(&cfg)

	// Assert
	assert.Equal(t, duration, cfg.cacheTTL)
}

func TestWithDisabledCacheConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}

	// Act
	WithDisabledCache()(&cfg)

	// Assert
	assert.Equal(t, 0*time.Second, cfg.cacheTTL)
}

func TestWithTimeoutStartConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}

	// Act
	WithTimeout(5 * time.Hour)(&cfg)

	// Assert
	assert.Equal(t, 5*time.Hour, cfg.timeout)
}

func TestWithDisabledDetailsConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}

	// Act
	WithDisabledDetails()(&cfg)

	// Assert
	assert.True(t, cfg.detailsDisabled)
}

func TestWithMiddlewareConfig(t *testing.T) {
	// Arrange
	cfg := HandlerConfig{}
	mw := func(MiddlewareFunc) MiddlewareFunc {
		return func(r *http.Request) Result {
			return Result{nil, StatusUp, nil}
		}
	}

	// Act
	WithMiddleware(mw)(&cfg)

	// Assert
	assert.Len(t, cfg.middleware, 1)
}

func TestWithInterceptorConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}
	interceptor := func(InterceptorFunc) InterceptorFunc {
		return func(ctx context.Context, name string, state CheckState) CheckState {
			return CheckState{}
		}
	}

	// Act
	WithInterceptors(interceptor)(&cfg)

	// Assert
	assert.Len(t, cfg.interceptors, 1)
}

func TestWithResultWriterConfig(t *testing.T) {
	// Arrange
	cfg := HandlerConfig{}
	w := resultWriterMock{}

	// Act
	WithResultWriter(&w)(&cfg)

	// Assert
	assert.Equal(t, &w, cfg.resultWriter)
}

func TestWithStatusChangeListenerConfig(t *testing.T) {
	// Arrange
	cfg := checkerConfig{}

	// Act
	// Use of non standard AvailabilityStatus codes.
	WithStatusListener(func(ctx context.Context, state State) {})(&cfg)

	// Assert
	assert.NotNil(t, cfg.statusChangeListener)
	// Not possible in Go to compare functions.
}

func TestNewWithDefaults(t *testing.T) {
	// Arrange
	configApplied := false
	opt := func(config *checkerConfig) { configApplied = true }

	// Act
	checker := NewChecker(opt)

	// Assert
	ckr, _ := checker.(*defaultChecker)
	assert.Equal(t, 1*time.Second, ckr.cfg.cacheTTL)
	assert.Equal(t, 10*time.Second, ckr.cfg.timeout)
	assert.True(t, configApplied)
}

func TestNewCheckerWithDefaults(t *testing.T) {
	// Arrange
	configApplied := false
	opt := func(config *checkerConfig) { configApplied = true }

	// Act
	checker := NewChecker(opt)

	// Assert
	ckr, _ := checker.(*defaultChecker)
	assert.Equal(t, 1*time.Second, ckr.cfg.cacheTTL)
	assert.Equal(t, 10*time.Second, ckr.cfg.timeout)
	assert.True(t, configApplied)
}

func TestCheckerAutostartConfig(t *testing.T) {
	// Arrange + Act
	c := NewChecker()

	// Assert
	assert.True(t, c.IsStarted())
}

func TestCheckerAutostartDisabledConfig(t *testing.T) {
	// Arrange
	c := NewChecker(WithDisabledAutostart())

	// Assert
	assert.False(t, c.IsStarted())
}

func TestWithChecks(t *testing.T) {
	// Arrange
	check := Check{Name: "test"}

	// Act
	checker := NewChecker(WithChecks(check))

	// Assert
	ckr, _ := checker.(*defaultChecker)
	assert.Len(t, ckr.cfg.checks, 1)
	assert.Contains(t, ckr.cfg.checks, check.Name)
}
