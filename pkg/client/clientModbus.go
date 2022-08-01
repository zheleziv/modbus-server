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
	connectionType     string
	ip                 string
	port               int
	slaveId            uint8
	connectionAttempts int
	connectionTimeout  time.Duration
	log                logger.Logger
	state              string
	tags               []tag.TagInterface
	sender             modbus.Client
}

func (c *сlientModbus) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string
		Ip      string
		SlaveId uint8
		State   string
		Tags    []tag.TagInterface
	}{
		Name:    c.name,
		Ip:      c.ip,
		SlaveId: c.slaveId,
		State:   c.state,
		Tags:    c.tags,
	})
}

// ======================инициализация========================{

// =================================Name
func (thisModbusClient *сlientModbus) Name() string {
	return thisModbusClient.name
}
func (thisModbusClient *сlientModbus) SetName(name string) error {
	if strings.TrimSpace(name) == "" {
		return myerr.New("invalid client name. {setter client name}")
	}
	thisModbusClient.name = name
	return nil
}

// =================================Type
func (thisModbusClient *сlientModbus) Type() string {
	return thisModbusClient.connectionType
}
func (thisModbusClient *сlientModbus) SetType() {
	thisModbusClient.connectionType = MODBUS_TCP
}

// ===================================IP
func (thisModbusClient *сlientModbus) Ip() string {
	return thisModbusClient.ip
}

// using net.parseIp
func (thisModbusClient *сlientModbus) SetIp(ip string) error {
	ipAddr := net.ParseIP(strings.TrimSpace(ip))
	if ipAddr == nil {
		return myerr.New("invalid client Ip. {setter client Ip}")
	} else {
		thisModbusClient.ip = ip
		return nil
	}
}

// ===================ConnectionAttempts
func (thisModbusClient *сlientModbus) ConnectionAttempts() int {
	return thisModbusClient.connectionAttempts
}
func (thisModbusClient *сlientModbus) SetConnectionAttempts(ca int) error {
	if ca <= 0 {
		return myerr.New("invalid client connection attempts. {setter client connection attempts}")
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

// func (c *сlientModbus) SetTags(tags []tag.TagInterface) error {
// 	for id := range tags {
// 		if _, err := c.TagByName(tags[id].Name()); err != nil {
// 			return err
// 		}
// 	}
// 	c.tags = tags
// 	return nil
// }

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
		adr16, err := checkModbusAddress(address, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}
		t, err := tag.NewTag(name, adr16, scanPeriod, dataType)
		if err != nil {
			return myerr.New(err.Error())
		}
		thisModbusClient.tags = append(thisModbusClient.tags, t)
		return nil
	}
	return myerr.New("client tag name already exists")
}

func checkModbusAddress(address uint32, dataType string) (uint16, error) {
	if address >= tag.UINT16_MAX_VALUE {
		t := FUNCTION__TAG_TYPE[dataType]
		if t == nil {
			return 0, myerr.New("invalid tag data type")
		}

		tmpINT := int(address / 100000.0)
		isDefined := false
		for i := range t {
			if t[i] == tmpINT {
				isDefined = isDefined || true
			}
		}
		if !isDefined {
			return 0, myerr.New("invalid tag address")
		}

		tmpUINT16 := uint16(address - uint32(tmpINT*100000))
		if tmpUINT16 >= tag.UINT16_MAX_VALUE {
			return 0, myerr.New("invalid tag address")
		}
		return tmpUINT16 - 1, nil
	}
	return uint16(address - 1), nil
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
func NewClinetModbus(ip string, port int, slaveID uint8, name string, debug bool, ConnectionAttempts int, ConnectionTimeout float64) (*сlientModbus, error) {
	var c сlientModbus

	err := c.SetName(name)
	if err != nil {
		return nil, myerr.New(err.Error())
	}
	c.SetType()
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
		IsLogOutput:    debug,
	}

	return &c, nil
}

//======================инициализация========================}

func (thisModbusClient *сlientModbus) Connect() error {
	provider := modbus.NewTCPClientProvider(
		thisModbusClient.ip + ":" + fmt.Sprint(thisModbusClient.port))
	thisModbusClient.sender = modbus.NewClient(provider)
	defer thisModbusClient.sender.Close()

	//устанавливаем соединение
	err := thisModbusClient.sender.Connect()
	if err != nil {
		return myerr.New("Ошибка соединения! " + err.Error())
	} else {
		return nil
	}
}

