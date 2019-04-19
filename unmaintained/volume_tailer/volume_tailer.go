package main

import (
	"flag"
	"log"
	"time"

	"github.com/chrislusf/seaweedfs/weed/operation"
	"github.com/chrislusf/seaweedfs/weed/security"
	weed_server "github.com/chrislusf/seaweedfs/weed/server"
	"github.com/chrislusf/seaweedfs/weed/storage"
	"github.com/spf13/viper"
	"golang.org/x/tools/godoc/util"
)

var (
	master         = flag.String("master", "localhost:9333", "master server host and port")
	volumeId       = flag.Int("volumeId", -1, "a volume id")
	rewindDuration = flag.Duration("rewind", -1, "rewind back in time. -1 means from the first entry. 0 means from now.")
	timeoutSeconds = flag.Int("timeoutSeconds", 0, "disconnect if no activity after these seconds")
	showTextFile   = flag.Bool("showTextFile", false, "display textual file content")
)

func main() {
	flag.Parse()

	weed_server.LoadConfiguration("security", false)
	grpcDialOption := security.LoadClientTLS(viper.Sub("grpc"), "client")

	vid := storage.VolumeId(*volumeId)

	var sinceTimeNs int64
	if *rewindDuration == 0 {
		sinceTimeNs = time.Now().UnixNano()
	} else if *rewindDuration == -1 {
		sinceTimeNs = 0
	} else if *rewindDuration > 0 {
		sinceTimeNs = time.Now().Add(-*rewindDuration).UnixNano()
	}

	err := storage.TailVolume(*master, grpcDialOption, vid, uint64(sinceTimeNs), *timeoutSeconds, func(n *storage.Needle) (err error) {
		if n.Size == 0 {
			println("-", n.String())
			return nil
		} else {
			println("+", n.String())
		}

		if *showTextFile {

			data := n.Data
			if n.IsGzipped() {
				if data, err = operation.UnGzipData(data); err != nil {
					return err
				}
			}
			if util.IsText(data) {
				println(string(data))
			}

			println("-", n.String(), "compressed", n.IsGzipped(), "original size", len(data))
		}
		return nil
	})

	if err != nil {
		log.Printf("Error VolumeTail volume %d: %v", vid, err)
	}

}
