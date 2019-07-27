package server

import (
	"strings"
)

// bootcfgHead const, this is the basis for the configuration that will be modified per use-case for the boot.cfg
// The main modifications with regards to this template is a replace-all on the module and kernel path
const bootcfg67u2 = `bootstate=0
title=Loading Plunder ESXi installer
timeout=5
prefix=
kernel=/b.b00
kernelopt=cdromBoot runweasel
modules=/jumpstrt.gz --- /useropts.gz --- /features.gz --- /k.b00 --- /chardevs.b00 --- /user.b00 --- /procfs.b00 --- /uc_intel.b00 --- /uc_amd.b00 --- /uc_hygon.b00 --- /vmx.v00 --- /vim.v00 --- /sb.v00 --- /s.v00 --- /ata_liba.v00 --- /ata_pata.v00 --- /ata_pata.v01 --- /ata_pata.v02 --- /ata_pata.v03 --- /ata_pata.v04 --- /ata_pata.v05 --- /ata_pata.v06 --- /ata_pata.v07 --- /block_cc.v00 --- /bnxtnet.v00 --- /bnxtroce.v00 --- /brcmfcoe.v00 --- /char_ran.v00 --- /ehci_ehc.v00 --- /elxiscsi.v00 --- /elxnet.v00 --- /hid_hid.v00 --- /i40en.v00 --- /iavmd.v00 --- /igbn.v00 --- /ima_qla4.v00 --- /ipmi_ipm.v00 --- /ipmi_ipm.v01 --- /ipmi_ipm.v02 --- /iser.v00 --- /ixgben.v00 --- /lpfc.v00 --- /lpnic.v00 --- /lsi_mr3.v00 --- /lsi_msgp.v00 --- /lsi_msgp.v01 --- /lsi_msgp.v02 --- /misc_cni.v00 --- /misc_dri.v00 --- /mtip32xx.v00 --- /ne1000.v00 --- /nenic.v00 --- /net_bnx2.v00 --- /net_bnx2.v01 --- /net_cdc_.v00 --- /net_cnic.v00 --- /net_e100.v00 --- /net_e100.v01 --- /net_enic.v00 --- /net_fcoe.v00 --- /net_forc.v00 --- /net_igb.v00 --- /net_ixgb.v00 --- /net_libf.v00 --- /net_mlx4.v00 --- /net_mlx4.v01 --- /net_nx_n.v00 --- /net_tg3.v00 --- /net_usbn.v00 --- /net_vmxn.v00 --- /nfnic.v00 --- /nhpsa.v00 --- /nmlx4_co.v00 --- /nmlx4_en.v00 --- /nmlx4_rd.v00 --- /nmlx5_co.v00 --- /nmlx5_rd.v00 --- /ntg3.v00 --- /nvme.v00 --- /nvmxnet3.v00 --- /nvmxnet3.v01 --- /ohci_usb.v00 --- /pvscsi.v00 --- /qcnic.v00 --- /qedentv.v00 --- /qfle3.v00 --- /qfle3f.v00 --- /qfle3i.v00 --- /qflge.v00 --- /sata_ahc.v00 --- /sata_ata.v00 --- /sata_sat.v00 --- /sata_sat.v01 --- /sata_sat.v02 --- /sata_sat.v03 --- /sata_sat.v04 --- /scsi_aac.v00 --- /scsi_adp.v00 --- /scsi_aic.v00 --- /scsi_bnx.v00 --- /scsi_bnx.v01 --- /scsi_fni.v00 --- /scsi_hps.v00 --- /scsi_ips.v00 --- /scsi_isc.v00 --- /scsi_lib.v00 --- /scsi_meg.v00 --- /scsi_meg.v01 --- /scsi_meg.v02 --- /scsi_mpt.v00 --- /scsi_mpt.v01 --- /scsi_mpt.v02 --- /scsi_qla.v00 --- /shim_isc.v00 --- /shim_isc.v01 --- /shim_lib.v00 --- /shim_lib.v01 --- /shim_lib.v02 --- /shim_lib.v03 --- /shim_lib.v04 --- /shim_lib.v05 --- /shim_vmk.v00 --- /shim_vmk.v01 --- /shim_vmk.v02 --- /smartpqi.v00 --- /uhci_usb.v00 --- /usb_stor.v00 --- /usbcore_.v00 --- /vmkata.v00 --- /vmkfcoe.v00 --- /vmkplexe.v00 --- /vmkusb.v00 --- /vmw_ahci.v00 --- /xhci_xhc.v00 --- /elx_esx_.v00 --- /btldr.t00 --- /esx_dvfi.v00 --- /esx_ui.v00 --- /esxupdt.v00 --- /weaselin.t00 --- /lsu_hp_h.v00 --- /lsu_inte.v00 --- /lsu_lsi_.v00 --- /lsu_lsi_.v01 --- /lsu_lsi_.v02 --- /lsu_lsi_.v03 --- /lsu_smar.v00 --- /native_m.v00 --- /qlnative.v00 --- /rste.v00 --- /vmware_e.v00 --- /vsan.v00 --- /vsanheal.v00 --- /vsanmgmt.v00 --- /tools.t00 --- /xorg.v00 --- /imgdb.tgz --- /imgpayld.tgz
build=
updated=0`

// kickstart67u2 const, this is the template for the actual installation of ESXi
const kickstart67u2 = `accepteula 
install --firstdisk --overwritevmfs 
rootpw P@ssw0rd
reboot
vmserialnum --esx=PUT IN YOUR LICENSE KEY
 
#network configuration 
network --bootproto=dhcp --device=vmnic0
 
# run the following command only on the firstboot
%firstboot --interpreter=busybox
 
 
# enable & start remote ESXi Shell (SSH)
vim-cmd hostsvc/enable_ssh
vim-cmd hostsvc/start_ssh
 
# enable & start ESXi Shell (TSM)
vim-cmd hostsvc/enable_esx_shell
vim-cmd hostsvc/start_esx_shell
 
# enable High Performance
# http://www.virtuallyghetto.com/2012/08/configuring-esxi-power-management.html
esxcli system settings advanced set --option=/Power/CpuPolicy --string-value="High Performance" 
 
# supress ESXi Shell shell warning - Thanks to Duncan (http://www.yellow-bricks.com/2011/07/21/esxi-5-suppressing-the-localremote-shell-warning/)
esxcli system settings advanced set -o /UserVars/SuppressShellWarning -i 1
 
#Disable ipv6
esxcli network ip set --ipv6-enabled=0
 
# NTP Configuration (thanks to http://www.virtuallyghetto.com)
cat > /etc/ntp.conf << __NTP_CONFIG__
restrict default kod nomodify notrap noquerynopeer
restrict 127.0.0.1
server 129.6.15.28
server 129.6.15.29
server 129.6.15.30
 
__NTP_CONFIG__
 
/sbin/chkconfig ntpd on
`

//BuildESXiConfig - Creates a new presseed configuration using the passed data
func (config *HostConfig) BuildESXiConfig(mountPath string) string {

	parsedBootCfg := strings.Replace(bootcfg67u2, "/", mountPath, -1)

	return parsedBootCfg
}
