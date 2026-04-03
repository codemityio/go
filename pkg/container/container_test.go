package container

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ServiceContainer struct {
	*DefaultContainer
}

func newContainer(t *testing.T) (*ServiceContainer, ShutdownFunc) {
	t.Helper()

	c, shutdown := New(
		WithInfoFunc(func(info string) {
			t.Log(info)
		}),
		WithWarnFunc(func(warn string) {
			t.Log(warn)
		}),
		WithShutdownTimeout(10*time.Second),
	)

	return &ServiceContainer{c}, shutdown
}

func (c *ServiceContainer) ServiceOne() ServiceOne {
	id := "service-one"

	s, ok := c.Get(id).(ServiceOne)
	if ok {
		return s
	}

	ss := ServiceOne{Name: "one"}

	c.Set(id, ss, func(any) {
		time.Sleep(2 * time.Second)
	})

	return ss
}

func (c *ServiceContainer) ServiceTwo() ServiceTwo {
	id := "service-two"

	s, ok := c.Get(id).(ServiceTwo)
	if ok {
		return s
	}

	ss := ServiceTwo{
		Name: "two",
		Dep:  c.ServiceOne(),
	}

	c.Set(id, ss, func(any) {
		time.Sleep(1 * time.Second)
	})

	return ss
}

// services

type ServiceOne struct {
	Name string
}

type ServiceTwo struct {
	Dep  ServiceOne
	Name string
}

func TestNew(t *testing.T) {
	t.Parallel()

	cont, shutdown := newContainer(t)
	defer shutdown()

	assert.NotNil(t, cont)
	assert.NotNil(t, shutdown)
}

func TestGet(t *testing.T) {
	t.Parallel()

	cont, shutdown := newContainer(t)
	defer shutdown()

	assert.Equal(t, "one", cont.ServiceOne().Name)
	assert.Equal(t, "one", cont.ServiceOne().Name) // second call returns cached

	assert.Nil(t, cont.Get("nonexistent"))
}

func TestSet(t *testing.T) {
	t.Parallel()

	cont, shutdown := newContainer(t)
	defer shutdown()

	_ = cont.ServiceOne()
	_ = cont.ServiceTwo()

	assert.Equal(t, "one", cont.ServiceOne().Name)
	assert.Equal(t, "two", cont.ServiceTwo().Name)
	assert.Equal(t, "one", cont.ServiceTwo().Dep.Name)
}

func TestKeys(t *testing.T) {
	t.Parallel()

	cont, shutdown := newContainer(t)
	defer shutdown()

	assert.Empty(t, cont.Keys())

	_ = cont.ServiceOne()
	_ = cont.ServiceTwo()

	keys := cont.Keys()

	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "service-one")
	assert.Contains(t, keys, "service-two")
}

func TestShutdown(t *testing.T) {
	t.Parallel()

	cont, shutdown := newContainer(t)

	_ = cont.ServiceOne()
	_ = cont.ServiceTwo()

	require.NotEmpty(t, cont.Keys())

	start := time.Now()

	shutdown()

	elapsed := time.Since(start)

	// services sleep 2s and 1s concurrently so total should be ~2s not 3s
	assert.Less(t, elapsed, 3*time.Second)
	assert.Empty(t, cont.Keys())
}

func TestShutdownIdempotent(t *testing.T) {
	t.Parallel()

	_, shutdown := newContainer(t)

	shutdown()
	shutdown() // must not panic or deadlock
}

func TestShutdownTimeout(t *testing.T) {
	t.Parallel()

	c, shutdown := New(
		WithInfoFunc(func(info string) { t.Log(info) }),
		WithWarnFunc(func(warn string) { t.Log(warn) }),
		WithShutdownTimeout(1*time.Second),
	)

	c.Set("slow-service", struct{}{}, func(any) {
		time.Sleep(5 * time.Second) // longer than timeout
	})

	start := time.Now()

	shutdown()

	elapsed := time.Since(start)

	// should not wait full 5s, timeout at 1s
	assert.Less(t, elapsed, 3*time.Second)
}
