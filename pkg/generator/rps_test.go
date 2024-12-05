package generator

import (
	"github.com/vhive-serverless/loader/pkg/common"
	"math"
	"testing"
)

func TestWarmStartMatrix(t *testing.T) {
	tests := []struct {
		testName               string
		experimentDuration     int
		rpsTarget              float64
		expectedIAT            common.IATArray
		expectedPerMinuteCount []int
	}{
		{
			testName:               "2min_0rps",
			experimentDuration:     2,
			rpsTarget:              0,
			expectedIAT:            []float64{},
			expectedPerMinuteCount: []int{0, 0},
		},
		{
			testName:           "2min_1rps",
			experimentDuration: 2,
			rpsTarget:          1,
			expectedIAT: []float64{
				// minute 1
				0, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				// minute 2
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
				1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000, 1_000_000,
			},
			expectedPerMinuteCount: []int{60, 60},
		},
		{
			testName:           "2min_0.5rps",
			experimentDuration: 2,
			rpsTarget:          0.5,
			expectedIAT: []float64{
				// minute 1
				0, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				// minute 2
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
				2_000_000, 2_000_000, 2_000_000, 2_000_000, 2_000_000,
			},
			expectedPerMinuteCount: []int{30, 30},
		},
		{
			testName:           "2min_0.125rps",
			experimentDuration: 2,
			rpsTarget:          0.125,
			expectedIAT: []float64{
				// minute 1
				0, 8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000,
				// minute 2
				8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000, 8_000_000,
			},
			expectedPerMinuteCount: []int{8, 7},
		},
		{
			testName:           "6min_0.01rps",
			experimentDuration: 6,
			rpsTarget:          0.01,
			expectedIAT: []float64{
				0, 100_000_000, 100_000_000, 100_000_000,
			},
			expectedPerMinuteCount: []int{1, 1, 0, 1, 0, 1},
		},
	}

	epsilon := 0.01

	for _, test := range tests {
		t.Run("warm_start_"+test.testName, func(t *testing.T) {
			matrix, minuteCount := GenerateWarmStartFunction(test.experimentDuration, test.rpsTarget)

			if len(matrix) != len(test.expectedIAT) {
				t.Errorf("Unexpected IAT array size - got: %d, expected: %d", len(matrix), len(test.expectedIAT))
			}
			if len(minuteCount) != len(test.expectedPerMinuteCount) {
				t.Errorf("Unexpected count array size - got: %d, expected: %d", len(minuteCount), len(test.expectedPerMinuteCount))
			}

			sum := 0.0
			count := 0

			for i := 0; i < len(matrix); i++ {
				if math.Abs(matrix[i]-test.expectedIAT[i]) > epsilon {
					t.Error("Unexpected IAT value.")
				}

				sum += matrix[i]
				count++
			}

			for i := 0; i < len(minuteCount); i++ {
				if test.expectedPerMinuteCount[i] != minuteCount[i] {
					t.Error("Unexpected per minute count.")
				}
			}
		})
	}
}

