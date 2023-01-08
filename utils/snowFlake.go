package utils

import (
	"file-server/setting"
	"sync"
	"time"
)

type snowFlake struct {
	sync.Mutex
	timeStamp  int64
	machineId  int
	sequenceId int
}

const (
	timeStampBit  = uint(41)
	machineIdBit  = uint(10)
	sequenceIdBit = uint(12)

	timeStampLeftBit = machineIdBit + sequenceIdBit
	machineIdLeftBit = sequenceIdBit

	sequenceIdMax = int64(1 ^ (-1 << sequenceIdBit))
)

func ConstructSnowFlake(machineId int) *snowFlake {
	return &snowFlake{
		timeStamp:  time.Now().UnixMicro(),
		machineId:  machineId,
		sequenceId: 0,
	}
}

var sF snowFlake

func init() {
	sF = *ConstructSnowFlake(int(setting.Conf.MachineID))
}

func GenSFID() int64 {
	sF.Lock()
	nowStamp := time.Now().UnixMilli()
	if nowStamp == sF.timeStamp {
		sF.sequenceId = (sF.sequenceId + 1) & int(sequenceIdMax)
		if sF.sequenceId == 0 {
			for nowStamp <= sF.timeStamp {
				nowStamp = time.Now().UnixMilli()
			}
		}
	} else {
		sF.sequenceId = 0
	}
	sF.timeStamp = nowStamp
	result := int64(sF.timeStamp<<int64(timeStampLeftBit) | int64(sF.machineId)<<int64(machineIdLeftBit) | int64(sF.sequenceId))
	sF.Unlock()
	return result
}
