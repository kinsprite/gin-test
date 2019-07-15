package main

// UserInfo  from user test service
type UserInfo struct {
	id   int
	name string
}

// ID   get user's ID
func (userInfo *UserInfo) ID() int {
	return userInfo.id
}

// Name   get user's name
func (userInfo *UserInfo) Name() string {
	return userInfo.name
}
