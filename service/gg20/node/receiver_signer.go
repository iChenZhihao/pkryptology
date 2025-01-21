package node

// RecvAskCandidateInfo 接收到需要签名的请求，将workId与自身id一起返回
func (s *SignOperator) RecvAskCandidateInfo(workId string) *CandidateInfo {
	return &CandidateInfo{Id: s.id, WorkId: workId}
}
