package bootstraps

import "fmt"

// This initial template will be modifiable based upon the build requirements
const kickstart = `
install
cdrom
lang en_US.UTF-8
keyboard us
unsupported_hardware
network --bootproto=dhcp --hostname centos-7.pelmet.loc
rootpw vagrant
firewall --disabled
selinux --permissive
timezone Europe/Prague
unsupported_hardware
bootloader --location=mbr
text
skipx
zerombr
clearpart --all --initlabel

#Disk partitioning information
part /boot --fstype ext4 --size=2048
part swap  --asprimary   --size=8192
part /     --fstype ext4 --size=1 --grow

auth --enableshadow --passalgo=sha512 --kickstart
firstboot --disabled
eula --agreed
services --enabled=NetworkManager,sshd
reboot
user --name=vagrant --plaintext --password vagrant --groups=vagrant,wheel

repo --name=base --baseurl=http://mirror.centos.org/centos/7.3.1611/os/x86_64/
repo --name=epel-release --baseurl=http://anorien.csc.warwick.ac.uk/mirrors/epel/7/x86_64/
repo --name=elrepo-kernel --baseurl=http://elrepo.org/linux/kernel/el7/x86_64/
repo --name=elrepo-release --baseurl=http://elrepo.org/linux/elrepo/el7/x86_64/
repo --name=elrepo-extras --baseurl=http://elrepo.org/linux/extras/el7/x86_64/

%packages --ignoremissing --excludedocs
@Base
@Core
@Development Tools
kernel-ml
kernel-ml-devel
kernel-ml-tools
kernel-ml-tools-libs
kernel-ml-headers
openssh-clients
expect
make
perl
patch
dkms
gcc
bzip2
sudo
openssl-devel
readline-devel
zlib-devel
net-tools
vim
wget
curl
rsync
epel-release
ansible
libselinux-python
-abrt-libs
-abrt-tui
-abrt-cli
-abrt
-abrt-addon-python
-abrt-addon-ccpp
-abrt-addon-kerneloops
-kernel
-kernel-devel
-kernel-tools-libs
-kernel-tools
-kernel-headers
-aic94xx-firmware
-atmel-firmware
-b43-openfwwf
-bfa-firmware
-ipw2100-firmware
-ipw2200-firmware
-ivtv-firmware
-iwl100-firmware
-iwl105-firmware
-iwl135-firmware
-iwl1000-firmware
-iwl2000-firmware
-iwl2030-firmware
-iwl3160-firmware
-iwl3945-firmware
-iwl4965-firmware
-iwl5000-firmware
-iwl5150-firmware
-iwl6000-firmware
-iwl6000g2a-firmware
-iwl6000g2b-firmware
-iwl6050-firmware
-iwl7260-firmware
-libertas-usb8388-firmware
-libertas-sd8686-firmware
-libertas-sd8787-firmware
-ql2100-firmware
-ql2200-firmware
-ql23xx-firmware
-ql2400-firmware
-ql2500-firmware
-rt61pci-firmware
-rt73usb-firmware
-xorg-x11-drv-ati-firmware
-zd1211-firmware
-iprutils
-fprintd-pam
-intltool

# unnecessary firmware
-aic94xx-firmware
-atmel-firmware
-b43-openfwwf
-bfa-firmware
-ipw2100-firmware
-ipw2200-firmware
-ivtv-firmware
-iwl100-firmware
-iwl1000-firmware
-iwl3945-firmware
-iwl4965-firmware
-iwl5000-firmware
-iwl5150-firmware
-iwl6000-firmware
-iwl6000g2a-firmware
-iwl6050-firmware
-libertas-usb8388-firmware
-ql2100-firmware
-ql2200-firmware
-ql23xx-firmware
-ql2400-firmware
-ql2500-firmware
-rt61pci-firmware
-rt73usb-firmware
-xorg-x11-drv-ati-firmware
-zd1211-firmware
%end

%post
yum update -y
yum install -y sudo
echo "vagrant        ALL=(ALL)       NOPASSWD: ALL" >> /etc/sudoers.d/vagrant
sed -i "s/^.*requiretty/#Defaults requiretty/" /etc/sudoers
/bin/echo 'UseDNS no' >> /etc/ssh/sshd_config
yum clean all

/bin/mkdir /home/vagrant/.ssh
/bin/chmod 700 /home/vagrant/.ssh
/bin/echo -e 'ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ== vagrant insecure public key' > /home/vagrant/.ssh/authorized_keys
/bin/chown -R vagrant:vagrant /home/vagrant/.ssh
/bin/chmod 0400 /home/vagrant/.ssh/*

%end
`

// BuildKickStartConfig - Creates a new presseed configuration using the passed data
func (config *ServerConfig) BuildKickStartConfig() string {
	return fmt.Sprintf("%s%s%s%s%s%s", preseed, preseedDisk, preseedNet, preseedPkg, preseedUsers, preseedCmd)
}
