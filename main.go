package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Title     string
	Completed bool
	Deadline  string
	Priority  int
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

		//バリデーションエラー
		if cleanTitle == "" {
			fmt.Println("⚠️エラー：タスク名が空です！文字を入力してください！")
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
			//表示する前に「優先度が高い順」に並び替える
			sort.Slice(tasks, func(i, j int) bool {
				//i番目とj番目を比較してiのほうが大きければ「iを前にして」という意味
				return tasks[i].Priority > tasks[j].Priority
			})

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
				//優先度表示
				stars := "⭐" //低
				if t.Priority == 2 {
					stars = "⭐⭐" //中
				} else if t.Priority == 3 {
					stars = "⭐⭐⭐" //高
				}

				//期限切れチェック
				deadlineTime, err := time.Parse("2006-01-02", t.Deadline)
				expiredTag := ""
				if err == nil && deadlineTime.Before(now) && !t.Completed {
					expiredTag = "⚠️　期限切れ！"
				}

				fmt.Printf("%d: %s %-15s %-10s (期限: %s)%s\n",
					i, mark, t.Title, "【"+stars+"】", t.Deadline, expiredTag)

			}
			fmt.Println("==================")
			continue
		}

		//６．期限を聞く
		fmt.Print("期限を入力（例：2026-01-01) > ")
		dateInput, _ := reader.ReadString('\n')
		cleanDeadline := strings.TrimSpace(dateInput)

		//日付のチェック機能
		//形チェック：「2026-01-01」の形になっているか
		deadlineTime, err := time.Parse("2006-01-02", cleanDeadline)

		if err != nil {
			//解析に失敗したらエラー（変な文字や存在しない数字）
			fmt.Println("⚠️　エラー：日付は'2026-01-01'の形で入力してください！")
			continue
		}

		//過去チェック
		//time.Now()で今の時間を取る
		//truncate(24*time.Hour)は「今日の0時0分」に合わせるおまじない
		now := time.Now().Truncate(24 * time.Hour)

		if deadlineTime.Before(now) {
			fmt.Println("⚠️　エラー：過去の日付は入力できません！未来に向かって生きよう！")
			continue
		}

		//優先度を聞く
		fmt.Print("優先度を入力（3:高, 2:中, 1:低）> ")
		priorityInput, _ := reader.ReadString('\n')
		cleanPriority := strings.TrimSpace(priorityInput)

		//文字を数字に変換(Atoi)
		priority, err := strconv.Atoi(cleanPriority)

		//エラーチェック(数字じゃない、もしくは1～3以外)
		if err != nil || priority < 1 || priority > 3 {
			fmt.Println("⚠️　エラー：優先度は1,2,3の数字で入力してください！")
			continue
		}

		//７．Deadlineにデータを入力
		newTask := Task{
			Title:     cleanTitle,
			Completed: false,
			Deadline:  cleanDeadline,
			Priority:  priority,
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
