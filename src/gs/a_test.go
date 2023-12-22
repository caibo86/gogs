// -------------------------------------------
// @file      : a_test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/20 下午10:39
// -------------------------------------------

package gs

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"gogs/pb"
	"testing"
)

func init() {
	gob.Register(&Car{})
	gob.Register(&Student{})
}

func BenchmarkCar_DeepCopy(b *testing.B) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = NewStudent()
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"

	for n := 0; n < b.N; n++ {
		car.DeepCopy()
	}
}

func BenchmarkCar_DeepCopy1(b *testing.B) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = NewStudent()
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	for n := 0; n < b.N; n++ {
		car.DeepCopy1()
	}
}

func BenchmarkCar_Marshal(b *testing.B) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = &Student{}
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = &Student{}
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	for n := 0; n < b.N; n++ {
		car.Marshal()
	}
}

func BenchmarkCar_JsonMarshal(b *testing.B) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = NewStudent()
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	for n := 0; n < b.N; n++ {
		_, _ = json.Marshal(car)
	}
}

func BenchmarkCar_GobMarshal(b *testing.B) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = NewStudent()
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"

	for n := 0; n < b.N; n++ {
		var buff bytes.Buffer
		encoder := gob.NewEncoder(&buff)
		_ = encoder.Encode(car)
	}
}

func BenchmarkCar_PBMarshal(b *testing.B) {
	car := &pb.Car{
		Attrs: map[string]string{},
	}
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = pb.Color_Red
	car.Owner = &pb.Student{}
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = pb.Err_Success
	car.Size_ = []int32{1, 2, 3, 4}
	car.Drivers = make([]*pb.Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = &pb.Student{}
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	for n := 0; n < b.N; n++ {
		_, err := car.Marshal()
		if err != nil {
			b.Failed()
		}
	}
}

func BenchmarkCar_PBDeepCopy(b *testing.B) {
	car := &pb.Car{
		Attrs: map[string]string{},
	}
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = pb.Color_Red
	car.Owner = &pb.Student{}
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = pb.Err_Success
	car.Size_ = []int32{1, 2, 3, 4}
	car.Drivers = make([]*pb.Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = &pb.Student{}
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	for n := 0; n < b.N; n++ {
		data, err := car.Marshal()
		if err != nil {
			b.Failed()
		}
		ncar := &pb.Car{}
		err = ncar.Unmarshal(data)
		if err != nil {
			b.Failed()
		}
	}
}

func TestCar_Marshal(t *testing.T) {
	car := NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = ColorRed
	car.Owner = NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = make([]*Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = NewStudent()
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}

	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	fmt.Println("gs序列化后长", len(car.Marshal()))
}

func TestCar_PBMarshal(t *testing.T) {
	car := &pb.Car{
		Attrs: map[string]string{},
	}
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = pb.Color_Red
	car.Owner = &pb.Student{}
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = pb.Err_Success
	car.Size_ = []int32{1, 2, 3, 4}
	car.Drivers = make([]*pb.Student, 1000)
	for i := 0; i < 1000; i++ {
		car.Drivers[i] = &pb.Student{}
		car.Drivers[i].Age = 18
		car.Drivers[i].Name = "张三三"
		car.Drivers[i].ID = int64(i)
	}
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"
	data, err := car.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pb序列化后长", len(data))
}
