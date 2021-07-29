package dbops

import "database/sql"

type TableFile struct {
	FileHash string         `json:"file_hash"`
	FileName sql.NullString `json:"file_name"`
	FileSize sql.NullInt64  `json:"file_size"`
	FileAddr sql.NullString `json:"file_addr"`
}
