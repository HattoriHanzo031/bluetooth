package main

import (
	"tinygo.org/x/bluetooth"
	"tinygo.org/x/bluetooth/rawterm"
)

var (
	serviceUUID = bluetooth.ServiceUUIDNordicUART
	rxUUID      = bluetooth.CharacteristicUUIDUARTRX
	txUUID      = bluetooth.CharacteristicUUIDUARTTX
)

func main() {
	println("starting")
	adapter := bluetooth.DefaultAdapter
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Moonboar",
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	}))
	must("start adv", adv.Start())

	var rxChar bluetooth.Characteristic
	var txChar bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &rxChar,
				UUID:   rxUUID,
				Flags:  bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					//txChar.Write(value)
					for _, c := range value {
						rawterm.Putchar(c)
					}
				},
			},
			{
				Handle: &txChar,
				UUID:   txUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))

	rawterm.Configure()
	defer rawterm.Restore()
	print("NUS console enabled, use Ctrl-X to exit\r\n")
	var line []byte
	for {
		ch := rawterm.Getchar()
		rawterm.Putchar(ch)
		line = append(line, ch)

		// Send the current line to the central.
		if ch == '\x18' {
			// The user pressed Ctrl-X, exit the terminal.
			break
		} else if ch == '\n' {
			// Reset the slice while keeping the buffer in place.
			line = line[:0]
			// TODO: parse custom commands
		}
	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
