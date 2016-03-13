package core

import (
  "math"
)

type StrSlice []string

func (s StrSlice) Pos(value string) int {
  for p, v := range s {
    if v == value {
      return p
    }
  }
  return -1
}

func MinMaxFloat32(v float32, min float32, max float32) float32 {
  return float32(math.Min(math.Max(float64(v), float64(min)), float64(max)))
}

func MinInt(x, y int) int {
  if x < y {
    return x
  }
  return y
}

func MaxInt(x, y int) int {
  if x > y {
    return x
  }
  return y
}
