package node

import (
	"github.com/golang/glog"
)

// CandidateInfo 候选的签名节点信息
type CandidateInfo struct {
	Id     uint32 `json:"id"`
	WorkId string `json:"workId"` // 签名工作编号id (目前使用uuid)
}

// DoAskForCosignerCandidate 向其它节点请求获取其id，最先返回的前threshold-1个节点，即与s本身一起作为本次签名的cosigner
func (s *SignOperator) DoAskForCosignerCandidate(workId string) {
	for _, nodeId := range s.otherParticipants {
		address := s.participantAddrs[nodeId-1]
		url := GetAskCosignerCandidateUrl(address)
		go func() {
			candidate, err := DoSendAskForCosignerCandidate(url, workId)
			if err != nil {
				glog.Errorf("请求获取Cosigner失败: %v\n", err)
			}
			info, exist := s.signerMap[candidate.WorkId[1:len(candidate.WorkId)-1]]
			//反序列化响应数据，string的WorkId会前后各多一个"号，需要将其去掉
			if !exist {
				glog.Error("当前WorkId对应的SignerInfo不存在：", candidate.WorkId)
			}
			info.candidateInfoChan <- candidate
		}()
	}
}
