#!/bin/bash

sudo groupadd ec2-group
sudo useradd -s /bin/false -g ec2-group ec2-user

sudo cp ./startup-scripts/webapp.service /etc/systemd/system

sudo chown -R ec2-user:ec2-group /home/ec2-user/webapp
sudo chmod 744 /home/ec2-user/webapp

sudo systemctl daemon-reload
sudo systemctl enable webapp
sudo systemctl start webapp
sudo systemctl restart webapp