func (thisModbusClient *сlientModbus) Send(id int) error {
	if thisModbusClient.sender == nil {
		thisModbusClient.tags[id].SetState(false)
		return myerr.New("sender nil")
	}
	defer thisModbusClient.sender.Close()

	//новые типы должны быть указаны здесь
	switch thisModbusClient.tags[id].DataType() {
	case tag.COIL_TYPE:
		{
			resp, err := thisModbusClient.sender.ReadDiscreteInputs(thisModbusClient.slaveId, thisModbusClient.tags[id].Address(), 1)

			if err != nil {
				thisModbusClient.tags[id].SetState(false)
				return myerr.New(err.Error())
			}

			if len(resp) > 0 {
				thisModbusClient.tags[id].(*tag.CoilTag).SetValue(resp[0])
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(resp[0]))+".")
				return nil
			} else {
				thisModbusClient.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	case tag.WORD_TYPE:
		{
			resp, err := thisModbusClient.sender.ReadHoldingRegisters(thisModbusClient.slaveId, thisModbusClient.tags[id].Address(), 1)

			if err != nil {
				thisModbusClient.tags[id].SetState(false)
				return myerr.New(err.Error())
			}

			if len(resp) > 0 {
				thisModbusClient.tags[id].(*tag.WordTag).SetValue(resp[0])
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(resp[0]))+".")
				return nil
			} else {
				thisModbusClient.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	case tag.DWORD_TYPE:
		{
			resp, err := thisModbusClient.sender.ReadHoldingRegisters(thisModbusClient.slaveId, thisModbusClient.tags[id].Address(), 2)

			if err != nil {
				thisModbusClient.tags[id].SetState(false)
				return myerr.New(err.Error())
			}

			if len(resp) > 1 {
				var tmp uint32 = (uint32(resp[0]) << 16) + uint32(resp[1])
				thisModbusClient.tags[id].(*tag.DWordTag).SetValue(tmp)
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[id].Name(), "Значение: "+strconv.Itoa(int(tmp))+".")
				return nil
			} else {
				thisModbusClient.tags[id].SetState(false)
				// c.log.WriteWithTag(logger.WARNING, c.state, c.tags[id].Name(), "Значение не было считано!")
				return nil
			}
		}
	default:
		return myerr.New("resp nil")
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
			{
				close(quit)
				wgi.Wait()
				return
			}
		case cb := <-connection: // канал снизу. Плохое подключение => реконект
			{
				if cb {
					close(quit)
					wgi.Wait()
					quit = make(chan struct{})
					for tagId := range thisModbusClient.tags {
						wgi.Add(1)
						go thisModbusClient.startSender(tagId, quit, &wgi, connection)
					}
				} else {
					close(quit)
					wgi.Wait()
					quit = make(chan struct{})
					wgi.Add(1)
					go thisModbusClient.tryConnect(stop, connection, &wgi)
				}
			}
		}
	}
}

func (thisModbusClient *сlientModbus) startSender(tagId int, quit chan struct{}, wg *sync.WaitGroup, connect chan bool) {
	// c.log.WriteWithTag(logger.INFO, c.state, c.tags[tagId].Name(), "Запущен опрос тега!")

	ticker := time.NewTicker(time.Duration(thisModbusClient.tags[tagId].ScanPeriod()) * time.Second)

	defer wg.Done()

	for {
		select {
		case <-quit:
			{
				// c.log.WriteWithTag(logger.INFO, c.state, c.tags[tagId].Name(), "Завершен опрос тега!")
				return
			}
		case <-ticker.C:
			{
				err := thisModbusClient.Send(tagId)
				if err != nil {
					thisModbusClient.log.WriteWithTag(logger.ERROR, thisModbusClient.state, thisModbusClient.tags[tagId].Name(), err.Error())
					connect <- false
				}
			}
		}
	}
}

func (thisModbusClient *сlientModbus) tryConnect(quit chan struct{}, connection chan bool, wg *sync.WaitGroup) { /// connection day out
	defer wg.Done()
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-quit:
			{
				return
			}
		case <-ticker.C:
			{
				ticker.Stop()
				for i := 1; i <= thisModbusClient.connectionAttempts; i++ {
					select {
					case <-quit:
						{
							return
						}
					default:
						{
							err := thisModbusClient.Connect()
							if err == nil {
								if isChanged := thisModbusClient.setState(GOOD); isChanged {
									thisModbusClient.log.Write(logger.INFO, thisModbusClient.state, "Подключенно!")
								}
								connection <- true
								return
							}
						}
					}
				}
				if isChanged := thisModbusClient.setState(BAD); isChanged {
					thisModbusClient.log.Write(logger.WARNING, thisModbusClient.state, "Не удалось подключиться!")
				}

				ticker.Reset(thisModbusClient.connectionTimeout)
			}
		}
	}
}
