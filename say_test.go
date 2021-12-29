package gitStudy

import (
	"fmt"
	"reflect"
	"testing"
)

type MisCellIndex int32

const (
	MisCellTestWeekClear      MisCellIndex = 0
	MisCellGuideIndex         MisCellIndex = 1
	MisCellTestDayClear       MisCellIndex = 2
	MisCellSoldierTrainCount  MisCellIndex = 3
	MisCellEquipForgeCount    MisCellIndex = 4
	MisCellEquipRefineCount   MisCellIndex = 5
	MisCellCastCommanderSkill MisCellIndex = 6
	MisCellUpgradeTalentSkill MisCellIndex = 7
)

func TestSay1(t *testing.T) {

	fmt.Println(reflect.ValueOf(MisCellTestWeekClear))
}
