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

func (thisC *Commander) Setup(nt configuration.NodeTag, th *tag.TagInterface) error {
	var err error

	thisC.name, err = verifyName(nt.Name) // проверка из cdApp
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.cheker, err = NewChecker(nt.ValueCondition, (*th).DataType())
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.logic, err = verifyLogic(nt.Logic)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.stateCondition, err = verifyStateCondition(nt.StateCondition)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.action, err = verifyAction(nt.Action)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.actionTimeout, err = verifyActionTimeout(nt.ActionTimeout)
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.scanPeriod, err = verifyScanPeriod(nt.ScanPeriod) // проверка из cdApp
	if err != nil {
		return myerr.New(err.Error())
	}

	thisC.th_ptr = th

	thisC.log = Logger{
		ParentNodeName: thisC.name,
		IsLogOutput:    nt.Log,
	}
	return nil
}

// инициализация=========================================== -- }

// логика================================================== -- {

func (dwc *Commander) StartChecking(quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(dwc.scanPeriod)
	var condition chan bool = make(chan bool)

	wg.Add(1)
	go dwc.startCommand(condition, quit, wg)
	for {
		select {
		case <-quit:
			{
				return
			}
		case <-ticker.C:
			{
				ticker.Stop()
				condition <- dwc.cheker.CheckValues(*dwc.th_ptr)
				ticker.Reset(dwc.scanPeriod)
			}
		}
	}
}

func (wc *Commander) startCommand(condition chan bool, quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	timeBetweenTick := wc.actionTimeout / 5
	tickerToCommand := time.NewTicker(wc.actionTimeout)
	tickerToCommand.Stop()
	tickCount := 0

	var lastCondition bool
	for {
		select {
		case <-quit:
			{
				if lastCondition {
					wc.log.Write(INFO, "Таймер команды остановлен из-за смены конфига!")
				}
				return
			}
		case <-tickerToCommand.C:
			{
				tickerToCommand.Stop()
				tickCount++
				if tickCount != 5 {
					timeToCommand := wc.actionTimeout - time.Duration(tickCount)*timeBetweenTick
					wc.log.Write(INFO, "Команда "+wc.action+", до завершения таймера: "+timeToCommand.String()+".")

				} else {
					tickCount = 0
					wc.log.Write(INFO, "Запущена команда!")
					err := command(wc.action)
					if err != nil {
						wc.log.Write(ERROR, err.Error())
					}
				}
				tickerToCommand.Reset(timeBetweenTick)
			}
		case v := <-condition:
			{
				if lastCondition != v {
					lastCondition = v
					tickCount = 0
					if v {
						wc.log.Write(INFO, "Запущен таймер команды "+wc.action+", до завершения: "+wc.actionTimeout.String()+".")
						tickerToCommand.Reset(timeBetweenTick)
					} else {
						wc.log.Write(INFO, "Таймер команды остановлен по значению!")
						tickerToCommand.Stop()
					}
				}
			}
		}
	}
}

func command(c string) error {
	var flag string
	exe := strings.Split(c, " ")

	switch exe[0] {
	case SHUTDOWN:
		{
			flag = "/s"
		}
	case RESTART:
		{
			flag = "/r"
		}
	case RUN_PROGRAM:
		{
			if len(exe) != 2 {
				return myerr.New("len(exe) != 2")
			}
			cmd := exec.Command("./" + exe[1])
			err := cmd.Run()
			if err != nil {
				return myerr.New(err.Error())
			}
			return nil
		}
	default:
		return myerr.New("invalid command")
	}

	if err := exec.Command("cmd", "/C", "shutdown "+flag+" /t 1").Run(); err != nil {
		return myerr.New(err.Error())
	}
	return nil
}

// логика================================================== -- }
