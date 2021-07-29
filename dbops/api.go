package dbops

import (
	"File-Uploader-go/dbops/mysql"
	"database/sql"
	"errors"
	"log"
)

var dbConn *sql.DB

func init() {
	dbConn = mysql.DBConn()
}

// OnFileUploadFinished 文件上传，保存 Meta 到表中
func OnFileUploadFinished(fileHash, fileName string, fileSize int64, fileAddr string) error {
	stmtIns, err := dbConn.Prepare("INSERT IGNORE INTO tbl_file (file_sha1, file_name, file_size, file_addr, status) VALUES (?, ?, ?, ?, 1)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	ret, err := stmtIns.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		return err
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Printf("File with hash: %s has been upload before", fileHash)
		}
		return nil
	}
	return errors.New("known db exec error")
}

// GetFileMeta 从表单获取文件元信息
func GetFileMeta(fileHash string) (*TableFile, error) {
	stmtOut, err := dbConn.Prepare("SELECT file_addr, file_name, file_size FROM tbl_file where file_sha1 = ? and status = 1 LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	tFile := TableFile{}
	err = stmtOut.QueryRow(fileHash).Scan(&tFile.FileName, &tFile.FileName, &tFile.FileSize)
	if err != nil {
		return nil, err
	}
	tFile.FileHash = fileHash

	return &tFile, nil
}
