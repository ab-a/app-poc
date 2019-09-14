# Project
Automated deployment of a small web application easily scalable and high available using among other things **HAProxy**, **Docker**, **Ansible**, **Galera**...
### Architecture example
![alt text](https://i.imgur.com/tZRjs8z.png)
### Requirements
- 4 nodes up and running on CentOS
- ssh connection accessible without password between all nodes
- each server reachable from its custom hostname (`node-00X`)
- ansible installed on the frontend server
### Ansible actions
- install and configure docker [*frontend and backend*]
- deploy an haproxy container [*frontend*]
- deploy a strict set of iptables rules [*backend*]
- deploy a mySQL cluster [*backend*]
- bootstrap the database for the application [*backend*]
- build and launch the application in a container [*backend*]
# Before deploy
#### You need to replace some IP :
- backends servers in `haproxy/haproxy.cfg` from the line 23
- subnet server `ansible/roles/iptables/tasks/main.yml` line 27
- frontend server in `ansible/roles/iptables/tasks/main.yml` line 35 and line 44

You can replace some easily with `sed` (`$FRONTEND_IP` and `$SUBNET` variables).
#### Hostname resolution
You need that all your servers are accessible from the name configured on the host file, so don't forget to edit your `/etc/hosts`, if necessary and check if you can ping all of your hosts (`ansible all -m ping -i hosts`)
# Deployment
Simply do :
```bash
ansible-playbook -i hosts app.yml
```
You can take advantage of the Ansible tag system :
```bash
ansible-playbook -i hosts app.yml --tags app # only deploy the app
ansible-playbook -i hosts app.yml --skip-tags pre # skip docker installation
ansible-playbook -i hosts app.yml --tags database # only deploy the Galera Cluster
```
Note : Once the galera cluster is deployed you must use `--skip-tags database` !
# Scale app and database
Simply add the server in the `hosts` file. Ideally name it `node-00X` incrementaly; this allows you to modify the `hosts` file only once. Exemple if you have 3 servers for the app and 3 dedicated servers for the database :
```bash
[app_cluster]
node-[001:010]

[galera_cluster]
node-[010:013]
```
And all you have to do is to run again the ansible playbook : 
```bash
ansible-playbook -i hosts app.yml
```
# Note
- The build of the container has a problem when adding the git package, caused by a dns issue for me, fixed by set docker dns in `/etc/docker/daemon.json`
- Synchronisation problem between Galera Node fixed by changing the `wsrep_sst_method`.
- Generate your own `ssl.pem` if you need, the file include here is a generic self signed pem
- The HAProxy frontend server is the SSL termination. The idea is to make a protocol break wich is great if we assume that our backend servers are really protected and isolated from the outside like on a traditional cloud (vxlan private network for example) to be able to implement a WAF between the front server and the application server for example.
- Unlike a cloud deployment where each instance has a security group and where rules can be made at this level, we set up filtering at server level here because all our servers are public.
- Thanks to the `iptables` rules that Ansible will apply, the frontend server (`node-000`) will be the only entry point to the backend servers and will act as a bastion to access the 3 backend servers (`node-001`, `node-002`, `node-003`).
- All operations are performed on the front server through `root` user at this time.

## Resources
Galera Cluster Ansible rôle : https://github.com/adfinis-sygroup/mariadb-ansible-galera-cluster (some modifications and cleaning) \
Golang application based on : https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html

