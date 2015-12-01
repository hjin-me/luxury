package uuid

import (
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/*
 4 byte timestamp
 1 byte machine
 1 byte pic%251
 2 byte counter
*/
type UUID []byte

var (
	firstDay  time.Time
	machineId uint8
	processId uint8
	counter   uint32
	once      sync.Once
)

func New() UUID {
	once.Do(func() {
		firstDay = time.Unix(1194451200, 0)
		processId = uint8(os.Getpid() % 251)
		name, err := os.Hostname()
		if err != nil {
			t := rand.New(rand.NewSource(time.Now().Unix())).Uint32()
			machineId = uint8(t & 0x00FF)
		} else {
			h := sha1.New()
			h.Write([]byte(name))
			t := h.Sum(nil)
			machineId = uint8((t[0]<<8 | t[1]) & 0x000000FF)
		}
		counter = 0
	})
	var t uint64
	var u = make([]byte, 8)
	atomic.AddUint32(&counter, 1)
	// log.Printf("%d %x %x %x", uint64(time.Now().Sub(firstDay)/time.Second), uint64(machineId)<<24, uint64(processId)<<16, uint64(counter&0x00FF))
	t = (uint64(time.Now().Sub(firstDay)/time.Second) << 32) +
		(uint64(machineId) << 24) + (uint64(processId) << 16) + uint64(counter&0x00FF)
	binary.BigEndian.PutUint64(u, t)
	return UUID(u)
}

func Parse(id uint64) (UUID, error) {
	var u = make([]byte, 8)
	binary.BigEndian.PutUint64(u, id)
	return UUID(u), nil
}

func (u UUID) SetTime(t time.Time) {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(t.Sub(firstDay)/time.Second))
	// log.Println("time to byte", b)
	for i := 0; i < 4; i++ {
		u[i] = b[i]
	}
}

func (u UUID) Time() time.Time {
	// log.Println("uuid time byte", u[0:4])
	d := binary.BigEndian.Uint32(u[0:4])
	// log.Println("byte to duraing", d)

	return firstDay.Add(time.Duration(d) * time.Second)
}

func (u UUID) String() string {
	return strconv.FormatUint(u.Uint64(), 10)
}

func (u UUID) Uint64() uint64 {
	return binary.BigEndian.Uint64(u[:])
}
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(u.String()), nil
}
