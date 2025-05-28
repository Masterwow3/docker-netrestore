# docker-netrestore

Restores Docker HNS networks after a Windows restart, based on the Docker local-kv.db.<br/>
This plugin needs to be installed as a cli-plugin on the respective Windows server.<br/>
Important: The command must be executed before the Docker daemon starts.<br/>
I resolved the issue by setting the Docker service on Windows to start with a delay using the following command:
````cmd
sc.exe config docker start= delayed-auto
````

## 🔧 Function

- Executes `docker netrestore` to fix network problems after reboots.
- Optional: automatically registers itself in the Windows Task Scheduler (`docker netrestore --register-autorun`).

## Install
cmd:
````cmd
REM Download docker-netrestore.exe to "C:\ProgramData\docker\cli-plugins\docker-netrestore.exe"
docker netrestore --register-autorun
````

### Example Output
````
PS C:\Users\SSHAdministrator> docker netrestore version
Version: 1.5.0
CommitSHA: 69db042
````

````
PS C:\Users\SSHAdministrator> docker netrestore --register-autorun
Task created successfully.
````

````
sshadministrator@ULM-DEV-CI06 C:\Users\SSHAdministrator>docker netrestore
nat (ID: 82cb42b9d3bfbe5847c3736cac7494d9915e474a7d1a3c5068a0d35938456f72) exists: true
default_network (ID: fbe919774495beaca42fb5c5cb4f9a0b913a3e462dc3d604773a2455cbc08e99) exists: true
````

