package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/things-go/go-modbus"
	myerr "zheleznovux.com/modbus-console/pkg"
	"zheleznovux.com/modbus-console/pkg/client/logger"
	tag "zheleznovux.com/modbus-console/pkg/client/tag"
)

type сlientModbus struct {
	name               string
	ip                 string
	port               int
	slaveId            uint8
	connectionAttempts uint
	connectionTimeout  time.Duration
	log                logger.Logger
	state              string
	tags               []tag.TagInterface
	mutex              sync.RWMutex
	modbus.Client
}

func (c *сlientModbus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name           string
		ConnectionType string
		Ip             string
		SlaveId        uint8
		State          string
		Tags           []tag.TagInterface
	}{
		Name:           c.name,
		ConnectionType: c.Type(),
		Ip:             c.ip,
		SlaveId:        c.slaveId,
		State:          c.state,
		Tags:           c.tags,
	})
}

// ======================инициализация========================{

// =================================Name
func (thisModbusClient *сlientModbus) Name() string {
	return thisModbusClient.name
}
func (thisModbusClient *сlientModbus) SetName(name string) error {
	if strings.TrimSpace(name) == "" {
		return myerr.New("invalid client name")
	}
	thisModbusClient.name = name
	return nil
}

// =================================Type
func (thisModbusClient *сlientModbus) Type() string {
	return MODBUS_TCP
}

// ===================================IP
func (thisModbusClient *сlientModbus) Ip() string {
	return thisModbusClient.ip
}

// using net.parseIp
func (thisModbusClient *сlientModbus) SetIp(ip string) error {
	ipAddr := net.ParseIP(strings.TrimSpace(ip))
	if ipAddr == nil {
		return myerr.New("invalid client Ip")
	} else {
		thisModbusClient.ip = ip
		return nil
	}
}

// ===================ConnectionAttempts
func (thisModbusClient *сlientModbus) ConnectionAttempts() uint {
	return thisModbusClient.connectionAttempts
}
func (thisModbusClient *сlientModbus) SetConnectionAttempts(ca uint) error {
	if ca <= 0 {
		return myerr.New("invalid client connection attempts")
	}
	thisModbusClient.connectionAttempts = ca
	return nil
}

// =================================Port
func (thisModbusClient *сlientModbus) Port() int {
	return thisModbusClient.port
}
func (thisModbusClient *сlientModbus) SetPort(port int) error {
	if (port > 0xFFFF) || (port < 0) {
		thisModbusClient.port = 502
		return myerr.New("invalid client port")
	} else {
		thisModbusClient.port = port
		return nil
	}
}

// ==============================SlaveID
func (thisModbusClient *сlientModbus) SalveId() uint8 {
	return thisModbusClient.slaveId
}
func (thisModbusClient *сlientModbus) SetSalveId(sid uint8) error {
	if sid > 0xFF {
		return myerr.New("invalid client slaveID")
	}
	thisModbusClient.slaveId = sid
	return nil
}

// =================================Tags
func (thisModbusClient *сlientModbus) Tags() []tag.TagInterface {
	return thisModbusClient.tags
}

// ============================TagById
func (thisModbusClient *сlientModbus) TagById(id int) (tag.TagInterface, error) {
	if (id >= len(thisModbusClient.tags)) || (id < 0) {
		return nil, myerr.New("invalid id client tag")
	}
	return thisModbusClient.tags[id], nil
}

func (thisModbusClient *сlientModbus) TagByName(name string) (tag.TagInterface, error) {
	for id := range thisModbusClient.tags {
		if thisModbusClient.tags[id].Name() == name {
			return thisModbusClient.tags[id], nil
		}
	}
	return nil, myerr.New("invalid client tag name")
}

