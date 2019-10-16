Windows 10 Pro - HyperV enabled
PowerShell with admin right

1) Verifier les ports que windows se réserve...
PS> C:\Program Files\Docker\Docker> netsh interface ipv4 show excludedportrange protocol=tcp

If you see that one of port ranges include port 2375 then you have the same issue.
Protocole tcp Plages d’exclusion de ports

Port de début    Port de fin
-------------    -----------
        80          80
      1587        1686
      1787        1886
      1887        1986
      2180        2279
      2280        2379
      2380        2479
      2480        2579
      2580        2679
      2699        2798
      5357        5357
     11870       11969
     15855       15954
     50000       50059     *

* - Exclusions de ports administrés.

Disable Hyper-V and reboot:
PS> dism.exe /Online /Disable-Feature:Microsoft-Hyper-V

Then reserve port 2375:
netsh int ipv4 add excludedportrange protocol=tcp startport=2375 numberofports=1 store=persistent

Enable Hyper-V and reboot again:

dism.exe /Online /Enable-Feature:Microsoft-Hyper-V /All

`PS> DISM /online /enable-feature /all /featureName:Microsoft-Hyper-V`

PowerShell Manager HyperV (optional) `PS> mmc virtmgmt.msc`

PS> C:\Program Files\Docker\Docker> docker-machine create --driver hyperv --hyperv-cpu-count 2  --hyperv-virtual-switch "Default Switch" default
Running pre-create checks...
Creating machine...
(default) Copying C:\Users\Pascal\.docker\machine\cache\boot2docker.iso to C:\Users\Pascal\.docker\machine\machines\default\boot2docker.iso...
(default) Creating SSH key...
(default) Creating VM...
(default) Using switch "Default Switch"
(default) Creating VHD
(default) Starting VM...
(default) Waiting for host to start...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Error creating machine: Error detecting OS: Too many retries waiting for SSH to be available.  Last error: Maximum number of retries (60) exceeded

PS> C:\Program Files\Docker\Docker> docker-machine ls
NAME      ACTIVE   DRIVER   STATE     URL                        SWARM   DOCKER    ERRORS
default   -        hyperv   Running   tcp://192.168.1.196:2376           Unknown   Unable to query docker version: Get https://192.168.1.196:2376/v1.15/version: x509: certificate is valid for 127.0.0.1, 127.0.0.1, ::1, 172.17.194.231, not 192.168.1.196

Probleme avec openssh
PS> ssh -vvv -F /dev/null -o PasswordAuthentication=no -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -o ConnectionAttempts=3 -o ConnectTimeout=10 -o ControlMaster=no -o ControlPath=none docker@192.168.1.196 -o IdentitiesOnly=yes -i /c/Users/Pascal/.docker/machine/machines/default/id_rsa -p 22

Modification des droits de id_rsa dans C:\Users\Pascal\.docker\machine\machines\default
Selectionner Priopriétés -> Sécurité -> Avancé -> Pascal -> Désactiver l'héritage

PS> C:\Program Files\Docker\Docker> docker-machine restart default
Restarting "default"...
(default) Waiting for host to stop...
(default) Waiting for host to start...
Waiting for SSH to be available...
Detecting the provisioner...
Restarted machines may have new IP addresses. You may need to re-run the `docker-machine env` command.

PS> C:\Program Files\Docker\Docker> docker-machine env
Error checking TLS connection: Error checking and/or regenerating the certs: There was an error validating certificates for host "192.168.1.196:2376": x509: certificate is valid for 127.0.0.1, 127.0.0.1, ::1, 172.17.194.231, not 192.168.1.196
You can attempt to regenerate them using 'docker-machine regenerate-certs [name]'.
Be advised that this will trigger a Docker daemon restart which might stop running containers.

PS> C:\Program Files\Docker\Docker> docker-machine regenerate-certs default
Regenerate TLS machine certs?  Warning: this is irreversible. (y/n): y
Regenerating TLS certificates
Waiting for SSH to be available...
Detecting the provisioner...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...

PS> Set-NetConnectionProfile -interfacealias "vEthernet (DockerExterne)" -NetworkCategory Private



#####Download _GoLang 1.13_ 
   https://golang.org/doc/install?download=go1.13.1.windows-amd64.msi#extra_versions

#####Download _Docker Desktop 2.1_ 
   https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe

#####Download _IntelliJ Goland 2019.2.3_ 