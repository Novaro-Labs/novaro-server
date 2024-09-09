package service

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
	"sync"
)

type PostService struct {
	dao              *dao.PostsDao
	rdb              *redis.Client
	collectionDao    *dao.CollectionsDao
	tagsDao          *dao.TagsDao
	imgsDao          *dao.ImgsDao
	userDao          *dao.UsersDao
	pointsHistoryDao *dao.PointsHistoryDao
	commentsDao      *dao.CommentsDao
	tagRecordDao     *dao.TagsRecordDao
}

func NewPostService() *PostService {
	db := model.GetDB()
	return &PostService{
		dao:              dao.NewPostsDao(db),
		rdb:              model.GetRedisCli(),
		collectionDao:    dao.NewCollectionsDao(db),
		tagsDao:          dao.NewTagsDao(db),
		imgsDao:          dao.NewImgsDao(db),
		pointsHistoryDao: dao.NewPointsHistoryDao(db),
		commentsDao:      dao.NewCommentsDao(db),
		userDao:          dao.NewUsersDao(db),
		tagRecordDao:     dao.NewTagsRecordDao(db),
	}
}

func (s *PostService) GetList(p *model.PostsQuery) (resp []model.Posts, err error) {
	resp, err = s.dao.GetPostsList(p)

	var wg sync.WaitGroup
	// 使用缓冲通道作为信号量来限制并发 goroutine 的数量
	semaphore := make(chan struct{}, 10) // 最多10个并发 goroutine

	for i := range resp {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 获取标签
			tags, total, _ := s.tagRecordDao.GetTagsRecordsByPostId(&resp[i])
			resp[i].Tags = tags
			resp[i].TagsAmount = total
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

func (s *PostService) GetById(id string) (*model.Posts, error) {
	resp, err := s.dao.GetPostsById(id)
	tags, total, err := s.tagRecordDao.GetTagsRecordsByPostId(resp)
	resp.Tags = tags
	resp.TagsAmount = total
	return resp, err
}

func (s *PostService) PostExists(id string) (bool, error) {
	return s.dao.PostExists(id)
}

func (s *PostService) GetByUserId(userId string) ([]model.Posts, error) {
	resp, err := s.dao.GetPostsByUserId(userId)

	// 处理标签
	for i := range resp {
		tags, total, err := s.tagRecordDao.GetTagsRecordsByPostId(&resp[i])
		if err != nil {
			resp[i].Tags = nil
		}
		resp[i].Tags = tags
		resp[i].TagsAmount = total
	}
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

		// 计算积分
		count, err := s.dao.GetCountByUserId(posts.UserId)
		if err != nil {
			return fmt.Errorf("failed to get post count: %w", err)
		}
		point := s.calculatePoints(count)

		// 创建积分历史记录
		if user.WalletPublicKey != nil {
			history := model.PointsHistory{
				Wallet: user.WalletPublicKey,
				Points: float64(point),
				Status: 0,
			}

			if err := s.pointsHistoryDao.Create(tx, &history); err != nil {
				return fmt.Errorf("failed to create points history: %w", err)
			}
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
	return s.dao.Delete(id)
}

func (s *PostService) calculatePoints(value int64) int {
	switch value {
	case 0:
		return 20
	case 1:
		return 10
	default:
		return 0
	}
}
