package commander

import (
	"os/exec"
	"strings"
	"sync"
	"time"

	"zheleznovux.com/modbus-console/internal/configuration"
	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client"
	"zheleznovux.com/modbus-console/pkg/client/tag"
)

type Commander struct {
	name           string
	stateCondition bool
	cheker         Checker
	logic          string
	action         string
	actionTimeout  time.Duration
	scanPeriod     time.Duration
	log            Logger
	th_ptr         *tag.TagInterface
}

// инициализация=========================================== -- {

func verifyName(name string) (string, error) { // эта функция также используется для верификации имени тэга
	rtn := strings.TrimSpace(name)

	if rtn == "" {
		return rtn, myerr.New("did not have Name")
	}
	return rtn, nil
}

func verifyLogic(logic string) (string, error) {
	str := strings.TrimSpace(strings.ToLower(logic))

	switch logic {
	case AND:
		return AND, nil
	case OR:
		return OR, nil
	case "&&":
		return AND, nil
	case "||":
		return OR, nil
	default:
		return str, myerr.New("did not have logic")
	}
}

func verifyStateCondition(state string) (bool, error) {
	str := strings.TrimSpace(strings.ToLower(state))

	switch str {
	case client.BAD:
		return false, nil
	case client.GOOD:
		return true, nil
	default:
		return true, myerr.New("did not have state")
	}
}

func verifyAction(logic string) (string, error) { // regexp
	str := strings.TrimSpace(strings.ToLower(logic))

	switch str {
	case SHUTDOWN, "/s", "s", "-s", "--s":
		return SHUTDOWN, nil
	case RESTART, "/r", "r", "-r", "--r":
		return RESTART, nil
	default:
		{
			if len(str) != 0 {
				return str, nil
			} else {
				return str, myerr.New("did not have action")
			}
		}
	}
}

func verifyActionTimeout(t float64) (time.Duration, error) {
	if t < MIN_ACTION_TIMEOUT {
		return 0, myerr.New("action timeout < 0")
	}
	return time.Duration(t * float64(time.Second)), nil
}

func verifyScanPeriod(sp float64) (time.Duration, error) {
	if sp < MIN_SCAN_PERIOD {
		return 0, myerr.New("scan period < 0.001")
	}
	return time.Duration(sp * float64(time.Second)), nil
}

func (thisCommander *Commander) Setup(nt configuration.NodeTag, th *tag.TagInterface) error {
	var err error

	thisCommander.name, err = verifyName(nt.Name)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.cheker, err = NewChecker(nt.ValueCondition, (*th).DataType())
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.logic, err = verifyLogic(nt.Logic)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.stateCondition, err = verifyStateCondition(nt.StateCondition)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.action, err = verifyAction(nt.Action)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.actionTimeout, err = verifyActionTimeout(nt.ActionTimeout)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.scanPeriod, err = verifyScanPeriod(nt.ScanPeriod)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisCommander.th_ptr = th

	thisCommander.log = Logger{
		ParentNodeName: thisCommander.name,
		IsLogOutput:    nt.Log,
	}

	return nil
}

// инициализация=========================================== -- }

// логика================================================== -- {

func (thisCommander *Commander) StartChecking(quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(thisCommander.scanPeriod) // timer
	condition := make(chan bool)

	wg.Add(1)
	go thisCommander.startCommand(condition, quit, wg)
	for {
		select {
		case <-quit:
			{
				return
			}
		case <-ticker.C:
			{
				ticker.Stop()
				condition <- thisCommander.generalTagCheck()

				ticker.Reset(thisCommander.scanPeriod)
			}
		}
	}
}

func (thisCommander *Commander) generalTagCheck() bool {
	if thisCommander.stateCondition == (*thisCommander.th_ptr).State() {
		return thisCommander.cheker.CheckValues(*thisCommander.th_ptr)
	}
	return false
}

func (thisCommander *Commander) startCommand(condition chan bool, quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	timeBetweenTick := thisCommander.actionTimeout / 5
	timer := &time.Timer{}
	tickCount := 0

	var lastCondition bool
	for {
		select {
		case <-quit:
			if lastCondition {
				thisCommander.log.Write(INFO, "Таймер команды остановлен")
			}
			return
		case <-timer.C:
			timer.Stop()

			if tickCount < 5 {
				timeToCommand := thisCommander.actionTimeout - time.Duration(tickCount)*timeBetweenTick
				thisCommander.log.Write(INFO, "Команда "+thisCommander.action+", до завершения таймера: "+timeToCommand.String())
				tickCount++
			} else {
				tickCount = 0
				err := command(thisCommander.action)
				thisCommander.log.Write(INFO, "Запущена команда")
				if err != nil {
					thisCommander.log.Write(ERROR, err.Error())
				}
			}

			timer.Reset(timeBetweenTick)
		case v := <-condition:
			if lastCondition != v {
				lastCondition = v
				tickCount = 0
				if v {
					thisCommander.log.Write(INFO, "Запущен таймер команды "+thisCommander.action)
					timer = time.NewTimer(1)
					continue
				}
				thisCommander.log.Write(INFO, "Таймер команды остановлен по значению")
				timer.Stop()
			}
		}
	}
}

func command(c string) error {
	var flag string
	exe := strings.Split(c, " ")

	switch exe[0] {
	case SHUTDOWN:
		flag = "/s"
	case RESTART:
		flag = "/r"
	case RUN_PROGRAM:
		if len(exe) != 2 {
			return myerr.New("len(exe) != 2")
		}
		cmd := exec.Command("./" + exe[1])
		err := cmd.Run()
		if err != nil {
			return myerr.New(err.Error())
		}
		return nil
	default:
		return myerr.New("invalid command")
	}

	if err := exec.Command("cmd", "/C", "shutdown "+flag+" /t 1").Run(); err != nil {
		return myerr.New(err.Error())
	}
	return nil
}

// логика================================================== -- }
