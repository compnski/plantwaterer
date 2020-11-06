
# Watering

## Emitters
![3d render of emitter](https://csg.tinkercad.com/things/kaySP06Uaiy/t725.png?rev=1568820315916000000&s=&v=1)
[Side-In Waterer](https://www.tinkercad.com/things/kaySP06Uaiy-side-in-waterer)

## Tubing
This project is using 5/16" ID hardware.
I'm currently using Vinyl tubing but it's very inflexible and hard to set up. Next time I set the system up, I plan to try the Silicone based tubing.

## Pump and Water Source
I use a submersible 12v pump, linked below. It delivers 280L/H, which is ~3.9L in the 50 seconds I run it.

I put this in the tank of an auto-filling resivior (toliet) so it has an "infinite" source of water without manual filling.
Preventing the siphon effect is very important with this setup. I have a few holes in the tube just inside the tank to allow in air. This prevents a vacuum from forming and pulling water past the pump after it is stopped.

## Parts I used
* Pump: https://www.amazon.com/gp/product/B00JWJIC0K
* Vinyl Tubing: https://www.amazon.com/gp/product/B000E62TCC
* T-Connector: https://www.amazon.com/gp/product/B0040CRX0Y/
* Silicone Tubing: https://www.amazon.com/Silicon-JoyTube-Silicone-Brewing-Winemaking/dp/B07T8SBX3Y/

## Important Notes
* Anti-Siphon hole


# Pi <-> Arduino link

## Purpose

## Optocoupling
The Rasberry Pi uses 3.3v for TTL serial and the Arduino uses 5v. Additionally, they have separate power sources and may not share a ground. Thus, the serial link between them must be isolated. I used the <PART> for this, along with a transitor because the output is inverted.
There are two of these, one for each direction.
The RST line of the Arduino is also isolated and connected to a GPIO pin on the Pi in case a reset is needed. This uses a slower, but simpler <PART>.
[KiCad schematic]()
[Breadboard]()

## Moisture Sensing
### Multiplexing
### Direct Connection

## Relay Board (Work In Progress)
Eventu
Relay Board:
* Shift Register - https://www.ti.com/lit/ds/symlink/sn74hc595.pdf
* Transistors    - https://www.ti.com/lit/ds/symlink/uln2803a.pdf
* LED Display    - https://optoelectronics.liteon.com/upload/download/DS-30-92-0810/LTA-1000HR.pdf
* Resistor Netwrk- https://www.bourns.com/docs/Product-Datasheets/4600x.pdf
* Relay          - https://content.kemet.com/datasheets/KEM_R7002_EC2_EE2.pdf

