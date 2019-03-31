package post

import (
    "github.com/go-gin-demo/entity"
    "github.com/go-gin-demo/utils/collections/set"
    UserService "github.com/go-gin-demo/service/user"
    "strconv"
    PostModel "github.com/go-gin-demo/model/post"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
)

func RenderPosts(posts []*entity.Post) ([]*entity.PostEntity, error) {
   uidSet := set.Make()
   for _, post := range posts {
       if post != nil {
           uidSet.Add(post.Uid)
       }
   }
   rawUids := uidSet.ToArray()
   uids := make([]uint64, len(rawUids))
   for i, e := range rawUids {
       uids[i] = e.(uint64)
   }

   userEntities, err := UserService.RenderUsersById(uids)
   if err != nil {
       return nil, err
   }
   userMap := make(map[uint64]*entity.UserEntity)
   for _, user := range userEntities {
       if user != nil {
           userMap[user.ID] = user
       }
   }

   entities := make([]*entity.PostEntity, len(posts))
   for i, post := range posts {
       user, ok := userMap[post.Uid]
       if !ok {
           continue
       }
       entities[i] = &entity.PostEntity{Id:strconv.FormatUint(post.ID, 10), CreatedAt:post.CreatedAt, User:user}
   }
   return entities, nil
}

func RenderPost(post *entity.Post) (*entity.PostEntity, error) {
    posts := []*entity.Post{post}
    entities, err := RenderPosts(posts)
    if err != nil {
        return nil, err
    }
    if len(entities) == 0 {
        return nil, nil
    }
    return entities[0], nil
}

func CreatePost(uid uint64, text string) (*entity.PostEntity, error) {
    post := &entity.Post{Uid:uid, Text:text, Valid:true}
    err := PostModel.Create(post)
    if err != nil {
        return nil, err
    }
    return RenderPost(post)
}

func GetPost(pid uint64) (*entity.PostEntity, error) {
    post, err := PostModel.Get(pid)
    if err != nil {
        return nil, err
    }
    if post == nil {
        return nil,  BizError.NotFound(fmt.Sprintf("post not found: %d", pid))
    }
    return RenderPost(post)
}

func DeletePost(currentUid uint64, pid uint64) error {
    post, err := PostModel.Get(pid)
    if err != nil {
        return err
    }
    if post == nil {
        return BizError.NotFound(fmt.Sprintf("post not found: %d", pid))
    }
    if post.Uid != currentUid {
        return BizError.Forbidden("not your own post")
    }
    return PostModel.Delete(post.ID)
}