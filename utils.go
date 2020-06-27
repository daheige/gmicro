package gmicro

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/peer"
)

// GetGRPCClientIP get client ip address from context
func GetGRPCClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[GetGRPCClientIp] invoke FromContext() failed")
	}

	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[GetGRPCClientIp] peer.Addr is nil")
	}

	addSlice := strings.Split(pr.Addr.String(), ":")

	return addSlice[0], nil
}

// RndUUID realizes unique uuid based on time ns and random number
// There is no duplication of uuid on a single machine
// If you want to generate non-duplicate uuid on a distributed architecture
// Just add some custom strings in front of rndStr
// Return format: eba1e8cd-0460-4910-49c6-44bdf3cf024d
func RndUUID() string {
	s := RndUUIDMd5()
	return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:])
}

// RndUUIDMd5 make an uuid
func RndUUIDMd5() string {
	ns := time.Now().UnixNano()
	rndStr := strings.Join([]string{
		strconv.FormatInt(ns, 10), strconv.FormatInt(RandInt64(1000, 9999), 10),
	}, "")

	return Md5(rndStr)
}

// RandInt64 crete a num [m,n]
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// Md5 md5 func
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
