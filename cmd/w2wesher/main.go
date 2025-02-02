package main

import (
	"context"
	"flag"
	"os"

	"github.com/derlaft/w2wesher/config"
	"github.com/derlaft/w2wesher/networkstate"
	"github.com/derlaft/w2wesher/p2p"
	"github.com/derlaft/w2wesher/runnergroup"
	"github.com/derlaft/w2wesher/wg"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("w2wesher")

var (
	configFile = flag.String("config", ".w2wesher.ini", "configuration file")
)

func main() {
	flag.Parse()

	if os.Getuid() <= 0 || os.Getegid() <= 0 {
		log.Fatal("w2wesher should never be started from root")
	}

	state := networkstate.New()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	adapter, err := wg.New(cfg, state)
	if err != nil {
		log.Fatal(err)
	}

	node, err := p2p.New(cfg, state, adapter)
	if err != nil {
		log.Fatal(err)
	}

	err = runnergroup.New(context.TODO()).
		Go(node.Run).
		Go(adapter.Run).
		Go(runnergroup.AbortOnSignal).
		Wait()
	if err != nil {
		log.Error(err)
	}
}
