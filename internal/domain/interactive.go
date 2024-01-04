package domain

// Interactive 这个是总体交互的计数
type Interactive struct {
	BizId      int64 `json:"biz_id"`
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`
}
