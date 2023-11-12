package sharding

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAlgo(t *testing.T) {
	testcases := []struct {
		algo          AlgorithmType
		expectedType  Algorithm
		expectedError error
	}{
		{
			algo:         NaiveAlgorithm,
			expectedType: &naive{},
		},
		{
			algo:         RendezvousAlgorithm,
			expectedType: &rendezvous{},
		},
		{
			algo:         ConsistentAlgorithm,
			expectedType: &consistent{},
		},
		{
			algo:          AlgorithmType("unknown"),
			expectedType:  nil,
			expectedError: fmt.Errorf("unknown sharding algorithm: %s", "unknown"),
		},
	}

	for _, tc := range testcases {
		t.Run(string(tc.algo), func(t *testing.T) {
			algo, err := NewAlgo(tc.algo, nil, nil)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			require.IsType(t, tc.expectedType, algo)
		})
	}
}
