package gitStudy

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func Hello(sub string) {
	_ = mongo.ClientStream
	var version string = "1.4.1"
	fmt.Println(version, sub)
}
