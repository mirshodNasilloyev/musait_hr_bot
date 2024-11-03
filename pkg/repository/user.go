package repository

type Bucket string

const (
	ApiKey        Bucket = "api_key"
	SpreadSheetId Bucket = "spreadsheet_id"
)

type UserRepository interface {
	Save(userID int64, vakue string, bucket Bucket) error
	Get(userID int64, bucket Bucket) (string, error)
}
