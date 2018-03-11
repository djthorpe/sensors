
In order to set up RFM69 for mihome, use the following parameters:

```
rfm69 -spi.slave=1 \
  -mode standby -sequencer -modulation fsk \
  -bitrate 4.8 -freq_carrier 434300 -freq_dev 30 \
  -afc_mode on -afc_routine standard \
  -datamode packet -packet_format variable -packet_coding manchester -packet_filter off -packet_crc off \
  -preamble_size 3 -payload_size 66 -sync_word D42D -sync_tol 3 \
  -aes_key "" \
  -fifo_threshold 1
```

The following will then wait for up to 10 seconds for a payload which is captured:

```
rfm69 -spi.slave=1 -mode standby -timeout 10s ClearFIFO ReadPayload
```



