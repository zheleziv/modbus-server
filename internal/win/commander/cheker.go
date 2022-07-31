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

type DWCondition struct {
	operator       string
	valueCondition uint32
}
type WCondition struct {
	operator       string
	valueCondition uint16
}

type CCondition struct {
	valueCondition bool
}

func NewChecker(s string, dataType string) (Checker, error) {
	var rtn Checker
	var err error
	switch dataType {
	case tag.COIL_TYPE:
		{
			var tmp CCondition
			tmp, err = newCoilCondition(s)
			if err != nil {
				return rtn, myerr.New(err.Error())
			}
			rtn.condition = append(rtn.condition, &tmp)
		}
	case tag.WORD_TYPE:
		{
			var tmp []ConditionInterface
			tmp, err = newWordCondition(s)
			if err != nil {
				return rtn, myerr.New(err.Error())
			}
			rtn.condition = tmp
		}
	case tag.DWORD_TYPE:
		{
			var tmp []ConditionInterface
			tmp, err = newDWordCondition(s)
			if err != nil {
				return rtn, myerr.New(err.Error())
			}
			rtn.condition = tmp
		}
	default:
		return Checker{}, myerr.New("no such data type")
	}
	return rtn, nil
}

func newCoilCondition(s string) (CCondition, error) {
	var rtn CCondition
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
func (cc *CCondition) checkValue(t tag.TagInterface) bool {
	return (t.(*tag.CoilTag).Value() == 1) && cc.valueCondition
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
			&WCondition{
				operator:       m[1],
				valueCondition: tmpNumber})
	}
	return rtn, nil
}

func (wc *WCondition) checkValue(t tag.TagInterface) bool {
	var condition bool
	switch wc.operator {
	case MORE:
		condition = (t.(*tag.WordTag).Value() > wc.valueCondition)
	case LESS:
		condition = (t.(*tag.WordTag).Value() < wc.valueCondition)
	case EQUAL:
		condition = (t.(*tag.WordTag).Value() == wc.valueCondition)
	case NOT_EQUAL:
		condition = (t.(*tag.WordTag).Value() != wc.valueCondition)
	case MORE_EQUAL:
		condition = (t.(*tag.WordTag).Value() >= wc.valueCondition)
	case LESS_EQUAL:
		condition = (t.(*tag.WordTag).Value() <= wc.valueCondition)
	case BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(wc.valueCondition)))) != 0)
	case NOT_BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(wc.valueCondition)))) == 0)
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
			&DWCondition{
				operator:       m[1],
				valueCondition: tmpNumber})
	}
	return rtn, nil
}

func (dwc *DWCondition) checkValue(t tag.TagInterface) bool {
	var condition bool
	switch dwc.operator {
	case MORE:
		condition = (t.(*tag.DWordTag).Value() > dwc.valueCondition)
	case LESS:
		condition = (t.(*tag.DWordTag).Value() < dwc.valueCondition)
	case EQUAL:
		condition = (t.(*tag.DWordTag).Value() == dwc.valueCondition)
	case NOT_EQUAL:
		condition = (t.(*tag.DWordTag).Value() != dwc.valueCondition)
	case MORE_EQUAL:
		condition = (t.(*tag.DWordTag).Value() >= dwc.valueCondition)
	case LESS_EQUAL:
		condition = (t.(*tag.DWordTag).Value() <= dwc.valueCondition)
	case BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(dwc.valueCondition)))) != 0)
	case NOT_BIT:
		condition = ((t.(*tag.WordTag).Value() & uint16(math.Pow(2, float64(dwc.valueCondition)))) == 0)
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
