# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

ShellPetは、コマンド履歴を食べて育つターミナルベースのペット育成ゲームです。git、docker、vimなどのシェルコマンドを解析し、ペットのパラメータや見た目を進化させます。CC BY-NC-SA 4.0ライセンス下の個人・学習用プロジェクトです（商用利用禁止）。

**技術スタック:**
- Go 1.25.0
- Bubble Tea（TUIフレームワーク - 追加予定）
- Lip Gloss（スタイリング - 追加予定）

## 開発コマンド

### 基本的なGoコマンド
```bash
# 依存関係の初期化・更新
go mod tidy

# プロジェクトのビルド
go build -o shellpet .

# アプリケーションの実行
go run .

# テストの実行
go test ./...

# カバレッジ付きでテスト実行
go test -cover ./...

# 特定パッケージのテスト実行
go test ./pkg/packagename

# 特定のテストのみ実行
go test -run TestName ./...

# コードフォーマット
go fmt ./...

# コードの潜在的な問題をチェック
go vet ./...
```