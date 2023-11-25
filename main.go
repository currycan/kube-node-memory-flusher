package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
)

var interval = flag.Int("interval", 360, "Timeout in seconds")
var resetCacheOption = flag.Int("reset_cache_option", 1, "Option to reset")
var buffersLimitString = flag.String("buffers_limit", "10 MB", "Maximum cache buffer size")
var cachedLimitString = flag.String("cached_limit", "900 MB", "Maximum cached memory size")
var dropCachesFilePath = flag.String("drop_caches_file_path", "/var/host_sys_vm/drop_caches", "Mounted host file path")
var debug = flag.Bool("debug_log", false, "print debug logs")

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	// 设置 format
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	flag.Parse()

	var buffersLimit datasize.ByteSize
	var cachedLimit datasize.ByteSize

	err := buffersLimit.UnmarshalText([]byte(*buffersLimitString))
	checkErr(err)
	err = cachedLimit.UnmarshalText([]byte(*cachedLimitString))
	checkErr(err)

	for {
		v, _ := mem.VirtualMemory()

		log.Infof("Buffers: %v", (datasize.ByteSize(v.Buffers) * datasize.B).String())
		log.Infof("Cached: %v", (datasize.ByteSize(v.Cached) * datasize.B).String())

		if datasize.ByteSize(v.Buffers) > buffersLimit || datasize.ByteSize(v.Cached) > cachedLimit {
			log.Debugf("Cleaning memory...")

			err := os.WriteFile(*dropCachesFilePath, []byte(strconv.Itoa(*resetCacheOption)), 0644)
			checkErr(err)
		}

		// Sleep until interval
		log.Debugf("Sleeping for %v\n seconds", *interval)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
