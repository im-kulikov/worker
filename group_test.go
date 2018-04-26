package worker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chapsuk/worker"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGroup(t *testing.T) {
	Convey("Given empty workers group", t, func() {
		wk := worker.NewGroup()
		So(wk, ShouldNotBeNil)

		Convey("When add 3 workers", func() {
			var counter int32

			wk.Add(
				worker.New(createFakeJob(&counter)),
				worker.New(createFakeJob(&counter)),
				worker.New(createFakeJob(&counter)),
			)

			Convey("workers should not be started", func() {
				So(atomic.LoadInt32(&counter), ShouldEqual, 0)
			})

			Convey("When run group with 3 workers", func() {
				wk.Run()
				time.Sleep(time.Second)

				Convey("all workers should be started", func() {
					So(atomic.LoadInt32(&counter), ShouldEqual, 3)
				})

				Convey("When add worker after group run", func() {
					wk.Add(worker.New(createFakeJob(&counter)))
					time.Sleep(time.Second)

					Convey("added worker should be executed", func() {
						So(atomic.LoadInt32(&counter), ShouldEqual, 4)
					})
				})

				Convey("Stop workers call", func() {
					ch := make(chan struct{})
					go func() {
						wk.Stop()
						ch <- struct{}{}
					}()

					Convey("should not be blocking", func() {
						select {
						case <-ch:
							So("non-blocking", ShouldEqual, "non-blocking")
						case <-time.NewTicker(time.Second).C:
							So("blocking", ShouldEqual, "non-blocking")
						}
					})
				})
			})
		})
	})
}

func createFakeJob(counter *int32) worker.Job {
	return func(ctx context.Context) {
		atomic.AddInt32(counter, 1)
	}
}
