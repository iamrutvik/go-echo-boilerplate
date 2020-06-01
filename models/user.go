package models

import (
	_ "github.com/go-playground/validator/v10"
)

type UserType string
type LoginSource string

type User struct {
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"firstName" validate:"required,alpha,min=2,max=100"`
	LastName       string `json:"lastName" validate:"required,alpha,min=2,max=100"`
	Mobile         string `json:"mobile" validate:"required,numeric,len=10"`
	Birthday       string `json:"birthday" validate:"datetime"`
	Password       string `json:"password" validate:"required,containsany=!@#?,min=6"`
	ProfilePicture string `json:"profilePicture" validate:"omitempty,url"`
	UserType       UserType `json:"userType" validate:"required,oneof=PUBLISHER ADVERTISER USER ADMIN ADMIN_RESTRICTED"`
	LoginSource    LoginSource `json:"loginSource" validate:"required,oneof=FACEBOOK GOOGLE APPLE LOCAL"`
	IsVerified     bool   `json:"isVerified" validate:"omitempty"`
	IsActive       bool   `json:"isActive" validate:"omitempty"`
	Otp            int32  `json:"otp" validate:"omitempty,numeric,len=6"`
}

type LoginUser struct {
	Email          string `json:"email" validate:"required_without_all=Mobile,omitempty,email"`
	Password       string `json:"password" validate:"required,containsany=!@#?,min=6"`
}

type ChangePassword struct {
	Password       string `json:"password" validate:"required,containsany=!@#?,min=6"`
}

