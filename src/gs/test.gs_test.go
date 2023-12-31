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
	"gogs/gss"
	"gogs/pb"
	"math"
	"testing"
)

var car = &Car{
	VarEnum:    ColorBlue,
	VarUint8:   'B',
	VarString:  "我这个是中文和English混合的string",
	VarBool:    true,
	VarByte:    'C',
	VarInt8:    -100,
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
		"成都": {
			ID:   1,
			Name: "成都扛把子",
			Age:  15,
		},
		"资中": {
			ID:   2,
			Name: "资中绕王",
			Age:  25,
		},
	},
	VarSubject:      gss.SubjectBiology,
	VarTeacher:      &gss.Teacher{ID: 1, Name: "石老师", Age: 15},
	VarTeachers:     []*gss.Teacher{{ID: 1, Name: "石老师", Age: 15}, {ID: 2, Name: "李老师", Age: 25}},
	VarSubjects:     []gss.Subject{gss.SubjectBiology, gss.SubjectChemistry, gss.SubjectChinese},
	VarMap2:         map[gss.Subject]gss.Subject{gss.SubjectBiology: gss.SubjectChemistry},
	VarArray:        [3]int32{3, 6, 9},
	VarStructArray:  [3]*gss.Teacher{{ID: 1, Name: "石老师", Age: 15}, nil, {ID: 2, Name: "李老师", Age: 25}},
	VarEnumArray:    [4]gss.Subject{gss.SubjectBiology, gss.SubjectChemistry, gss.SubjectChinese, gss.SubjectChinese},
	VarStructArray1: [4]*Student{{ID: 904088, Name: "蔡波", Age: 18}, nil, nil, {ID: 111, Name: "桌子", Age: 22}},
	VarStructSlice:  []*Table{{}, nil, {}},
	VarStructMap:    map[string]*Table{"aaa": {}, "bbb": nil},
	VarData:         []byte("asdsadsadsadsadsadsadsa"),
	VarBytes:        []byte("bbbdf1221231321"),
	VarArrayBytes:   [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	VarMap3: map[gss.Subject]*gss.Teacher{gss.SubjectBiology: {ID: 1, Name: "石老师", Age: 15},
		gss.SubjectChemistry: {ID: 2, Name: "李老师", Age: 25},
		gss.SubjectChinese:   nil},
	VarSliceBytes: [][]byte{[]byte("abcdef"), []byte("ffff123456"), []byte("这是中文不是英文^&%&^%")},
}

