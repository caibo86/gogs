package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	SignBits      = 1                                       // 1位是符号位，也就是最高位，始终是0，没有任何意义，因为要是唯一计算机二进制补码中就是负数，0才是正数。
	TimeUnitBits  = 41                                      // 41位是时间戳，具体到毫秒，41位的二进制可以使用69年，因为时间理论上永恒递增，所以根据这个排序是可以的。
	ServerIDBits  = 10                                      // 10位是机器标识，可以全部用作机器ID，也可以用来标识机房ID + 机器ID，10位最多可以表示1024台机器。
	SequenceBits  = 12                                      // 12位是计数序列号，也就是同一台机器上同一时间，理论上还可以同时生成不同的ID，12位的序列号能够区分出4096个ID。
	TimeUnitShift = ServerIDBits + SequenceBits             // 时间单元左移位数
	TimeUnit      = int64(1 * time.Millisecond)             // 时间单位 1毫秒
	MaxTimeUnits  = (1 << TimeUnitBits) - 1                 // 最大时间单元
	MaxSeqID      = (1 << SequenceBits) - 1                 // 最大序号数
	MaxServerID   = (1 << ServerIDBits) - 1                 // 最大机器ID
	CustomEpoch   = int64(1675161600000 * time.Millisecond) // 开始纪元 2024-01-01 00:00:00 UTC
)

var (
	ErrClockTurnedBack  = errors.New("clock turned back")
	ErrTimeUnitOverflow = errors.New("time unit overflow")
	ErrIDGenOverflow    = errors.New("id gen overflow")
)

// Snowflake 标准雪花算法
type Snowflake struct {
	sync.Mutex         // 互斥锁
	seq          int64 // 当前序列号
	serverID     int64 // 机器号
	lastTimeUnit int64 // 上一次取号的时间单元
	lastID       int64 // 上一次的取的号
}

func NewSnowflake(serverID int64) *Snowflake {
	if serverID > MaxServerID || serverID < 0 {
		return nil
	}
	return &Snowflake{
		serverID:     serverID,
		lastTimeUnit: currentTimeUnit(),
	}
}

func currentTimeUnit() int64 {
	return (time.Now().UnixNano() - CustomEpoch) / TimeUnit
}

func (s *Snowflake) Next() (int64, error) {
	s.Lock()
	defer s.Unlock()
	curTimeUnit := currentTimeUnit()
	if curTimeUnit > MaxTimeUnits {
		return 0, ErrTimeUnitOverflow
	}
	if curTimeUnit < s.lastTimeUnit {
		// 时钟回拨
		return 0, ErrClockTurnedBack
	}
	if curTimeUnit == s.lastTimeUnit {
		s.seq++
		if s.seq > MaxSeqID {
			if tu, err := s.waitTilNext(curTimeUnit); err != nil {
				return 0, err
			} else {
				curTimeUnit = tu
				s.seq = 0
			}
		}
	} else {
		s.seq = 0
	}
	s.lastTimeUnit = curTimeUnit
	id := (s.lastTimeUnit << TimeUnitShift) | (s.serverID << SequenceBits) | s.seq
	if id <= s.lastID {
		return 0, ErrIDGenOverflow
	}
	s.lastID = id
	return id, nil

}

func (s *Snowflake) MaxID() int64 {
	return MaxTimeUnits<<TimeUnitShift | s.serverID<<SequenceBits | MaxSeqID
}

func (s *Snowflake) waitTilNext(lastTU int64) (int64, error) {
	var now int64
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Microsecond)
		now = currentTimeUnit()
		if now < lastTU {
			return 0, ErrClockTurnedBack
		}
		if now == lastTU {
			continue
		}
		break
	}
	return now, nil
}
