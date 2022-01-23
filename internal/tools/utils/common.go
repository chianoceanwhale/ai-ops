package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetCurrntTime() time.Time {
	return time.Now()
}

func GetCurrntTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetCurrntTimeStr2() string {
	return time.Now().Format("20060102_1504")
}

func GetCurrntTimeStr3() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrntTimeStr4() string {
	return time.Now().Format("200601021504")
}

// 获取当前时间戳
func GetTimeUnix() int64 {
	return time.Now().Unix()
}

// MD5 方法
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func StrToInt(err error, index string) int {
	i, err := strconv.Atoi(index)
	return i
}

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func GetFileAsString(name string) string {
	buf, err := ioutil.ReadFile("../resources/" + name)
	if err != nil {
		panic(err)
	}

	return string(buf)
}

func TimestampToString(str string) (string, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return "", err
	}

	t := time.Unix(0, i*int64(time.Millisecond)).Format("2006-01-02 15:04:00.000")
	return t, nil
}

func TimestampIntToString(str int64) string {
	t := time.Unix(0, str*int64(time.Millisecond)).Format("2006-01-02 15:04:00.000")
	return t
}

func StringToTimestamp(str string) int64 {
	t, _ := time.Parse(str, "2006-01-02 15:04:00.000")
	return t.Unix()
}

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//  E: Exabyte
//  P: Petabyte
//  T: Terabyte
//  G: Gigabyte
//  M: Megabyte
//  K: Kilobyte
//  B: Byte
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "E"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "P"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

func ByteToMb(bytes int64) string {
	value := float64(bytes) / MEGABYTE

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func TimestampSubNow(tm string) (sb float64, err error) {
	i, err := strconv.ParseInt(tm, 10, 64)
	if err != nil {
		return 0, err
	}

	t := time.Unix(i/1000, 0)
	f := t.Sub(time.Now())
	sb = math.Abs(f.Hours())
	return
}

//get abs path according the path and filename
func GetFileAbsPath(dir, filename string) string {
	var fileAbsPath string
	sysType := runtime.GOOS
	switch sysType {
	case "linux":
		fileAbsPath = filepath.ToSlash(filepath.Join(dir, filename))
	case "windows":
		fileAbsPath = filepath.Join(dir, filename)
	}
	return fileAbsPath
}
