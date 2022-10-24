package mspi

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsPi
// ###########################################################################
// ###########################################################################

// MsPi Encapsulates the ms-pi data
type MsPi struct {
	microservice.MicroService

	seededRand *rand.Rand
	piMutex    sync.Mutex
}

// ###########################################################################

// InitMsPiFromArgs ...
func InitFromArgs(ms *MsPi, args []string, flagset *flag.FlagSet) *MsPi {
	var cfg Configuration

	if flagset == nil {
		flagset = flag.NewFlagSet("ms-pi", flag.PanicOnError)
	}

	InitConfigurationFromArgs(&cfg, args, flagset)
	microservice.Init(&ms.MicroService, &cfg.Configuration, nil)
	ms.seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	piHandler := ms.DefaultHandler()
	piHandler.Get = ms.httpGetPi
	ms.AddHandler("/pi", piHandler)
	return ms
}

// ---------------------------------------------------------------------------

var deviceMutex sync.Mutex

func (ms *MsPi) httpGetPi(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	ms.piMutex.Lock()
	value := ms.seededRand.Intn(100)
	ms.piMutex.Unlock()
	time.sleep(10)
	status = http.StatusOK
	name := r.Header.Get("X-Cid")
	version := r.Header.Get("X-Version")
	environment := r.Header.Get("X-Environment")
	msg = fmt.Sprintf("'v%s' in '%s' Generated pi number '%d' for client '%s@%s'.", ms.GetVersion(), environment, value, name, version)
	response := NewPiResponse(status, msg, ms)
	response.Value = value
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}
