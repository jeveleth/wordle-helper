package utils

import (
	"reflect"
	"strings"
	"testing"
)

func TestSetIncludeLetters(t *testing.T) {
	type args struct {
		letters string
		wrds    []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test1", args{"ar", []string{"foo", "bar"}}, []string{"bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetIncludeLetters(tt.args.letters, tt.args.wrds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIncludeLetters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAntiPattern(t *testing.T) {
	type args struct {
		antiPattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"-n-i--"}, "([a-z])([^n])([a-z])([^i])([a-z])([a-z])"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseAntiPattern(tt.args.antiPattern); got != tt.want {
				t.Errorf("ParseAntiPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{"-n-i--"}, "([a-z])[n]([a-z])[i]([a-z])([a-z])"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePattern(tt.args.pattern); got != tt.want {
				t.Errorf("ParsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAntiPattern(t *testing.T) {
	type args struct {
		pattern string
		wrds    []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test1", args{"([a-z])([^n])([a-z])([^i])([a-z])", []string{"afoot", "bar"}}, []string{"afoot"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetAntiPattern(tt.args.pattern, tt.args.wrds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetAntiPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetExcludedLetters(t *testing.T) {
	type args struct {
		exclude string
		wrds    []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test1", args{"abcd", []string{"foo", "bar"}}, []string{"foo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetExcludedLetters(tt.args.exclude, tt.args.wrds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetExcludedLetters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetComboPattern(t *testing.T) {
	type args struct {
		pattern string
		words   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"test1", args{"([a-z])([^n])([a-z])([^i])([a-z])", []string{"afoot", "bar"}}, []string{"afoot"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetComboPattern(tt.args.pattern, tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetComboPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCommonLettersCount(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{"test1", args{
			[]string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"}},
			map[string]int{"a": 3, "b": 2, "e": 4, "f": 3, "h": 1, "j": 1, "l": 1, "n": 1, "o": 3, "r": 3, "t": 2, "u": 1, "z": 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCommonLettersCount(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommonLettersCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport(t *testing.T) {
	type args struct {
		wrds []string
		c    Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "passnoconfig",
			args: args{
				[]string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"},
				Config{
					Include:     "",
					Exclude:     "",
					Pattern:     "",
					AntiPattern: "",
				},
			},
			want: "7 words found",
		},
		{
			name: "passincludeonly",
			args: args{
				[]string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"},
				Config{
					Include:     "ab",
					Exclude:     "",
					Pattern:     "",
					AntiPattern: "",
				},
			},
			want: "2 words found that include 'ab'. The words are:\nzebra\nbar\n\nThe most common letters are:\n\na: 2 times\nb: 2 times\nr: 2 times\ne: 1 times\nz: 1 times",
		},
		{
			name: "passexcludeonly",
			args: args{
				[]string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"},
				Config{
					Include:     "",
					Exclude:     "ab",
					Pattern:     "",
					AntiPattern: "",
				},
			},
			want: "4 words found that exclude 'ab'.The words are:\nfoot\nhotel\nfun\nfree\n\nThe most common letters are:\n\ne: 3 times\nf: 3 times\no: 3 times\nt: 2 times\nh: 1 times\nl: 1 times\nn: 1 times\nr: 1 times\nu: 1 times",
		},
		{
			name: "passexcludeandinclude",
			args: args{
				[]string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"},
				Config{
					Include:     "u",
					Exclude:     "a",
					Pattern:     "",
					AntiPattern: "",
				},
			},
			want: "1 words found that exclude 'a' and include 'u'. The words are:\nfun\n\nThe most common letters are:\n\nf: 1 times\nn: 1 times\nu: 1 times",
		},
		{
			name: "passexcludeandincludeandantipattern",
			args: args{
				[]string{"zebra", "jazzy", "hotel", "armor", "aroma", "carom"},
				Config{
					Include:     "z",
					Exclude:     "h",
					Pattern:     "",
					AntiPattern: "-e---",
				},
			},
			want: "1 words found that exclude 'h' and include 'z'. The words are:\njazzy\n\nThe most common letters are:\n\nz: 2 times\na: 1 times\nj: 1 times\ny: 1 times",
		},
		{
			name: "passexcludeandincludeandantipatternandpattern",
			args: args{
				[]string{"zebra", "jazzy", "joker", "hotel", "armor", "aroma", "carom"},
				Config{
					Include:     "z",
					Exclude:     "h",
					Pattern:     "j----",
					AntiPattern: "-e-e-",
				},
			},
			want: "1 words found that exclude 'h', include 'z', and match the pattern 'j----'. The words are:\njazzy\n\nThe most common letters are:\n\nz: 2 times\na: 1 times\nj: 1 times\ny: 1 times",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Report(tt.args.wrds, tt.args.c); got != tt.want {
				t.Errorf("Report() = %v, want %v", got, tt.want)
			}
		})
	}
}

func IncludeOnly(t *testing.T) {
	var testB strings.Builder
	var testWrds = []string{"zebra", "bar", "foot", "hotel", "fun", "free", "jazz"}
	var include = "ab"
	type args struct {
		b    strings.Builder
		wrds []string
		c    Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"passincludeonly",
			args{testB, testWrds, Config{Include: include}},
			"2 words found that include 'ab'. The words are:\nzebra\nbar\n\nThe most common letters are:\n\nb: 2 times\nr: 2 times\na: 2 times\nz: 1 times\ne: 1 times"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := includeOnly(tt.args.b, tt.args.wrds, tt.args.c); got != tt.want {
				t.Errorf("includeOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortRankedPairsAlaphabetically(t *testing.T) {
	type args struct {
		pl PairList
	}
	tests := []struct {
		name string
		args args
		want PairList
	}{
		{"pass", args{PairList{{"h", 3}, {"c", 3}, {"y", 2}, {"b", 2}, {"a", 2}}}, PairList{{"c", 3}, {"h", 3}, {"a", 2}, {"b", 2}, {"y", 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortRankedPairsAlaphabetically(tt.args.pl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortRankedPairsAlaphabetically() = %v, want %v", got, tt.want)
			}
		})
	}
}
