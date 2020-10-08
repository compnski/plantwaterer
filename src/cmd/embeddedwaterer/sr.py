#!/usr/bin/env python

import serial
readMode = False
writeMode = True
with serial.Serial('/dev/ttyS1', 115200, bytesize=8, parity='N', stopbits=1, timeout=1) as ser:
    if writeMode:
        ser.write("Hi")
    if readMode:
        while True:
            data = ser.read()
            if data != b"\0":
                print "%s" % data,


