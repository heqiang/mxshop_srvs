package handle

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"strings"
	"time"
)

type UserService struct{}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func ModelToResponse(user model.User) *proto.UserInfo {
	//
	userInfoRsp := proto.UserInfo{
		Id:       user.ID,
		Password: user.Password,
		NickNmae: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return &userInfoRsp
}
func (user *UserService) GetUserList(ctc context.Context, req *proto.PageInfo) (res *proto.UserListReponse, err error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListReponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userinfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, userinfoRsp)
	}
	return rsp, nil
}

// GetUserByMobile 通过手机号码查询用户
func (user *UserService) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (res *proto.UserInfo, err error) {
	var userinfo model.User
	result := global.DB.Where(&model.User{
		Mobile: req.Mobile,
	}).Find(&userinfo)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "不存在该用户")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userinforsp := ModelToResponse(userinfo)
	return userinforsp, nil
}

// GetUserById 通过id查询用户
func (user *UserService) GetUserById(ctx context.Context, req *proto.IdRequest) (res *proto.UserInfo, err error) {
	var userinfo model.User
	result := global.DB.First(&userinfo, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "不存在该用户")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userinforsp := ModelToResponse(userinfo)
	return userinforsp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (res *proto.UserInfo, err error) {
	var userinfo model.User
	result := global.DB.Where(&model.User{
		Mobile: req.Mobile,
	}).Find(&userinfo)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "该手机已被注册")
	}
	userinfo.Mobile = req.Mobile
	userinfo.NickName = req.NickName
	//密码加密
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	userinfo.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	resu := global.DB.Create(&userinfo)
	if resu.Error != nil {
		return nil, status.Errorf(codes.Internal, resu.Error.Error())
	}

	protouserinfo := ModelToResponse(userinfo)
	return protouserinfo, nil
}

// UpdateUser 个人中心用户信息更新
func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (res *proto.Empty, err error) {
	var userinfo model.User
	result := global.DB.First(&userinfo, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthday := time.Unix(int64(req.Birthday), 0)
	userinfo.NickName = req.NickName
	userinfo.Birthday = &birthday
	userinfo.Gender = req.Gender

	result = global.DB.Save(&userinfo)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &proto.Empty{}, nil

}

func (s *UserService) CheckPassword(ctx context.Context, req *proto.CheckPwd) (res *proto.CheckResponse, err error) {
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}

	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{
		Success: check,
	}, nil
}
