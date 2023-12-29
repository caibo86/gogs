// -------------------------------------------
// @file      : test.gs_test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/27 下午2:04
// -------------------------------------------

package gs

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"
	"gogs/pb"
	"math"
	"testing"
)

var car = &Car{
	VarEnum:    ColorBlue,
	VarString:  "我这个是中文和English混合的string",
	VarBool:    true,
	VarByte:    'C',
	VarUint8:   'B',
	VarInt16:   math.MaxInt16,
	VarUint16:  math.MaxUint16,
	VarInt32:   math.MaxInt32,
	VarUint32:  math.MaxUint32,
	VarInt64:   math.MaxInt64,
	VarUint64:  123456789,
	VarFloat32: math.MaxFloat32,
	VarFloat64: math.MaxFloat64,
	VarStruct:  &Student{ID: 904088, Name: "蔡波", Age: 18},
	VarList:    []int32{1, 2, 3, 4, 5},
	VarStructs: []*Student{{ID: 123, Name: "汪汪队", Age: 33},
		{ID: 999, Name: "超人", Age: 888}},
	VarBools:    []bool{true, false, true, false},
	VarStrings:  []string{"我是中文", "I am English"},
	VarFloat32s: []float32{math.MaxFloat32, 3.14125},
	VarFloat64s: []float64{math.MaxFloat64, 9999.99, 33, 24.1287},
	VarEnums:    []Color{ColorRed, ColorBlue, ColorGreen},
	VarMap:      map[string]string{"key1": "我是中文", "键2": "I am English"},
	VarMap1: map[string]*Student{
		"成都": &Student{
			ID:   1,
			Name: "成都扛把子",
			Age:  15,
		},
		"资中": &Student{
			ID:   2,
			Name: "资中绕王",
			Age:  25,
		},
	},
}

var pbCar = &pb.Car{
	VarEnum:    pb.Color_Blue,
	VarString:  "我这个是中文和English混合的string",
	VarBool:    true,
	VarByte:    true,
	VarSbyte:   true,
	VarInt16:   math.MaxInt16,
	VarUint16:  math.MaxUint16,
	VarInt32:   math.MaxInt32,
	VarUint32:  math.MaxUint32,
	VarInt64:   math.MaxInt64,
	VarUint64:  123456789,
	VarFloat32: math.MaxFloat32,
	VarFloat64: math.MaxFloat64,
	VarStruct:  &pb.Student{ID: 904088, Name: "蔡波", Age: 18},
	VarArray:   []int32{1, 2, 3, 4, 5},
	VarStructs: []*pb.Student{{ID: 123, Name: "汪汪队", Age: 33},
		{ID: 999, Name: "超人", Age: 888}},
	VarBools:    []bool{true, false, true, false},
	VarStrings:  []string{"我是中文", "I am English"},
	VarFloat32S: []float32{math.MaxFloat32, 3.14125},
	VarFloat64S: []float64{math.MaxFloat64, 9999.99, 33, 24.1287},
	VarEnums:    []pb.Color{pb.Color_Red, pb.Color_Blue, pb.Color_Green},
	VarMap:      map[string]string{"key1": "我是中文", "键2": "I am English"},
	VarMap1: map[string]*pb.Student{
		"成都": &pb.Student{
			ID:   1,
			Name: "成都扛把子",
			Age:  15,
		},
		"资中": &pb.Student{
			ID:   2,
			Name: "资中绕王",
			Age:  25,
		},
	},
}

func TestGSLangMarshal(t *testing.T) {
	Convey("测试gslang的序列化和反序列化", t, func() {
		data := car.Marshal()
		newCar := &Car{}
		err := newCar.Unmarshal(data)
		fmt.Println(*car, *newCar)
		So(err, ShouldBeNil)
		So(newCar.VarEnum, ShouldEqual, ColorBlue)
		So(newCar.VarString, ShouldEqual, "我这个是中文和English混合的string")
		So(newCar.VarBool, ShouldEqual, true)
		So(newCar.VarByte, ShouldEqual, 'C')
		So(newCar.VarUint8, ShouldEqual, 'B')
		So(newCar.VarInt16, ShouldEqual, math.MaxInt16)
		So(newCar.VarUint16, ShouldEqual, math.MaxUint16)
		So(newCar.VarInt32, ShouldEqual, math.MaxInt32)
		So(newCar.VarUint32, ShouldEqual, math.MaxUint32)
		So(newCar.VarInt64, ShouldEqual, math.MaxInt64)
		So(newCar.VarUint64, ShouldEqual, 123456789)
		So(newCar.VarFloat32, ShouldEqual, math.MaxFloat32)
		So(newCar.VarFloat64, ShouldEqual, math.MaxFloat64)
		So(newCar.VarStruct.ID, ShouldEqual, 904088)
		So(newCar.VarStruct.Name, ShouldEqual, "蔡波")
		So(newCar.VarStruct.Age, ShouldEqual, 18)
		So(newCar.VarList, ShouldResemble, []int32{1, 2, 3, 4, 5})
		So(newCar.VarStructs[0].ID, ShouldEqual, 123)
		So(newCar.VarStructs[0].Name, ShouldEqual, "汪汪队")
		So(newCar.VarStructs[0].Age, ShouldEqual, 33)
		So(newCar.VarStructs[1].ID, ShouldEqual, 999)
		So(newCar.VarStructs[1].Name, ShouldEqual, "超人")
		So(newCar.VarStructs[1].Age, ShouldEqual, 888)
		So(newCar.VarBools, ShouldResemble, []bool{true, false, true, false})
		So(newCar.VarStrings, ShouldResemble, []string{"我是中文", "I am English"})
		So(newCar.VarFloat32s, ShouldResemble, []float32{math.MaxFloat32, 3.14125})
		So(newCar.VarFloat64s, ShouldResemble, []float64{math.MaxFloat64, 9999.99, 33, 24.1287})
		So(newCar.VarEnums, ShouldResemble, []Color{ColorRed, ColorBlue, ColorGreen})
		So(newCar.VarMap["key1"], ShouldEqual, "我是中文")
		So(newCar.VarMap["键2"], ShouldEqual, "I am English")
		So(newCar.VarMap1["成都"].ID, ShouldEqual, 1)
		So(newCar.VarMap1["成都"].Name, ShouldEqual, "成都扛把子")
		So(newCar.VarMap1["成都"].Age, ShouldEqual, 15)
		So(newCar.VarMap1["资中"].ID, ShouldEqual, 2)
		So(newCar.VarMap1["资中"].Name, ShouldEqual, "资中绕王")
		So(newCar.VarMap1["资中"].Age, ShouldEqual, 25)
	})
}

