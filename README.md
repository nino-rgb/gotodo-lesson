# gotodo-lesson
参考記事[https://qiita.com/hrk_ym/items/c73c5ad41c92688c3b94]

- リポジトリ層の引数
1件取得 (GetTodoByID)

表現方法: nil（ポインタ）

意味: 「対象のデータが存在しない」

一覧取得 (GetTodos)

表現方法: []（空スライス）

意味: 「条件に該当するデータが0件」

## test

### service層 
- GetTodos 正常系    
- GetTodoByID 正常系  
- GetTodoByID エラー系
- CreateTodo 正常系   
- UpdateTodo 正常系   
- UpdateTodo エラー系 
- DeleteTodo 正常系   