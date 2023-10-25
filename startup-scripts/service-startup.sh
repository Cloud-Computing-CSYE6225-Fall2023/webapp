#!/bin/bash


sudo groupadd ec2-user
sudo useradd -s /bin/false -g ec2-user ec2-user

sudo cp ./startup-scripts/webapp.service /etc/systemd/system

sudo chown -R ec2-user:ec2-user /home/webapp
sudo chown -R ec2-user:ec2-user /etc/systemd/system/webapp.service
sudo u+x /etc/systemd/system/webapp.service

sudo systemctl daemon-reload
sudo -u ec2-user systemctl enable webapp
sudo -u ec2-user systemctl start webapp
sudo -u ec2-user systemctl restart webapp