````json
PS C:\Users\SSHAdministrator> docker netrestore
nat (ID: 82cb42b9d3bfbe5847c3736cac7494d9915e474a7d1a3c5068a0d35938456f72) exists: false
Response: {"Success":true,"Output":{"ActivityId":"F7289430-273B-4B29-BF97-5DAC479AE3DC","AdditionalParams":{},"CurrentEndpointCount":0,"Extensions":[{"Id":"E7C3B2F0-F3C5-48DF-AF2B-10FED6D72E7A","IsEnabled":false,"Name":"Microsoft Windows Filtering Platform"},{"Id":"F74F241B-440F-4433-BB28-00F89EAD20D8","IsEnabled":false,"Name":"Microsoft Azure VFP Switch Extension"},{"Id":"430BDADD-BAB0-41AB-A369-94B67FA5BE0A","IsEnabled":true,"Name":"Microsoft NDIS Capture"}],"Flags":8,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875495496471},"ID":"404947AC-FC8F-4229-A4B5-595B9161F3A8","IPv6":false,"LayeredOn":"6252BE31-13AF-492C-B0CE-4FED3E16B231","MacPools":[{"EndMacAddress":"00-15-5D-51-8F-FF","StartMacAddress":"00-15-5D-51-80-00"}],"MaxConcurrentEndpoints":0,"Name":"82cb42b9d3bfbe5847c3736cac7494d9915e474a7d1a3c5068a0d35938456f72","NatName":"NAT8177E9BF-266D-4BBD-8B97-0340A82921E2","Policies":[{"Type":"VLAN","VLAN":1}],"State":1,"Subnets":[{"AdditionalParams":{},"AddressPrefix":"172.21.16.0/20","Flags":0,"GatewayAddress":"172.21.16.1","Health":{"LastErrorCode":0,"LastUpdateTime":133928875495496471},"ID":"96C9F144-2A30-414A-87C3-C39F43E4211E","IpSubnets":[{"AdditionalParams":{},"Flags":3,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875495496471},"ID":"A1BFAFB0-2E2E-428B-ADC2-CA71CE8C5DFC","IpAddressPrefix":"172.21.16.0/20","ObjectType":6,"Policies":[],"State":0}],"ObjectType":5,"Policies":[],"State":0}],"TotalEndpoints":0,"Type":"nat","Version":55834574851,"Resources":{"AdditionalParams":{},"AllocationOrder":2,"Allocators":[{"AdapterNetCfgInstanceId":"{8177E9BF-266D-4BBD-8B97-0340A82921E2}","AdditionalParams":{},"AllocationOrder":0,"CompartmendId":0,"Connected":true,"DNSFirewallRules":true,"DevicelessNic":true,"DhcpDisabled":true,"EndpointNicGuid":"337253B0-09D4-4A3F-BBDB-126ECDB6EEFF","EndpointPortGuid":"ACE214FD-4A7C-462A-AF6B-6119010C5BD2","Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875496545976},"ID":"C4700EC5-0531-4BE1-8081-F89CEC3E59CE","InterfaceGuid":"8177E9BF-266D-4BBD-8B97-0340A82921E2","IsPolicy":false,"IsolationId":1,"MacAddress":"00-15-5D-51-86-A4","ManagementPort":true,"NcfHidden":false,"NicFriendlyName":"82cb42b9d3bfbe5","NlmHidden":true,"PreferredPortFriendlyName":"Container NIC c4700ec5","State":3,"SwitchId":"6670A4DB-E650-4E9A-8C36-CFB3CE9A78DB","Tag":"Host Vnic","WaitForIpv6Interface":false,"nonPersistentPort":false},{"AdditionalParams":{},"AllocationOrder":1,"Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875497327044,"NATState":2},"ID":"E9F9170C-5407-4334-B5A0-434332BB814A","IsPolicy":false,"Prefix":20,"PrivateInterfaceGUID":"8177E9BF-266D-4BBD-8B97-0340A82921E2","State":3,"SubnetIPAddress":"172.21.16.0","Tag":"NAT"}],"CompartmentOperationTime":0,"Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875496545976},"ID":"F7289430-273B-4B29-BF97-5DAC479AE3DC","PortOperationTime":0,"State":1,"SwitchOperationTime":0,"VfpOperationTime":0,"parentId":"08DB5E96-31B1-4406-AB68-94EA6883C706"}}}        
default_network (ID: fbe919774495beaca42fb5c5cb4f9a0b913a3e462dc3d604773a2455cbc08e99) exists: false
Response: {"Success":true,"Output":{"ActivityId":"CA3B3110-78E8-4B0D-B6EA-3D7D30369322","AdditionalParams":{},"CurrentEndpointCount":0,"Extensions":[{"Id":"E7C3B2F0-F3C5-48DF-AF2B-10FED6D72E7A","IsEnabled":false,"Name":"Microsoft Windows Filtering Platform"},{"Id":"F74F241B-440F-4433-BB28-00F89EAD20D8","IsEnabled":false,"Name":"Microsoft Azure VFP Switch Extension"},{"Id":"430BDADD-BAB0-41AB-A369-94B67FA5BE0A","IsEnabled":true,"Name":"Microsoft NDIS Capture"}],"Flags":8,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875497594435},"ID":"C8026314-FC13-42CC-B41B-4DF9873B1D11","IPv6":false,"LayeredOn":"6252BE31-13AF-492C-B0CE-4FED3E16B231","MacPools":[{"EndMacAddress":"00-15-5D-EC-3F-FF","StartMacAddress":"00-15-5D-EC-30-00"}],"MaxConcurrentEndpoints":0,"Name":"fbe919774495beaca42fb5c5cb4f9a0b913a3e462dc3d604773a2455cbc08e99","NatName":"NATEBB494F4-4929-4CBA-A426-472E8CEB3E93","Policies":[{"Type":"VLAN","VLAN":2}],"State":1,"Subnets":[{"AdditionalParams":{},"AddressPrefix":"192.168.0.0/22","Flags":0,"GatewayAddress":"192.168.0.1","Health":{"LastErrorCode":0,"LastUpdateTime":133928875497594435},"ID":"FAD22C8B-3E63-463E-9938-C66A1115B0D2","IpSubnets":[{"AdditionalParams":{},"Flags":3,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875497599495},"ID":"0C82562C-1E63-46E0-A0FC-066C43DF05D3","IpAddressPrefix":"192.168.0.0/22","ObjectType":6,"Policies":[],"State":0}],"ObjectType":5,"Policies":[],"State":0}],"TotalEndpoints":0,"Type":"nat","Version":55834574851,"Resources":{"AdditionalParams":{},"AllocationOrder":2,"Allocators":[{"AdapterNetCfgInstanceId":"{EBB494F4-4929-4CBA-A426-472E8CEB3E93}","AdditionalParams":{},"AllocationOrder":0,"CompartmendId":0,"Connected":true,"DNSFirewallRules":true,"DevicelessNic":true,"DhcpDisabled":true,"EndpointNicGuid":"2ECCADDC-9043-49ED-B369-E24C7D1F5853","EndpointPortGuid":"EB8C9E76-D609-4BA3-A6C2-FDDEB8445360","Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875497600269},"ID":"50FFB17A-56D2-4862-8895-ED5CDED769B8","InterfaceGuid":"EBB494F4-4929-4CBA-A426-472E8CEB3E93","IsPolicy":false,"IsolationId":2,"MacAddress":"00-15-5D-EC-38-47","ManagementPort":true,"NcfHidden":false,"NicFriendlyName":"fbe919774495bea","NlmHidden":true,"PreferredPortFriendlyName":"Container NIC 50ffb17a","State":3,"SwitchId":"6670A4DB-E650-4E9A-8C36-CFB3CE9A78DB","Tag":"Host Vnic","WaitForIpv6Interface":false,"nonPersistentPort":false},{"AdditionalParams":{},"AllocationOrder":1,"Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875498262943,"NATState":2},"ID":"26363919-CE35-4F0C-A6D6-98601B21BEFA","IsPolicy":false,"Prefix":22,"PrivateInterfaceGUID":"EBB494F4-4929-4CBA-A426-472E8CEB3E93","State":3,"SubnetIPAddress":"192.168.0.0","Tag":"NAT"}],"CompartmentOperationTime":0,"Flags":0,"Health":{"LastErrorCode":0,"LastUpdateTime":133928875497600269},"ID":"CA3B3110-78E8-4B0D-B6EA-3D7D30369322","PortOperationTime":0,"State":1,"SwitchOperationTime":0,"VfpOperationTime":0,"parentId":"08DB5E96-31B1-4406-AB68-94EA6883C706"}}}        
````

### cli-plugin Examples
* https://github.com/docker/cli/blob/master/cli-plugins/examples/helloworld/main.go
* https://github.com/docker/compose/blob/main/cmd/main.go