package container

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// DefaultContainer is a service container structure.
type DefaultContainer struct {
	infoFunc        LogFunc
	warnFunc        LogFunc
	signalChannel   chan os.Signal
	shutdownTimeout time.Duration
	services        sync.Map
	shutdownOnce    sync.Once
}

// New factory function.
func New(options ...Option) (*DefaultContainer, ShutdownFunc) {
	container := &DefaultContainer{
		infoFunc:        func(string) {},
		warnFunc:        func(string) {},
		signalChannel:   make(chan os.Signal, 1),
		shutdownTimeout: 0,
		services:        sync.Map{},
		shutdownOnce:    sync.Once{},
	}

	for i := range options {
		options[i](container)
	}

	container.handleOSSignals()

	return container, container.shutdown
}

// Get method to get the service.
func (c *DefaultContainer) Get(key string) any {
	value, ok := c.services.Load(key)
	if !ok {
		return nil
	}

	rec, ok := value.(ServiceRecord)
	if !ok {
		c.warnFunc("invalid assertion: unable to get `ServiceRecord`")

		return nil
	}

	return rec.Service
}

// Set method to set the service in the container.
func (c *DefaultContainer) Set(key string, service any, ss func(service any)) {
	c.services.Store(key, ServiceRecord{
		Service:  service,
		Shutdown: ss,
	})

	c.infoFunc(fmt.Sprintf("container %s service stored", key))
}

// Keys list all service keys.
func (c *DefaultContainer) Keys() []string {
	var ids []string

	c.services.Range(func(key, _ any) bool {
		id, ok := key.(string)
		if !ok {
			c.warnFunc("invalid assertion, unable to assert that a key is a string")

			return false
		}

		ids = append(ids, id)

		return true
	})

	return ids
}

func (c *DefaultContainer) handleOSSignals() {
	signal.Notify(c.signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		osCall, ok := <-c.signalChannel
		if !ok {
			return
		}

		c.infoFunc(fmt.Sprintf("container detected system call %+v", osCall))

		signal.Stop(c.signalChannel)
		close(c.signalChannel)

		c.shutdown()
	}()
}

func (c *DefaultContainer) shutdown() {
	c.shutdownOnce.Do(func() {
		c.infoFunc(fmt.Sprintf(
			"container shutdown started with %s timeout duration",
			c.shutdownTimeout.String(),
		))

		ctx, cancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
		defer cancel()

		var wg sync.WaitGroup

		c.services.Range(func(key, value any) bool {
			wg.Go(func() {
				record, ok := value.(ServiceRecord)
				if !ok {
					c.warnFunc("invalid service record")

					return
				}

				c.infoFunc(fmt.Sprintf("container service %s shutdown started", key))

				if record.Shutdown != nil {
					record.Shutdown(record.Service)
				}

				c.services.Delete(key)

				c.infoFunc(fmt.Sprintf("container service %s shutdown completed", key))
			})

			return true
		})

		done := make(chan struct{})

		go func() {
			defer close(done)

			wg.Wait()
		}()

		select {
		case <-done:
			c.infoFunc("container shutdown completed - now exit")
		case <-ctx.Done():
			c.warnFunc("container shutdown timed out")
		}
	})
}
