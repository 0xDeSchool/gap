package aws

import (
	"context"
	"path"
	"strings"

	"github.com/0xDeSchool/gap/errx"
	"github.com/0xDeSchool/gap/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elastictranscoder"
	"github.com/aws/aws-sdk-go-v2/service/elastictranscoder/types"
	"github.com/google/uuid"
)

// TransStatus section 的转换状态，如Completed、Pedding、Processing、Error
type TransStatus string

const (
	TransCompleted  TransStatus = "Completed"
	TransPending    TransStatus = "Pending"
	TransProcessing TransStatus = "Processing"
	TransFailed     TransStatus = "Error"
)

type VideoTranscoder struct {
	client *elastictranscoder.Client
	opts   *AwsOptions
}

func NewVideoTranscoder(cfg *aws.Config, opts *AwsOptions) *VideoTranscoder {
	client := elastictranscoder.NewFromConfig(*cfg)
	return &VideoTranscoder{
		client: client,
		opts:   opts,
	}
}

func (vt *VideoTranscoder) Trans(key string, watermarkKey string, userData map[string]string) string {
	key = GetS3Path(key)
	watermarkKey = GetS3Path(watermarkKey)
	dirs := strings.Split(key, "/")
	dirs[0] = "v/1080p"
	dir := strings.Join(dirs[:len(dirs)-1], "/") + "/"
	ext := path.Ext(dirs[len(dirs)-1])
	outputKey := uuid.New().String() + ext
	output, err := vt.client.CreateJob(context.Background(), &elastictranscoder.CreateJobInput{
		PipelineId: &vt.opts.PipelineId,
		Input: &types.JobInput{
			Key: &key,
		},
		OutputKeyPrefix: &dir,
		Output: &types.CreateJobOutput{
			PresetId: &vt.opts.PresetId,
			Key:      &outputKey,
			Watermarks: []types.JobWatermark{
				{
					InputKey:          &watermarkKey,
					PresetWatermarkId: &vt.opts.PresetWatermarkId,
				},
			},
		},
		UserMetadata: userData,
	})
	errx.CheckError(err)
	log.Info("created video job: " + *output.Job.Id)
	return *output.Job.Id
}

func (vt *VideoTranscoder) ReadStatus(id string) (TransStatus, string) {
	job, err := vt.client.ReadJob(context.Background(), &elastictranscoder.ReadJobInput{
		Id: &id,
	})
	errx.CheckError(err)
	status := *job.Job.Status
	key := *job.Job.OutputKeyPrefix + *job.Job.Output.Key
	if status == "Submitted" {
		return TransPending, ""
	} else if status == "Progressing" {
		return TransProcessing, ""
	} else if status == "Complete" {
		return TransCompleted, key
	} else if status == "Error" {
		return TransFailed + ": " + TransStatus(*job.Job.Output.StatusDetail), ""
	}
	return TransStatus(status), ""
}

func GetS3Path(p string) string {
	p = strings.TrimPrefix(p, "https://s3.us-east-1.amazonaws.com/deschool/")
	p = strings.TrimPrefix(p, "https://deschool.s3.amazonaws.com/")
	p = strings.TrimPrefix(p, "https://deschooldev.s3.amazonaws.com/")
	return p
}

type S3Key struct {
	Bucket string
	Key    string
}

func GetS3Key(p string) S3Key {
	path := strings.TrimPrefix(p, "https://")
	paths := strings.Split(path, "/")
	bucket := strings.Split(paths[0], ".")[0]
	key := strings.Join(paths[1:], "/")
	if bucket == "s3" {
		bucket = paths[1]
		key = strings.Join(paths[2:], "/")
	}
	return S3Key{
		bucket,
		key,
	}
}
