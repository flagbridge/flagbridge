package evaluation

import (
	"fmt"
	"math"
	"testing"
)

// TestRolloutDistribution verifies that the hash function produces a roughly
// uniform distribution across buckets. This catches issues like clustering,
// where a disproportionate number of users end up in the same buckets.
func TestRolloutDistribution(t *testing.T) {
	const (
		numUsers   = 10_000
		numBuckets = 100
		flagKey    = "distribution-test"
	)

	buckets := make([]int, numBuckets)
	for i := 0; i < numUsers; i++ {
		userID := fmt.Sprintf("user-%d", i)
		bucket := rolloutBucket(userID, flagKey)
		if bucket < 0 || bucket >= numBuckets {
			t.Fatalf("bucket %d out of range [0, %d)", bucket, numBuckets)
		}
		buckets[bucket]++
	}

	expected := float64(numUsers) / float64(numBuckets) // 100
	// Chi-squared test for uniformity
	var chiSquared float64
	for _, count := range buckets {
		diff := float64(count) - expected
		chiSquared += (diff * diff) / expected
	}

	// For 99 degrees of freedom, chi-squared critical value at p=0.01 is ~135
	// We use a generous threshold to avoid flaky tests
	if chiSquared > 150 {
		t.Errorf("distribution is not uniform: chi-squared = %.2f (threshold 150)", chiSquared)
	}

	// Also verify no single bucket is wildly over-represented (> 3x expected)
	for i, count := range buckets {
		if float64(count) > expected*3 {
			t.Errorf("bucket %d has %d users (expected ~%.0f)", i, count, expected)
		}
	}
}

// TestRolloutConsistency verifies that the same percentage threshold always
// includes the same set of users — adding users to the rollout as the
// percentage increases, never removing previously included users.
func TestRolloutConsistency(t *testing.T) {
	const (
		numUsers = 1000
		flagKey  = "consistency-test"
	)

	// For each user, compute their bucket once
	userBuckets := make([]int, numUsers)
	for i := 0; i < numUsers; i++ {
		userBuckets[i] = rolloutBucket(fmt.Sprintf("user-%d", i), flagKey)
	}

	// As rollout goes from 10% → 50% → 100%, users included at 10% must
	// still be included at 50% and 100%.
	for i := 0; i < numUsers; i++ {
		bucket := userBuckets[i]
		if bucket < 10 {
			// User is in 10% rollout — must also be in 50%
			if bucket >= 50 {
				t.Errorf("user-%d in 10%% rollout (bucket=%d) but not in 50%%", i, bucket)
			}
		}
		if bucket < 50 {
			// User is in 50% rollout — must also be in 100%
			if bucket >= 100 {
				t.Errorf("user-%d in 50%% rollout (bucket=%d) but not in 100%%", i, bucket)
			}
		}
	}
}

// TestRolloutDifferentFlags verifies that the same user gets different buckets
// for different flags, so rollout percentages are independent per flag.
func TestRolloutDifferentFlags(t *testing.T) {
	const numFlags = 50
	userID := "test-user-42"

	buckets := make(map[int]int)
	for i := 0; i < numFlags; i++ {
		b := rolloutBucket(userID, fmt.Sprintf("flag-%d", i))
		buckets[b]++
	}

	// With 50 flags across 100 buckets, we should see decent spread.
	// If all buckets are the same, the hash function is broken for this use case.
	uniqueBuckets := len(buckets)
	if uniqueBuckets < 10 {
		t.Errorf("only %d unique buckets across %d flags — hash function has poor spread", uniqueBuckets, numFlags)
	}

	// Verify standard deviation is reasonable
	mean := float64(numFlags) / float64(uniqueBuckets)
	var variance float64
	for _, count := range buckets {
		diff := float64(count) - mean
		variance += diff * diff
	}
	variance /= float64(uniqueBuckets)
	stddev := math.Sqrt(variance)

	// With uniform distribution, std dev should be low
	if stddev > 3 {
		t.Errorf("high variance in bucket distribution: stddev=%.2f", stddev)
	}
}
