package canal


import (
    "github.com/go-gin-demo/lib/canal"
    PostService "github.com/go-gin-demo/app/service/post"
    UserService "github.com/go-gin-demo/app/service/user"
)

func GetConsumerMap() map[string]canal.ConsumerFunc {
    consumerMap := make(map[string]canal.ConsumerFunc)

    consumerMap["go.posts:insert"] = PostService.AfterCreatePost
    consumerMap["go.posts:update"] = PostService.AfterUpdatePost

    consumerMap["go.users:insert"] = UserService.AfterCreateUser
    consumerMap["go.users:update"] = UserService.AfterUpdateUser

    consumerMap["go.follows:insert"] = UserService.AfterCreateFollow
    consumerMap["go.follows:update"] = UserService.AfterUpdateFollow

    return consumerMap
}