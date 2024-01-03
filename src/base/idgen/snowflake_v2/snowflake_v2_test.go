package snowflake_v2

import (
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNewSnowflakeV2(t *testing.T) {
	Convey("创建雪花算法V2 IDGen", t, func() {
		Convey("参数非法检查", func() {
			Convey("serverID大于最大值", func() {
				serverID := MaxServerID + 1
				idGen := NewSnowflakeV2(int64(serverID))
				So(idGen, ShouldBeNil)
			})
			Convey("serverID小于零", func() {
				serverID := -1
				idGen := NewSnowflakeV2(int64(serverID))
				So(idGen, ShouldBeNil)
			})
		})
		Convey("正确返回", func() {
			serverID := 1
			idGen := NewSnowflakeV2(int64(serverID))
			So(idGen, ShouldNotBeNil)
		})
	})
}

func TestSnowflakeV2_Next(t *testing.T) {
	Convey("雪花算法", t, func() {
		serverID := 1
		idGen := NewSnowflakeV2(int64(serverID))
		Convey("支持时钟回拨", func() {
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Truncate(10 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			id1, _ := idGen.Next()
			id2, err := idGen.Next()
			So(err, ShouldBeNil)
			So(id2, ShouldBeGreaterThan, id1)
		})
		Convey("支持时钟回拨达到最大次数报错", func() {
			idGen.clockBacks = MaxClockBack
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Truncate(1 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			_, _ = idGen.Next()
			_, err := idGen.Next()
			So(err, ShouldEqual, ErrClockTurnedBack)
		})
		Convey("新的时间单元序号正确置0", func() {
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Add(1 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			id1, _ := idGen.Next()
			id2, err := idGen.Next()
			So(err, ShouldBeNil)
			So(idGen.seq, ShouldEqual, 0)
			So(id2, ShouldBeGreaterThan, id1)
		})
		Convey("序列号到最大等待到下一个时间单元时遇上时钟回拨", func() {
			idGen.seq = MaxSeqID - 1
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Truncate(1 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			id1, _ := idGen.Next()
			id2, err := idGen.Next()
			So(err, ShouldBeNil)
			So(idGen.seq, ShouldEqual, 0)
			So(id2, ShouldBeGreaterThan, id1)
		})
		Convey("序列号到最大等待到下一个时间单元时遇上时钟回拨,但回拨次数达到最大", func() {
			idGen.seq = MaxSeqID - 1
			idGen.clockBacks = MaxClockBack
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Truncate(1 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			_, _ = idGen.Next()
			_, err := idGen.Next()
			So(err, ShouldEqual, ErrClockTurnedBack)
		})
		Convey("序列号到最大能正确等待到下一个时间单元", func() {
			idGen.seq = MaxSeqID - 1
			id1, _ := idGen.Next()
			id2, err := idGen.Next()
			So(err, ShouldBeNil)
			So(id2, ShouldBeGreaterThan, id1)
		})

		Convey("时间单元溢出", func() {
			now, _ := time.Parse(time.DateTime, "2200-01-01 00:00:00")
			patch := gomonkey.ApplyFuncReturn(time.Now, now)
			defer patch.Reset()
			_, err := idGen.Next()
			So(err, ShouldEqual, ErrTimeUnitOverflow)
		})
		Convey("发号器溢出", func() {
			idGen.lastID = idGen.MaxID()
			_, err := idGen.Next()
			So(err, ShouldEqual, ErrIDGenOverflow)
		})
		Convey("正确返回", func() {
			_, err := idGen.Next()
			So(err, ShouldBeNil)
		})
		Convey("ID递增", func() {
			id1, _ := idGen.Next()
			id2, _ := idGen.Next()
			So(id2, ShouldBeGreaterThan, id1)
		})

	})
}
