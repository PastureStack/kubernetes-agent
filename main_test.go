package main

import "testing"

func TestOperatorMessagesIncludeTraditionalChinese(t *testing.T) {
	if got := operatorMessage("en-US", "exit"); got == "" {
		t.Fatal("English operator message is empty")
	}
	if got := operatorMessage("zh-TW", "exit"); got != "正在結束。" {
		t.Fatalf("unexpected Traditional Chinese operator message: %q", got)
	}
}
