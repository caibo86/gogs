package snowflake_v2

import (
	"errors"
	"sync"
	"time"
)

// 改进版雪花算法
// 2位时钟回拨位,支持时钟被回拨3次
// 时间戳缩短到37位,但是时间单元设置位10毫秒
// 12位机器号,最大机器ID=4095
// 12位序列号,每个时间单位最大4095个号

const (
	SignBits       = 1                            // 1位是符号位，也就是最高位，始终是0，没有任何意义，因为要是唯一计算机二进制补码中就是负数，0才是正数。
	ClockBackBits  = 2                            // 2位时钟回拨标记
	TimeUnitBits   = 37                           // 37位是时间戳，具体到厘秒，37位的二进制可以使用43年，因为时间理论上永恒递增，所以根据这个排序是可以的。
	ServerIDBits   = 12                           // 12位是机器标识，可以全部用作机器ID，也可以用来标识机房ID + 机器ID，12位最多可以表示4086台机器。
	SequenceBits   = 12                           // 12位是计数序列号，也就是同一台机器上同一时间，理论上还可以同时生成不同的ID，12位的序列号能够区分出4096个ID。
	TimeUnitShift  = ServerIDBits + SequenceBits  // 时间单元左移位数
	TimeUnit       = int64(10 * time.Millisecond) // 时间单位 10毫秒
	MaxTimeUnits   = (1 << TimeUnitBits) - 1      // 最大时间单元
	MaxSeqID       = (1 << SequenceBits) - 1      // 最大序号数
	ClockBackShift = TimeUnitBits + ServerIDBits + SequenceBits
	MaxClockBack   = (1 << ClockBackBits) - 1                // 最大时钟回拨次数
	CustomEpoch    = int64(1675161600000 * time.Millisecond) // 开始纪元 2024-01-01 00:00:00 UTC
	MaxServerID    = (1 << ServerIDBits) - 1
)

var (
	ErrClockTurnedBack  = errors.New("clock turned back")
	ErrTimeUnitOverflow = errors.New("time unit overflow")
	ErrIDGenOverflow    = errors.New("id gen overflow")
)

// SnowflakeV2 改进版雪花算法
type SnowflakeV2 struct {
	sync.Mutex         // 互斥锁
	seq          int64 // 当前序列号
	serverID     int64 // 机器号
	lastTimeUnit int64 // 上一次取号的时间单元
	lastID       int64 // 上一次的取的号
	clockBacks   int64 // 时钟回拨的次数
}

func NewSnowflakeV2(serverID int64) *SnowflakeV2 {
	if serverID > MaxServerID || serverID < 0 {
		return nil
	}
	return &SnowflakeV2{
		serverID:     serverID,
		lastTimeUnit: currentTimeUnit(),
	}
}

func currentTimeUnit() int64 {
	return (time.Now().UnixNano() - CustomEpoch) / TimeUnit
}

// 标记时钟被回拨了
func (s *SnowflakeV2) clockBack() error {
	if s.clockBacks >= MaxClockBack {
		// 达到最大回拨次数,不再发号
		return ErrClockTurnedBack
	}
	s.clockBacks++
	return nil
}

func (s *SnowflakeV2) Next() (int64, error) {
	s.Lock()
	defer s.Unlock()
	curTimeUnit := currentTimeUnit()
	if curTimeUnit > MaxTimeUnits {
		return 0, ErrTimeUnitOverflow
	}
	if curTimeUnit < s.lastTimeUnit {
		// 时钟回拨
		if err := s.clockBack(); err != nil {
			return 0, err
		}
		s.seq = 0
	} else if curTimeUnit == s.lastTimeUnit {
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
	id := (s.clockBacks << ClockBackShift) | (s.lastTimeUnit << TimeUnitShift) | (s.serverID << SequenceBits) | s.seq
	if id <= s.lastID {
		return 0, ErrIDGenOverflow
	}
	s.lastID = id
	return id, nil

}

func (s *SnowflakeV2) MaxID() int64 {
	return MaxClockBack<<ClockBackShift | MaxTimeUnits<<TimeUnitShift | s.serverID<<SequenceBits | MaxSeqID
}

func (s *SnowflakeV2) waitTilNext(lastTU int64) (int64, error) {
	var prevUT int64
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Millisecond)
		now := currentTimeUnit()
		if now > lastTU {
			return now, nil
		}
		if now < lastTU || now < prevUT {
			if err := s.clockBack(); err != nil {
				return 0, err
			} else {
				return now, nil
			}
		}
		prevUT = now
	}
	return 0, ErrClockTurnedBack
}
