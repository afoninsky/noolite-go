### Protocol specification
https://www.noo.com.by/assets/files/PDF/MTRF-64-USB.pdf

### Running on macOS
`libusb: bad access [code -3]` fix:
```
kextstat | grep -i ftdi
sudo kextunload -b com.FTDI.driver.FTDIUSBSerialDriver
sudo kextunload -b com.apple.driver.AppleUSBFTDI
```

### Running on linux
To have access on device from common user add the next rule to udev. For example to /etc/udev/rules.d/50-noolite.rules next line:
```
ATTRS{idVendor}=="1027", ATTRS{idProduct}=="24577", SUBSYSTEMS=="usb", ACTION=="add", MODE="0666", GROUP="noolite"
```
Then add your user to noolite group:
```
sudo usermod <user> -aG noolite
```

/dev/ttyUSB0
Bus 001 Device 004: ID 0403:6001 Future Technology Devices International, Ltd FT232 USB-Serial (UART) IC

### Misc
https://www.home-assistant.io/components/light.mqtt/

docker run -p 1883:1883 eclipse-mosquitto
mosquitto_sub -v -t "#"
mosquitto_pub -t "home/noolite_test/set" -m ON
mosquitto_pub -h 192.168.1.3 -t "home/noolitef/1/command" -m ON


docker build -t vkfont/pi3-noolite:latest .
docker run \
    -e MQTT_HOST=192.168.1.3:1883 \
    -e DEVICE_PORT=/dev/ttyUSB0 \
    --device /dev/ttyUSB0 \
    vkfont/pi3-noolite:latest