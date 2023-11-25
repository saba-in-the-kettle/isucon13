package isuutil

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"path/filepath"
)

const limit = 100

// TableSchema は本番ではテーブル構造に応じてよしなに書き換える。
type TableSchema struct {
	ID      int    `db:"id"`
	Imgdata []byte `db:"imgdata"`
	Mime    string `db:"mime"`
}

// ExportImages はDBに保存されている画像をファイルに出力する。
// 本番では適切にクエリを書き換える必要がある。
// この関数は一度実行すれば、initializeのタイミングで再実行する必要はない。
// initializeでは MakeSymbolicLinks 関数を使って、実際に使うディレクトリにシンボリックリンクを張る。
func ExportImages(db *sqlx.DB, dirName string) error {
	if err := os.MkdirAll(dirName, 0777); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	var rows []TableSchema

	seekID := 0
	for {
		fmt.Println("seekID:", seekID)

		// LIMIT OFFSETだと重いのでシーク法でいく
		query := `SELECT id, imgdata, mime FROM posts WHERE id > ? ORDER BY id LIMIT ?`
		err := db.Select(&rows, query, seekID, limit)
		if err != nil {
			return err
		}

		for _, row := range rows {
			if err := saveImage(dirName, row); err != nil {
				return fmt.Errorf("failed to save image: %w", err)
			}
		}

		fmt.Println("len(rows):", len(rows))
		if len(rows) < limit {
			break
		}
		seekID = rows[len(rows)-1].ID
	}

	return nil
}

// MakeSymbolicLinks は、 ExportImages で書き出した画像を実際に使うディレクトリにシンボリックリンクを張る。
// srcDir は ExportImages で書き出したディレクトリ。何度ベンチを実行しても中身は変わらない想定。
// distDir は実際にNginxから参照されたり、ベンチ時に追加で書き込みが行われるディレクトリ。
func MakeSymbolicLinks(srcDir string, distDir string) error {
	// すでにシンボリックリンクが張られている場合は削除する
	if err := os.RemoveAll(distDir); err != nil {
		return fmt.Errorf("failed to remove all: %w", err)
	}

	if err := os.MkdirAll(distDir, 0777); err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}

	// srcDirに含まれているすべてのファイルごとにシンボリックリンクを張る
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)

		if err != nil {
			return fmt.Errorf("failed to walk: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		// srcDirからの相対パスを取得する
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// distDirにシンボリックリンクを張る
		if err := os.Symlink(filepath.Join("../", srcDir, relPath), filepath.Join(distDir, relPath)); err != nil {
			return fmt.Errorf("failed to symlink: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk: %w", err)
	}

	return nil
}

func saveImage(dirName string, row TableSchema) error {
	filename := filepath.Join(dirName, fmt.Sprintf("%d.jpg", row.ID))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)

	}
	defer f.Close()

	reader := bytes.NewReader(row.Imgdata)

	_, err = io.Copy(f, reader)
	if err != nil {

	}
	return nil
}
