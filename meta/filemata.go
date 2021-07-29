package meta

import "File-Uploader-go/dbops"

// FileMeta 文件元信息
type FileMeta struct {
	FileSha1 string `json:"file_sha_1"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Location string `json:"location"`
	UploadAt string `json:"upload_at"`
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 新增或更新文件员信息
func UpdateFileMeta(meta FileMeta) {
	fileMetas[meta.FileSha1] = meta
}

// UpdateFileMetaDB 更新文件元数据到数据库中
func UpdateFileMetaDB(meta FileMeta) error {
	err := dbops.OnFileUploadFinished(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
	if err != nil {
		return err
	}
	return nil
}

// GetFileMeta 通过Sha获取文件的元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB 从数据库获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tFile, err := dbops.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}

	fMeta := FileMeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return fMeta, nil
}

// RemoveFileMeta 删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
