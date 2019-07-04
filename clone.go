/*
@author : Aeishen
@data :  19/07/04, 12:20

@description : 实现深拷贝与浅拷贝
*/

//深拷贝和浅拷贝是针对复杂数据类型来说的，浅拷贝只拷贝一层，而深拷贝是层层拷贝

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Color int

const (
	black = iota
	green
	blue
)

type Cat struct {
	Name string
	Age  int
	Leg Legs
	Eye  map[string]Eye
	Tail Tails
}

type Eye struct {
	Colors Color
}

type Tails struct {
	Length float64
}

type Legs struct {
	Count  int
	Length []int
}

// 修改猫
func ModifyCat(cat Cat) {
	fmt.Printf("enter value cat:%v\n", cat)
	fmt.Printf("in modify: cat:%p, cat.name:%p, cat.tail:%p\n", &cat, &cat.Name, &cat.Tail)
	fmt.Printf("in modify: cat.legs.length:%p, cat.legs.count:%p\n", &cat.Leg.Length, &cat.Leg.Count)
	cat.Name = "Ben"
	cat.Eye["left"] = Eye{blue}
	cat.Tail = Tails{234.56}
	cat.Leg.Count = 3
	cat.Leg.Length[0] = 0
	fmt.Printf("exit value cat:%v\n\n", cat)
}

//深拷贝最简单的方式是基于序列化和反序列化来实现对象的深度复制（当对象是结构体且结构体中有小写成员变量时，该方式无效）:
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

/*
    Gob 和 bytes.Buffer 简单组合就搞定了.当然,Gob 的底层也是基于 reflect 包完成的.
    在内存中序列化,反序列化对象实体 来完成对象实体的深拷贝这个是大部分语言实现对象深拷贝的惯用法。
*/


//测试浅拷贝
func test_shallowCopy(){

	//创建一只猫
	catA := Cat{
		Name: "tom",
		Age:  1,
		Leg: Legs{Count: 4, Length: []int{10, 10, 10, 10}},
		Eye:  map[string]Eye{"left":{black}, "right":{green}},
		Tail: Tails{123.45},
	}

	//直接将catA赋值给catB（浅拷贝）
	catB := catA
	fmt.Printf("value catA:%v, catB:%v\n", catA, catB)
	fmt.Println()

	//将catB当参数传入ModifyCat
	ModifyCat(catB)
	fmt.Printf("value catA:%v, catB:%v\n", catA, catB)

	// 在外部对catB做改变
	catB.Name = "Ben"
	catB.Eye["right"] = Eye{black}
	catB.Tail = Tails{234.56}
	catB.Leg.Count = 3
	catB.Leg.Length[1] = 0

	fmt.Printf("value catA:%v, catB:%v\n", catA, catB)
}

func test_deepCopy(){

	//创建一只猫
	catA := &Cat{
		Name: "tom",
		Age:  1,
		Leg: Legs{Count: 4, Length: []int{10, 10, 10, 10}},
		Eye:  map[string]Eye{"left":{black}, "right":{green}},
		Tail: Tails{123.45},
	}

	//深度拷贝
	catB := &Cat{}
	if err := deepCopy(catB, catA);err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("value catA:%v, catB:%v\n\n", catA, catB)

	//将catB当参数传入ModifyCat
	ModifyCat(*catB)
	fmt.Printf("value catA:%v, catB:%v\n", catA, catB)

	// 在外部对catB做改变
	catB.Name = "Ben"
	catB.Eye["right"] = Eye{black}
	catB.Tail = Tails{234.56}
	catB.Leg.Count = 3
	catB.Leg.Length[1] = 0

	fmt.Printf("value catA:%v, catB:%v\n", catA, catB)


}

func main() {
	fmt.Printf("\n=====================shallowCopy===============\n")
	test_shallowCopy()
	fmt.Printf("\n=====================deepCopy===============\n")
	test_deepCopy()
}

/*
    测试浅拷贝：
		根据输出可以看到几点：
		1.catB进行赋值之后，不管是cat本身，还是cat内部的各个属性的地址都已经改变了。
		2.在函数内打印cat的几个属性的地址，可以看到和传入之前的catB是不同的。
		3.在函数内对cat的名字修改后，并没有影响catB，cat是传值的。
		4.通过赋值得到的catB，直接修改catB的名字，也没有影响catA的名字。

		几点意外：
		1.我们在函数内对cat的左眼的修改，影响到了外部的catB，甚至影响到了catA！
		2.在外部对catB的右眼的修改，影响到了catA
		3.在函数内外对cat的腿部的长度(slice类型)的修改，都和眼睛有一样的效果。但是对腿部的数量(int类型)的修改则没有。

		意外的原因：
		slice并不是简单的传值关系，就像指针和chan这些引用类型一样，map和slice在复制的时候复制的是引用，
		复制出来的引用与原来的引用都是指向同一个地址，共享底层元素，对复制出来的引用进行操作会影响到底层元素，
		进而影响到其他也指向这些底层元素的引用

    测试深拷贝：
        根据输出可以看到几点：
		1.在函数内打印cat的几个属性的地址，可以看到和传入之前的catB是不同的
        2.无论对catB做任何操作都不会影响到原来的catA

        几点意外：
		1.我们在函数内对cat的左眼的修改，影响到了外部的catB，不会影响到了catA！
		2.在外部对catB的右眼的修改，也不会影响到了catA
		3.在函数内外对cat的腿部的长度(slice类型)的修改，都和眼睛有一样的效果。但是对腿部的数量(int类型)的修改则没有。

    总结：
        golang 完全是按值传递，所以正常的赋值都是值拷贝，当然如果类型里面嵌套的有指针，也是指针值的拷贝，此时就会出现两个类型变量的内部有一部分是共享的。
		深拷贝：
			深拷贝复制变量值，对于非基础类型的变量，则递归至基本变量后再复制，
			深拷贝后的对象与原来的对象是完全隔离的互不影响，对一个对象的修改不会影响另一个

		浅拷贝：
			浅拷贝是将对象的属性依次复制一遍，但当对象的属性值是引用类型时，实际复制的是其引用，
			即新对象和原对象的该属性值指向的是同一块内存地址，当引用指向的值改变时相互影响。而
			当为基本类型的属性值改变时互相不影响
*/


