// Copyright © 2015 Matthias Benkort <matthias.benkort@gmail.com>
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

// ------------------------------------------------------------
// ------------------------- Parse (raw []byte) (error, Packet)
// ------------------------------------------------------------

// Parse() with valid raw data and no payload (PUSH_ACK)
func TestParse1(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_ACK}
	packet, err := Parse(raw)

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PUSH_ACK {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload != nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.GatewayId != nil {
		t.Errorf("Invalid parsed gateway id: % x", packet.GatewayId)
	}
}

// Parse() with valid raw data and stat payload
func TestParse2(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`{
        "stat": {
            "time":"2014-01-12 08:59:28 GMT",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)
	packet, err := Parse(append(append(raw, gatewayId...), payload...))

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if !bytes.Equal(gatewayId, packet.GatewayId) {
		t.Errorf("Invalid parsed gatewayId: % x", packet.GatewayId)
	}

	if packet.Identifier != PUSH_DATA {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload == nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
		return
	}

	if !bytes.Equal(payload, packet.Payload.Raw) {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.Payload.Stat == nil {
		t.Errorf("Invalid parsed payload Stat: %#v", packet.Payload.Stat)
	}

	statTime, _ := time.Parse(time.RFC3339, "2014-01-12T08:59:28.000Z")

	stat := Stat{
		Time: statTime,
		Lati: 46.24000,
		Long: 3.25230,
		Alti: 145,
		Rxnb: 2,
		Rxok: 2,
		Rxfw: 2,
		Ackr: 100.0,
		Dwnb: 2,
		Txnb: 2,
	}

	if !reflect.DeepEqual(stat, *packet.Payload.Stat) {
		t.Errorf("Invalid parsed payload Stat: %#v", packet.Payload.Stat)
	}

}

// Parse() with valid raw data and rxpk payloads
func TestParse3(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`{
        "rxpk":[
            {
                "chan":2,
                "codr":"4/6",
                "data":"-DS4CGaDCdG+48eJNM3Vai-zDpsR71Pn9CPA9uCON84",
                "datr":"SF7BW125",
                "freq":866.349812,
                "lsnr":5.1,
                "modu":"LORA",
                "rfch":0,
                "rssi":-35,
                "size":32,
                "stat":1,
                "time":"2013-03-31T16:21:17.528002Z",
                "tmst":3512348611
            },{
                "chan":9,
                "data":"VEVTVF9QQUNLRVRfMTIzNA==",
                "datr":50000,
                "freq":869.1,
                "modu":"FSK",
                "rfch":1,
                "rssi":-75,
                "size":16,
                "stat":1,
                "time":"2013-03-31T16:21:17.530974Z",
                "tmst":3512348514
            },{
                "chan":0,
                "codr":"4/7",
                "data":"ysgRl452xNLep9S1NTIg2lomKDxUgn3DJ7DE+b00Ass",
                "datr":"SF10BW125",
                "freq":863.00981,
                "lsnr":5.5,
                "modu":"LORA",
                "rfch":0,
                "rssi":-38,
                "size":32,
                "stat":1,
                "time":"2013-03-31T16:21:17.532038Z",
                "tmst":3316387610
            }
        ]
    }`)

	packet, err := Parse(append(append(raw, gatewayId...), payload...))

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PUSH_DATA {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if !bytes.Equal(gatewayId, packet.GatewayId) {
		t.Errorf("Invalid parsed gatewayId: % x", packet.GatewayId)
	}

	if packet.Payload == nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
		return
	}

	if !bytes.Equal(payload, packet.Payload.Raw) {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.Payload.RXPK == nil || len(*packet.Payload.RXPK) != 3 {
		t.Errorf("Invalid parsed payload RXPK: %#v", packet.Payload.RXPK)
	}

	rxpkTime, _ := time.Parse(time.RFC3339, "2013-03-31T16:21:17.530974Z")

	rxpk := RXPK{
		Chan: 9,
		Data: "VEVTVF9QQUNLRVRfMTIzNA==",
		Datr: "50000",
		Freq: 869.1,
		Modu: "FSK",
		Rfch: 1,
		Rssi: -75,
		Size: 16,
		Stat: 1,
		Time: rxpkTime,
		Tmst: 3512348514,
	}

	if !reflect.DeepEqual(rxpk, (*packet.Payload.RXPK)[1]) {
		t.Errorf("Invalid parsed payload RXPK: %#v", packet.Payload.RXPK)
	}
}

// Parse() with valid raw data and rxpk payloads + stat
func TestParse4(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`{
        "rxpk":[
            {
                "chan":2,
                "codr":"4/6",
                "data":"-DS4CGaDCdG+48eJNM3Vai-zDpsR71Pn9CPA9uCON84",
                "datr":"SF7BW125",
                "freq":866.349812,
                "lsnr":5.1,
                "modu":"LORA",
                "rfch":0,
                "rssi":-35,
                "size":32,
                "stat":1,
                "time":"2013-03-31T16:21:17.528002Z",
                "tmst":3512348611
            },{
                "chan":9,
                "data":"VEVTVF9QQUNLRVRfMTIzNA==",
                "datr":50000,
                "freq":869.1,
                "modu":"FSK",
                "rfch":1,
                "rssi":-75,
                "size":16,
                "stat":1,
                "time":"2013-03-31T16:21:17.530974Z",
                "tmst":3512348514
            },{
                "chan":0,
                "codr":"4/7",
                "data":"ysgRl452xNLep9S1NTIg2lomKDxUgn3DJ7DE+b00Ass",
                "datr":"SF10BW125",
                "freq":863.00981,
                "lsnr":5.5,
                "modu":"LORA",
                "rfch":0,
                "rssi":-38,
                "size":32,
                "stat":1,
                "time":"2013-03-31T16:21:17.532038Z",
                "tmst":3316387610
            }
        ],
        "stat": {
            "time":"2014-01-12 08:59:28 GMT",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)

	packet, err := Parse(append(append(raw, gatewayId...), payload...))

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PUSH_DATA {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload == nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
		return
	}

	if !bytes.Equal(payload, packet.Payload.Raw) {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.Payload.RXPK == nil || len(*packet.Payload.RXPK) != 3 {
		t.Errorf("Invalid parsed payload RXPK: %#v", packet.Payload.RXPK)
	}

	if packet.Payload.Stat == nil {
		t.Errorf("Invalid parsed payload Stat: %#v", packet.Payload.Stat)
	}

	rxpkTime, _ := time.Parse(time.RFC3339, "2013-03-31T16:21:17.530974Z")

	rxpk := RXPK{
		Chan: 9,
		Data: "VEVTVF9QQUNLRVRfMTIzNA==",
		Datr: "50000",
		Freq: 869.1,
		Modu: "FSK",
		Rfch: 1,
		Rssi: -75,
		Size: 16,
		Stat: 1,
		Time: rxpkTime,
		Tmst: 3512348514,
	}

	if !reflect.DeepEqual(rxpk, (*packet.Payload.RXPK)[1]) {
		t.Errorf("Invalid parsed payload RXPK: %#v", packet.Payload.RXPK)
	}

	statTime, _ := time.Parse(time.RFC3339, "2014-01-12T08:59:28.000Z")

	stat := Stat{
		Time: statTime,
		Lati: 46.24000,
		Long: 3.25230,
		Alti: 145,
		Rxnb: 2,
		Rxok: 2,
		Rxfw: 2,
		Ackr: 100.0,
		Dwnb: 2,
		Txnb: 2,
	}

	if !reflect.DeepEqual(stat, *packet.Payload.Stat) {
		t.Errorf("Invalid parsed payload Stat: %#v", packet.Payload.Stat)
	}
}

// Parse() with valid raw data and txpk payload
func TestParse5(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PULL_RESP}

	payload := []byte(`{
        "txpk":{
            "imme":true,
            "freq":864.123456,
            "rfch":0,
            "powe":14,
            "modu":"LORA",
            "datr":"SF11BW125",
            "codr":"4/6",
            "ipol":false,
            "size":32,
            "data":"H3P3N2i9qc4yt7rK7ldqoeCVJGBybzPY5h1Dd7P7p8v"
        }
    }`)

	packet, err := Parse(append(raw, payload...))

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PULL_RESP {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload == nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
		return
	}

	if !bytes.Equal(payload, packet.Payload.Raw) {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.Payload.TXPK == nil {
		t.Errorf("Invalid parsed payload TXPK: %#v", packet.Payload.TXPK)
	}

	txpk := TXPK{
		Imme: true,
		Freq: 864.123456,
		Rfch: 0,
		Powe: 14,
		Modu: "LORA",
		Datr: "SF11BW125",
		Codr: "4/6",
		Ipol: false,
		Size: 32,
		Data: "H3P3N2i9qc4yt7rK7ldqoeCVJGBybzPY5h1Dd7P7p8v",
	}

	if !reflect.DeepEqual(txpk, *packet.Payload.TXPK) {
		t.Errorf("Invalid parsed payload TXPK: %#v", packet.Payload.TXPK)
	}
}

// Parse() with an invalid version number
func TestParse6(t *testing.T) {
	raw := []byte{0x00, 0x14, 0x14, PUSH_ACK}
	_, err := Parse(raw)

	if err == nil {
		t.Errorf("Successfully parsed an incorrect version number")
	}
}

// Parse() with an invalid raw message
func TestParse7(t *testing.T) {
	raw1 := []byte{VERSION}
	var raw2 []byte
	_, err1 := Parse(raw1)
	_, err2 := Parse(raw2)

	if err1 == nil {
		t.Errorf("Successfully parsed an raw message")
	}

	if err2 == nil {
		t.Errorf("Successfully parsed a nil raw message")
	}
}

// Parse() with an invalid identifier
func TestParse8(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, 0xFF}
	_, err := Parse(raw)

	if err == nil {
		t.Errorf("Successfully parsed an incorrect identifier")
	}
}

// Parse() with an invalid payload
func TestParse9(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`wrong`)
	_, err := Parse(append(append(raw, gatewayId...), payload...))
	if err == nil {
		t.Errorf("Successfully parsed an incorrect payload")
	}
}

