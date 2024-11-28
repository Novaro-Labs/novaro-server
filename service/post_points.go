package service

import (
	"fmt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"novaro-server/dao"
	"novaro-server/model"
	"sync"
	"time"
)

type PostPointsService struct {
	dao              *dao.PostPointsDao
	tagRecordDao     *dao.TagsRecordDao
	postPointsDao    *dao.PostPointsDao
	userDao          *dao.UsersDao
	pointsHistoryDao *dao.PointsHistoryDao
	postDao          *dao.PostsDao
}

func NewPostPointsService() *PostPointsService {
	db := model.GetDB()
	return &PostPointsService{
		dao:              dao.NewPostPointsDao(db),
		tagRecordDao:     dao.NewTagsRecordDao(db),
		postPointsDao:    dao.NewPostPointsDao(db),
		userDao:          dao.NewUsersDao(db),
		pointsHistoryDao: dao.NewPointsHistoryDao(db),
		postDao:          dao.NewPostsDao(db),
	}
}

func (s *PostPointsService) Save(m *model.PostPoints) error {
	return s.dao.Save(nil, m)
}

func (s *PostPointsService) Delete(postId string) error {
	return s.dao.Delete(nil, postId)
}

func (s *PostPointsService) SyncData() error {
	records, err := s.tagRecordDao.GetYesterdayTagsRecords()
	if err != nil {
		return fmt.Errorf("failed to get yesterday's tag records: %w", err)
	}

	if len(records) == 0 {
		return nil
	}

	postPoints, userPoints, err := s.processRecords(records)

	if err != nil {
		return fmt.Errorf("failed to process records: %w", err)
	}

	if err := s.savePostPoints(postPoints); err != nil {
		return fmt.Errorf("failed to save post points: %w", err)
	}

	if err := s.pointsHistoryDao.BatchSave(userPoints); err != nil {
		return fmt.Errorf("failed to save user points history: %w", err)
	}
	logger.Debug("sync data success")
	return nil
}

func (s *PostPointsService) processRecords(records []model.TagsRecords) ([]model.PostPoints, []model.PointsHistory, error) {
	postPoints := make([]model.PostPoints, len(records))
	var userPoints []model.PointsHistory
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, item := range records {
		postPoints[i] = model.PostPoints{
			PostId:    item.PostId,
			Points:    item.PostPoints,
			CreatedAt: item.CreatedAt,
		}
		now := time.Now()

		wg.Add(1)
		go func(i int, item *model.TagsRecords) {
			defer wg.Done()
			user, err := s.userDao.GetById(item.UserId)
			if err == nil && user.WalletPublicKey != "" {
				mu.Lock()
				userPoints = append(userPoints, model.PointsHistory{
					Wallet:   user.NftInfo.Wallet,
					Points:   item.Points,
					Status:   0,
					CreateAt: &now,
				})
				mu.Unlock()
			}
		}(i, &item)
	}

	wg.Wait()
	return postPoints, userPoints, nil
}
func (s *PostPointsService) savePostPoints(postPoints []model.PostPoints) error {
	ok := s.postPointsDao.BatchSave(postPoints)
	if ok {
		s.processPostHistory()
	}
	return nil
}

func (s *PostPointsService) processPostHistory() error {
	histories, err := s.postPointsDao.GetYesterdayPostHistory()

	if err != nil {
		return fmt.Errorf("failed to get yesterday's post history: %w", err)
	}

	var userPoints []model.PointsHistory

	for _, history := range histories {

		postUser := s.postDao.GetPostIdByUser(history.PostId)
		if postUser.User.WalletPublicKey == "" {
			continue
		}
		now := time.Now()
		userPoints = append(userPoints, model.PointsHistory{
			Wallet:   postUser.User.WalletPublicKey,
			Points:   history.Points,
			Status:   0,
			CreateAt: &now,
		})
	}

	return s.pointsHistoryDao.BatchSave(userPoints)
}
