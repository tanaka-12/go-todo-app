package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type Task struct {
	Title     string
	Completed bool
	Deadline  string
}

func main() {
	//１．箱（スライス）を作る
	//ここはループの外
	//タスクを保存するリスト
	tasks := []Task{}

	//２．読み込み準備
	//ここもループの外で1回やればおっけー
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("TODOアプリを開始します（exitと入力すると終了） ")

	//３．無限ループ開始
	for{
		fmt.Print("\nタスクを入力 > ")

		//４．入力を受け取る&お掃除
		input, _ := reader.ReadString('\n')
		cleanTitle := strings.TrimSpace(input)

		//５．脱出チェック
		if cleanTitle == "exit" {
			fmt.Println("アプリを終了します・・・")
			break
		}

		if cleanTitle == "" {
			continue
		}

		//６．期限を聞く
		fmt.Print("期限を入力 > ")
		dateInput, _ := reader.ReadString('\n')
		cleanDeadline := strings.TrimSpace(dateInput)

		//７．Deadlineにデータを入力
		newTask := Task{
			Title: cleanTitle,
			Completed: false, Deadline: cleanDeadline,
		}

		//８．リストに追加
		tasks = append(tasks, newTask)

		//９．現在のリストを表示
		fmt.Println("=== 現在のタスク ===")
		for i, t := range tasks {
			fmt.Printf("%d: [ ] %s (期限: %s)\n", i, t.Title, t.Deadline)
		}
		fmt.Println("==================")
	}
}






	
