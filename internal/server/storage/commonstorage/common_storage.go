package commonstorage

import (
	log "github.com/sirupsen/logrus"
)

type CommonStorage struct {
	Events   map[string][]chan func()
	OnUpdate []func()
}

func New() *CommonStorage {
	defaultStore := CommonStorage{
		Events: map[string][]chan func(){
			"OnUpdate": {
				make(chan func()),
			},
		},
		OnUpdate: []func(){},
	}

	return &defaultStore
}

func (s *CommonStorage) ListenEvents() {
	go func() {
		chs := s.Events["OnUpdate"]
		for _, ch := range chs {
			for {
				handler := <-ch
				handler()
				log.Info("Event: OnUpdate")
			}
		}
	}()
}

func (s *CommonStorage) Emit(event string) {
	if _, ok := s.Events[event]; ok {
		for _, handler := range s.Events[event] {
			go func(handler chan func()) {
				for _, h := range s.OnUpdate {
					handler <- h
				}
			}(handler)
		}
	}
}

func (s *CommonStorage) AddOnUpdateListener(fn func()) {
	s.OnUpdate = append(s.OnUpdate, fn)
}
