package core

type StrSlice []string

func (s StrSlice) Pos(value string) int {
  for p, v := range s {
    if v == value {
      return p
    }
  }
  return -1
}
