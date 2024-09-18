package database

import "testing"

func TestVideosTable_getSeconds(t *testing.T) {
	type fields struct {
		Table *Table
	}
	type args struct {
		duration string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"", fields{nil}, args{""}, 0},
		{"LIVE", fields{nil}, args{"LIVE"}, 0},
		{"ERROR", fields{nil}, args{"ERROR"}, 0},
		{"MEMBER", fields{nil}, args{"MEMBER"}, 0},
		{"112:34:56", fields{nil}, args{"112:34:56"}, 405296},
		{"12:34:56", fields{nil}, args{"12:34:56"}, 45296},
		{"2:34:56", fields{nil}, args{"2:34:56"}, 9296},
		{"13:28", fields{nil}, args{"13:28"}, 808},
		{"3:27", fields{nil}, args{"3:27"}, 207},
		{"45", fields{nil}, args{"45"}, 45},
		{"8", fields{nil}, args{"8"}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VideosTable{
				Table: tt.fields.Table,
			}
			if got := v.getSeconds(tt.args.duration); got != tt.want {
				t.Errorf("getSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
