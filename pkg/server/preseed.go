package server

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Preseed const, this is the basis for the configuration that will be modified per use-case
const preseedHead = `
# Force debconf priority to critical.
debconf debconf/priority select critical
# Override default frontend to Noninteractive
debconf debconf/frontend select Noninteractive

# Preseeding only locale sets language, country and locale.
d-i debian-installer/locale string en_US

# Disable automatic (interactive) keymap detection.
d-i console-setup/ask_detect boolean false
d-i keyboard-configuration/layoutcode string us

### Clock and time zone setup
d-i clock-setup/utc boolean true
d-i time/zone string Europe/GMT
d-i clock-setup/ntp boolean true
d-i clock-setup/ntp-server string 1.pl.pool.ntp.org

### Preseed Early
d-i preseed/early_command string kill-all-dhcp; netcfg
`

const preseedNet = `
### Network configuration
d-i netcfg/wireless_wep string

# Set network interface or 'auto'
d-i netcfg/choose_interface select %s

# Any hostname and domain names assigned from dhcp take precedence over
d-i netcfg/get_gateway string %s
d-i netcfg/get_ipaddress string %s
d-i netcfg/get_nameservers string %s
d-i netcfg/get_netmask string %s
d-i netcfg/use_dhcp string
d-i netcfg/disable_dhcp boolean true

d-i netcfg/get_hostname string ubuntu
d-i netcfg/get_domain string internal

d-i netcfg/hostname string %s`

const preseedDisk = `
### Partitions
d-i partman/mount_style select label

### Boot loader installation
d-i grub-installer/only_debian boolean true
d-i grub-installer/with_other_os boolean true

### Finishing up the installation
d-i finish-install/reboot_in_progress note
d-i cdrom-detect/eject boolean true

### Preseeding other packages
popularity-contest popularity-contest/participate boolean false

### GRUB
grub-pc grub-pc/hidden_timeout  boolean true
grub-pc grub-pc/timeout string  0
d-i grub-installer/bootdev string /dev/sda

### Regular, primary partitions
d-i partman-auto/disk string /dev/sda

#d-i partman/alignment string cylinder
d-i partman/confirm_write_new_label boolean true

d-i partman-basicfilesystems/choose_label string gpt
d-i partman-basicfilesystems/default_label string gpt

d-i partman-partitioning/choose_label string gpt
d-i partman-partitioning/default_label string gpt
d-i partman/choose_label string gpt
d-i partman/default_label string gpt
#d-i partman-partitioning/confirm_write_new_label boolean true
d-i partman-basicfilesystems/no_swap boolean false
d-i partman/choose_partition select finish
d-i partman/confirm boolean true
d-i partman/confirm_nooverwrite boolean true

d-i partman-auto/method string regular

d-i partman-auto/choose_recipe select parlayfs
d-i partman-auto/expert_recipe string         \
   parlayfs ::                      \
      1 1 1 free                              \
         $bios_boot{ }                        \
         method{ biosgrub } .                 \
      200 200 200 fat32                       \
         $primary{ }                          \
         method{ efi } format{ } .            \
      512 512 512 ext3                        \
         $primary{ } $bootable{ }             \
         method{ format } format{ }           \
         use_filesystem{ } filesystem{ ext3 } \
         mountpoint{ /boot } .                \
      1000 20000 -1 ext4                      \
         $primary{ }                          \
         method{ format } format{ }           \
         use_filesystem{ } filesystem{ ext4 } \
         mountpoint{ / } .                    \
`

const swap = `      65536 65536 65536 linux-swap            \
$primary{ }                          \
method{ swap } format{ } .`

const noswap = `
partman-basicfilesystems partman-basicfilesystems/no_swap boolean false
`

const preseedUsers = `
### Account setup
d-i passwd/root-login boolean false
d-i passwd/make-user boolean true
d-i passwd/user-fullname string %s
d-i passwd/username string %s

d-i passwd/user-password password %s
d-i passwd/user-password-again password %s
d-i user-setup/allow-password-weak boolean true
d-i user-setup/encrypt-home boolean false
`

const preseedPkg = `
### Apt setup
d-i apt-setup/restricted boolean true
d-i apt-setup/universe boolean true
d-i mirror/http/hostname string %s
d-i mirror/http/directory string %s
d-i mirror/country string manual
d-i mirror/http/proxy string

### Base system installation
d-i base-installer/install-recommends boolean false

### Package selection
tasksel tasksel/first multiselect
tasksel/skip-tasks multiselect server
d-i pkgsel/ubuntu-standard boolean false

# Allowed values: none, safe-upgrade, full-upgrade
d-i pkgsel/upgrade select none
d-i pkgsel/ignore-incomplete-language-support boolean true
d-i pkgsel/include string %s

# Language pack selection
d-i pkgsel/install-language-support boolean false
d-i pkgsel/language-pack-patterns string
d-i pkgsel/language-packs multiselect
# or ...
#d-i pkgsel/language-packs multiselect en, pl

# Policy for applying updates. May be "none" (no automatic updates),
# "unattended-upgrades" (install security updates automatically), or
# "landscape" (manage system with Landscape).
d-i pkgsel/update-policy select unattended-upgrades
d-i pkgsel/updatedb boolean false
`

const preseedCmd = `
d-i preseed/late_command string \
    in-target sed -i 's/^%%sudo.*$/%%sudo ALL=(ALL:ALL) NOPASSWD: ALL/g' /etc/sudoers; \
    in-target /bin/sh -c "echo 'Defaults env_keep += \"SSH_AUTH_SOCK\" >> /etc/sudoers"; \
    in-target mkdir -p /home/ubuntu/.ssh; \
    in-target /bin/sh -c "echo '%s' >> /home/ubuntu/.ssh/authorized_keys"; \
    in-target chown -R ubuntu:ubuntu /home/ubuntu/; \
	in-target chmod -R go-rwx /home/ubuntu/.ssh/authorized_keys;
`

//BuildPreeSeedConfig - Creates a new presseed configuration using the passed data
func (config *HostConfig) BuildPreeSeedConfig() string {

	var key string
	var err error

	// This will populate anything missing from the global configuration
	config.PopulateConfiguration()

	if config.SSHKeyPath != "" {
		key, err = config.ReadKeyFromFile(config.SSHKeyPath)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	var parsedDisk string

	if config.SwapEnable == true {
		parsedDisk = preseedDisk + swap
	} else {
		parsedDisk = preseedDisk + noswap
	}
	parsedNet := fmt.Sprintf(preseedNet, config.Adapter, config.Gateway, config.IPAddress, config.NameServer, config.Subnet, config.ServerName)
	parsedPkg := fmt.Sprintf(preseedPkg, config.RepositoryAddress, config.MirrorDirectory, config.Packages)
	parsedCmd := fmt.Sprintf(preseedCmd, key)
	parsedUsr := fmt.Sprintf(preseedUsers, config.Username, config.Username, config.Password, config.Password)
	return fmt.Sprintf("%s%s%s%s%s%s", preseedHead, parsedDisk, parsedNet, parsedPkg, parsedUsr, parsedCmd)
}
