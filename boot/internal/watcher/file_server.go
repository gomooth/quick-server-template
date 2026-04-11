package watcher

import (
	"context"
	"log/slog"

	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/xerror"

	"github.com/fsnotify/fsnotify"
)

type server struct {
	ctx      context.Context
	filename string
	f        func(ev fsnotify.Event) error

	watcher *fsnotify.Watcher
}

func NewFileServer(ctx context.Context, filename string, cb func(ev fsnotify.Event) error) app.IApp {
	return &server{
		ctx:      ctx,
		filename: filename,
		f:        cb,
	}
}

func (s *server) Start(_ context.Context) error {
	if len(s.filename) == 0 {
		return xerror.New("filename is empty")
	}

	var err error
	if s.watcher, err = fsnotify.NewWatcher(); err != nil {
		return xerror.Wrap(err, "file watch failed")
	}

	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}

				slog.Debug("file event", "name", event.Name, "op", event.Op)
				if event.Name == s.filename {
					if err := s.f(event); nil != err {
						slog.Error("file charged, watch handle failed", "name", event.Name, "err", err)
					}
				}
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				slog.Error("file watch error", "err", err)
			}
		}
	}()

	if err := s.watcher.Add(s.filename); nil != err {
		return err
	}

	slog.Info("file watch server starting...")
	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	return s.watcher.Close()
}
