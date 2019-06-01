package post

import (
    "github.com/go-gin-demo/app/entity"
    "github.com/go-gin-demo/lib/collections/set"
    UserService "github.com/go-gin-demo/app/service/user"
    "strconv"
    PostModel "github.com/go-gin-demo/app/model/post"
    BizError "github.com/go-gin-demo/lib/errors"
    "fmt"
    UserTimelineModel "github.com/go-gin-demo/app/model/timeline/user"
    "github.com/go-gin-demo/lib/canal"
    "github.com/mitchellh/mapstructure"
    "github.com/go-gin-demo/app/context/context"
)

func RenderPosts(currentUid uint64, posts []*entity.Post) ([]*entity.PostEntity, error) {
   if len(posts) == 0 {
       return []*entity.PostEntity{}, nil
   }
   uidSet := set.MakeUint64Set()
   for _, post := range posts {
       if post != nil {
           uidSet.Add(post.Uid)
       }
   }
   uids := uidSet.ToArray()

   userEntities, err := UserService.RenderUsersById(currentUid, uids)
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
       entities[i] = &entity.PostEntity{
           Id:strconv.FormatUint(post.ID, 10),
           CreatedAt:post.CreatedAt,
           User:user,
           Text:post.Text,
       }
   }
   return entities, nil
}

func RenderPost(currentUid uint64, post *entity.Post) (*entity.PostEntity, error) {
    posts := []*entity.Post{post}
    entities, err := RenderPosts(currentUid, posts)
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
    if !context.EnableCanal() {
        UserTimelineModel.Push(post)
    }
    return RenderPost(uid, post)
}

// canal handler
func AfterCreatePost(former canal.Row, row canal.Row) error {
    var post entity.Post
    err := mapstructure.Decode(row, &post)
    if err != nil {
        return err
    }
    if !post.Valid {
        return nil
    }
    err = PostModel.AfterCreate(&post)
    if err != nil {
        return err
    }
    return UserTimelineModel.Push(&post)
}

func GetPost(currentUid uint64, pid uint64) (*entity.PostEntity, error) {
    post, err := PostModel.Get(pid)
    if err != nil {
        return nil, err
    }
    if post == nil {
        return nil,  BizError.NotFound(fmt.Sprintf("post not found: %d", pid))
    }
    return RenderPost(currentUid, post)
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
    err = PostModel.Delete(post)
    if err != nil {
        return err
    }
    if !context.EnableCanal() {
        err = UserTimelineModel.Remove(post)
    }
    return err
}

func AfterUpdatePost(former canal.Row, row canal.Row) error {
    var post entity.Post
    err := mapstructure.Decode(row, &post)
    if err != nil {
        return err
    }
    if !post.Valid {
        err = PostModel.AfterDelete(&post)
        if err != nil {
            return err
        }
        err = UserTimelineModel.Remove(&post)
        return err
    }
    return nil
}