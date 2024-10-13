package main

//
//import (
//	"fmt"
//	"kubeS3/internal/controller"
//)
//
//// this main is for testing only
//func main() {
//
//	sess, err := controller.CreateSession("us-west-1")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//log session
//
//	fmt.Println(controller.IfBucketExistsOnS3(sess, "testbucket"))
//
//	fmt.Println("Creating bucket")
//
//	controller.CreateBucket(sess, "testbucket")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	//fmt.Println(sess)
//
//}
