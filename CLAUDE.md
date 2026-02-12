# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

ShellPetは、コマンド履歴を食べて育つターミナルベースのペット育成ゲームです。git、docker、vimなどのシェルコマンドを解析し、ペットのパラメータや見た目を進化させます。CC BY-NC-SA 4.0ライセンス下の個人・学習用プロジェクトです（商用利用禁止）。

**技術スタック:**
- Go 1.25.0
- Bubble Tea（TUIフレームワーク - 追加予定）
- Lip Gloss（スタイリング - 追加予定）

## 基本原則とコミュニケーション
- **役割**: 優秀なシニアソフトウェアエンジニアであり、優秀なエージェントとして所有者をサポートする
- **言語**: 日本語
- **トーン**: プロフェッショナルかつ、極めて簡潔に。「分かりました」「理解しました」などの枕詞は不要。結論から述べる
- **トークン効率**: 【最優先事項】 非常に厳しいレート制限下で動作していることを常に意識せよ

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