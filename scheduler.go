package pullbuddy

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type scheduler struct {
	images                 []image
	lock                   sync.RWMutex
	processing             bool
	scheduleChan           chan imageID
	processingFinishedChan chan processingImageResult
	puller                 imagePuller
	logger                 logrus.FieldLogger
}

func newScheduler() *scheduler {
	return &scheduler{
		images:                 make([]image, 0),
		scheduleChan:           make(chan imageID),
		processingFinishedChan: make(chan processingImageResult),
	}
}

func (sch *scheduler) list() []image {
	sch.logger.Info("list images")
	locker := sch.lock.RLocker()
	locker.Lock()
	defer locker.Unlock()
	images := make([]image, 0, len(sch.images))
	copy(images, sch.images)
	return images
}

func (sch *scheduler) schedule(id string) {
	sch.logger.WithField("image_id", id).Info("schedule image")
	sch.scheduleChan <- imageID(id)
}

func (sch *scheduler) run() {
	for {
		select {
		case id := <-sch.scheduleChan:
			sch.scheduleFromChannel(id)
		case result := <-sch.processingFinishedChan:
			sch.handleProcessResult(result)
		}
		sch.nextImage()
	}
}

func (sch *scheduler) scheduleFromChannel(id imageID) {
	sch.lock.Lock()
	defer sch.lock.Unlock()
	for i := range sch.images {
		if sch.images[i].id == id && !sch.images[i].status.Done() {
			sch.logger.WithField("image_id", id).Info("image already scheduled")
			return
		}
	}
	sch.images = append(
		sch.images,
		image{
			id:     id,
			status: pending,
		},
	)
	sch.logger.WithField("image_id", id).Info("image scheduled")
}

func (sch *scheduler) nextImage() {
	if sch.processing {
		return
	}
	sch.lock.Lock()
	defer sch.lock.Unlock()
	for i := range sch.images {
		if sch.images[i].status == pending {
			sch.logger.WithField("image_id", sch.images[i].id).Info("pull image")
			go sch.processImage(sch.images[i].id)
			sch.images[i].status = pulling
			sch.processing = true
			return
		}
	}
}

func (sch *scheduler) processImage(id imageID) {
	sch.processingFinishedChan <- processingImageResult{
		id:  id,
		err: sch.puller.pull(string(id)),
	}
}

func (sch *scheduler) handleProcessResult(result processingImageResult) {
	logger := sch.logger.WithField("image_id", result.id)
	if result.err != nil {
		logger = logger.WithField("pull_error", result.err)
	}
	logger.Info("handled image")
	sch.processing = false
	sch.lock.Lock()
	defer sch.lock.Unlock()
	for i := range sch.images {
		if sch.images[i].id == result.id && !sch.images[i].status.Done() {
			if result.err == nil {
				sch.images[i].status = finished
			} else {
				sch.images[i].status = failed
				sch.images[i].err = result.err
			}
		}
	}
}

type processingImageResult struct {
	id  imageID
	err error
}

type image struct {
	id     imageID
	status imageStatus
	err    error
}

type imageID string

type imageStatus string

const (
	pending  imageStatus = "pending"
	pulling              = "pulling"
	finished             = "finished"
	failed               = "failed"
)

func (status imageStatus) Done() bool {
	return status == finished || status == failed
}

type imagePuller interface {
	pull(id string) error
}
