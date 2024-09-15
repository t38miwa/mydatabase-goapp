package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Album 構造体はアルバムの情報を格納します
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

// PostgreSQLに接続し、クエリを実行する関数
func main() {
	// PostgreSQLに接続するためのDSN (Data Source Name)
	connStr := "user=miwatakuma password=kinako0325 dbname=mydb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open the database: ", err)
	}
	defer db.Close()

	// データベースへの接続を確認
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	// John Coltraneのアルバムを検索する
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}

	// 検索結果を表示
	fmt.Printf("Albums found: %v\n", albums)

	// ID 2のアルバムを検索する
	alb, err := albumByID(2)
    if err != nil {
        log.Fatal(err)
    }

    // 検索結果を表示
    fmt.Printf("Album found: %v\n", alb)

	titlealb, err := albumByTitle("Jeru")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album found: %v\n", titlealb)

	// 新しいアルバムをデータベースに追加
    albID, err := addAlbum(Album{
        Title:  "The Modern Sound of Betty Carter",
        Artist: "Betty Carter",
        Price:  49.99,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 新しいアルバムIDを表示
    fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist 関数は指定されたアーティスト名のアルバムを検索します
func albumsByArtist(name string) ([]Album, error) {
	// 検索結果を格納するAlbum型のスライス
	var albums []Album

	// SQLクエリを実行してアルバムを取得
	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	// 検索結果の行を処理
	for rows.Next() {
		var alb Album
		// 各行のデータを構造体のフィールドにマッピング
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	// クエリ中のエラーをチェック
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumByID(id int64) (Album, error) {
    // 取得結果を格納するAlbum構造体の変数を宣言
    var alb Album

    // 単一行をクエリする
    row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)

    // 取得した行のデータを構造体のフィールドにスキャンする
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumByID %d: no such album", id)
        }
        return alb, fmt.Errorf("albumByID %d: %v", id, err)
    }

    return alb, nil
}

func albumByTitle(title string) (Album, error) {
	var tialb Album

	row := db.QueryRow("SELECT * FROM album WHERE title = $1", title)

	if err := row.Scan(&tialb.ID, &tialb.Title, &tialb.Artist, &tialb.Price); err != nil {
        if err == sql.ErrNoRows {
            return tialb, fmt.Errorf("albumByTitle %q: no such album", title)
        }
		return tialb, fmt.Errorf("albumByTitle %q: %v", title, err)
    }

	return tialb, nil
}

// addAlbum は指定されたアルバムをデータベースに追加し、
// 新しいエントリのアルバムIDを返します
func addAlbum(alb Album) (int64, error) {
    // RETURNING id を使って挿入後にIDを返す
    var id int64
    err := db.QueryRow(
        "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
        alb.Title, alb.Artist, alb.Price).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("addAlbum: %v", err)
    }
    return id, nil
}

