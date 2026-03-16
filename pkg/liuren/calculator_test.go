package liuren

import (
	"testing"
)

func TestParsePillar(t *testing.T) {
	tests := []struct {
		input    string
		expected Sexagenary
		wantErr  bool
	}{
		{"甲子", Sexagenary{Stem: StemJia, Branch: Zi}, false},
		{"乙丑", Sexagenary{Stem: StemYi, Branch: Chou}, false},
		{"丙寅", Sexagenary{Stem: StemBing, Branch: Yin}, false},
		{"丁卯", Sexagenary{Stem: StemDing, Branch: Mao}, false},
		{"戊辰", Sexagenary{Stem: StemWu, Branch: Chen}, false},
		{"己巳", Sexagenary{Stem: StemJi, Branch: Si}, false},
		{"庚午", Sexagenary{Stem: StemGeng, Branch: Wu}, false},
		{"辛未", Sexagenary{Stem: StemXin, Branch: Wei}, false},
		{"壬申", Sexagenary{Stem: StemRen, Branch: Shen}, false},
		{"癸酉", Sexagenary{Stem: StemGui, Branch: You}, false},
		{"甲戌", Sexagenary{Stem: StemJia, Branch: Xu}, false},
		{"乙亥", Sexagenary{Stem: StemYi, Branch: Hai}, false},
		{"己丑", Sexagenary{Stem: StemJi, Branch: Chou}, false}, // 這個是問題案例
		{"", Sexagenary{}, true},
		{"甲", Sexagenary{}, true},
		{"invalid", Sexagenary{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePillar(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePillar(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.Stem != tt.expected.Stem || got.Branch != tt.expected.Branch) {
				t.Errorf("ParsePillar(%q) = {Stem:%v, Branch:%v}, want {Stem:%v, Branch:%v}",
					tt.input, got.Stem, got.Branch, tt.expected.Stem, tt.expected.Branch)
			}
		})
	}
}

func TestStemNames(t *testing.T) {
	t.Logf("StemNames: %v", StemNames)
	for i, name := range StemNames {
		t.Logf("Stem %d: %s (rune: %v)", i, name, []rune(name))
	}
}

func TestBranchNames(t *testing.T) {
	t.Logf("BranchNames: %v", BranchNames)
	for i, name := range BranchNames {
		t.Logf("Branch %d: %s (rune: %v)", i, name, []rune(name))
	}
}