var pbCar = &pb.Car{
	VarEnum:    pb.Color_Blue,
	VarUint8:   'B',
	VarString:  "我这个是中文和English混合的string",
	VarBool:    true,
	VarByte:    'C',
	VarInt8:    -100,
	VarInt16:   math.MaxInt16,
	VarUint16:  math.MaxUint16,
	VarInt32:   math.MaxInt32,
	VarUint32:  math.MaxUint32,
	VarInt64:   math.MaxInt64,
	VarUint64:  123456789,
	VarFloat32: math.MaxFloat32,
	VarFloat64: math.MaxFloat64,
	VarStruct:  &pb.Student{ID: 904088, Name: "蔡波", Age: 18},
	VarList:    []int32{1, 2, 3, 4, 5},
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
	VarSubject:      pb.Subject_Biology,
	VarTeacher:      &pb.Teacher{ID: 1, Name: "石老师", Age: 15},
	VarTeachers:     []*pb.Teacher{{ID: 1, Name: "石老师", Age: 15}, {ID: 2, Name: "李老师", Age: 25}},
	VarSubjects:     []pb.Subject{pb.Subject_Biology, pb.Subject_Chemistry, pb.Subject_Chinese},
	VarMap2:         map[int32]pb.Subject{10000000: pb.Subject_Chemistry},
	VarArray:        []int32{3, 6, 9},
	VarStructArray:  []*pb.Teacher{{ID: 1, Name: "石老师", Age: 15}, {}, {ID: 2, Name: "李老师", Age: 25}},
	VarEnumArray:    []pb.Subject{pb.Subject_Biology, pb.Subject_Chemistry, pb.Subject_Chinese, pb.Subject_Chinese},
	VarStructArray1: []*pb.Student{{ID: 904088, Name: "蔡波", Age: 18}, {}, {}, {ID: 111, Name: "桌子", Age: 22}},
	VarStructSlice:  []*pb.Table{{}, {}, {}},
	VarStructMap:    map[string]*pb.Table{"aaa": {}, "bbb": {}},
	VarData:         []byte("asdsadsadsadsadsadsadsa"),
	VarBytes:        []byte("bbbdf1221231321"),
	VarArrayBytes:   []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	VarMap3: map[int32]*pb.Teacher{int32(pb.Subject_Biology): {ID: 1, Name: "石老师", Age: 15},
		int32(pb.Subject_Chemistry): {ID: 2, Name: "李老师", Age: 25},
		int32(pb.Subject_Chinese):   {}},
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
		So(newCar.VarSubject, ShouldEqual, gss.SubjectBiology)
		So(newCar.VarTeacher.ID, ShouldEqual, 1)
		So(newCar.VarTeacher.Name, ShouldEqual, "石老师")
		So(newCar.VarTeacher.Age, ShouldEqual, 15)
		So(newCar.VarTeachers[0].ID, ShouldEqual, 1)
		So(newCar.VarTeachers[0].Name, ShouldEqual, "石老师")
		So(newCar.VarTeachers[0].Age, ShouldEqual, 15)
		So(newCar.VarTeachers[1].ID, ShouldEqual, 2)
		So(newCar.VarTeachers[1].Name, ShouldEqual, "李老师")
		So(newCar.VarTeachers[1].Age, ShouldEqual, 25)
		So(newCar.VarSubjects[0], ShouldEqual, gss.SubjectBiology)
		So(newCar.VarSubjects[1], ShouldEqual, gss.SubjectChemistry)
		So(newCar.VarSubjects[2], ShouldEqual, gss.SubjectChinese)
		So(newCar.VarMap2[gss.SubjectBiology], ShouldEqual, gss.SubjectChemistry)
		So(newCar.VarArray[0], ShouldEqual, 3)
		So(newCar.VarArray[1], ShouldEqual, 6)
		So(newCar.VarArray[2], ShouldEqual, 9)
		So(newCar.VarStructArray[0].ID, ShouldEqual, 1)
		So(newCar.VarStructArray[0].Name, ShouldEqual, "石老师")
		So(newCar.VarStructArray[0].Age, ShouldEqual, 15)
		So(newCar.VarStructArray[1], ShouldBeNil)
		So(newCar.VarStructArray[2].ID, ShouldEqual, 2)
		So(newCar.VarStructArray[2].Name, ShouldEqual, "李老师")
		So(newCar.VarStructArray[2].Age, ShouldEqual, 25)
		So(newCar.VarEnumArray[0], ShouldEqual, gss.SubjectBiology)
		So(newCar.VarEnumArray[1], ShouldEqual, gss.SubjectChemistry)
		So(newCar.VarEnumArray[2], ShouldEqual, gss.SubjectChinese)
		So(newCar.VarEnumArray[3], ShouldEqual, gss.SubjectChinese)
		So(newCar.VarStructArray1[0].ID, ShouldEqual, 904088)
		So(newCar.VarStructArray1[0].Name, ShouldEqual, "蔡波")
		So(newCar.VarStructArray1[0].Age, ShouldEqual, 18)
		So(newCar.VarStructArray1[1], ShouldBeNil)
		So(newCar.VarStructArray1[2], ShouldBeNil)
		So(newCar.VarStructArray1[3].ID, ShouldEqual, 111)
		So(newCar.VarStructArray1[3].Name, ShouldEqual, "桌子")
		So(newCar.VarStructArray1[3].Age, ShouldEqual, 22)
		So(newCar.VarStructSlice[0], ShouldNotBeNil)
		So(newCar.VarStructSlice[1], ShouldBeNil)
		So(newCar.VarStructSlice[2], ShouldNotBeNil)
		So(newCar.VarStructMap["aaa"], ShouldNotBeNil)
		So(newCar.VarStructMap["bbb"], ShouldBeNil)
		So(newCar.VarData, ShouldResemble, []byte("asdsadsadsadsadsadsadsa"))
		So(newCar.VarBytes, ShouldResemble, []byte("bbbdf1221231321"))
		So(newCar.VarArrayBytes[0], ShouldEqual, 1)
		So(newCar.VarArrayBytes[1], ShouldEqual, 2)
		So(newCar.VarArrayBytes[2], ShouldEqual, 3)
		So(newCar.VarArrayBytes[3], ShouldEqual, 4)
		So(newCar.VarArrayBytes[4], ShouldEqual, 5)
		So(newCar.VarArrayBytes[5], ShouldEqual, 6)
		So(newCar.VarArrayBytes[6], ShouldEqual, 7)
		So(newCar.VarArrayBytes[7], ShouldEqual, 8)
		So(newCar.VarArrayBytes[8], ShouldEqual, 9)
		So(newCar.VarArrayBytes[9], ShouldEqual, 0)
		So(newCar.VarMap3[gss.SubjectBiology].ID, ShouldEqual, 1)
		So(newCar.VarMap3[gss.SubjectBiology].Name, ShouldEqual, "石老师")
		So(newCar.VarMap3[gss.SubjectBiology].Age, ShouldEqual, 15)
		So(newCar.VarMap3[gss.SubjectChemistry].ID, ShouldEqual, 2)
		So(newCar.VarMap3[gss.SubjectChemistry].Name, ShouldEqual, "李老师")
		So(newCar.VarMap3[gss.SubjectChemistry].Age, ShouldEqual, 25)
		So(newCar.VarMap3[gss.SubjectChinese], ShouldBeNil)
		So(newCar.VarSliceBytes[0], ShouldResemble, []byte("abcdef"))
		So(newCar.VarSliceBytes[1], ShouldResemble, []byte("ffff123456"))
		So(newCar.VarSliceBytes[2], ShouldResemble, []byte("这是中文不是英文^&%&^%"))
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

func BenchmarkPBCar_CopyByMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		car1 := &pb.Car{}
		data, _ := pbCar.Marshal()
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
			VarStructArray1: [4]*Student{stu},
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

func TestUse(t *testing.T) {
	Convey("测试使用", t, func() {
		data := []byte("aaa")
		s := string(data)
		fmt.Println(data)
		fmt.Println(s)
		data[1] = 100
		fmt.Println(data)
		fmt.Println(s)
	})
}
