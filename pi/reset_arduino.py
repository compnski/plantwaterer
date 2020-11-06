# This pulses GPIO2 for 1 second.
# It's used to force the RST line low on the Arduino so it resets.
# The lack of +exec and she-bang is so I don't run it by accident.
import gpiozero
import time
l = gpiozero.LED(2)
l.on()
time.sleep(1)
l.off()
