package service

import (
	"fmt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
	"sync"
	"time"
)

type PostService struct {
	dao              *dao.PostsDao
	collectionDao    *dao.CollectionsDao
	tagsDao          *dao.TagsDao
	imgsDao          *dao.ImgsDao
	userDao          *dao.UsersDao
	pointsHistoryDao *dao.PointsHistoryDao
	commentsDao      *dao.CommentsDao
	tagRecordDao     *dao.TagsRecordDao
	postPointsDao    *dao.PostPointsDao
}

func NewPostService() *PostService {
	db := model.GetDB()
	return &PostService{
		dao:              dao.NewPostsDao(db),
		collectionDao:    dao.NewCollectionsDao(db),
		tagsDao:          dao.NewTagsDao(db),
		imgsDao:          dao.NewImgsDao(db),
		pointsHistoryDao: dao.NewPointsHistoryDao(db),
		commentsDao:      dao.NewCommentsDao(db),
		userDao:          dao.NewUsersDao(db),
		tagRecordDao:     dao.NewTagsRecordDao(db),
		postPointsDao:    dao.NewPostPointsDao(db),
	}
}

func (s *PostService) GetLikeByUser(userId string) ([]model.Posts, error) {

	return s.dao.GetLikeByUser(userId)
}

func (s *PostService) GetList(p *model.PostsQuery) (resp []model.Posts, err error) {
	resp, err = s.dao.GetPostsList(p)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 100)

	for i := range resp {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 获取标签
			//tags, total, _ := s.tagRecordDao.GetTagsRecordsByPostId(&resp[i])
			//resp[i].TagResp = tags
			//resp[i].TagsAmount = total
			//list, _ := s.tagsDao.GetTagsList()
			//resp[i].Tags = list

			count, _ := s.commentsDao.GetCommentCount(resp[i].Id)
			resp[i].CommentsAmount = count

			// 获取图片
			source, err2 := s.imgsDao.GetImgsBySourceId(resp[i].SourceId)
			if err2 != nil {
				resp[i].Imgs = nil
			} else {
				resp[i].Imgs = source
			}

		}(i)
	}

	wg.Wait()

	return resp, err
}

func (s *PostService) GetById(id string) (model.Posts, error) {
	resp, err := s.GetList(&model.PostsQuery{
		Id:   id,
		Page: 1,
		Size: 10,
	})
	//tags, total, err := s.tagRecordDao.GetTagsRecordsByPostId(&resp)
	//resp.TagResp = tags
	//resp.TagsAmount = total
	//list, _ := s.tagsDao.GetTagsList()
	//resp.Tags = list
	return resp[0], err
}

func (s *PostService) PostExists(id string) (bool, error) {
	return s.dao.PostExists(id)
}

func (s *PostService) GetByUserId(userId string) ([]model.Posts, error) {
	resp, err := s.dao.GetPostsByUserId(userId)

	// 处理标签
	//for i := range resp {
	//	tags, total, err := s.tagRecordDao.GetTagsRecordsByPostId(&resp[i])
	//	if err != nil {
	//		resp[i].Tags = nil
	//	}
	//	resp[i].TagResp = tags
	//	resp[i].TagsAmount = total
	//}
	return resp, err
}

func (s *PostService) Save(posts *model.Posts) error {
	user, err := s.userDao.GetById(posts.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = model.GetDB().Transaction(func(tx *gorm.DB) error {
		// 保存帖子
		if err := s.dao.Save(tx, posts); err != nil {
			return fmt.Errorf("failed to save post: %w", err)
		}

		points, err2 := s.dao.CalculateCommission(posts.UserId)
		if err2 != nil {
			points = 0
		}

		if user.WalletPublicKey != nil {
			err2 := s.postPointsDao.Save(tx, &model.PostPoints{
				PostId:    posts.Id,
				Points:    float64(points),
				CreatedAt: time.Now(),
			})
			logger.Errorf("failed to save post points: %v", err2)
		}

		return nil
	})
	return err
}

func (s *PostService) Update(posts *model.Posts) error {
	return s.dao.Update(posts)
}

func (s *PostService) UpdateBatch(posts []model.Posts) error {
	return s.dao.UpdateBatch(posts)
}

func (s *PostService) Delete(id string) error {
	err := model.GetDB().Transaction(func(tx *gorm.DB) error {
		err := s.dao.Delete(tx, id)
		if err != nil {
			logger.Errorf("failed to delete post: %v", err)
			return err
		}

		err = s.postPointsDao.Delete(tx, id)
		if err != nil {
			logger.Errorf("failed to delete post_points: %v", err)
			return err
		}

		return nil
	})
	return err
}
