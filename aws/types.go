package aws

type AwsOptions struct {
	Url               string
	AccessKeyId       string
	SecretAccessKey   string
	SessionToken      string
	Region            string
	PipelineId        string
	PresetId          string
	PresetWatermarkId string
}
