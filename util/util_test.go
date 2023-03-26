package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Recover(t *testing.T) {
	defer Recover("test")
	panic("test")
}

func Test_GetGeoDistance(t *testing.T) {
	// https://jingwei.supfree.net/
	pm := map[string][]float64{
		"cn-beijing":  []float64{39.90, 116.40},
		"cn-tianjin":  []float64{39.12, 117.20},
		"cn-shenzhen": []float64{22.55, 114.05},
	}
	da := GetGeoDistance(pm["cn-beijing"][0], pm["cn-beijing"][1], pm["cn-tianjin"][0], pm["cn-tianjin"][1])
	t.Logf("beijing --> tianjin: %v", da)
	db := GetGeoDistance(pm["cn-tianjin"][0], pm["cn-tianjin"][1], pm["cn-beijing"][0], pm["cn-beijing"][1])
	t.Logf("tianjin --> beijing: %v", db)
	assert.Equal(t, da, db)

}