func (thisModbusClient *сlientModbus) SetTag(name string, address uint32, scanPeriod float64, dataType string) error {
	if _, err := thisModbusClient.TagByName(name); err != nil {
		adr16, function_number, err := checkModbusAddress(address, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}
		t, err := tag.NewTag(name, adr16, scanPeriod, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}

		switch v := t.(type) {
		case *tag.CoilTag:
			var function func(byte, uint16, uint16) ([]byte, error)
			if function_number == FUNCTION_1 {
				function = thisModbusClient.ReadCoils
			} else if function_number == FUNCTION_2 {
				function = thisModbusClient.ReadDiscreteInputs
			} else {
				return myerr.New("invalid function number")
			}

			v.ReadFunc = func() (byte, error) {
				thisModbusClient.mutex.Lock()
				resp, err := function(thisModbusClient.slaveId, v.Address(), 1)
				thisModbusClient.mutex.Unlock()
				if err != nil {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ с ошибкой")
					return 0, myerr.New(err.Error())
				}

				if len(resp) == 0 {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен пустой ответ")
					return 0, myerr.New("empty response")
				}
				thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ:"+fmt.Sprintf("%d", resp[0]))
				return resp[0], nil
			}
		case *tag.WordTag:
			var function func(byte, uint16, uint16) ([]uint16, error)
			if function_number == FUNCTION_3 {
				function = thisModbusClient.ReadHoldingRegisters
			} else if function_number == FUNCTION_4 {
				function = thisModbusClient.ReadInputRegisters
			} else {
				return myerr.New("invalid function number")
			}
			v.ReadFunc = func() (uint16, error) {
				thisModbusClient.mutex.Lock()
				resp, err := function(thisModbusClient.slaveId, v.Address(), 1)
				thisModbusClient.mutex.Unlock()
				if err != nil {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ с ошибкой")
					return 0, myerr.New(err.Error())
				}

				if len(resp) == 0 {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен пустой ответ")
					return 0, myerr.New("empty response")
				}
				thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ:"+fmt.Sprintf("%d", resp[0]))
				return resp[0], nil
			}
		case *tag.DWordTag:
			var function func(byte, uint16, uint16) ([]uint16, error)
			if function_number == FUNCTION_3 {
				function = thisModbusClient.ReadHoldingRegisters
			} else if function_number == FUNCTION_4 {
				function = thisModbusClient.ReadInputRegisters
			} else {
				return myerr.New("invalid function number")
			}
			v.ReadFunc = func() (uint32, error) {
				thisModbusClient.mutex.Lock()
				resp, err := function(thisModbusClient.slaveId, v.Address(), 2)
				thisModbusClient.mutex.Unlock()
				if err != nil {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ с ошибкой")
					return 0, myerr.New(err.Error())
				}

				if len(resp) < 2 {
					thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен пустой ответ")
					return 0, myerr.New("empty response")
				}
				thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, v.Name(), "Получен ответ:"+fmt.Sprintf("%d", (uint32(resp[0])<<16)+uint32(resp[1])))
				return (uint32(resp[0]) << 16) + uint32(resp[1]), nil
			}
		}

		thisModbusClient.tags = append(thisModbusClient.tags, t)
		return nil
	}
	return myerr.New("client tag name already exists")
}

func checkModbusAddress(address uint32, dataType string) (uint16, int, error) {
	// if address >= tag.UINT16_MAX_VALUE {
	t := FUNCTION__TAG_TYPE[dataType]
	if t == nil {
		return 0, 0, myerr.New("invalid tag data type")
	}

	tmpINT := int(address / 100000.0)
	fmt.Println(tmpINT)
	isDefined := false
	for i := range t {
		if t[i] == tmpINT {
			isDefined = isDefined || true
		}
	}
	if !isDefined {
		return 0, 0, myerr.New("invalid tag address")
	}

	tmpUINT16 := uint16(address - uint32(tmpINT*100000))
	if tmpUINT16 >= tag.UINT16_MAX_VALUE {
		return 0, 0, myerr.New("invalid tag address")
	}
	return tmpUINT16 - 1, tmpINT, nil
	// }
	// return uint16(address - 1), nil
}

// ====================ConnectionTimeout
func (thisModbusClient *сlientModbus) ConnectionTimeout() time.Duration {
	return thisModbusClient.connectionTimeout
}
func (thisModbusClient *сlientModbus) SetConnectionTimeout(s float64) error {
	if s < 0 {
		return myerr.New("client connection timeout < 0")
	}
	thisModbusClient.connectionTimeout = time.Duration(s) * time.Second
	return nil
}

// // ==============================State
func (thisModbusClient *сlientModbus) setState(state string) (isChanged bool) {
	if thisModbusClient.state != state {
		thisModbusClient.state = state
		isChanged = true
	} else {
		isChanged = false
	}
	return
}