func TestLength(t *testing.T) {
	Convey("看下序列化后数据长度对比", t, func() {
		fmt.Println("car size:", car.Size())
		data := car.Marshal()
		fmt.Println("gslang序列化后:", len(data))
		data1, _ := pbCar.Marshal()
		fmt.Println("protobuf序列化后:", len(data1))
	})
}

func BenchmarkCar_Marshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		car.Marshal()
	}
}

func BenchmarkPBCar_Marshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = pbCar.Marshal()
	}
}

func BenchmarkCar_Copy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		car.Copy()
	}
}

func BenchmarkCar_DeepCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		car.DeepCopy()
	}
}

func BenchmarkCar_CopyByMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		car1 := &Car{}
		data := car.Marshal()
		car1.Unmarshal(data)
	}
}

func TestPB(t *testing.T) {
	pbCar := &pb.Car{
		// VarStructs: make([]*pb.Student, 3),
		VarMap1: make(map[string]*pb.Student, 3),
	}
	pbCar.VarMap1["成都"] = nil
	fmt.Println(pbCar.VarMap1)
	data, err := proto.Marshal(pbCar)
	if err != nil {
		t.Error(err)
	}
	newPbCar := &pb.Car{}
	err = proto.Unmarshal(data, newPbCar)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("结果:", newPbCar.VarMap1 == nil)
	stu := new(Student)
	fmt.Println("stu的size:", stu.Size())
	car := &Car{
		VarStructs:   make([]*Student, 3),
		VarStructMap: make(map[string]*Table, 3),
	}
	car.VarStructs[1] = new(Student)
	car.VarStructMap["成都"] = nil
	data = car.Marshal()
	newCar := &Car{}
	err = newCar.Unmarshal(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("第二个结果:", newCar.VarStructMap, newCar.VarStructs)
}

func TestDeepCopy(t *testing.T) {
	Convey("测试深拷贝", t, func() {
		stu := &Student{
			ID:   1,
			Name: "成都扛把子",
		}
		car := &Car{
			VarStructArray1: [10]*Student{stu},
		}
		fmt.Println("car:", car.VarStructArray1[0])
		newCar := car.DeepCopy()
		newCar1 := car.Copy()
		fmt.Println("newCar:", newCar.VarStructArray1[0])
		fmt.Println("newCar1:", newCar1.VarStructArray1[0])
		car.VarStructArray1[0].Name = "改了名字了"
		fmt.Println("car:", car.VarStructArray1[0])
		fmt.Println("newCar:", newCar.VarStructArray1[0])
		fmt.Println("newCar1:", newCar1.VarStructArray1[0])
	})
}

func BenchmarkCar_MarshalBytes(b *testing.B) {
	data := "sdnhfdsfhjiupwefhniuweqsdnhfdsfhjiupwefhniuweqfhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui12983217sdnhfdsfhjiupwefhniuweqfhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui129832173892176321876381270632188fjiewnfuiwenc82u3739823193892176321876381270632188fjiewnfuiwenc82u373982319fhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui129832173892176321876381270632188fjiewnfuiwenc82u373982319-42-0=40-32843290u432=094u3294234n3i2u4n32uinfds8fdfekwi9mfmew98fuy328374234j2m8f9ewhf78ehdfuydshfdusfhdsuiohf78023y378473204h432u8hfd9uerhfuehf7reh79823"
	car := &Car{
		VarBytes: []byte(data),
	}
	for i := 0; i < b.N; i++ {
		car.Marshal()
	}
}

func BenchmarkCar_MarshalByteSlice(b *testing.B) {
	data := "sdnhfdsfhjiupwefhniuweqsdnhfdsfhjiupwefhniuweqfhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui12983217sdnhfdsfhjiupwefhniuweqfhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui129832173892176321876381270632188fjiewnfuiwenc82u3739823193892176321876381270632188fjiewnfuiwenc82u373982319fhnniwuehfwehfweiofhewiopfhweiofnjkdsnfjkdsnfjksdnfdjkshfndsjfhwuifhwefhweifhuwehfnewui129832173892176321876381270632188fjiewnfuiwenc82u373982319-42-0=40-32843290u432=094u3294234n3i2u4n32uinfds8fdfekwi9mfmew98fuy328374234j2m8f9ewhf78ehdfuydshfdusfhdsuiohf78023y378473204h432u8hfd9uerhfuehf7reh79823"
	car := &Car{
		VarData: []byte(data),
	}
	for i := 0; i < b.N; i++ {
		car.Marshal()
	}
}
