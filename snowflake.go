package snowflake

import (
	"errors"
	"sync"
	"time"
)

var (
	SequenceBit   uint = 12
	MachineBit    uint = 5
	DataCenterBit uint = 5

	MaxSequence      int64 = -1 ^ (-1 << SequenceBit)
	MaxMachineNum    int64 = -1 ^ (-1 << MachineBit)
	MaxDataCenterNum int64 = -1 ^ (-1 << DataCenterBit)

	MachineLeft    = SequenceBit
	DataCenterLeft = SequenceBit + MachineBit
	TimestampLeft  = SequenceBit + MachineBit + DataCenterBit

	StartTimeStamp int64 = 1330808767000
)

type IDGenerator struct {
	dataCenterId  int64
	machineId     int64
	sequence      int64
	lastTimeStamp int64
	mutex         sync.Mutex
}

func (g *IDGenerator) NextId() (int64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	currTimeStamp := getNewTimeStamp()
	if currTimeStamp < g.lastTimeStamp {
		return 0, errors.New("Clock moved backwards.  Refusing to generate id")
	}

	if currTimeStamp == g.lastTimeStamp {
		g.sequence = (g.sequence + 1) & MaxSequence
		if g.sequence == 0 {
			currTimeStamp = getNextMill(g.lastTimeStamp)
		}
	} else {
		g.sequence = 0
	}

	g.lastTimeStamp = currTimeStamp

	return ((currTimeStamp - StartTimeStamp) << TimestampLeft) | (g.dataCenterId << DataCenterLeft) | (g.machineId << MachineLeft) | g.sequence, nil
}

func getNewTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getNextMill(lastTimeStamp int64) int64 {
	mill := getNewTimeStamp()
	for mill <= lastTimeStamp {
		mill = getNewTimeStamp()
	}
	return mill
}

func NewIDGenerator(dataCenterId, machineId int64) (*IDGenerator, error) {
	if dataCenterId > MaxDataCenterNum || dataCenterId < 0 {
		return nil, errors.New("DtaCenterId can't be greater than MAX_DATA_CENTER_NUM or less than 0！")
	}
	if machineId > MaxMachineNum || machineId < 0 {
		return nil, errors.New("MachineId can't be greater than MAX_MACHINE_NUM or less than 0！")
	}

	return &IDGenerator{
		dataCenterId: dataCenterId,
		machineId:    machineId,
	}, nil
}
