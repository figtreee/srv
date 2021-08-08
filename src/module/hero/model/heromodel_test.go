package heromodel

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDemoHero(t *testing.T) {
	heroArr := make([]Hero, 0)

	for i := 0; i < 10; i++ {
		heroArr = append(heroArr, Hero{
			ID:   i,
			Name: fmt.Sprintf("hero %d", i),
		})
	}

	// Encode Go object to JSON string
	var err error
	var buf []byte

	//if buf, err = json.Marshal(heroArr); err != nil {
	if buf, err = json.MarshalIndent(heroArr, "", "  "); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("buf: %s\n", buf)
}

func TestReturnFirstThreeCharacters(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test case 1",
			args: args{
				name: "abcdefg",
			},
			want: "abc",
		},
		{
			name: "test case 2",
			args: args{
				name: "aaaaaa",
			},
			want: "aaa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReturnFirstThreeCharacters(tt.args.name); got != tt.want {
				t.Errorf("ReturnFirstThreeCharacters() = %v, want %v", got, tt.want)
			}
		})
	}
}
