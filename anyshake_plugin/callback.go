package main

/*
#cgo CFLAGS: -I ../
#cgo LDFLAGS: -L ./ -lplugin
#include <stdlib.h>
#include "plugin.h"
*/
import "C"
import (
	"log"
	"time"
	"unsafe"
)

func seisCompDaemonCallback(message Message) {
	data := (*C.int)(unsafe.Pointer(&message.Data[0]))
	sampleRate := C.int(message.SampleRate)
	pTime := &C.struct_ptime{
		year:   C.int(message.Time.Year()),
		yday:   C.int(message.Time.YearDay()),
		hour:   C.int(message.Time.Hour()),
		minute: C.int(message.Time.Minute()),
		second: C.int(message.Time.Second()),
		usec:   C.int(message.Time.Nanosecond() / 1000),
	}

	station := C.CString(message.Station)
	defer C.free(unsafe.Pointer(station))

	channel := C.CString(message.Channel)
	defer C.free(unsafe.Pointer(channel))

	C.send_raw3(station, channel, pTime, C.int(0), C.int(100), data, sampleRate)
	log.Printf(
		"1 message sent, station %s, channel: %s, sample rate: %d Hz, time: %s",
		message.Station, message.Channel, message.SampleRate, message.Time.Format(time.RFC3339Nano),
	)
}
