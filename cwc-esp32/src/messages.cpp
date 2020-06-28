/*
Copyright (C) 2020 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
#include <sys/time.h>
#include <inttypes.h>
#include <string.h>
#include <WiFi.h>
#include "timestamp.h"
#include "messages.h"

uint64_t
ntoh64(const uint64_t *input)
{
    uint64_t rval;
    uint8_t *data = (uint8_t *)&rval;

    data[0] = *input >> 56;
    data[1] = *input >> 48;
    data[2] = *input >> 40;
    data[3] = *input >> 32;
    data[4] = *input >> 24;
    data[5] = *input >> 16;
    data[6] = *input >> 8;
    data[7] = *input >> 0;

    return rval;
}

uint64_t
hton64(const uint64_t *input)
{
    return (ntoh64(input));
}

/**
 * Send message via UDP
 */
void msg_send(char * pkt, int len) {

}

// Message encode and decode functions

void enumerate_channels() {
    EnumerateChannelsPayload ecp;
    msg_send((char *)&ecp, sizeof(ecp));
}

void time_sync() {
    TimeSyncPayload tsp;

    uint64_t now = timestamp_64_now();
    tsp.current_time = hton64(&now);
     
    msg_send((char *)&tsp, sizeof(tsp));
}

void listen_request(int channel, char * callsign) {
    ListenRequestPayload lrp;
    
    lrp.channel = (ChannelIdType) htons(channel);
    strncpy(lrp.callsign, callsign, 16);

    msg_send((char *)&lrp, sizeof(lrp));
}

void unlisten(int channel, unsigned short key) {
    UnlistenPayload up;
    
    up.channel = (ChannelIdType) htons(channel);
    up.carrier_key = (CarrierKeyType) htons(key);

    msg_send((char *)&up, sizeof(up));
} 

void key_value(char * key, char * value) {
    KeyValuePayload kvp;
    
    strncpy(kvp.key, key, KEY_SIZE);
    strncpy(kvp.value, value, VALUE_SIZE);

    msg_send((char *)&kvp, sizeof(kvp));
}

void carrier_event() {



}