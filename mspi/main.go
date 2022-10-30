package mspi

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
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

var requestRegexp *regexp.Regexp

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
	ms.AddHandler("/pi/", piHandler)
	requestRegexp = regexp.MustCompile("^.*/(?P<iterations>\\d+)(?:[.@/#-](?P<precision>\\d+))$")
	return ms
}

// ---------------------------------------------------------------------------

func (ms *MsPi) httpGetPi(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var iterations uint
	var precision uint
	var value string
	var err error

	iterations = 100
	precision = 10

	if r.URL.Path != "/pi" && r.URL.Path != "/pi/" {

		match := requestRegexp.FindStringSubmatch(r.URL.Path)

		if len(match) == 0 {
			status = http.StatusInternalServerError
			msg = "URL has bad format"
			response := NewPiResponse(status, msg, ms)
			w.WriteHeader(status)
			contentLen = ms.Reply(w, response)
			return status, contentLen, msg
		}

		if len(match) >= 1 {
			var i int
			i, err = strconv.Atoi(match[1])

			if err != nil {
				msg = fmt.Sprintf("Failed to read the number of iteration, error was '%s'!", err.Error())
				http.Error(w, msg, http.StatusBadRequest)
				return http.StatusBadRequest, 0, msg
			}

			iterations = uint(i)
		}

		if len(match) >= 3 {
			var i int
			i, err = strconv.Atoi(match[2])

			if err != nil {
				msg = fmt.Sprintf("Failed to read the precision, error was '%s'!", err.Error())
				http.Error(w, msg, http.StatusBadRequest)
				return http.StatusBadRequest, 0, msg
			}

			precision = uint(i)
		}
	}

	p := Pi(iterations, precision)
	value = p.Text('f', int(precision))

	status = http.StatusOK
	name := r.Header.Get("X-Cid")
	version := r.Header.Get("X-Version")
	environment := r.Header.Get("X-Environment")
	msg = fmt.Sprintf("'v%s' in '%s' Generated pi with '%d' iteration for client '%s@%s'.", ms.GetVersion(), environment, iterations, name, version)
	response := NewPiResponse(status, msg, ms)
	response.Value = value
	response.Iterations = int(iterations)
	response.Precision = int(precision)
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}

func Pi(iterations uint, precision uint) *big.Float {

	k := uint(0)
	pi := new(big.Float).SetPrec(precision).SetFloat64(0)
	k1k2k3 := new(big.Float).SetPrec(precision).SetFloat64(0)
	k4k5k6 := new(big.Float).SetPrec(precision).SetFloat64(0)
	temp := new(big.Float).SetPrec(precision).SetFloat64(0)
	minusOne := new(big.Float).SetPrec(precision).SetFloat64(-1)
	total := new(big.Float).SetPrec(precision).SetFloat64(0)

	two2Six := math.Pow(2, 6)
	two2SixBig := new(big.Float).SetPrec(precision).SetFloat64(two2Six)

	for k < iterations {
		t1 := float64(float64(1) / float64(10*k+9))
		k1 := new(big.Float).SetPrec(precision).SetFloat64(t1)
		t2 := float64(float64(64) / float64(10*k+3))
		k2 := new(big.Float).SetPrec(precision).SetFloat64(t2)
		t3 := float64(float64(32) / float64(4*k+1))
		k3 := new(big.Float).SetPrec(precision).SetFloat64(t3)
		k1k2k3.Sub(k1, k2)
		k1k2k3.Sub(k1k2k3, k3)

		t4 := float64(float64(4) / float64(10*k+5))
		k4 := new(big.Float).SetPrec(precision).SetFloat64(t4)
		t5 := float64(float64(4) / float64(10*k+7))
		k5 := new(big.Float).SetPrec(precision).SetFloat64(t5)
		t6 := float64(float64(1) / float64(4*k+3))
		k6 := new(big.Float).SetPrec(precision).SetFloat64(t6)
		k4k5k6.Add(k4, k5)
		k4k5k6.Add(k4k5k6, k6)
		k4k5k6 = k4k5k6.Mul(k4k5k6, minusOne)
		temp.Add(k1k2k3, k4k5k6)

		k7temp := new(big.Int).Exp(big.NewInt(-1), big.NewInt(int64(k)), nil)
		k8temp := new(big.Int).Exp(big.NewInt(1024), big.NewInt(int64(k)), nil)

		k7 := new(big.Float).SetPrec(precision).SetFloat64(0)
		k7.SetInt(k7temp)
		k8 := new(big.Float).SetPrec(precision).SetFloat64(0)
		k8.SetInt(k8temp)

		t9 := float64(256) / float64(10*k+1)
		k9 := new(big.Float).SetPrec(precision).SetFloat64(t9)
		k9.Add(k9, temp)
		total.Mul(k9, k7)
		total.Quo(total, k8)
		pi.Add(pi, total)

		k = k + 1
	}
	pi.Quo(pi, two2SixBig)
	return pi
}
