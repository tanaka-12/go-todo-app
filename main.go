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
	//タスクを保存するリスト
	tasks := []Task{}

	fmt.Println("タスクを入力してください: ")

	//入力の準備
	reader := bufio.NewReader(os.Stdin)

	//入力待ち（エンターキーが押されるまで）
	input, _ := reader.ReadString('\n')

	//お掃除
	cleanTitle := strings.TrimSpace(input)

	//タスクの作成
	newTask := Task{
		Title: cleanTitle,
		Completed: false,
		Deadline: "",
	}

	//リストの追加
	tasks = append(tasks, newTask)

	//結果発表
	fmt.Printf("現在のリスト: %+v\n", tasks)
}






	
