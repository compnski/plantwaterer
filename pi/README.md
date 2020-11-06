
### Serial Port Setup
First, read up on using the serial port of your rasberry pi.
https://www.raspberrypi.org/documentation/configuration/uart.md
I opted to switch the UARTs using /boot/config.txt. it's not clear if that was needed.

It took a bunch of trial and error to get these settings. I pulled them after using the Arduino programmer to flash a program, and they continue to work after the Arduino boots.

``` sh
ls -l /dev/serial0
lrwxrwxrwx 1 root root 7 Oct  7 23:17 /dev/serial0 -> ttyAMA0

stty -F /dev/serial0
speed 115200 baud; line = 0;
eof = <undef>; start = <undef>; stop = <undef>; min = 1; time = 0;
ignbrk -brkint -icrnl -imaxbel
-opost -onlcr
-isig -icanon -iexten -echo -echoe -echok -echoctl -echoke
```

To set all these, use the command below
``` sh
stty -F /dev/serial0 1:0:18b2:0:3:1c:7f:15:0:0:1:0:0:0:1a:0:12:f:17:16:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0:0
```
You can add this to `/etc/rc.local` to have it done on startup. Be careful, that file often ends with `exit 0`, so you can't blindly append to it.
