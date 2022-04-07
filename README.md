run blockchain node
git clone https://github.com/chandanoodles/scriptblockchaincode.git usr/local/go/src/github.com/scripttoken/script
export SCRIPT_HOME=/usr/local/go/src/github.com/scripttoken/script
#sudo apt-get install build-essential
#sudo snap install go --classic
cd $SCRIPT_HOME
export GO111MODULE=on

make install

cd $SCRIPT_HOME

cp -r ./integration/scriptnet ../scriptnet
mkdir ~/.scriptcli
cp -r ./integration/scriptnet/scriptcli/* ~/.scriptcli/

Sudo chmod 700 ~/.scriptcli/keys/encrypted

Sudo script start --config=../scriptnet/node
