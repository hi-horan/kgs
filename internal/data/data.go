package data

import (
	"fmt"
	"io"
	"os"
	"strconv"
    "errors"

	v1 "kgs/api/v1"
	"kgs/internal/conf"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var (
	ErrHaveNoKeys = kerrors.NotFound(v1.ErrorReason_HAVE_NO_KEYS.String(), "have no keys")
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewKGSRepo)

// Data .
type Data struct {
	key chan string

	currentKeyIndex     int64
	currentKeyIndexFile *os.File

	keyfile *os.File

	conf *conf.Data

	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(logger)

	keyfile, err := os.Open(c.FilePath)
	if err != nil {
		l.Errorf("failed to open file=%s err=%w", c.FilePath, err)
		return nil, nil, err
	}

	keyIdxFile, err := os.OpenFile(c.CurrentKeyIndex, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		l.Errorf("failed to open index file=%s err=%w", c.CurrentKeyIndex, err)
		keyfile.Close()
		return nil, nil, err
	}

	var idx int64
	n, err := fmt.Fscanf(keyIdxFile, "%d", &idx)
	if err != nil {
		l.Warnf("failed to read index file=%s err=%v", c.CurrentKeyIndex, err)
	}
	l.Infof("read index file. n=%d", n)

	if idx != 0 {
		ret, err := keyfile.Seek(idx, io.SeekStart)
		if err != nil {
			l.Errorf("failed to seek idx file. err=%w", err)
			keyfile.Close()
			keyIdxFile.Close()
			return nil, nil, err
		}
		if ret != idx {
			l.Warnf("ret=%d idx=%d", ret, idx)
			idx = ret
		}
	}

	keyChan := make(chan string, c.KeyChanLen) // 带缓冲chan 可以做到分配key和从文件获取key并行

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data file resources")
		keyfile.Close()
		keyIdxFile.Close()
	}
	d := Data{
		key: keyChan,

		currentKeyIndex:     idx,
		currentKeyIndexFile: keyIdxFile,

		keyfile: keyfile,

		conf: c,

		log: l,
	}

	go d.getKeys()

	return &d, cleanup, nil
}

func (d *Data) Getkeys(count int64) ([]string, error) {
	keys := make([]string, 0, count)
	if count == 0 {
		return keys, nil
	}

	for i := 0; i < int(count); i++ {
		key, ok := <-d.key
		if !ok {
			msg := "have no keys."
			d.log.Errorf(msg)
			break
		}
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		// NOTE: 已经没有key，停止服务
		return keys, ErrHaveNoKeys
	}
	return keys, nil
}

func (d *Data) getKeys() {
	for {
		keys, e := d.readKeys()

		for idx := range keys {
			d.key <- keys[idx]
		}
        if e != nil {
            break
        }
	}
}

func (d *Data) readKeys() ([]string, error) {
	keys := make([]string, 0, d.conf.PreAllocCount)
    var e error
	for int64(len(keys)) < d.conf.PreAllocCount {
        key := make([]byte, d.conf.EachKeyLen)
        n, err := d.keyfile.Read(key)
        if err != nil {
            d.log.Errorf("failed to read key. err=%v", err)
            break
        }
        if n != int(d.conf.EachKeyLen) {
            d.log.Errorf("failed to read key. err=%v", err)
            break
        }
		if key[len(key)-1] == '\n' {
			key = key[0 : len(key)-1]
		}
		keys = append(keys, string(key))
	}
	if len(keys) == 0 {
		msg := "have no keys."
		d.log.Errorf(msg)
		close(d.key) // NOTE: 将chan中的key消耗完
        e = errors.New("have no keys")
		// os.Exit(-1)
	}
	idx, err := d.keyfile.Seek(0, io.SeekCurrent)
	if err != nil {
		d.log.Errorf("failed to seek new idx. lastKey=%s err=%w", keys[len(keys)-1], err)
		// os.Exit(-1)
		close(d.key)
        e = errors.New("failed to seek new idx")
	}
	d.currentKeyIndex = idx
	d.currentKeyIndexFile.Truncate(0)
	d.currentKeyIndexFile.WriteAt([]byte(strconv.FormatInt(idx, 10)), 0)
	d.currentKeyIndexFile.Sync()
	return keys, e
}
