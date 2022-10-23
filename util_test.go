package render

import "testing"

func Test_max(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "happy path",
			args: args{
				x: 10,
				y: 20,
			},
			want: 20,
		},
		{
			name: "happy path",
			args: args{
				x: 10,
				y: 11,
			},
			want: 11,
		},
		{
			name: "happy path",
			args: args{
				x: 54,
				y: 60,
			},
			want: 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "happy path",
			args: args{
				x: 10,
				y: 20,
			},
			want: 10,
		},
		{
			name: "happy path swap values",
			args: args{
				x: 20,
				y: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_totalPages(t *testing.T) {
	type args struct {
		size  int
		total int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "happy path",
			args: args{
				size:  25,
				total: 100,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := totalPages(tt.args.size, tt.args.total); got != tt.want {
				t.Errorf("totalPages() = %v, want %v", got, tt.want)
			}
		})
	}
}
