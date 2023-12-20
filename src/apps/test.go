// -------------------------------------------
// @file      : test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/20 下午8:42
// -------------------------------------------

package main

import (
	"fmt"
	"gogs/gs"
)

func main() {
	car := gs.NewCar()
	car.ID = 1
	car.Name = "BMW"
	car.Price = 1000000
	car.Color = gs.ColorRed
	car.Owner = gs.NewStudent()
	car.Owner.Age = 18
	car.Owner.Name = "张三"
	car.Owner.ID = 1
	car.Code = gs.ErrSuccess
	car.Size = [4]int32{1, 2, 3, 4}
	car.Drivers = []*gs.Student{gs.NewStudent(), gs.NewStudent()}
	car.Drivers[0].Age = 18
	car.Drivers[0].Name = "张三三"
	car.Drivers[0].ID = 2
	car.Drivers[1].Age = 18
	car.Drivers[1].Name = "张三三三"
	car.Drivers[1].ID = 3
	car.Attrs["颜色"] = "黑色"
	car.Attrs["车牌"] = "京A88888"
	car.Attrs["车型"] = "宝马X6"

	newCar := car.DeepCopy()

	fmt.Printf("%+v %+v \n%+v\n", car, *car.Owner, car.Drivers)
	fmt.Printf("%+v %+v \n%+v\n", newCar, *newCar.Owner, newCar.Drivers)

	newCar.ID = 2
	newCar.Name = "奔驰"
	newCar.Price = 2000000
	newCar.Color = gs.ColorBlue
	newCar.Owner.Age = 19
	newCar.Owner.Name = "李四"
	newCar.Owner.ID = 12
	newCar.Code = gs.ErrOOS
	newCar.Size = [4]int32{4, 3, 2, 1}
	newCar.Drivers[0].Age = 19
	newCar.Drivers[0].Name = "李四四"
	newCar.Drivers[0].ID = 22
	newCar.Drivers[1].Age = 19
	newCar.Drivers[1].Name = "李四四四"
	newCar.Drivers[1].ID = 32
	newCar.Attrs["颜色"] = "白色"
	newCar.Attrs["车牌"] = "京A66666"
	newCar.Attrs["车型"] = "奔驰S600"
	fmt.Printf("%+v %+v \n%+v\n", car, *car.Owner, car.Drivers)
	fmt.Printf("%+v %+v \n%+v\n", newCar, *newCar.Owner, newCar.Drivers)

}
