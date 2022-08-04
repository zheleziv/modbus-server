package commander

import (
	"fmt"
	"math"
	"regexp"

	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client/tag"
)

type Checker struct {
	condition []ConditionInterface
}

type ConditionInterface interface {
	checkValue(tag.TagInterface) bool
}

type DWordCondition struct {
	operator       string
	valueCondition uint32
}
type WordCondition struct {
	operator       string
	valueCondition uint16
}

type CoilCondition struct {
	valueCondition bool
}

func NewChecker(s string, dataType string) (Checker, error) {
	var rtn Checker
	var err error
	switch dataType {
	case tag.COIL_TYPE:
		var tmp CoilCondition
		tmp, err = newCoilCondition(s)
		if err != nil {
			return rtn, myerr.New(err.Error())
		}
		rtn.condition = append(rtn.condition, &tmp)
	case tag.WORD_TYPE:
		var tmp []ConditionInterface
		tmp, err = newWordCondition(s)
		if err != nil {
			return rtn, myerr.New(err.Error())
		}
		rtn.condition = tmp
	case tag.DWORD_TYPE:
		var tmp []ConditionInterface
		tmp, err = newDWordCondition(s)
		if err != nil {
			return rtn, myerr.New(err.Error())
		}
		rtn.condition = tmp
	default:
		return Checker{}, myerr.New("no such data type")
	}
	return rtn, nil
}

func newCoilCondition(s string) (CoilCondition, error) {
	var rtn CoilCondition
	if s == "true" {
		rtn.valueCondition = true
		return rtn, nil
	} else if s == "false" {
		rtn.valueCondition = false
		return rtn, nil
	} else {
		return rtn, myerr.New("did not have coil value condition")
	}
}
func (thisCoilCondition *CoilCondition) checkValue(t tag.TagInterface) bool {
	return (t.(*tag.CoilTag).Value() == 1) && thisCoilCondition.valueCondition
}

func newWordCondition(s string) ([]ConditionInterface, error) {
	var rtn []ConditionInterface

	re := regexp.MustCompile(`(?P<OPERATOR>>|>=|<|<=|!=|==|bit|!bit)\(?(?P<VALUE>[0-9]+)\)?`)
	matcher := re.FindAllStringSubmatch(s, -1)

	if len(matcher) == 0 {
		return nil, myerr.New("regexp found no match")
	}

	for _, m := range matcher {
		var tmpNumber uint16
		_, err := fmt.Sscanf(m[2], "%d", &tmpNumber)
		if err != nil {
			return nil, myerr.New(err.Error())
		}
		rtn = append(
			rtn,
			&WordCondition{
				operator:       m[1],
				valueCondition: tmpNumber})
	}
	return rtn, nil
}

func (thisWordCondition *WordCondition) checkValue(t tag.TagInterface) bool {
	var condition bool
	switch thisWordCondition.operator {
	case MORE:
		condition = (t.(*tag.WordTag).Value() > thisWordCondition.valueCondition)
	case LESS:
		condition = (t.(*tag.WordTag).Value() < thisWordCondition.valueCondition)
	case EQUAL:
		condition = (t.(*tag.WordTag).Value() == thisWordCondition.valueCondition)
	case NOT_EQUAL:
		condition = (t.(*tag.WordTag).Value() != thisWordCondition.valueCondition)
	case MORE_EQUAL:
		condition = (t.(*tag.WordTag).Value() >= thisWordCondition.valueCondition)
	case LESS_EQUAL:
		condition = (t.(*tag.WordTag).Value() <= thisWordCondition.valueCondition)
	case BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(thisWordCondition.valueCondition)))) != 0)
	case NOT_BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(thisWordCondition.valueCondition)))) == 0)
	default:
		return false
	}
	return condition
}

func newDWordCondition(s string) ([]ConditionInterface, error) {
	var rtn []ConditionInterface

	re := regexp.MustCompile(`(?P<OPERATOR>>|>=|<|<=|!=|==|bit|!bit)\(?(?P<VALUE>[0-9]+)\)?`)
	matcher := re.FindAllStringSubmatch(s, -1)

	if len(matcher) == 0 {
		return nil, myerr.New("regexp found no match")
	}

	for _, m := range matcher {
		var tmpNumber uint32
		_, err := fmt.Sscanf(m[2], "%d", &tmpNumber)
		if err != nil {
			return nil, myerr.New(err.Error())
		}
		rtn = append(
			rtn,
			&DWordCondition{
				operator:       m[1],
				valueCondition: tmpNumber})
	}
	return rtn, nil
}

func (thisDWordCondition *DWordCondition) checkValue(t tag.TagInterface) bool {
	var condition bool
	switch thisDWordCondition.operator {
	case MORE:
		condition = (t.(*tag.DWordTag).Value() > thisDWordCondition.valueCondition)
	case LESS:
		condition = (t.(*tag.DWordTag).Value() < thisDWordCondition.valueCondition)
	case EQUAL:
		condition = (t.(*tag.DWordTag).Value() == thisDWordCondition.valueCondition)
	case NOT_EQUAL:
		condition = (t.(*tag.DWordTag).Value() != thisDWordCondition.valueCondition)
	case MORE_EQUAL:
		condition = (t.(*tag.DWordTag).Value() >= thisDWordCondition.valueCondition)
	case LESS_EQUAL:
		condition = (t.(*tag.DWordTag).Value() <= thisDWordCondition.valueCondition)
	case BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(thisDWordCondition.valueCondition)))) != 0)
	case NOT_BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(thisDWordCondition.valueCondition)))) == 0)
	default:
		return false
	}
	return condition
}

func (thisCheker *Checker) CheckValues(t tag.TagInterface) (rtn bool) {
	rtn = true
	for i := range thisCheker.condition {
		rtn = rtn && thisCheker.condition[i].checkValue(t)
	}
	return
}
