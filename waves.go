package main

import (
	"time"
)

type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		time.Sleep(time.Second * WavesMonitorTick)
	}
}

func initWavesMonitor() {
	wm := &WavesMonitor{}
	go wm.start()
}