// Parse() with an invalid date
func TestParse10(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`{
        "stat": {
            "time":"null",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)
	_, err := Parse(append(append(raw, gatewayId...), payload...))
	if err == nil {
		t.Errorf("Successfully parsed an incorrect payload time")
	}
}

// Parse() with valid raw data but a useless payload
func TestParse11(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PUSH_ACK}
	payload := []byte(`{
        "stat": {
            "time":"2014-01-12 08:59:28 GMT",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)
	packet, err := Parse(append(raw, payload...))

	if err != nil {
		t.Errorf("Failed to parse a valid PUSH_ACK packet")
	}

	if packet.Payload != nil {
		t.Errorf("Parsed payload on a PUSH_ACK packet")
	}
}

// Parse() with valid raw data but a useless payload
func TestParse12(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PULL_ACK}
	payload := []byte(`{
        "stat": {
            "time":"2014-01-12 08:59:28 GMT",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)
	packet, err := Parse(append(raw, payload...))

	if err != nil {
		t.Errorf("Failed to parse a valid PULL_ACK packet")
	}

	if packet.Payload != nil {
		t.Errorf("Parsed payload on a PULL_ACK packet")
	}
}

// Parse() with valid raw data but a useless payload
func TestParse13(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PULL_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	payload := []byte(`{
        "stat": {
            "time":"2014-01-12 08:59:28 GMT",
            "lati":46.24000,
            "long":3.25230,
            "alti":145,
            "rxnb":2,
            "rxok":2,
            "rxfw":2,
            "ackr":100.0,
            "dwnb":2,
            "txnb":2
        }
    }`)
	packet, err := Parse(append(append(raw, gatewayId...), payload...))

	if err != nil {
		t.Errorf("Failed to parse a valid PULL_DATA packet")
	}

	if packet.Payload != nil {
		t.Errorf("Parsed payload on a PULL_DATA packet")
	}
}

// Parse() with valid raw data and no payload (PULL_ACK)
func TestParse14(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PULL_ACK}
	packet, err := Parse(raw)

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PULL_ACK {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload != nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if packet.GatewayId != nil {
		t.Errorf("Invalid parsed gateway id: % x", packet.GatewayId)
	}
}

// Parse() with valid raw data and no payload (PULL_DATA)
func TestParse15(t *testing.T) {
	raw := []byte{VERSION, 0x14, 0x14, PULL_DATA}
	gatewayId := []byte("qwerty1234")[0:8]
	packet, err := Parse(append(raw, gatewayId...))

	if err != nil {
		t.Errorf("Failed to parse with error: %#v", err)
	}

	if packet.Version != VERSION {
		t.Errorf("Invalid parsed version: %x", packet.Version)
	}

	if !bytes.Equal([]byte{0x14, 0x14}, packet.Token) {
		t.Errorf("Invalid parsed token: %x", packet.Token)
	}

	if packet.Identifier != PULL_DATA {
		t.Errorf("Invalid parsed identifier: %x", packet.Identifier)
	}

	if packet.Payload != nil {
		t.Errorf("Invalid parsed payload: % x", packet.Payload)
	}

	if !bytes.Equal(gatewayId, packet.GatewayId) {
		t.Errorf("Invalid parsed gateway id: % x", packet.GatewayId)
	}
}