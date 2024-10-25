package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/venil7/func/task"
)

func getItem(num int, db *SQLiteRepository) task.Task[int] {
	client := http.Client{}

	body := task.Task[io.ReadCloser](func() (io.ReadCloser, error) {
		// fmt.Printf("Fetching %d\n", num)
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", num)
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	})

	item := task.Then(body, func(body io.ReadCloser) (Item, error) {
		defer body.Close()
		decoder := json.NewDecoder(body)
		var item Item
		err := decoder.Decode(&item)
		return item, err
	})

	addRecord := task.From1(db.Add)

	item = task.Tap(item, addRecord)

	kids := task.FlatMap(item, func(item Item) task.Task[[]int] {
		return task.Traverse(item.Kids, func(id int) task.Task[int] {
			return getItem(id, db)
		})
	})

	return task.Then(kids, func(ids []int) (int, error) {
		return len(ids), nil
	})
}

func _main() {
	num := task.From1(strconv.Atoi)(os.Args[1])

	db := New("posts.db")
	migrate := task.From(db.Migrate)

	app := task.FlatMap(
		migrate,
		func(res sql.Result) task.Task[int] {
			return task.FlatMap(num, func(num int) task.Task[int] {
				return getItem(num, db)
			})
		})

	_, err := app()

	if err != nil {
		panic(err)
	}
}

func main() {
	db := New("posts.db")
	migrate := task.From(db.Migrate)

	nums := task.Traverse(os.Args[1:], task.From1(strconv.Atoi))
	items := task.FlatMap(nums, func(nums []int) task.Task[[]int] {
		return task.Traverse(nums, func(num int) task.Task[int] {
			return getItem(num, db)
		})
	})
	app := task.FlatMap(migrate, func(r sql.Result) task.Task[[]int] {
		return items
	})

	_, err := app()

	if err != nil {
		panic(err)
	}
}

/*
copy([...new Set($$('a').map(x => x.href).filter(x => x.includes("item?id=")).map(x => x.match(/id=(\d+)/)[1]).map(Number))])
*/
