package watcher

import (
	"context"
	"server-api/global"

	"github.com/gomooth/pkg/framework/app"
	"github.com/save95/xerror"

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

func (s *server) Start() error {
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

				global.Log.Debugf("file(%s) %s...", event.Name, event.Op)
				if event.Name == s.filename {
					if err := s.f(event); nil != err {
						global.Log.Errorf("file(%s) charged, watch handle failed: %+v", event.Name, err)
					}
				}
			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				global.Log.Errorf("file watch error: %+v", err)
			}
		}
	}()

	if err := s.watcher.Add(s.filename); nil != err {
		return err
	}

	global.Log.Info("file watch server starting...")
	return nil
}

func (s *server) Shutdown() error {
	return s.watcher.Close()
}
