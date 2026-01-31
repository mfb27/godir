package redis

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"godir/internal/common/logger"
	"godir/internal/common/svc"
	"godir/internal/model"

	minioLib "github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

// ThumbnailTask 表示生成缩略图的任务
type ThumbnailTask struct {
	MaterialID  uint   `json:"material_id"`
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	ContentType string `json:"content_type"`
}

// PushThumbnailTask 将缩略图生成任务推送到队列
func PushThumbnailTask(task *ThumbnailTask) error {
	ctx := context.Background()
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 推送到Redis队列
	return svc.Redis().LPush(ctx, "thumbnail_tasks", data).Err()
}

// StartThumbnailWorker 启动处理缩略图任务的工作进程
func StartThumbnailWorker() {
	ctx := context.Background()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Logger.Error("处理缩略图任务发生错误", r)
			}
		}()

		for {
			// 从队列中取出任务
			result, err := svc.Redis().BRPop(ctx, 5*time.Second, "thumbnail_tasks").Result()
			if err != nil && err != redis.Nil {
				logger.Logger.Error("从Redis队列获取任务失败", err)
				<-time.After(5 * time.Second) // 等待5秒后继续
				continue
			}

			// 如果有任务则处理
			if len(result) > 1 {
				var task ThumbnailTask
				err := json.Unmarshal([]byte(result[1]), &task)
				if err != nil {
					logger.Logger.Error("解析任务失败", err)
					continue
				}

				// 处理任务
				processThumbnailTask(&task)
			}
		}
	}()
}

// processThumbnailTask 处理缩略图生成任务
func processThumbnailTask(task *ThumbnailTask) {
	logger.Logger.Info("开始处理缩略图任务", "material_id", task.MaterialID, "key", task.Key)

	minioClient := svc.Minio()
	if minioClient == nil {
		logger.Logger.Error("MinIO客户端未初始化")
		return
	}

	ctx := context.Background()

	// 从MinIO下载对象到临时文件
	obj, err := minioClient.GetObject(ctx, task.Bucket, task.Key, minioLib.GetObjectOptions{})
	if err != nil {
		logger.Logger.Error("获取MinIO对象失败", err)
		return
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "material-*")
	if err != nil {
		logger.Logger.Error("创建临时文件失败", err)
		return
	}

	_, err = io.Copy(tmpFile, obj)
	if err != nil {
		logger.Logger.Error("复制对象到临时文件失败", err)
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return
	}

	_ = tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// 确保ffmpeg存在
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		logger.Logger.Error("找不到ffmpeg命令", err)
		return
	}

	// 生成缩略图
	thumbPath := tmpFile.Name() + ".thumb.jpg"
	var cmd *exec.Cmd

	if strings.HasPrefix(task.ContentType, "image/") {
		cmd = exec.Command("ffmpeg", "-y", "-i", tmpFile.Name(), "-vf", "scale=640:-1", "-vframes", "1", "-q:v", "2", thumbPath)
	} else {
		// 视频文件
		cmd = exec.Command("ffmpeg", "-y", "-i", tmpFile.Name(), "-ss", "00:00:01", "-vframes", "1", "-q:v", "2", thumbPath)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Error("ffmpeg执行失败", err, "output", string(out))
		return
	}

	// 上传缩略图到MinIO
	tf, err := os.Open(thumbPath)
	if err != nil {
		logger.Logger.Error("打开缩略图文件失败", err)
		return
	}

	fi, _ := tf.Stat()
	thumbKey := task.Key + ".thumb.jpg"
	_, err = minioClient.PutObject(ctx, task.Bucket, thumbKey, tf, fi.Size(), minioLib.PutObjectOptions{ContentType: "image/jpeg"})
	_ = tf.Close()

	if err != nil {
		logger.Logger.Error("上传缩略图到MinIO失败", err)
		_ = os.Remove(thumbPath)
		return
	}

	// 更新数据库中的封面信息
	db := svc.DB()
	material := model.GodirMaterial{}
	result := db.First(&material, task.MaterialID)
	if result.Error != nil {
		logger.Logger.Error("查找素材记录失败", result.Error)
		_ = os.Remove(thumbPath)
		return
	}

	material.CoverOssFilePath = thumbKey
	result = db.Model(&material).Updates(map[string]interface{}{"cover_oss_file_path": material.CoverOssFilePath})
	if result.Error != nil {
		logger.Logger.Error("更新素材封面信息失败", result.Error)
	}

	_ = os.Remove(thumbPath)
	logger.Logger.Info("缩略图任务处理完成", "material_id", task.MaterialID)
}
