sudo groupadd ec2-user
sudo useradd -s /bin/false -g ec2-user ec2-user
sudo cp ./startup-scripts/webapp.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable webapp
sudo systemctl start webapp
sudo systemctl restart webapp
