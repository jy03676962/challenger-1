package core

type Room struct {
  MaxNum int      `json:"max"`
  Hoster string   `json:"hoster"`
  Member []string `json:"member"`
}

type ArenaPosition struct {
  X int
  Y int
}

type P ArenaPosition

type ArenaWall struct {
  P1 P
  P2 P
}

type W ArenaWall

type GameVar struct {
  ArenaWidth    int     `json:"arenaWidth"`    // 场地宽度
  ArenaHeight   int     `json:"arenaHeight"`   // 场地高度
  ArenaCellSize int     `json:"arenaCellSize"` // 场地格子大小
  ArenaBorder   int     `json:"arenaBorder"`   // 场地边框宽度
  Warmup        float32 `json:"warmup"`        // 预热时间
  ArenaWallList []W     `json:"walls"`         // 墙列表
}

func DefaultGameVar() *GameVar {
  v := GameVar{}
  v.ArenaWidth = 8
  v.ArenaHeight = 6
  v.ArenaCellSize = 60
  v.ArenaBorder = 4
  w := []W{
    W{P{4, 0}, P{5, 0}},
    W{P{1, 0}, P{1, 1}},
    W{P{6, 0}, P{6, 1}},
    W{P{0, 1}, P{1, 1}},
    W{P{2, 1}, P{3, 1}},
    W{P{3, 1}, P{4, 1}},
    W{P{6, 1}, P{7, 1}},
    W{P{2, 1}, P{2, 2}},
    W{P{5, 1}, P{5, 2}},
    W{P{2, 2}, P{3, 2}},
    W{P{3, 2}, P{4, 2}},
    W{P{4, 2}, P{5, 2}},
    W{P{1, 2}, P{1, 3}},
    W{P{6, 2}, P{6, 3}},
    W{P{0, 3}, P{1, 3}},
    W{P{2, 3}, P{3, 3}},
    W{P{4, 3}, P{5, 3}},
    W{P{6, 3}, P{7, 3}},
    W{P{2, 3}, P{2, 4}},
    W{P{3, 3}, P{3, 4}},
    W{P{0, 4}, P{1, 4}},
    W{P{1, 4}, P{2, 4}},
    W{P{4, 4}, P{5, 4}},
    W{P{5, 4}, P{6, 4}},
    W{P{6, 4}, P{7, 4}},
    W{P{4, 4}, P{4, 5}},
    W{P{2, 5}, P{3, 5}},
    W{P{5, 5}, P{6, 5}},
    W{P{6, 5}, P{7, 5}},
  }
  v.ArenaWallList = w
  return &v
}
