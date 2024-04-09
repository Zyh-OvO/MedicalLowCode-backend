package model

import (
	"database/sql"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var UserFileRootDirPath = "./file/userFile/"

type Directory struct {
	DirId   int `gorm:"primaryKey;autoIncrement"`
	UserId  int
	DirName string
	IsRoot  bool
}

func (d Directory) TableName() string {
	return "directory"
}

type DirectoryPath struct {
	AncestorId   int
	DescendantId int
	Depth        int
}

func (d DirectoryPath) TableName() string {
	return "directory_path"
}

type File struct {
	FileId   int `gorm:"primaryKey;autoIncrement"`
	DirId    int
	UserId   int
	FileName string
}

func (f File) TableName() string {
	return "file"
}

type queryDirPathResult struct {
	DirId   int    `gorm:"column:ancestor_id"`
	DirName string `gorm:"column:dir_name"`
	Depth   int    `gorm:"column:depth"`
}

func QueryDirPath(userId, dirId int) (string, error) {
	var results []queryDirPathResult
	sql1 := "select ancestor_id, dir_name, depth from directory_path left join directory on directory_path.ancestor_id = directory.dir_id where user_id = ? and descendant_id = ? order by depth desc"
	if err := DB.Raw(sql1, userId, dirId).Scan(&results).Error; err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", os.ErrNotExist
	}
	dirPath := UserFileRootDirPath
	for _, r := range results {
		dirPath = filepath.Join(dirPath, r.DirName)
	}
	return dirPath, nil
}

func QueryFilePath(userId, fileId int) (string, error) {
	var file File
	if err := DB.Where("user_id = ? and file_id = ?", userId, fileId).First(&file).Error; err != nil {
		return "", err
	}
	dirPath, err := QueryDirPath(userId, file.DirId)
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, file.FileName), nil
}

func NewDirectory(userId int, parentDirId int, dirName string) *Directory {
	dir := Directory{
		UserId:  userId,
		DirName: dirName,
		IsRoot:  false,
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		//sql1
		if err := tx.Create(&dir).Error; err != nil {
			return err
		}
		//sql2
		insertPathSql := "insert into directory_path (ancestor_id, descendant_id, depth) select dp.ancestor_id, @newDirId, dp.depth + 1 from directory_path as dp where dp.descendant_id = @parentDirId union all select @newDirId, @newDirId, 0"
		if err := tx.Exec(insertPathSql, sql.Named("newDirId", dir.DirId), sql.Named("parentDirId", parentDirId)).Error; err != nil {
			return err
		}
		//sql3
		if parentDirPath, err := QueryDirPath(userId, parentDirId); err != nil {
			return err
		} else {
			targetDirPath := filepath.Join(parentDirPath, dirName)
			//检查目录是否存在
			if _, err := os.Stat(targetDirPath); os.IsNotExist(err) {
				// 使用 MkdirAll 函数递归创建目录
				if err := os.MkdirAll(targetDirPath, os.ModePerm); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				return os.ErrExist
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &dir
}

func DeleteDirectory(userId int, dirId int) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if dirPath, err := QueryDirPath(userId, dirId); err != nil {
			return err
		} else {
			if result := tx.Where("user_id = ? and dir_id = ? and is_root = ?", userId, dirId, false).Delete(&Directory{}); result.Error != nil {
				return result.Error
			} else if result.RowsAffected == 0 {
				return os.ErrNotExist
			}
			if err := os.RemoveAll(dirPath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RenameDirectory(userId int, dirId int, dirName string) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if dirPath, err := QueryDirPath(userId, dirId); err != nil {
			return err
		} else {
			newDirPath := filepath.Join(filepath.Dir(dirPath), dirName)
			if result := tx.Model(&Directory{}).Where("user_id = ? and dir_id = ? and is_root = ?", userId, dirId, false).Update("dir_name", dirName); result.Error != nil {
				return result.Error
			} else if result.RowsAffected == 0 {
				return os.ErrNotExist
			}
			if err := os.Rename(dirPath, newDirPath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func UploadFile(userId, dirId int, srcFile multipart.File, fileName string) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if dirPath, err := QueryDirPath(userId, dirId); err != nil {
			return err
		} else {
			if err := tx.Create(&File{
				DirId:    dirId,
				UserId:   userId,
				FileName: fileName,
			}).Error; err != nil {
				return err
			}
			targetFilePath := filepath.Join(dirPath, fileName)
			if _, err := os.Stat(targetFilePath); err == nil {
				return os.ErrExist
			}
			tgtFile, err := os.Create(targetFilePath)
			if err != nil {
				return err
			}
			defer tgtFile.Close()
			if _, err := io.Copy(tgtFile, srcFile); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func DeleteFile(userId, fileId int) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if filePath, err := QueryFilePath(userId, fileId); err != nil {
			return err
		} else {
			if err := tx.Where("user_id = ? and file_id = ?", userId, fileId).Delete(&File{}).Error; err != nil {
				return err
			}
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func RenameFile(userId, fileId int, fileName string) {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if filePath, err := QueryFilePath(userId, fileId); err != nil {
			return err
		} else {
			newFilePath := filepath.Join(filepath.Dir(filePath), fileName)
			if err := tx.Model(&File{}).Where("user_id = ? and file_id = ?", userId, fileId).Update("file_name", fileName).Error; err != nil {
				return err
			}
			if err := os.Rename(filePath, newFilePath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func GetFilesUnderDir(userId, dirId int) []File {
	var files []File
	if err := DB.Where("user_id = ? and dir_id = ?", userId, dirId).Find(&files).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return files
}

type GetDirContentResult struct {
	DirId   int    `gorm:"column:descendant_id"`
	DirName string `gorm:"column:dir_name"`
}

func GetDirsUnderDir(userId, dirId int) []GetDirContentResult {
	var results []GetDirContentResult
	sql1 := "select descendant_id, dir_name from directory_path left join directory on directory_path.descendant_id = directory.dir_id where user_id = ? and ancestor_id = ? and depth = 1"
	if err := DB.Raw(sql1, userId, dirId).Scan(&results).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return results
}

func GetRootDir(userId int) *Directory {
	var dir Directory
	if err := DB.Where("user_id = ? and is_root = ?", userId, true).First(&dir).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return &dir
}
