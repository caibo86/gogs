package snowflake

import (
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNewSnowflake(t *testing.T) {
	Convey("创建雪花算法IDGen", t, func() {
		Convey("参数非法检查", func() {
			Convey("serverID大于最大值", func() {
				serverID := MaxServerID + 1
				idGen := NewSnowflake(int64(serverID))
				So(idGen, ShouldBeNil)
			})
			Convey("serverID小于零", func() {
				serverID := -1
				idGen := NewSnowflake(int64(serverID))
				So(idGen, ShouldBeNil)
			})
		})
		Convey("正确返回", func() {
			serverID := 1
			idGen := NewSnowflake(int64(serverID))
			So(idGen, ShouldNotBeNil)
		})
	})
}

func TestSnowflake_Next(t *testing.T) {
	Convey("雪花算法", t, func() {
		serverID := 1
		idGen := NewSnowflake(int64(serverID))
		Convey("时钟回拨报错", func() {
			now := time.Now()
			o := []gomonkey.OutputCell{
				{Values: gomonkey.Params{now}, Times: 1},
				{Values: gomonkey.Params{now.Truncate(1 * time.Second)}, Times: 1},
			}
			patch := gomonkey.ApplyFuncSeq(time.Now, o)
			defer patch.Reset()
			_, err := idGen.Next()
			_, err = idGen.Next()
			So(err, ShouldEqual, ErrClockTurnedBack)
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
