package dao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"novaro-server/model"
	"strconv"
	"sync"
	"time"
)

type TagsRecordDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewTagsRecordDao(db *gorm.DB) *TagsRecordDao {

	return &TagsRecordDao{
		db:  db,
		rdb: model.GetRedisCli(),
	}
}

func (d *TagsRecordDao) AddTagsRecords(t *model.TagsRecords) error {
	err := d.db.Create(&t).Error
	return err
}

func (d *TagsRecordDao) Delete(t *model.TagsRecords) error {
	err := d.db.Where("user_id = ? and post_id = ? and tag_id = ?", t.UserId, t.PostId, t.TagId).Delete(&t).Error
	return err
}

func (d *TagsRecordDao) GetRecord(tagId, postId, userId string) int64 {
	var count int64
	d.db.Where("tag_id = ? and post_id = ? and user_id = ?", tagId, postId, userId).Count(&count)
	return count
}

func (d *TagsRecordDao) GetCountByUserId(userId string) (int64, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	var count int64
	err := d.db.Where("user_id = ? and created_at >= ? AND created_at < ?", userId, startOfDay, endOfDay).Count(&count).Error
	return count, err
}

func (d *TagsRecordDao) GetTagsRecordsByPostId(post *model.Posts) ([]model.TagRecordResponse, int, error) {
	var wg sync.WaitGroup
	resultMap := make(map[string]*model.TagRecordResponse)
	var mysqlErr, redisErr error
	var redisResult map[string]string

	var sqlResults []model.TagRecordResponse

	wg.Add(1)
	go func() {
		defer wg.Done()
		mysqlErr = d.db.Model(&model.TagsRecords{}).
			Select("tag_id as id, count(tag_id) as count").
			Where("post_id = ?", post.Id).
			Group("tag_id").
			Scan(&sqlResults).Error

		for _, result := range sqlResults {
			resultMap[result.Id] = &result
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		key := fmt.Sprintf("tags:count:%s", post.Id)
		redisResult, redisErr = d.rdb.HGetAll(context.Background(), key).Result()
	}()

	// 等待两个协程完成
	wg.Wait()

	// 检查错误
	if mysqlErr != nil {
		return nil, 0, fmt.Errorf("error querying MySQL: %w", mysqlErr)
	}
	if redisErr != nil {
		return nil, 0, fmt.Errorf("error querying Redis: %w", redisErr)
	}

	// 合并 Redis 结果
	for field, value := range redisResult {

		count, err := strconv.Atoi(value)
		if err != nil {
			count = 0
		}

		if record, exists := resultMap[field]; exists {

			record.Count += count
		} else {
			resultMap[field] = &model.TagRecordResponse{
				Id:    field,
				Count: count,
			}
		}
	}

	// 转换为切片
	result := make([]model.TagRecordResponse, 0, len(resultMap))
	var total int
	for _, record := range resultMap {
		total += record.Count
		result = append(result, *record)
	}

	return result, total, nil
}
func (d *TagsRecordDao) Points(wattle *string, nftLevel int) int64 {
	if wattle == nil {
		return 0
	}
	defaultPoints := 5
	rewards := nftLevel

	return int64((nftLevel * defaultPoints) + rewards)
}
