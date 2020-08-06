# plantwaterer

# Plant Waterer v2.
 A robotic plant watering system in two parts.

embeddedwaterer: Tinygo program for a microcontroller that waters on a fixed schedule.
watererserver: Go program for rasberryPI that communicates with the embedded application to report stat and control remotely.


## Big Features:
* Soil humidity sensors
* Multiple relays, allows different water circuits
* Read stats / update config via web

## Maybe Features:
* Webcam for monitoring, maybe time-lapse?
* Awair integration?
* Humidity based control, water when soil dries past certain threshold


## Arduino
Turns relays on/off for specified on/off periods.
Takes commands, sends stats back over serial.
Monitors humidity sensors, reports back to Pi
Store customization data in EEPROM?

## Raspberry Pi
Simple go service, hosts website.
Polls Arduino via serial link on some interval. Keeps state of last values to show deltas.
If Arduino resets, Pi should automatically reset the water on/off durations if they’ve been changed.

## Serial Link
4-byte request
Request: `<command> <8bit number> <16bit number>`
Response: Echo command back. Stats is special and will then send stats followed by a terminator.

Commands
| command                	| value 	| args?                                	| Notes                                                	|   	|
|------------------------	|-------	|--------------------------------------	|------------------------------------------------------	|---	|
| stats                  	| 0x01  	| -                                    	|                                                      	|   	|
| water_on               	| 0x02  	| 8bit water id                        	| Resets water timer                                   	|   	|
| water_off              	| 0x03  	| 8bit water id                        	| Resets water time                                    	|   	|
| set_water_on_duration  	| 0x04  	| 8bit water id 16bit duration seconds 	|                                                      	|   	|
| set_water_off_duration 	| 0x05  	| 8bit water id 16bit duration seconds 	|                                                      	|   	|
| set_next_water_at      	| 0x06  	| 8bit sensor id 16bit value           	| Sets wait time to desired value. Useful after resets 	|   	|
| zero_stats             	| 0x07  	| -                                    	|                                                      	|   	|


Stats Return
Also 4-byte frames
| name                  	| value 	| args                                 	| notes                                       	|   	|
|-----------------------	|-------	|--------------------------------------	|---------------------------------------------	|---	|
| total_water_secs      	| 0x01  	| 8bit water id 16bit duration seconds 	| Per-water loop                              	|   	|
| total_water_days      	| 0x02  	| 8bit water id 16bit duration seconds 	| I think numbers are max 16bit               	|   	|
| total_wait_secs       	| 0x03  	| 8bit water id 16bit duration days    	|                                             	|   	|
| total_wait_days       	| 0x04  	| 8bit water id 16bit duration days    	| I think numbers are max 16bit               	|   	|
| since_last_water_secs 	| 0x05  	| 8bit water id 16bit duration seconds 	| Per-water loop                              	|   	|
| current_water_secs    	| 0x06  	| 8bit water id 16bit duration seconds 	| Per-water loop, 0 if not currently watering 	|   	|
| moisture_current      	| 0x07  	| 8bit sensor id 16bit value           	| Per sensor                                  	|   	|
| moisture_max          	| 0x08  	| 8bit sensor id 16bit value           	| Per sensor                                  	|   	|
| moisture_min          	| 0x09  	| 8bit sensor id 16bit value           	| Per sensor                                  	|   	|
|                       	|       	|                                      	|                                             	|   	|



References
* https://github.com/huin/goserial
* https://www.digikey.com/product-detail/en/adafruit-industries-llc/997/1528-2003-ND/6827136
* https://tinygo.org/
* https://github.com/d2r2/go-i2c
* http://www.krekr.nl/wp-content/uploads/2013/08/Arduino-uno.pdf
* https://www.raspberrypi.org/documentation/hardware/raspberrypi/mechanical/rpi_MECH_3b_1p2.pdf
* https://www.youtube.com/watch?v=ioSYlxHlYdI (For solenoid board design)
