package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
)

var (
	totalKeyCount     int64
	falsePositiveRate float64

	bloomFilterFile          string
	newBloomFilterFileSaveTo string

	keyFile string

	keyLen      int64
	genKeyCount int64
)


func main() {
	flag.Parse()

    filter := NewBloomFilter(uint(totalKeyCount), falsePositiveRate, bloomFilterFile)
    if filter == nil {
        return
    }

	newKeyFile, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to open file=%s err=%v\n", keyFile, err)
		return
	}
	defer newKeyFile.Close()

    if !GenKeys(filter, newKeyFile, int(keyLen)) {
        return
    }

    if !SaveBloomFile(filter, newBloomFilterFileSaveTo) {
        return
    }
}

func RandomString(n int) string {
	// base64 将+改为-，将/改为_
	var letter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func NewBloomFilter(totalKeyCount uint, falsePositiveRate float64, bloomFilterFile string) *bloom.BloomFilter{
	filter := bloom.NewWithEstimates(totalKeyCount, falsePositiveRate)
	if len(bloomFilterFile) > 0 {
		file, err := os.Open(bloomFilterFile)
		if err != nil {
			fmt.Printf("failed to open file=%s err=%v\n", bloomFilterFile, err)
			return nil
		}
		defer file.Close()

		num, err := filter.ReadFrom(file)
		if err != nil {
			fmt.Printf("failed to read file=%s err=%v\n", bloomFilterFile, err)
			return nil
		}
		fmt.Printf("read file=%s success. num=%d\n", bloomFilterFile, num)
	}
    return filter
}

func GenKeys(filter *bloom.BloomFilter, newKeyFile *os.File, KeyLength int) bool {
	fmt.Printf("start gen keys\n")
	startT := time.Now()
	fileLen := 0
	for i := 0; i < int(genKeyCount); i++ {
		str := RandomString(KeyLength)
		filter.AddString(str)
		_, err := newKeyFile.WriteString(str)
		if err != nil {
			fmt.Printf("failed to write key. err=%v idx=%d\n", err, i)
			return false
		}
		_, err = newKeyFile.Write([]byte{'\n'})
		if err != nil {
			fmt.Printf("failed to write key. err=%v idx=%d\n", err, i)
			return false
		}
		fileLen += len(str) + 1
	}
	tc := time.Since(startT)
	fmt.Printf("save key_file=%s success. num=%d time_cost=%v fileLen=%dKB\n", keyFile, genKeyCount, tc, fileLen/1024)
    return true
}

func SaveBloomFile(filter *bloom.BloomFilter, filterFile string) bool {
	fmt.Printf("start to save new bloom filter file\n")
	newFile, err := os.OpenFile(filterFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to open file=%s err=%v\n", filterFile, err)
		return false
	}
	defer newFile.Close()
	num, err := filter.WriteTo(newFile)
	if err != nil {
		fmt.Printf("failed to write to file=%s err=%v\n", filterFile, err)
		return false
	}
	fmt.Printf("save to file=%s success. num=%d\n", filterFile, num)
    return true
}

func init() {
	flag.Int64Var(&totalKeyCount, "total_key_count",
		600_000_000,
		// 64*64*64*64*64*64,
		"预估的数据量") // default 6位64进制
	flag.Float64Var(&falsePositiveRate, "false_positive_rate", 0.01, "预估的错误率")

	flag.StringVar(&bloomFilterFile, "old_file",
		"", "已经生成过key文件, 则导入此文件的bloomfilter, 且上述的两个选项将被此文件替换")
	flag.StringVar(&newBloomFilterFileSaveTo, "save_to", "./new_bloom_filter_file", "新生成的bloomfilter文件")

	flag.StringVar(&keyFile, "key_file", "./key_file", "生成的key文件")

	flag.Int64Var(&keyLen, "key_len", 6, "key 的长度")
	flag.Int64Var(&genKeyCount, "gen_key_count", 1_000_000, "本次生成key数量")

	rand.Seed(time.Now().UnixNano())
}