func TestColdStartMatrix(t *testing.T) {
	tests := []struct {
		testName           string
		experimentDuration int
		rpsTarget          float64
		cooldownSeconds    int
		expectedIAT        []common.IATArray
		expectedCount      [][]int
	}{
		{
			testName:           "2min_0rps",
			experimentDuration: 2,
			rpsTarget:          0,
			cooldownSeconds:    10,
			expectedIAT:        []common.IATArray{},
			expectedCount:      [][]int{}, // empty since no functions
		},
		{
			testName:           "2min_1rps",
			experimentDuration: 2,
			rpsTarget:          1,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{0_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{1_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{2_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{3_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{4_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{5_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{6_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{7_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{8_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{9_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
			},
			expectedCount: [][]int{
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
				{6, 6},
			},
		},
		{
			testName:           "1min_0.25rps",
			experimentDuration: 1,
			rpsTarget:          0.25,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{0_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{4_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{8_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
			},
			expectedCount: [][]int{
				{5},
				{5},
				{5},
			},
		},
		{
			testName:           "2min_0.25rps",
			experimentDuration: 2,
			rpsTarget:          0.25,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{0_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{4_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{8_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
			},
			expectedCount: [][]int{
				{5, 5},
				{5, 5},
				{5, 5},
			},
		},
		{
			testName:           "1min_0.33rps",
			experimentDuration: 1,
			rpsTarget:          1.0 / 3,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{0_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{3_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{6_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
				{9_000_000, 12_000_000, 12_000_000, 12_000_000, 12_000_000},
			},
			expectedCount: [][]int{
				{5},
				{5},
				{5},
				{5},
			},
		},
		{
			testName:           "1min_5rps",
			experimentDuration: 1,
			rpsTarget:          5,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{1_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{1_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{1_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{1_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{1_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{2_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{2_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{2_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{2_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{2_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{3_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{3_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{3_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{3_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{3_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{4_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{4_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{4_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{4_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{4_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{5_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{5_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{5_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{5_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{5_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{6_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{6_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{6_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{6_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{6_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{7_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{7_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{7_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{7_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{7_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{8_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{8_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{8_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{8_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{8_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},

				{9_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{9_200_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{9_400_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{9_600_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
				{9_800_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000, 10_000_000},
			},
			expectedCount: [][]int{
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},

				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
				{6}, {6}, {6}, {6}, {6},
			},
		},
		{
			testName:           "1min_5rps_cooldown5s",
			experimentDuration: 1,
			rpsTarget:          5,
			cooldownSeconds:    5,
			expectedIAT: []common.IATArray{
				{000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{200_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{400_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{600_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{800_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},

				{1_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{1_200_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{1_400_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{1_600_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{1_800_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},

				{2_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{2_200_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{2_400_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{2_600_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{2_800_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},

				{3_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{3_200_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{3_400_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{3_600_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{3_800_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},

				{4_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{4_200_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{4_400_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{4_600_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
				{4_800_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000, 5_000_000},
			},
			expectedCount: [][]int{
				{12}, {12}, {12}, {12}, {12},
				{12}, {12}, {12}, {12}, {12},
				{12}, {12}, {12}, {12}, {12},
				{12}, {12}, {12}, {12}, {12},
				{12}, {12}, {12}, {12}, {12},
			},
		},
		{
			testName:           "6min_0.01rps_10s_cooldown",
			experimentDuration: 6,
			rpsTarget:          0.01,
			cooldownSeconds:    10,
			expectedIAT: []common.IATArray{
				{0, 100_000_000, 100_000_000, 100_000_000},
			},
			expectedCount: [][]int{
				{1, 1, 0, 1, 0, 1},
			},
		},
		{
			testName:           "6min_0.01rps_120s_cooldown",
			experimentDuration: 6,
			rpsTarget:          0.01,
			cooldownSeconds:    120,
			expectedIAT: []common.IATArray{
				{0, 200_000_000},
				{100_000_000, 200_000_000},
			},
			expectedCount: [][]int{
				{1, 0, 0, 1, 0, 0},
				{0, 1, 0, 0, 0, 1},
			},
		},
	}

	epsilon := 0.01

	for _, test := range tests {
		t.Run("cold_start_"+test.testName, func(t *testing.T) {
			matrix, minuteCounts := GenerateColdStartFunctions(test.experimentDuration, test.rpsTarget, test.cooldownSeconds)

			if len(matrix) != len(test.expectedIAT) {
				t.Errorf("Unexpected number of functions - got: %d, expected: %d", len(matrix), len(test.expectedIAT))
			}
			if len(minuteCounts) != len(test.expectedCount) {
				t.Errorf("Unexpected count array size - got: %d, expected: %d", len(minuteCounts), len(test.expectedCount))
			}

			for fIndex := 0; fIndex < len(matrix); fIndex++ {
				if len(matrix[fIndex]) != len(test.expectedIAT[fIndex]) {
					t.Errorf("Unexpected length of function %d IAT array - got: %d, expected: %d", fIndex, len(matrix[fIndex]), len(test.expectedIAT[fIndex]))
				}

				for i := 0; i < len(matrix[fIndex]); i++ {
					if math.Abs(matrix[fIndex][i]-test.expectedIAT[fIndex][i]) > epsilon {
						t.Errorf("Unexpected value fx %d val %d - got: %f; expected: %f", fIndex, i, matrix[fIndex][i], test.expectedIAT[fIndex][i])
					}
				}

				for i := 0; i < len(test.expectedCount[fIndex]); i++ {
					if test.expectedCount[fIndex][i] != minuteCounts[fIndex][i] {
						t.Error("Unexpected per minute count.")
					}
				}
			}
		})
	}
}
