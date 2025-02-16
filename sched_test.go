package sched

import (
	"crypto/rand"
	"slices"
	"testing"
)

func TestSched(t *testing.T) {
	buf := make([]int, 4096)
	for i := 0; i < len(buf); i++ {
		buf[i] = i
	}
	for i := 1; i < len(buf); i++ {
		sc := NewTask(buf[:i], func(_ int, t []int) ([]int, error) {
			return t, nil
		}, false, false)
		for j := 1; j < 256; j++ {
			r, err := sc.Collect(j, false, true)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(r, buf[:i]) {
				t.Fatal("expect", buf[:i], "got", r)
			}
		}
	}
}

func BenchmarkSched(b *testing.B) {
	arr := make([]uint8, 64*1024*1024) // 64M
	_, err := rand.Read(arr)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("single", func(b *testing.B) {
		b.SetBytes(int64(len(arr)))
		b.ResetTimer()
		for i := range arr {
			arr[i]++
		}
	})
	_, err = rand.Read(arr)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("para512K", func(b *testing.B) {
		b.SetBytes(int64(len(arr)))
		b.ResetTimer()
		_, _ = NewTask(arr, func(_ int, x []byte) ([]byte, error) {
			for i := range x {
				arr[i] /= 3
				arr[i]++
				arr[i] *= 5
				arr[i]++
				arr[i] /= 7
				arr[i]++
			}
			return arr, nil
		}, false, false).Collect(512*1024, true, true) // 512K batch
	})
	_, err = rand.Read(arr)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("para1M", func(b *testing.B) {
		b.SetBytes(int64(len(arr)))
		b.ResetTimer()
		_, _ = NewTask(arr, func(_ int, x []byte) ([]byte, error) {
			for i := range x {
				arr[i] /= 3
				arr[i]++
				arr[i] *= 5
				arr[i]++
				arr[i] /= 7
				arr[i]++
			}
			return arr, nil
		}, false, false).Collect(1024*1024, true, true) // 1M batch
	})
	_, err = rand.Read(arr)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("para2M", func(b *testing.B) {
		b.SetBytes(int64(len(arr)))
		b.ResetTimer()
		_, _ = NewTask(arr, func(_ int, x []byte) ([]byte, error) {
			for i := range x {
				arr[i] /= 3
				arr[i]++
				arr[i] *= 5
				arr[i]++
				arr[i] /= 7
				arr[i]++
			}
			return arr, nil
		}, false, false).Collect(2*1024*1024, true, true) // 2M batch
	})
}
