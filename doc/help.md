wget  https://github.com/ekzhang/bore/releases/download/v0.5.2/bore-v0.5.2-x86_64-unknown-linux-musl.tar.gz
tar xvf bore-v0.5.2-x86_64-unknown-linux-musl.tar.gz 
./bore local 80 --to bore.pub
apt install iproute2
apt install openssh-server
ss -l
service ssh start
ss -l
./bore local 22 --to bore.pub