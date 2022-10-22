package render

import (
	"math"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_approxDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "less than 1 second",
			args: args{
				d: time.Millisecond * 100,
			},
			want: "less than 1 second",
		},
		{
			name: "1 second",
			args: args{
				d: time.Second * 1,
			},
			want: "1 second",
		},
		{
			name: "more than 1 second",
			args: args{
				d: time.Second * 2,
			},
			want: "2 seconds",
		},
		{
			name: "1 minute",
			args: args{
				d: time.Minute * 1,
			},
			want: "1 minute",
		},
		{
			name: "more than minute",
			args: args{
				d: time.Minute * 2,
			},
			want: "2 minutes",
		},
		{
			name: "1 hour",
			args: args{
				d: time.Hour * 1,
			},
			want: "1 hour",
		},
		{
			name: "more than hour",
			args: args{
				d: time.Hour * 2,
			},
			want: "2 hours",
		},
		{
			name: "1 day",
			args: args{
				d: day,
			},
			want: "1 day",
		},
		{
			name: "more than 1 day",
			args: args{
				d: day * 2,
			},
			want: "2 days",
		},
		{
			name: "1 year",
			args: args{
				d: year,
			},
			want: "1 year",
		},
		{
			name: "more than 1 year",
			args: args{
				d: year * 2,
			},
			want: "2 years",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := approxDuration(tt.args.d); got != tt.want {
				t.Errorf("approxDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decr(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				i: 10,
			},
			want:    9,
			wantErr: false,
		},
		{
			name: "argument as string number",
			args: args{
				i: "10",
			},
			want:    9,
			wantErr: false,
		},
		{
			name: "illegal argument return 0",
			args: args{
				i: "10a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decr(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("decr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatFloat(t *testing.T) {
	type args struct {
		f  float64
		dp int
	}
	test := struct {
		name string
		args args
		want string
	}{
		name: "happy path",
		args: args{
			f:  math.Pi,
			dp: 2,
		},
		want: "3.14",
	}

	if got := formatFloat(test.args.f, test.args.dp); got != test.want {
		t.Errorf("formatFloat() = %v, want %v", got, test.want)
	}
}

func Test_formatInt(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				i: 10,
			},
			want:    "10",
			wantErr: false,
		},
		{
			name: "wrong input value",
			args: args{
				i: "10a",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatInt(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatTime(t *testing.T) {
	type args struct {
		format string
		t      time.Time
	}
	test := struct {
		name string
		args args
		want string
	}{

		name: "happy path",
		args: args{
			format: "Jan 2, 2006",
			t:      time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC),
		},
		want: "Oct 10, 2022",
	}

	if got := formatTime(test.args.format, test.args.t); got != test.want {
		t.Errorf("formatTime() = %v, want %v", got, test.want)
	}
}

func Test_incr(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				i: 10,
			},
			want:    11,
			wantErr: false,
		},
		{
			name: "argument as string number",
			args: args{
				i: "10",
			},
			want:    11,
			wantErr: false,
		},
		{
			name: "illegal argument return 0",
			args: args{
				i: "11a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := incr(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("incr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("incr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pluralize(t *testing.T) {
	type args struct {
		count    interface{}
		singular string
		plural   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				count:    1,
				singular: "computer",
				plural:   "computers",
			},
			want: "computer",
		},
		{
			name: "happy path plural",
			args: args{
				count:    2,
				singular: "computer",
				plural:   "computers",
			},
			want: "computers",
		},
		{
			name: "happy path",
			args: args{
				count:    "10a",
				singular: "computer",
				plural:   "computers",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pluralize(tt.args.count, tt.args.singular, tt.args.plural)
			if (err != nil) != tt.wantErr {
				t.Errorf("pluralize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pluralize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_slugify(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slugify(tt.args.s); got != tt.want {
				t.Errorf("slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toInt64(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toInt64(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("toInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlDelParam(t *testing.T) {
	type args struct {
		u   *url.URL
		key string
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := urlDelParam(tt.args.u, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("urlDelParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlSetParam(t *testing.T) {
	type args struct {
		u     *url.URL
		key   string
		value interface{}
	}
	test := struct {
		name string
		args args
		want *url.URL
	}{

		name: "happy path",
		args: args{
			u: &url.URL{
				Scheme: "http",
				Host:   "localhost",
			},
			key:   "page",
			value: 1,
		},
		want: &url.URL{
			Scheme:   "http",
			Host:     "localhost",
			RawQuery: "page=1",
		},
	}

	if got := urlSetParam(test.args.u, test.args.key, test.args.value); !reflect.DeepEqual(got, test.want) {
		t.Errorf("urlSetParam() = %v, want %v", got, test.want)
	}

}

func Test_yesno(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test yes",
			args: args{
				b: true,
			},
			want: "Yes",
		},
		{
			name: "test no",
			args: args{
				b: false,
			},
			want: "No",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yesno(tt.args.b); got != tt.want {
				t.Errorf("yesno() = %v, want %v", got, tt.want)
			}
		})
	}
}
