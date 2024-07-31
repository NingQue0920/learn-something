package main

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"math"
	"net/url"
	"strings"
	"sync"
)

var setMu sync.RWMutex
var storeMu sync.Mutex

type ShortUrlGenerator interface {
	Generate(input string) string
	Store(input, shorten string)
}

func UrlValidator(inputUrl string) error {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", inputUrl)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %v", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("missing URL host")
	}
	domains := strings.Split(u.Host, ".")
	if len(domains) < 2 {
		return fmt.Errorf("invalid URL host: %v", u.Host)
	}
	/*  check if domain is valid , e.g. not a reserved domain
	 *	var topLevelDomain = []string{"com", "org", "net", "io", "gov", "edu", "mil", "int", "arpa"}
	 *	if !slices.Contains(topLevelDomain, domains[len(domains)-1])
	 *    return fmt.Errorf("invalid URL host: %v", u.Host)
	 */

	return nil
}
func GenerateByMurmurHash(input string) string {

	// load from cache
	setMu.RLock()
	shorten, ok := urlMapCache[input]
	setMu.RUnlock()
	if ok {
		return shorten
	}
	// generate short code
	sum32 := murmur3.Sum64([]byte(input))
	// convert to string
	//shortCode := strconv.FormatInt(int64(sum32), 36)
	shortCode := FormatInt62(sum32)
	// check if short code is already present ,using db unique index or bloom filter
	shortCode = handleCollisions(input, shortCode)

	setMu.Lock()
	urlMapCache[input] = shortCode
	setMu.Unlock()
	storeMu.Lock()
	urlSet[shortCode] = struct{}{}
	storeMu.Unlock()
	return shortCode
}

func handleCollisions(input, shortCode string) string {
	for i := 0; i < maxRetries; i++ {

		setMu.RLock()
		_, ok := urlSet[shortCode]
		setMu.RUnlock()
		if !ok {
			return shortCode
		}
		shortCode = GenerateByMurmurHash(input + string(rune(i)))
	}
	return shortCode
}

// FormatInt62 使用62位表示法，0-9a-zA-Z来表示一个字符串
func FormatInt62(i uint64) string {
	format := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result string
	for i > 0 {
		remainder := math.Mod(float64(i), 62)
		i = i / 62
		result = string(format[int(remainder)]) + result
	}
	return result
}
