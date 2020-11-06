# plantwaterer

# Plant Waterer v2.
 A robotic plant watering system in two parts.

embeddedwaterer: Tinygo program for a microcontroller that waters on a fixed schedule.
watererserver: Go program for rasberryPI that communicates with the embedded application to report stat and control remotely.
Currently is a standalone parser that emits CSVs with stats data.

## Big Features:
* Soil humidity sensors
* Multiple relays, allows different water circuits
* Read stats / update config remotely
* Webcam for monitoring, takes snapshots

## Arduino
* Turns relays on/off for specified on/off periods.
* Takes commands, sends stats back over serial.
* Monitors humidity sensors, reports back to Pi

### Future Plans
* Humidity based control, water when soil dries past certain threshold
* Store customization data in EEPROM

## Raspberry Pi
* Currently just records serial data to disk for later analysis
* Takes webcam snapshot every 5 minutes

### Future Plans
* Simple go service, hosts website.
* Polls Arduino via serial link on some interval. Keeps state of last values to show deltas.
* If Arduino resets, Pi monitor should automatically reset the water on/off durations if they’ve been changed.
* Awair integration to pull in room humidity/temp/co2 levels

## Serial Link
More details on the linux side in the pi/README.md, and the hardware side in hardware/README.md
In general, standard arduino serial at 115200 baud.


3-byte header, 4-byte request
Request: `CMD <command> <8bit number> <16bit number>`
Response: Echo command back. Stats is special and will then send stats followed by a terminator.

### Commands
| command                	| value 	| args?                                	| Notes                                                	|   	|
|------------------------	|-------	|--------------------------------------	|------------------------------------------------------	|---	|
| water_on               	| 0x02  	| 8bit water id                        	| Resets water timer                                   	|   	|
| water_off              	| 0x03  	| 8bit water id                        	| Resets water time                                    	|   	|
| set_water_on_duration  	| 0x04  	| 8bit water id 16bit duration seconds 	|                                                      	|   	|
| set_water_off_duration 	| 0x05  	| 8bit water id 16bit duration seconds 	|                                                      	|   	|
| set_next_water_at      	| 0x06  	| 8bit sensor id 16bit value           	| Sets wait time to desired value. Useful after resets 	|   	|

To send a command, echo over the serial link.
`echo -e 'CMD\x03\x00\x00\x00' > /dev/serial0`


### Stats Return:
Stats are continually sent via serial every StatsIntervalSeconds, default of 5.

``` json
{"now": 1485775,
"sensors":[
 {"id": 0, "d":[ 1, 37440, 2, 37376, 3, 37376, 4, 37440, 5, 37440 ]},
 {"id": 1, "d":[ 1, 37376, 2, 37440, 3, 37376, 4, 37376, 5, 37376 ]}
],"sections":[
 {"id": 0, "on": false, "next": 1555513, "last": 1296313, 
 "onTime": 50, "offTime": 259200, 
 "onAcc": 306, "offAcc": 1296007 }
]}
```

Notes:
* Number of sensors and length of history is determined at compile-time. It is limited by the memory capacity of your chip. More powerful chips can support more sensors and longer histories.
* All timestamps are relative to the local epoch, when the board was turned on.
* The data in sensors.d are pairs of (timestamp, value) from the analog input. 
* onTime / offTime are current status values
* onAcc / offAcc accumulate time spent in on / off states over time


## References
* https://github.com/huin/goserial
* https://www.digikey.com/product-detail/en/adafruit-industries-llc/997/1528-2003-ND/6827136
* https://tinygo.org/
* https://github.com/d2r2/go-i2c
* http://www.krekr.nl/wp-content/uploads/2013/08/Arduino-uno.pdf
* https://www.raspberrypi.org/documentation/hardware/raspberrypi/mechanical/rpi_MECH_3b_1p2.pdf
* https://www.youtube.com/watch?v=ioSYlxHlYdI (For solenoid board design)
* https://jlcpcb.com/ - Cheap PCB!
* https://www.ti.com/lit/ds/symlink/sn74hc595.pdf
* https://www.ti.com/lit/ds/symlink/uln2803a.pdf?

