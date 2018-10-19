package main

const clientID = "noolite"

const modeTX = "tx"
const modeRX = "rx"
const modeFTX = "ftx"
const modeFRX = "frx"

const willTopicPattern = "%s/state"
const willOnlineMessage = "online"
const willOfflineMessage = "offline"

const setTopicPattern = "%s/%s/%s/command"
const stateTopicPattern = "%s/%s/%s/state"

/**
noolite/state: online|offline
noolite/{mode}/{channel}/command: bind|on|off|brightness|rgb
noolite/{mode}/{channel}/state: on|off

{mode}: tx|rx|ftx|frx
{channel}: 1..64
*/