// ==========================Constructor
// конструктор Modbus TCP/IP клиента с проверками
func NewClinetModbus(ip string, port int, slaveID uint8, name string, debug bool, ConnectionAttempts uint, ConnectionTimeout float64) (*сlientModbus, error) {
	var c сlientModbus

	err := c.SetName(name)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	err = c.SetIp(ip)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	err = c.SetPort(port)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	err = c.SetSalveId(slaveID)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	err = c.SetConnectionAttempts(ConnectionAttempts)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	err = c.SetConnectionTimeout(ConnectionTimeout)
	if err != nil {
		return nil, myerr.New(err.Error())
	}

	c.log = logger.Logger{
		ParentNodeName: c.name,
		ParentNodeIp:   c.ip,
		ParentNodeId:   c.slaveId,
		IsDebug:        debug,
	}

	provider := modbus.NewTCPClientProvider(
		c.ip + ":" + fmt.Sprint(c.port))
	c.Client = modbus.NewClient(provider)

	return &c, nil
}

//======================инициализация========================}

func (thisModbusClient *сlientModbus) checkConnect() error {
	defer thisModbusClient.Close()
	defer thisModbusClient.mutex.Unlock()
	//устанавливаем соединение
	thisModbusClient.mutex.Lock()
	err := thisModbusClient.Connect()
	if err != nil {
		return myerr.New(err.Error())
	} else {
		return nil
	}
}

func (thisModbusClient *сlientModbus) Start(stop chan struct{}, wg *sync.WaitGroup) {
	connection := make(chan bool)
	quit := make(chan struct{})

	defer wg.Done()

	var wgi sync.WaitGroup
	wgi.Add(1)
	go thisModbusClient.tryConnect(stop, connection, &wgi)
	for {
		select {
		case <-stop: //канал сверху. Завершение сессии
			close(quit)
			wgi.Wait()
			return
		case cb := <-connection: // канал снизу. Плохое подключение => реконект
			close(quit)
			wgi.Wait()
			quit = make(chan struct{})
			if cb {
				for tagId := range thisModbusClient.tags {
					wgi.Add(1)
					go thisModbusClient.startSender(tagId, quit, &wgi, connection)
				}
				continue
			}
			wgi.Add(1)
			go thisModbusClient.tryConnect(stop, connection, &wgi)
		}
	}
}

func (thisModbusClient *сlientModbus) startSender(tagId int, quit chan struct{}, wg *sync.WaitGroup, connect chan bool) {
	thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, thisModbusClient.tags[tagId].Name(), "Запущен опрос тега")
	defer wg.Done()

	timer := time.NewTimer(1)

	for {
		select {
		case <-quit:
			thisModbusClient.log.DebugWithTag(logger.INFO, thisModbusClient.state, thisModbusClient.tags[tagId].Name(), "Завершен опрос тега")
			return
		case <-timer.C:
			timer.Stop()
			if thisModbusClient.checkConnect() != nil {
				connect <- false
				return
			}
			err := thisModbusClient.tags[tagId].ReadDevice()
			if err != nil {
				thisModbusClient.log.WriteWithTag(logger.ERROR, thisModbusClient.state, thisModbusClient.tags[tagId].Name(), err.Error())
			}
			timer.Reset(time.Duration(thisModbusClient.tags[tagId].ScanPeriod()) * time.Second)
		}
	}
}

func (thisModbusClient *сlientModbus) tryConnect(quit chan struct{}, connection chan bool, wg *sync.WaitGroup) { /// connection day out
	defer wg.Done()
	timer := time.NewTimer(1)

	for {
		select {
		case <-quit:
			return
		case <-timer.C:
			timer.Stop()
			for i := 1; i <= int(thisModbusClient.connectionAttempts); i++ {
				select {
				case <-quit:
					return
				default:
					if err := thisModbusClient.checkConnect(); err == nil {
						if isChanged := thisModbusClient.setState(GOOD); isChanged {
							thisModbusClient.log.Write(logger.INFO, thisModbusClient.state, "Подключенно")
						}
						connection <- true
						return
					}
				}
			}
			if isChanged := thisModbusClient.setState(BAD); isChanged {
				thisModbusClient.log.Write(logger.WARNING, thisModbusClient.state, "Не удалось подключиться")
			}

			timer.Reset(thisModbusClient.connectionTimeout)
		}
	}
}
