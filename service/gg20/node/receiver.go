package node

// RecvDkgRound1 接收其它节点（称之为j）发来的第一轮数据
func (d *DkgOperator) RecvDkgRound1(recv DkgRound1Recv) {
	d.ChanRecvRound1 <- recv
}

// RecvDkgRound2 接收其它节点（称之为j）发来的第二轮数据
func (d *DkgOperator) RecvDkgRound2(recv DkgRound2Recv) {
	d.ChanRecvRound2 <- recv
}

// RecvDkgRound3 接收其它节点（称之为j）发来的第三轮数据
func (d *DkgOperator) RecvDkgRound3(recv DkgRound3Recv) {
	d.ChanRecvRound3 <- recv
}
