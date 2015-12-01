package main

import (
	"log"
	"time"

	"github.com/go-fsnotify/fsevents"
)

func main() {
	dev, _ := fsevents.DeviceForPath("/Users/huangjin02/Project/Haruhi/output")
	log.Print(dev)
	log.Println(fsevents.EventIDForDeviceBeforeTime(dev, time.Now()))

	es := &fsevents.EventStream{
		Paths:   []string{"/Users/huangjin02/Project/Haruhi/output"},
		Latency: 500 * time.Millisecond,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot,
	}

	es.Start()
	ec := es.Events

	go func() {
		for msg := range ec {
			for _, event := range msg {
				logEvent(event)
			}
		}

	}()

	time.Sleep(999999 * time.Second)
}

var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	fsevents.ItemCreated:       "Created",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

func logEvent(event fsevents.Event) {
	note := ""
	for bit, description := range noteDescription {
		if event.Flags&bit == bit {
			note += description + " "
		}
	}
	log.Printf("EventID: %d Path: %s Flags: %s", event.ID, event.Path, note)
}
