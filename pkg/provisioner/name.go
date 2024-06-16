package provisioner

import (
	"math/rand"
)

var randomNameList = []string{
	"apple",
	"banana",
	"cherry",
	"darjeeling",
	"elderberry",
	"fig",
	"grape",
	"honeydew",
	"imbe",
	"jackfruit",
	"kiwi",
	"lemon",
	"melon",
	"nectarine",
	"orange",
	"papaya",
	"quince",
	"raambutan",
	"strawberry",
	"tangerine",
	"ugli",
	"vanilla",
	"watermelon",
	"ximenia",
	"yuzu",
	"zucchini",
}

func GetRandomName() string {
	return randomNameList[rand.Intn(len(randomNameList))]
}
