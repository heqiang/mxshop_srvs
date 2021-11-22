package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_srvs/user_srv/proto"
)

func main() {
	address := fmt.Sprintf("127.0.0.1:8081")
	dial, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return
	}
	// 信息传输 metadata
	client := proto.NewUserClient(dial)
	//获取所有的用户
	pageInfo := proto.PageInfo{}
	list, err := client.GetUserList(context.Background(), &pageInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list)
	//通过id查询用户
	//idInfo := proto.IdRequest{
	//	Id: 21,
	//}
	//userbyId, err := client.GetUserById(context.Background(), &idInfo)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(userbyId)
	//user, err := client.CreateUser(context.Background(), &proto.CreateUserInfo{
	//	NickName: "hq",
	//	PassWord: "123456",
	//	Mobile:   "18282132023",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(user)
	////用户更新
	//_, err = client.UpdateUser(context.Background(), &proto.UpdateUserInfo{
	//	Id:       5,
	//	NickName: "更新",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//通过手机好查询用户
	//mobile, err := client.GetUserByMobile(context.Background(), &proto.MobileRequest{
	//	Mobile: "18282132023",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(mobile)
	////密码检查
	//checkBool, err := client.CheckPassword(context.Background(), &proto.CheckPwd{
	//	EncryptedPassword: mobile.Password,
	//	Password:          "123456",
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//if checkBool.Success {
	//	fmt.Println("ok")
	//} else {
	//	fmt.Println("no")
	//	fmt.Println(err)
	//}

}
