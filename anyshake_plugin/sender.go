package main

/*
#cgo CFLAGS: -I ../../libs/plugin
#cgo LDFLAGS: -L ./ -lplugin
#include <stdlib.h>
#include "plugin.h"
*/
import "C"
import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

func sendMessage(message Message) {
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

	networkStation := fmt.Sprintf("%s.%s", message.Network, message.Station)
	networkStationCString := C.CString(networkStation)
	defer C.free(unsafe.Pointer(networkStationCString))

	channelCString := C.CString(message.Channel)
	defer C.free(unsafe.Pointer(channelCString))

	C.send_raw3(networkStationCString, channelCString, pTime, C.int(0), C.int(100), data, sampleRate)
	log.Printf(
		"1 message sent, station %s, channel: %s, sample rate: %d Hz, time: %s",
		networkStation, message.Channel, message.SampleRate, message.Time.Format(time.RFC3339Nano),
	)
}
