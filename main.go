package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

	//起動時にファイルを読み込む
	//ファイルがあるか確認して読み込む
	bytes, err := os.ReadFile("tasks.json")
	if err == nil {
		//ファイルがあったらデータを戻す
		json.Unmarshal(bytes, &tasks)
		fmt.Println("📂 データを読み込みました！")
	}

	//２．読み込み準備
	//ここもループの外で1回やればおっけー
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("TODOアプリを開始します（exitと入力すると終了） ")

	//３．無限ループ開始
	for {
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

		//入力された文字をスペースで分割する
		parts := strings.Split(cleanTitle, " ")

		//「最初の単語」がdoneかどうかチェック
		if parts[0] == "done" {

			//番号が入力されているかチェック
			if len(parts) < 2 {
				fmt.Println("エラー : 番号を入力してください（例： done 0）")
				continue
			}

			//文字を数字に変換
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("正しい数字を入力してください")
				continue
			}

			//タスクの番号が存在するかチェック
			if index < 0 || index >= len(tasks) {
				fmt.Println("エラー：その番号のタスクはありません")
				continue
			}

			//ここに完了処理を書く
			tasks[index].Completed = true
			fmt.Println("🎉タスクを完了にしました！")

			//ここでcontinueすると期限入力をスキップして、次のループに戻る
			continue
		}

		//delete doneコマンド(一括削除)
		//もし「delete」かつ「done」と入力されたら・・・
		if parts[0] == "delete" && len(parts) > 1 && parts[1] == "done" {

			//引っ越し先の新しい箱(スライス)を用意
			newTasks := []Task{}

			//古いタスクをチェックする
			for _, t := range tasks {
				//完了してない(!t.Completed)のやつだけ選ぶ！
				if !t.Completed {
					//新しい箱に入れる(append)
					newTasks = append(newTasks, t)
				}
			}

			//古い箱を捨てて新しい箱に置き換える
			tasks = newTasks
			fmt.Println("🗑️ 完了済みのタスクを削除しました！スッキリ！")
			continue
		}

		if parts[0] == "delete" {

			//番号があるかチェック
			if len(parts) < 2 {
				fmt.Println("エラー : 削除する番号を入れてください（例： delete 0）")
				continue
			}

			//文字を数字に変換
			index, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("エラー：正しい数字を入力してください")
				continue
			}

			//タスクの番号が存在するかチェック
			if index < 0 || index >= len(tasks) {
				fmt.Println("エラー：その番号のタスクはありません")
				continue
			}

			//削除実行
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Println("🗑️ タスクを削除しました！")
			continue
		}

		//saveコマンド
		if parts[0] == "save" {
			//Goのデータ(tasks)をJSON(bytes)に変換
			bytes, err := json.Marshal(tasks)
			if err != nil {
				fmt.Println("変換に失敗しました...", err)
				continue
			}

			//成功したらファイルに書き込む
			//"tasks.json"はファイル名
			//bytesは書き込みデータ
			//0644は自分は読み書きおっけー、他人は見るだけ
			err = os.WriteFile("tasks.json", bytes, 0644)

			if err != nil {
				fmt.Println("保存に失敗しました...", err)
			} else {
				fmt.Println("💾 タスクを 'tasks.json' に保存しました！")
			}
			continue
		}

		//listを追加
		if parts[0] == "list" {
			fmt.Println("=== 現在のタスク ===")

			//現在の時間を取得(ループ外で1回だけやること)
			now := time.Now()

			//９からコピーしてくる
			for i, t := range tasks {
				//ここでマークを決める
				mark := "[]"
				if t.Completed == true {
					mark = "[x]"
				}

				//期限切れチェック！
				//t.Deadline(文字)を時間データに変換してみる
				deadlineTime, err := time.Parse("2006-01-02", t.Deadline)

				//変換が成功(err == nil)して、かつ
				//期限が過ぎていて(Before)、まだ完了していなければ(!t.Completed)
				if err == nil && deadlineTime.Before(now) && !t.Completed {
					//赤文字っぽく目立たせる(⚠️マーク)
					fmt.Printf("%d: %s %s (期限: %s) ⚠️ 期限切れ！\n", i, mark, t.Title, t.Deadline)
				} else {
					//通常表示
					fmt.Printf("%d: %s %s (期限: %s)\n", i, mark, t.Title, t.Deadline)
				}
			}
			fmt.Println("==================")
			continue
		}

		//６．期限を聞く
		fmt.Print("期限を入力 > ")
		dateInput, _ := reader.ReadString('\n')
		cleanDeadline := strings.TrimSpace(dateInput)

		//７．Deadlineにデータを入力
		newTask := Task{
			Title:     cleanTitle,
			Completed: false, Deadline: cleanDeadline,
		}

		//８．リストに追加
		tasks = append(tasks, newTask)

		//９．現在のリストを表示
		fmt.Println("=== 現在のタスク ===")
		for i, t := range tasks {
			//ここでマークを決める
			mark := "[]"
			if t.Completed == true {
				mark = "[x]"
			}

			//[]の代わりにmark変数を使う
			fmt.Printf("%d: %s %s (期限: %s)\n", i, mark, t.Title, t.Deadline)
		}
		fmt.Println("==================")
	}
}
