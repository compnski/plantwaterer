#!/bin/bash
/usr/bin/rsync -va /home/pi/serialStream/*.s pi_bastion:serialStream/
/usr/bin/rsync -va /home/pi/captures/ pi_bastion:captures/
