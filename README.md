##Homekit-daemon
A cli tool to expose sensor data stored in InfluxDB to apple home.
> Based on github.com/brutella/hc all credits go to him

##Running as a daemon

Create the file ``/etc/systemd/system/homekit-daemon.service`` with the following content:

    [Unit]
    Description=A deamon to expose local influxdb data to apple home
    Documentation=https://github.com/r0bot/homekit-daemon
    Wants=network.target
    After=network.target
    
    [Service]
    Type=simple
    DynamicUser=yes
    ExecStart=/home/pi/projects/homekitIAQ/homekit-daemon
    WorkingDirectory=/home/pi/projects/homekitIAQ
    Restart=always
    RestartSec=3
    
    [Install]
    WantedBy=multi-user.target
    
And then execute

    # sudo systemctl daemon-reload
    # sudo systemctl start homekit-deamon
    # sudo systemctl enable homekit-deamon
    # sudo journalctl -f -u homekit-deamon
    
