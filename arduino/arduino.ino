#include "DHT.h"
#include <ArduinoJson.h>

#define DHT_SENSOR_PIN 3
#define DHT_SENSOR_TYPE DHT11

DHT dht(DHT_SENSOR_PIN, DHT_SENSOR_TYPE);
StaticJsonDocument<64> doc;

void setup() {
  Serial.begin(9600);
  dht.begin();
}

void loop() {
  if (Serial.available()) {
    String command = Serial.readStringUntil('\n');
    if (command == "requestData") {
      float humi = dht.readHumidity();
      float tempC = dht.readTemperature();
      float tempF = dht.readTemperature(true);

      doc.clear();

      if (!isnan(humi) && !isnan(tempC) && !isnan(tempF)) {
        doc["celsius"] = tempC;
        doc["fahrenheit"] = tempF;
        doc["humidity"] = humi;

        serializeJson(doc, Serial);
        Serial.println();  // Envia a nova linha após o JSON
      } else {
        doc["errorCode"] = 500;
        doc["errorMessage"] = "Erro na leitura do sensor.";
        serializeJson(doc, Serial);
        Serial.println();
      }
    }
  }
}
