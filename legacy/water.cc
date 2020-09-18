#include "water.h"


#include <SoftwareSerial.h>



const int led_g = 2;
const int led_r = 9;
const int relay_1 = 4;
const int relay_2 = 5;
const int serial_tx = 11;
const int sw_top = 12;
const int sw_bottom = 13;

const int WAIT = 0;
const int WATER = 1;

int state = 0;
unsigned long lastWatered = 0;
unsigned long startedWatering = 0;
unsigned long totalWaterSec = 0;
unsigned long totalTimeSec = 0;
unsigned long lastRun = 0;
const int delayTime = 1000;

const long WATER_TIME_SEC = 60;
const long WAIT_TIME_SEC = (long)3600 * 72;

SoftwareSerial mySerial =  SoftwareSerial(0, serial_tx, 1);
const char clearScreen[ ] = {
  254, 1, 254, 128, 0
};


void setup() {
  // put your setup code here, to run once:
  pinMode(led_r, OUTPUT);
  pinMode(led_g, OUTPUT);
  pinMode(relay_1, OUTPUT);
  pinMode(relay_2, OUTPUT);

  digitalWrite(relay_1, LOW);
  digitalWrite(relay_2, LOW);

  digitalWrite(led_r, LOW);
  digitalWrite(led_g, LOW);

  pinMode(sw_top, INPUT_PULLUP);
  pinMode(sw_bottom, INPUT_PULLUP);

  digitalWrite(serial_tx, LOW);   // Stop bit state for inverted serial
  pinMode(serial_tx, OUTPUT);
  mySerial.begin(9600);    // Set the data rate
  lastRun = millis();
  status();
}

void setStatus(char* msg) {
  mySerial.print(clearScreen);
  mySerial.print(msg);
}

void startWatering() {
  digitalWrite(led_r, HIGH);
  startedWatering = millis();
  digitalWrite(relay_1, HIGH);
  digitalWrite(relay_2, HIGH);
  state = WATER;
}

void stopWatering() {
  digitalWrite(led_r, LOW);
  lastWatered = millis();
  totalWaterSec += max(lastWatered - startedWatering, 1000) / 1000;

  digitalWrite(relay_1, LOW);
  digitalWrite(relay_2, LOW);
  state = WAIT;
}

bool timeToWater() {
  return nextWaterSec() < 0;
}

bool doneWatering() {
  return (millis() > startedWatering + WATER_TIME_SEC * 1000);
}

unsigned long lastWateredSec() {
  return (millis() - lastWatered) / 1000;
}

long nextWaterSec() {
  return WAIT_TIME_SEC - lastWateredSec();
}

void updateScreen() {
  char message[256];
  //Last Watered: xx seconds ago
  //Next Water in: xx seconds
  //Total water time since restart??
  snprintf(message, 256, "Last: %ld:%02ld\xFE\xC0Next: %ld:%02ld", lastWateredSec() / 60, lastWateredSec() % 60, nextWaterSec() / 60, nextWaterSec() % 60);
  setStatus(message);
}

bool swTopState() {
  return digitalRead(sw_top) == LOW;
}
bool swBottomState() {
  return digitalRead(sw_bottom) == LOW;
}

void status() {
  char message[256];
  snprintf(message, 256, "Cycle: %lds/%02ldh\xFE\xC0Total: %lds/%01d.%01dh", WATER_TIME_SEC, WAIT_TIME_SEC / 3600, totalWaterSec, totalTimeSec / 3600, totalTimeSec / 360);

  setStatus(message);
  delay(1000);
}

void loop() {
  totalTimeSec += max(millis() - lastRun, delayTime) / 1000;
  lastRun = millis();
  // put your main code here, to run repeatedly:
  if (swTopState()) {
    if (swBottomState()) {
      status();
      return;
    }
    startWatering();
  }
  else if (swBottomState()) {
    stopWatering();
  }

  switch (state) {
    case WAIT:
      if (timeToWater()) {
        startWatering();
      }
      break;
    case WATER:
      if (doneWatering()) {
        stopWatering();
      }
      break;
  }
  updateScreen();
  delay(delayTime);
}
