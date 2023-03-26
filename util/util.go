package util

import (
	"demo/logger"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func GetRedisClusterClient(addrs []string, passwd string) *redis.ClusterClient {
	rcc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: passwd,
	})
	_, err := rcc.Ping(context.Background()).Result()
	if err != nil {
		logger.D.Error("new redis cluster client error: %v", err)
	}
	return rcc
}

func GetNowMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

func Recover(msg string) {
	if p := recover(); p != nil {
		logger.D.Error("[%v] Recoverd from panic: %v", msg, p)
		var buf [8192]byte
		n := runtime.Stack(buf[:], false)
		logger.D.Error("==> %v", string(buf[:n]))
	}
}

func ThisFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func GetFuncName(fc interface{}) string {
	pc := reflect.ValueOf(fc).Pointer()
	return runtime.FuncForPC(pc).Name()
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func GetRandomStr(n int) string {
	res := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx >= 0 && idx < len(letterBytes) {
			res[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(res)
}

func TickerLoop(fc func(), period int64, numLimit int64) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(period))
	fcName := ThisFuncName()
	go func() {
		defer Recover(fcName)
		var idx int64 = 0
		for range ticker.C {
			idx++
			if numLimit > 0 && idx > numLimit {
				logger.D.Info("hit loop number limit [%v]", numLimit)
				break
			}
			startTime := time.Now()
			logger.D.Trace("doing loop [%v] for func [%v]", idx, GetFuncName(fc))
			fc()
			logger.D.Debug("done  loop [%v] for func [%v] cost [%v]", idx, GetFuncName(fc), time.Now().Sub(startTime))
		}
	}()
}

func PostSleepLoop(fc func(), period int64, numLimit int64) {
	fcName := ThisFuncName()
	go func() {
		defer Recover(fcName)
		var idx int64 = 0
		for {
			idx++
			if numLimit > 0 && idx > numLimit {
				logger.D.Info("hit loop number limit [%v]", numLimit)
				break
			}
			startTime := time.Now()
			logger.D.Trace("doing loop [%v] for func [%v]", idx, GetFuncName(fc))
			fc()
			logger.D.Debug("done  loop [%v] for func [%v] cost [%v]", idx, GetFuncName(fc), time.Now().Sub(startTime))
			time.Sleep(time.Millisecond * time.Duration(period))
		}
	}()
}

func HandleSignal(fc func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer Recover(ThisFuncName())
		sig := <-sigs
		logger.D.Error("got signal [%d] [%+v]", sig, sig)
		fc()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

// meter
func GetGeoDistance(latitudeA, longitudeA, latitudeB, longitudeB float64) int {
	radLatA := float64(math.Pi * latitudeA / 180)
	radLatB := float64(math.Pi * latitudeB / 180)
	theta := float64(longitudeA - longitudeB)
	radTheta := float64(math.Pi * theta / 180)
	dist := math.Sin(radLatA)*math.Sin(radLatB) + math.Cos(radLatA)*math.Cos(radLatB)*math.Cos(radTheta)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344 * 1000
	return int(dist)
}
