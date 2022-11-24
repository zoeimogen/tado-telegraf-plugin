# Tado Telegraf Plugin

A [Telegraf](https://github.com/influxdata/telegraf) plugin to gather zone temperature settings and current readings from [Tado](https://www.tado.com/)

### Configuration

This is an [external plugin](https://github.com/influxdata/telegraf/blob/master/docs/EXTERNAL_PLUGINS.md)
which has to be integrated via Telegraf's [excecd plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd).

To use it, first create a config file such as **/etc/telegraf/tado.conf**:
```toml
[[inputs.tado]]
	# Generic Tado client ID and secret from https://my.tado.com/webapp/env.js
	client_id     = "tado-web-app"
	client_secret = "wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc"
    # Your username and password, as used to log in to the app/web site
	tado_username = "user@example.com"
	tado_password = "correct horse battery staple"

```

Then add the following section to your **telegraf.conf**
```toml
[[inputs.execd]]
  command = ["/usr/local/bin/telegraf/tado-telegraf-plugin", "-config", "/etc/telegraf/tado.conf", "-poll_interval", "120s"]
  signal = "none"
```

### Metrics

- tado
  - tags:
    - home - The name of the home
    - zone - The name of the zone
  - fields:
    - setting - Zone temperature setting in degrees Celsius
    - temperature - Current zone temperature in degrees Celsius
    - humidity - Current zone humidity as a percentage

### Example Output

```
tado,home=Our\ House,zone=Kitchen setting=16,temperature=18.77,humidity=47.9 1669328180677550000
tado,home=Our\ House,zone=Bedroom setting=18,temperature=18.44,humidity=52.9 1669328180761980000
tado,home=Our\ House,zone=Office setting=16,temperature=16.32,humidity=55.4 1669328180842656000
tado,home=Our\ House,zone=Living\ Room humidity=56.2,setting=16,temperature=15.82 1669328180917656000

```

### License
This project is subject to the the MIT License. See LICENSE information for details.